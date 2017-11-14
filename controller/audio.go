package controller

import (
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/vincent-petithory/dataurl"
)

// HandleAudioFileUpload saves an uploaded audio file on FS in "/audio". Returns 400 if
// Content Type not audio
func HandleAudioFileUpload(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	dataURL, err := dataurl.Decode(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if strings.Contains(dataURL.ContentType(), "audio/") {
		var stageDir = "audio/" + params.ByName("stageID")
		if _, err = os.Stat(stageDir); err != nil {
			if os.IsNotExist(err) {
				os.Mkdir(stageDir, 0777)
			} else {
				panic(err)
			}
		}
		ioutil.WriteFile(filepath.Join(stageDir, params.ByName("fileName")), dataURL.Data, 0644)
	} else {
		writeJSONError(w, http.StatusBadRequest, "Uploaded file is not audio")
		return
	}
	writeJSONResponse(w, http.StatusOK, "Audio file uploaded with success")
}

// GetAudioFile serves audio files
func GetAudioFile(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	http.ServeFile(w, r, "audio/"+params.ByName("stageID")+"/"+params.ByName("fileName"))
}
