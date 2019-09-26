package service

import (
	"encoding/base64"
	"encoding/json"
	"github.com/heaptracetechnology/gmail/result"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"net/http"
	"os"
	"strings"
)

//GmailArgument struct
type GmailArgument struct {
	From        string   `json:"from"`
	To          []string `json:"to"`
	Subject     string   `json:"subject"`
	Body        string   `json:"body"`
	AccessToken string   `json:"accessToken"`
}

//AuthURL struct
type AuthURL struct {
	URL     string `json:"url"`
	Message string `json:"message"`
}

//HealthCheck Google-Sheets
func HealthCheck(responseWriter http.ResponseWriter, request *http.Request) {

	bytes, _ := json.Marshal("OK")
	result.WriteJSONResponse(responseWriter, bytes, http.StatusOK)
}

//AccessToken Gmail
func AccessToken(responseWriter http.ResponseWriter, request *http.Request) {

	var base64CredentialsJSON = os.Getenv("CREDENTIAL_JSON")

	decodedJSON, err := base64.StdEncoding.DecodeString(base64CredentialsJSON)
	if err != nil {
		result.WriteErrorResponseString(responseWriter, err.Error())
		return
	}

	conf, confErr := google.ConfigFromJSON(decodedJSON, gmail.GmailSendScope)
	if confErr != nil {
		result.WriteErrorResponseString(responseWriter, confErr.Error())
		return
	}

	url := conf.AuthCodeURL("CSRF", oauth2.AccessTypeOffline)

	respURL := AuthURL{URL: url, Message: "Copy and paste the URL and get access token from browser"}
	bytes, _ := json.Marshal(respURL)
	result.WriteJSONResponse(responseWriter, bytes, http.StatusOK)
}

//SendMail Gmail
func SendMail(responseWriter http.ResponseWriter, request *http.Request) {

	var base64CredentialsJSON = os.Getenv("CREDENTIAL_JSON")

	decodedJSON, err := base64.StdEncoding.DecodeString(base64CredentialsJSON)
	if err != nil {
		result.WriteErrorResponseString(responseWriter, err.Error())
		return
	}

	conf, confErr := google.ConfigFromJSON(decodedJSON, gmail.GmailSendScope)
	if confErr != nil {
		result.WriteErrorResponseString(responseWriter, confErr.Error())
		return
	}

	decoder := json.NewDecoder(request.Body)
	var gmailArgument GmailArgument
	reqErr := decoder.Decode(&gmailArgument)
	if reqErr != nil {
		result.WriteErrorResponseString(responseWriter, reqErr.Error())
		return
	}

	// Exchange the auth code for an access token
	token, tokErr := conf.Exchange(oauth2.NoContext, gmailArgument.AccessToken)
	if tokErr != nil {
		result.WriteErrorResponseString(responseWriter, tokErr.Error())
		return
	}

	// Create the *http.Client using the access token
	client := conf.Client(oauth2.NoContext, token)

	// Create a new gmail service using the client
	gmailService, serviceErr := gmail.New(client)
	if serviceErr != nil {
		result.WriteErrorResponseString(responseWriter, serviceErr.Error())
		return
	}

	// New message for our gmail service to send
	var message gmail.Message

	var receiverList string

	for _, element := range gmailArgument.To {
		receiverList = receiverList + element + ","
	}

	receiverList = strings.TrimSuffix(receiverList, ",")

	// Compose the message
	messageStr := []byte(
		"To: " + receiverList + "\r\n" +
			"Subject: " + gmailArgument.Subject + "\r\n\r\n" +
			gmailArgument.Body)

	// Place messageStr into message.Raw in base64 encoded format
	message.Raw = base64.URLEncoding.EncodeToString(messageStr)

	// Send the message
	send, sendErr := gmailService.Users.Messages.Send(gmailArgument.From, &message).Do()
	if sendErr != nil {
		result.WriteErrorResponseString(responseWriter, sendErr.Error())
		return
	}

	bytes, _ := json.Marshal(send)
	result.WriteJSONResponse(responseWriter, bytes, http.StatusOK)
}
