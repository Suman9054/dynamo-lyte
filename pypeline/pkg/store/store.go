package store

import (
	
	"sync"
	"sync/atomic"
	"time"

	"google.golang.org/protobuf/types/known/structpb"
)

type Regmap[k comparable, v any] struct {
	Defaltv v
	M       sync.Map
	count   atomic.Int64
}
func (r *Regmap[k, v]) Get(key k) (v, bool) {
	value, ok := r.M.Load(key)
	if !ok {
		return r.Defaltv, false
	}
	return value.(v), true
}
func (r *Regmap[k, v]) Set(key k, value v) {
	r.M.Store(key, value)
	r.count.Add(1)
}

func (r *Regmap[k, v]) Count() int64 {
	return r.count.Load()
}
func (r *Regmap[k, v]) Delete(key k) {
	r.M.Delete(key)
	r.count.Add(-1)
}
func (r *Regmap[k, v]) All(f func(key k, value v) bool) {
	r.M.Range(func(key, value interface{}) bool {
		return f(key.(k), value.(v))
	})
}

func newstoremap() Stabel[string,*structpb.Struct] {
	return &Regmap[string, *structpb.Struct]{
		M: sync.Map{},
	}
}
func newexpiresmap() Stabel[*structpb.Struct, int64] {
	return &Regmap[*structpb.Struct, int64]{
		M: sync.Map{},
	}
}
type Store struct{
	store Stabel[string,*structpb.Struct]
	expires Stabel[*structpb.Struct, int64]
	numskey int
}

func Newstore() *Store{
	return &Store{
		store:newstoremap(),
		expires:newexpiresmap(),
	}
}


func (s *Store) Putdata(key string, data *structpb.Struct){
	curobj, ok := s.store.Get(key)
	if ok {
		expiretime, ok := s.expires.Get(curobj)
		if ok {
			s.expires.Set(data, expiretime)
			s.expires.Delete(curobj)
		}
	}
	
	s.store.Set(key, data)
	s.expires.Set(data, time.Now().Add(1*time.Minute).UnixNano())
	s.numskey++
}

func (s *Store) Getdata(key string) (*structpb.Struct, bool){
	
	data, ok := s.store.Get(key)
	
	if !ok {
		return nil, false
	}
	// expiretime, ok := s.expires.Get(data)
	// if !ok {
	// 	return data, true
	// }
	// if time.Now().UnixNano() > expiretime {
	// 	s.store.Delete(key)
	// 	s.expires.Delete(data)
	// 	return nil, false
	// }
	return data, true
}

func (s *Store) Expirewatcher(){
	ticker := time.NewTicker(1*time.Minute)
	for {
		<- ticker.C
		now := time.Now().UnixNano()
		s.expires.All(func(key *structpb.Struct, value int64) bool {
			if now > value {
				s.expires.Delete(key)
				// find the key in store
				s.store.All(func(k string, v *structpb.Struct) bool {
					if v == key {
						s.store.Delete(k)
						return false
					}
					return true
				})
				s.numskey--
			}
			return true
		})
	}
}