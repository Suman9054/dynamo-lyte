package store

import (
	"context"
	"encoding/json"
	"path"
	"sort"
	"strconv"
	"sync"

	"github.com/spaolacci/murmur3"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type Node struct {
	Adder    string `json:"adder"`
	Location string `json:"location"`
	Reader   bool   `json:"reader"`
	Writer   bool   `json:"writer"`
}

type Ring struct {
	sync.RWMutex
	Replicas int
	VNodes   map[uint32]Node
	Sorted   []uint32
}





func NewRing(replicas int) *Ring {
	return &Ring{
		Replicas: replicas,
		VNodes:   make(map[uint32]Node),
	}
}

func (r *Ring) addNode(nodeID string, node Node) {
	for i := 0; i < r.Replicas; i++ {
		h := vnodeHash(nodeID, i)
		r.VNodes[h] = node
		r.Sorted = append(r.Sorted, h)
	}
	sort.Slice(r.Sorted, func(i, j int) bool {
		return r.Sorted[i] < r.Sorted[j]
	})
}

func vnodeHash(nodeID string, replica int) uint32 {
	key := nodeID + "||" + strconv.Itoa(replica)
	return murmur3.Sum32([]byte(key))
}

func (r *Ring) removeNode(nodeID string) {
	for i := 0; i < r.Replicas; i++ {
		h := vnodeHash(nodeID, i)
		delete(r.VNodes, h)
	}

	r.Sorted = r.Sorted[:0]
	for h := range r.VNodes {
		r.Sorted = append(r.Sorted, h)
	}
	sort.Slice(r.Sorted, func(i, j int) bool {
		return r.Sorted[i] < r.Sorted[j]
	})
}

func WatchNodes(ctx context.Context, cli *clientv3.Client, ring *Ring) error {
	watch := cli.Watch(ctx, "/nodes/", clientv3.WithPrefix())

	for resp := range watch {
		for _, ev := range resp.Events {

			nodeID := path.Base(string(ev.Kv.Key))

			switch ev.Type {

			case clientv3.EventTypePut:
				var node Node
				if err := json.Unmarshal(ev.Kv.Value, &node); err != nil {
					continue
				}

				ring.Lock()
				ring.addNode(nodeID, node)
				ring.Unlock()

			case clientv3.EventTypeDelete:
				ring.Lock()
				ring.removeNode(nodeID)
				ring.Unlock()
			}
		}
	}

	return nil

}
