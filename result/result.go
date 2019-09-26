package result

import (
	"encoding/json"
	"log"
	"net/http"
)

//WriteErrorResponse Response
func WriteErrorResponse(responseWriter http.ResponseWriter, err error) {
	messageBytes, _ := json.Marshal(err)
	WriteJSONResponse(responseWriter, messageBytes, http.StatusBadRequest)
}

//WriteErrorResponseString Response
func WriteErrorResponseString(responseWriter http.ResponseWriter, err string) {
	messageBytes, _ := json.Marshal(err)
	WriteJSONResponse(responseWriter, messageBytes, http.StatusBadRequest)
}

//WriteJSONResponse Response
func WriteJSONResponse(responseWriter http.ResponseWriter, bytes []byte, statusCode int) {
	responseWriter.Header().Set("Content-Type", "application/json; charset=UTF-8")
	responseWriter.WriteHeader(statusCode)
	_, err := responseWriter.Write(bytes)
	if err != nil {
		log.Fatal(err)
	}
}
