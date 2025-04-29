package controllers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/ladecadence/EkiAPI/pkg/models"
)

func ApiGetMissions(writer http.ResponseWriter, request *http.Request) {
	missions, err := db.GetMissions()

	if err != nil || missions == nil {
		writer.WriteHeader(http.StatusNoContent)
		writer.Write([]byte(`{}\n`))
		return
	}

	res, _ := json.Marshal(missions)
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	writer.Write(res)
	writer.Write([]byte("\n"))
}

func ApiGetMission(writer http.ResponseWriter, request *http.Request) {
	// get ID
	name := request.PathValue("name")
	if name == "" {
		writer.WriteHeader(http.StatusNoContent)
		writer.Write([]byte(`{}\n`))
		return
	}
	mission, err := db.GetMission(name)

	if err != nil {
		writer.WriteHeader(http.StatusNoContent)
		writer.Write([]byte(`{}\n`))
		return
	}

	res, _ := json.Marshal(mission)
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	writer.Write(res)
	writer.Write([]byte("\n"))
}

func ApiNewMission(writer http.ResponseWriter, request *http.Request) {
	// check auth
	authOk := CheckAuth(request)
	if authOk {
		reqBody, _ := io.ReadAll(request.Body)
		request.Body.Close()
		// try to create new mission
		mission := models.Mission{}
		err := json.Unmarshal(reqBody, &mission)
		if err != nil {
			log.Printf("❌ Error decoding body: %v", err.Error())
			writer.WriteHeader(http.StatusBadRequest)
			writer.Write([]byte("{}\n"))
			return
		}
		err = db.InsertMission(mission)
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			writer.Write([]byte("{}\n"))
			return
		}
		data, _ := json.Marshal(mission)
		writer.WriteHeader(http.StatusOK)
		writer.Write(data)
		writer.Write([]byte("\n"))
	} else {
		writer.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
		http.Error(writer, "Unauthorized", http.StatusUnauthorized)
	}
}
