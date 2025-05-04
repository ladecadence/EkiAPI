package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/ladecadence/EkiAPI/pkg/models"
)

func ApiGetImageListMission(writer http.ResponseWriter, request *http.Request) {
	// get ID
	mission := request.PathValue("mission")
	if mission == "" {
		writer.WriteHeader(http.StatusNoContent)
		writer.Write([]byte(`{}\n`))
		return
	}

	data, err := db.GetImages(mission)

	if err != nil  || data == nil {
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

func ApiGetLastImageMission(writer http.ResponseWriter, request *http.Request) {
	// get ID
	mission := request.PathValue("mission")
	if mission == "" {
		writer.WriteHeader(http.StatusNoContent)
		writer.Write([]byte(`{}\n`))
		return
	}

	data, err := db.GetLastImage(mission)

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

func ApiUploadImage(writer http.ResponseWriter, request *http.Request) {
	// check auth
	authOk := CheckAuth(request)
	if authOk {
		// get the form data
		// 32 MB is the default used by FormFile() function
		if err := request.ParseMultipartForm(32 * 1024 * 1024); err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		// get the mission name
		mission := request.MultipartForm.Value["mission"]
		if len(mission) < 1 {
			writer.WriteHeader(http.StatusBadRequest)
			writer.Write([]byte("{}\n"))
			return
		}
		// Get a reference to the file.
		// They are accessible only after ParseMultipartForm is called
		files := request.MultipartForm.File["file"]
		if len(files) < 1 {
			writer.WriteHeader(http.StatusBadRequest)
			writer.Write([]byte("{}\n"))
			return
		}
		// Open the file
		file, err := files[0].Open()
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			writer.Write([]byte("{}\n"))
			return
		}
		defer file.Close()
		buff := make([]byte, 512)
		_, err = file.Read(buff)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			writer.Write([]byte("{}\n"))
			return
		}
		// checking the content type
		// so we don't allow files other than images
		filetype := http.DetectContentType(buff)
		if filetype != "image/jpeg" && filetype != "image/png" && filetype != "image/jpg" {
			writer.WriteHeader(http.StatusBadRequest)
			writer.Write([]byte("{}\n"))
			return
		}
		_, err = file.Seek(0, io.SeekStart)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			writer.Write([]byte("{}\n"))
			return
		}

		fileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), filepath.Ext(files[0].Filename))
		f, err := os.Create(conf.ImagePath + fileName)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			writer.Write([]byte("{}\n"))
			return
		}
		defer f.Close()
		_, err = io.Copy(f, file)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			writer.Write([]byte("{}\n"))
			return
		}

		image := models.Image{FileName: fileName, Mission: mission[0], DateTime: time.Now()}

		err = db.InsertImage(image)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			writer.Write([]byte("{}\n"))
			return
		}
		data, _ := json.Marshal(fileName)
		writer.WriteHeader(http.StatusOK)
		writer.Write(data)
		writer.Write([]byte("\n"))
	} else {
		writer.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
		http.Error(writer, "Unauthorized", http.StatusUnauthorized)
	}
}

func ApiDownloadImage(writer http.ResponseWriter, request *http.Request) {
	// get file name
	fileName := request.PathValue("name")
	if fileName == "" {
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte("{}\n"))
		return
	}

	// try to get file
	if _, err := os.Stat(conf.ImagePath + fileName); errors.Is(err, os.ErrNotExist) {
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte("{}\n"))
		return
	}
	http.ServeFile(writer, request, conf.ImagePath+fileName)
}
