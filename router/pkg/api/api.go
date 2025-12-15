package api

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type DataRequest struct {
	Primarykey  string                 `json:"primarykey" validate:"required"`
	Pertisunkey string                 `json:"secondarykey"`
	TenantId    string                 `json:"tenantId" validate:"required"`
	Collection  string                 `json:"collection" validate:"required"`
	Data        map[string]interface{} `json:"data"`
}

var validate = validator.New()

func StoredataHandler(w http.ResponseWriter, r *http.Request) {
	var dataReq DataRequest

	if err := json.NewDecoder(r.Body).Decode(&dataReq); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	if err := validate.Struct(dataReq); err != nil {
		http.Error(w, "Validation failed: "+err.Error(), http.StatusBadRequest)
		return
	}
	h := sha256.New()
	h.Write([]byte(dataReq.Primarykey + dataReq.Pertisunkey))
	hashBytes := h.Sum(nil)
    hashNum := binary.BigEndian.Uint32(hashBytes[:4])
	partitionkey := fmt.Sprintf("%d", hashNum%3)
	response := map[string]string{
		"tenantId":    dataReq.TenantId,
		"collection":  dataReq.Collection,
		"partitionkey": partitionkey,
		"data":       fmt.Sprintf("%v", dataReq.Data),
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

}
