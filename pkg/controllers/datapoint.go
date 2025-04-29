package controllers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/ladecadence/EkiAPI/pkg/models"
)

func ApiNewDatapoint(writer http.ResponseWriter, request *http.Request) {
	// check auth
	authOk := CheckAuth(request)
	if authOk {
		reqBody, _ := io.ReadAll(request.Body)
		request.Body.Close()
		// try to create new Datapoint
		dp := models.Datapoint{}
		err := json.Unmarshal(reqBody, &dp)
		if err != nil {
			log.Printf("❌ Error decoding body: %v", err.Error())
			writer.WriteHeader(http.StatusBadRequest)
			writer.Write([]byte("{}\n"))
			return
		}
		err = db.InsertDatapoint(dp)
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			writer.Write([]byte("{}\n"))
			return
		}
		data, _ := json.Marshal(dp)
		writer.WriteHeader(http.StatusOK)
		writer.Write(data)
		writer.Write([]byte("\n"))
	} else {
		writer.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
		http.Error(writer, "Unauthorized", http.StatusUnauthorized)
	}
}

func ApiGetDataMission(writer http.ResponseWriter, request *http.Request) {
	// get ID
	name := request.PathValue("mission")
	if name == "" {
		writer.WriteHeader(http.StatusNoContent)
		writer.Write([]byte(`{}\n`))
		return
	}

	data, err := db.GetMissionData(name)

	if err != nil {
		writer.WriteHeader(http.StatusNoContent)
		writer.Write([]byte(`{}\n`))
		return
	}

	res, _ := json.Marshal(data)
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	writer.Write(res)
	writer.Write([]byte("\n"))
}
