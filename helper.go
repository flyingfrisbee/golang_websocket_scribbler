package main

import (
	"encoding/json"
	"net/http"
)

func SendAppVersionResponse(w http.ResponseWriter, r *http.Request, httpstatus int) {
	resp := AppVersion{
		MandatoryVersion: AppVersionMandatory,
		OptionalVersion:  AppVersionOptional,
	}

	jsonBytes, _ := json.Marshal(resp)

	w.Header().Add("Content-type", "application/json")
	w.WriteHeader(httpstatus)
	w.Write(jsonBytes)
}

func SendResponse(w http.ResponseWriter, r *http.Request, success bool, httpstatus int) {
	resp := Resp{
		Success: success,
	}

	jsonBytes, _ := json.Marshal(resp)

	w.Header().Add("Content-type", "application/json")
	w.WriteHeader(httpstatus)
	w.Write(jsonBytes)
}
