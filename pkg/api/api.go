package api

import (
	"encoding/json"
	"net/http"
	"github.com/go-playground/validator/v10"
)

type DataRequest struct {
	Primarykey string `json:"primarykey" validate:"required"`
	Pertisunkey string `json:"secondarykey"`
	Data map[string]interface{} `json:"data"`
}

var validate = validator.New()

func StoredataHandler(w http.ResponseWriter, r *http.Request) {
     var dataReq DataRequest

	 if err:= json.NewDecoder(r.Body).Decode(&dataReq); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	 }
	if err:= validate.Struct(dataReq); err != nil {
		http.Error(w, "Validation failed: "+err.Error(), http.StatusBadRequest)
		return
	} 

  
	
}