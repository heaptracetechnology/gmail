package service

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/cloudevents/sdk-go"
	"github.com/heaptracetechnology/gmail/result"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

//MailContent struct
type MailContent struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

//GmailArgument struct
type GmailArgument struct {
	UserID                string   `json:"userId"`
	To                    []string `json:"to"`
	Subject               string   `json:"subject,omitempty"`
	Body                  string   `json:"body"`
	AuthorizationCode     string   `json:"authorizationCode"`
	AccessToken           string   `json:"accessToken"`
	TokenObj              Token    `json:"token"`
	LabelName             string   `json:"name"`
	BackgroundColor       string   `json:"backgroundColor,omitempty"`
	TextColor             string   `json:"textColor,omitempty"`
	LabelListVisibility   string   `json:"labelListVisibility,omitempty"`
	MessageListVisibility string   `json:"messageListVisibility,omitempty"`
	LabelID               string   `json:"labelID"`
	Scope                 string   `json:"scope"`
}

//GmailFilter struct
type GmailFilter struct {
	AccessToken    string `json:"accessToken"`
	UserID         string `json:"userId"`
	FilterID       string `json:"filterID"`
	RemoveLabelID  string `json:"removeLabelId,omitempty"`
	AddLabelID     string `json:"addLabelId,omitempty"`
	ExcludeChats   bool   `json:"excludeChats,omitempty"`
	From           string `json:"from,omitempty"`
	To             string `json:"to,omitempty"`
	Subject        string `json:"subject,omitempty"`
	HasAttachment  bool   `json:"hasAttachment,omitempty"`
	NegatedQuery   string `json:"negatedQuery,omitempty"`
	Query          string `json:"query,omitempty"`
	Size           int64  `json:"size,omitempty"`
	SizeComparison string `json:"sizeComparison,omitempty"`
}

//Token struct
type Token struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Expiry       string `json:"expiry"`
	TokenType    string `json:"token_type"`
}

//AuthURL struct
type AuthURL struct {
	URL     string `json:"url"`
	Message string `json:"message"`
}

//Subscribe struct
type Subscribe struct {
	Data      RequestParam `json:"data"`
	Endpoint  string       `json:"endpoint"`
	ID        string       `json:"id"`
	IsTesting bool         `json:"isTesting"`
}

//RequestParam struct
type RequestParam struct {
	AccessToken string `json:"accessToken"`
	UserID      string `json:"UserId"`
}

//Message struct
type Message struct {
	Success    bool   `json:"success"`
	Message    string `json:"message"`
	StatusCode int    `json:"statusCode"`
}

//Global Variables
var (
	Listener         = make(map[string]Subscribe)
	rtmStarted       bool
	gmailService     *gmail.Service
	currentMessageID string
	oldMessageID     string
)

//HealthCheck Google-Sheets
func HealthCheck(responseWriter http.ResponseWriter, request *http.Request) {

	bytes, _ := json.Marshal("OK")
	result.WriteJSONResponse(responseWriter, bytes, http.StatusOK)
}

//Authorization Gmail
func Authorization(responseWriter http.ResponseWriter, request *http.Request) {

	var base64CredentialsJSON = os.Getenv("CREDENTIAL_JSON")

	decodedJSON, err := base64.StdEncoding.DecodeString(base64CredentialsJSON)
	if err != nil {
		result.WriteErrorResponseString(responseWriter, err.Error())
		return
	}

	decoder := json.NewDecoder(request.Body)
	var gmailArgument GmailArgument
	reqErr := decoder.Decode(&gmailArgument)
	if reqErr != nil {
		result.WriteErrorResponseString(responseWriter, reqErr.Error())
		return
	}

	fmt.Println("gmailArgument.Scope :::", gmailArgument.Scope)
	conf, confErr := google.ConfigFromJSON(decodedJSON, gmailArgument.Scope)
	if confErr != nil {
		result.WriteErrorResponseString(responseWriter, confErr.Error())
		return
	}

	url := conf.AuthCodeURL("CSRF", oauth2.AccessTypeOffline)

	respURL := AuthURL{URL: url, Message: "Copy and paste the URL and get authorization code from browser"}
	bytes, _ := json.Marshal(respURL)
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

	decoder := json.NewDecoder(request.Body)
	var gmailArgument GmailArgument
	reqErr := decoder.Decode(&gmailArgument)
	if reqErr != nil {
		result.WriteErrorResponseString(responseWriter, reqErr.Error())
		return
	}

	conf, confErr := google.ConfigFromJSON(decodedJSON, gmail.MailGoogleComScope)
	if confErr != nil {
		result.WriteErrorResponseString(responseWriter, confErr.Error())
		return
	}

	// Exchange the auth code for an access token
	token, tokErr := conf.Exchange(oauth2.NoContext, gmailArgument.AuthorizationCode)
	if tokErr != nil {
		result.WriteErrorResponseString(responseWriter, tokErr.Error())
		return
	}

	bytes, _ := json.Marshal(token)
	result.WriteJSONResponse(responseWriter, bytes, http.StatusOK)
}

//RefreshToken Gmail
func RefreshToken(responseWriter http.ResponseWriter, request *http.Request) {

	var base64CredentialsJSON = os.Getenv("CREDENTIAL_JSON")

	decodedJSON, err := base64.StdEncoding.DecodeString(base64CredentialsJSON)
	if err != nil {
		result.WriteErrorResponseString(responseWriter, err.Error())
		return
	}

	decoder := json.NewDecoder(request.Body)
	var gmailArgument GmailArgument
	reqErr := decoder.Decode(&gmailArgument)
	if reqErr != nil {
		result.WriteErrorResponseString(responseWriter, reqErr.Error())
		return
	}

	conf, confErr := google.ConfigFromJSON(decodedJSON, gmail.MailGoogleComScope)
	if confErr != nil {
		result.WriteErrorResponseString(responseWriter, confErr.Error())
		return
	}

	expTime, _ := time.Parse(time.RFC3339, gmailArgument.TokenObj.Expiry)

	tok := oauth2.Token{
		AccessToken:  gmailArgument.TokenObj.AccessToken,
		RefreshToken: gmailArgument.TokenObj.RefreshToken,
		Expiry:       expTime,
		TokenType:    gmailArgument.TokenObj.TokenType,
	}

	updatedToken, updatedTokenErr := conf.TokenSource(context.TODO(), &tok).Token()
	if updatedTokenErr != nil {
		result.WriteErrorResponseString(responseWriter, updatedTokenErr.Error())
		return
	}

	bytes, _ := json.Marshal(updatedToken)
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

	conf, confErr := google.ConfigFromJSON(decodedJSON, gmail.MailGoogleComScope)
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

	token := oauth2.Token{
		AccessToken: gmailArgument.AccessToken,
	}

	// Create the *http.Client using the access token
	client := conf.Client(oauth2.NoContext, &token)

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
	sendMail, sendErr := gmailService.Users.Messages.Send(gmailArgument.UserID, &message).Do()
	if sendErr != nil {
		result.WriteErrorResponseString(responseWriter, sendErr.Error())
		return
	}

	bytes, _ := json.Marshal(sendMail)
	result.WriteJSONResponse(responseWriter, bytes, http.StatusOK)
}

//CreateLabel Gmail
func CreateLabel(responseWriter http.ResponseWriter, request *http.Request) {

	var base64CredentialsJSON = os.Getenv("CREDENTIAL_JSON")

	decodedJSON, err := base64.StdEncoding.DecodeString(base64CredentialsJSON)
	if err != nil {
		result.WriteErrorResponseString(responseWriter, err.Error())
		return
	}

	conf, confErr := google.ConfigFromJSON(decodedJSON, gmail.MailGoogleComScope)
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

	token := oauth2.Token{
		AccessToken: gmailArgument.AccessToken,
	}

	// Create the *http.Client using the access token
	client := conf.Client(oauth2.NoContext, &token)

	// Create a new gmail service using the client
	gmailService, serviceErr := gmail.New(client)
	if serviceErr != nil {
		result.WriteErrorResponseString(responseWriter, serviceErr.Error())
		return
	}

	labelprop := gmail.Label{
		Name:                  gmailArgument.LabelName,
		LabelListVisibility:   gmailArgument.LabelListVisibility,
		MessageListVisibility: gmailArgument.MessageListVisibility,
		Color: &gmail.LabelColor{
			BackgroundColor: gmailArgument.BackgroundColor,
			TextColor:       gmailArgument.TextColor,
		},
	}

	// Create new label
	newLabel, newLabelErr := gmailService.Users.Labels.Create(gmailArgument.UserID, &labelprop).Do()
	if newLabelErr != nil {
		result.WriteErrorResponseString(responseWriter, newLabelErr.Error())
		return
	}

	bytes, _ := json.Marshal(newLabel)
	result.WriteJSONResponse(responseWriter, bytes, http.StatusOK)
}

//PatchLabel Gmail
func PatchLabel(responseWriter http.ResponseWriter, request *http.Request) {

	var base64CredentialsJSON = os.Getenv("CREDENTIAL_JSON")

	decodedJSON, err := base64.StdEncoding.DecodeString(base64CredentialsJSON)
	if err != nil {
		result.WriteErrorResponseString(responseWriter, err.Error())
		return
	}

	conf, confErr := google.ConfigFromJSON(decodedJSON, gmail.MailGoogleComScope)
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

	token := oauth2.Token{
		AccessToken: gmailArgument.AccessToken,
	}

	// Create the *http.Client using the access token
	client := conf.Client(oauth2.NoContext, &token)

	// Create a new gmail service using the client
	gmailService, serviceErr := gmail.New(client)
	if serviceErr != nil {
		result.WriteErrorResponseString(responseWriter, serviceErr.Error())
		return
	}

	labelprop := gmail.Label{
		Id:                    gmailArgument.LabelID,
		Name:                  gmailArgument.LabelName,
		LabelListVisibility:   gmailArgument.LabelListVisibility,
		MessageListVisibility: gmailArgument.MessageListVisibility,
		Color: &gmail.LabelColor{
			BackgroundColor: gmailArgument.BackgroundColor,
			TextColor:       gmailArgument.TextColor,
		},
	}

	//Patch label
	patchLabel, patchLabelErr := gmailService.Users.Labels.Patch(gmailArgument.UserID, gmailArgument.LabelID, &labelprop).Do()
	if patchLabelErr != nil {
		result.WriteErrorResponseString(responseWriter, patchLabelErr.Error())
		return
	}

	bytes, _ := json.Marshal(patchLabel)
	result.WriteJSONResponse(responseWriter, bytes, http.StatusOK)
}

//ListLabel Gmail
func ListLabel(responseWriter http.ResponseWriter, request *http.Request) {

	var base64CredentialsJSON = os.Getenv("CREDENTIAL_JSON")

	decodedJSON, err := base64.StdEncoding.DecodeString(base64CredentialsJSON)
	if err != nil {
		result.WriteErrorResponseString(responseWriter, err.Error())
		return
	}

	conf, confErr := google.ConfigFromJSON(decodedJSON, gmail.MailGoogleComScope)
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

	token := oauth2.Token{
		AccessToken: gmailArgument.AccessToken,
	}

	// Create the *http.Client using the access token
	client := conf.Client(oauth2.NoContext, &token)

	// Create a new gmail service using the client
	gmailService, serviceErr := gmail.New(client)
	if serviceErr != nil {
		result.WriteErrorResponseString(responseWriter, serviceErr.Error())
		return
	}

	//List all label from mailbox
	listLabel, listLabelErr := gmailService.Users.Labels.List(gmailArgument.UserID).Do()
	if listLabelErr != nil {
		result.WriteErrorResponseString(responseWriter, listLabelErr.Error())
		return
	}

	bytes, _ := json.Marshal(listLabel)
	result.WriteJSONResponse(responseWriter, bytes, http.StatusOK)
}

//DeleteLabel Gmail
func DeleteLabel(responseWriter http.ResponseWriter, request *http.Request) {

	var base64CredentialsJSON = os.Getenv("CREDENTIAL_JSON")

	decodedJSON, err := base64.StdEncoding.DecodeString(base64CredentialsJSON)
	if err != nil {
		result.WriteErrorResponseString(responseWriter, err.Error())
		return
	}

	conf, confErr := google.ConfigFromJSON(decodedJSON, gmail.MailGoogleComScope)
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

	token := oauth2.Token{
		AccessToken: gmailArgument.AccessToken,
	}

	// Create the *http.Client using the access token
	client := conf.Client(oauth2.NoContext, &token)

	// Create a new gmail service using the client
	gmailService, serviceErr := gmail.New(client)
	if serviceErr != nil {
		result.WriteErrorResponseString(responseWriter, serviceErr.Error())
		return
	}

	//Delete label
	deleteLabel := gmailService.Users.Labels.Delete(gmailArgument.UserID, gmailArgument.LabelID)
	deleteLabelErr := deleteLabel.Do()
	if deleteLabelErr != nil {
		result.WriteErrorResponseString(responseWriter, deleteLabelErr.Error())
		return
	}

	message := Message{Success: true, Message: "Label deleted successfully", StatusCode: http.StatusOK}
	bytes, _ := json.Marshal(message)
	result.WriteJSONResponse(responseWriter, bytes, http.StatusOK)
}

//CreateFilter Gmail
func CreateFilter(responseWriter http.ResponseWriter, request *http.Request) {

	var base64CredentialsJSON = os.Getenv("CREDENTIAL_JSON")

	decodedJSON, err := base64.StdEncoding.DecodeString(base64CredentialsJSON)
	if err != nil {
		result.WriteErrorResponseString(responseWriter, err.Error())
		return
	}

	conf, confErr := google.ConfigFromJSON(decodedJSON, gmail.GmailSettingsBasicScope)
	if confErr != nil {
		result.WriteErrorResponseString(responseWriter, confErr.Error())
		return
	}

	decoder := json.NewDecoder(request.Body)
	var gmailFilter GmailFilter
	reqErr := decoder.Decode(&gmailFilter)
	if reqErr != nil {
		result.WriteErrorResponseString(responseWriter, reqErr.Error())
		return
	}

	token := oauth2.Token{
		AccessToken: gmailFilter.AccessToken,
	}

	// Create the *http.Client using the access token
	client := conf.Client(oauth2.NoContext, &token)

	// Create a new gmail service using the client
	gmailService, serviceErr := gmail.New(client)
	if serviceErr != nil {
		result.WriteErrorResponseString(responseWriter, serviceErr.Error())
		return
	}

	filterprops := gmail.Filter{
		Action: &gmail.FilterAction{
			AddLabelIds:    []string{gmailFilter.AddLabelID},
			RemoveLabelIds: []string{gmailFilter.RemoveLabelID},
		},
		Criteria: &gmail.FilterCriteria{
			ExcludeChats:   gmailFilter.ExcludeChats,
			From:           gmailFilter.From,
			HasAttachment:  gmailFilter.HasAttachment,
			NegatedQuery:   gmailFilter.NegatedQuery,
			Query:          gmailFilter.Query,
			Size:           gmailFilter.Size,
			SizeComparison: gmailFilter.SizeComparison,
			Subject:        gmailFilter.Subject,
			To:             gmailFilter.To,
		},
	}

	//Create filter form setting
	filter, filterErr := gmailService.Users.Settings.Filters.Create(gmailFilter.UserID, &filterprops).Do()
	if filterErr != nil {
		result.WriteErrorResponseString(responseWriter, filterErr.Error())
		return
	}

	bytes, _ := json.Marshal(filter)
	result.WriteJSONResponse(responseWriter, bytes, http.StatusOK)
}

//DeleteFilter Gmail
func DeleteFilter(responseWriter http.ResponseWriter, request *http.Request) {

	var base64CredentialsJSON = os.Getenv("CREDENTIAL_JSON")

	decodedJSON, err := base64.StdEncoding.DecodeString(base64CredentialsJSON)
	if err != nil {
		result.WriteErrorResponseString(responseWriter, err.Error())
		return
	}

	conf, confErr := google.ConfigFromJSON(decodedJSON, gmail.GmailSettingsBasicScope)
	if confErr != nil {
		result.WriteErrorResponseString(responseWriter, confErr.Error())
		return
	}

	decoder := json.NewDecoder(request.Body)
	var gmailFilter GmailFilter
	reqErr := decoder.Decode(&gmailFilter)
	if reqErr != nil {
		result.WriteErrorResponseString(responseWriter, reqErr.Error())
		return
	}

	token := oauth2.Token{
		AccessToken: gmailFilter.AccessToken,
	}

	// Create the *http.Client using the access token
	client := conf.Client(oauth2.NoContext, &token)

	// Create a new gmail service using the client
	gmailService, serviceErr := gmail.New(client)
	if serviceErr != nil {
		result.WriteErrorResponseString(responseWriter, serviceErr.Error())
		return
	}

	//Delete filter form setting
	deleteFilter := gmailService.Users.Settings.Filters.Delete(gmailFilter.UserID, gmailFilter.FilterID)
	deleteFilterErr := deleteFilter.Do()
	if deleteFilterErr != nil {
		result.WriteErrorResponseString(responseWriter, deleteFilterErr.Error())
		return
	}

	message := Message{Success: true, Message: "Filter deleted successfully", StatusCode: http.StatusOK}
	bytes, _ := json.Marshal(message)
	result.WriteJSONResponse(responseWriter, bytes, http.StatusOK)
}

//FilterList Gmail
func FilterList(responseWriter http.ResponseWriter, request *http.Request) {

	var base64CredentialsJSON = os.Getenv("CREDENTIAL_JSON")

	decodedJSON, err := base64.StdEncoding.DecodeString(base64CredentialsJSON)
	if err != nil {
		result.WriteErrorResponseString(responseWriter, err.Error())
		return
	}

	conf, confErr := google.ConfigFromJSON(decodedJSON, gmail.GmailSettingsBasicScope)
	if confErr != nil {
		result.WriteErrorResponseString(responseWriter, confErr.Error())
		return
	}

	decoder := json.NewDecoder(request.Body)
	var gmailFilter GmailFilter
	reqErr := decoder.Decode(&gmailFilter)
	if reqErr != nil {
		result.WriteErrorResponseString(responseWriter, reqErr.Error())
		return
	}

	token := oauth2.Token{
		AccessToken: gmailFilter.AccessToken,
	}

	// Create the *http.Client using the access token
	client := conf.Client(oauth2.NoContext, &token)

	// Create a new gmail service using the client
	gmailService, serviceErr := gmail.New(client)
	if serviceErr != nil {
		result.WriteErrorResponseString(responseWriter, serviceErr.Error())
		return
	}

	//List all filter form setting
	filterList, filterListErr := gmailService.Users.Settings.Filters.List(gmailFilter.UserID).Do()
	if filterListErr != nil {
		result.WriteErrorResponseString(responseWriter, filterListErr.Error())
		return
	}

	bytes, _ := json.Marshal(filterList)
	result.WriteJSONResponse(responseWriter, bytes, http.StatusOK)
}

//ReceiveMail Gmail
func ReceiveMail(responseWriter http.ResponseWriter, request *http.Request) {

	var key = os.Getenv("CREDENTIAL_JSON")

	decodedJSON, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		result.WriteErrorResponseString(responseWriter, err.Error())
		return
	}

	decoder := json.NewDecoder(request.Body)

	var sub Subscribe
	decodeError := decoder.Decode(&sub)
	if decodeError != nil {
		result.WriteErrorResponseString(responseWriter, decodeError.Error())
		return
	}

	conf, confErr := google.ConfigFromJSON(decodedJSON, gmail.MailGoogleComScope)
	if confErr != nil {
		result.WriteErrorResponseString(responseWriter, confErr.Error())
		return
	}

	token := oauth2.Token{
		AccessToken: sub.Data.AccessToken,
	}

	// Create the *http.Client using the access token
	client := conf.Client(oauth2.NoContext, &token)

	var serviceErr error

	// Create a new gmail service using the client
	gmailService, serviceErr = gmail.New(client)
	if serviceErr != nil {
		result.WriteErrorResponseString(responseWriter, serviceErr.Error())
		return
	}

	Listener[sub.Data.UserID] = sub
	if !rtmStarted {
		go GmailRTM()
		rtmStarted = true
	}

	bytes, _ := json.Marshal("Subscribed")
	result.WriteJSONResponse(responseWriter, bytes, http.StatusOK)
}

//GmailRTM func
func GmailRTM() {
	isTest := false
	for {
		if len(Listener) > 0 {
			for k, v := range Listener {
				go getNewEmail(k, v)
				isTest = v.IsTesting
			}
		} else {
			rtmStarted = false
			break
		}
		time.Sleep(10 * time.Second)
		if isTest {
			break
		}
	}
}

//getNewEmail func
func getNewEmail(userID string, sub Subscribe) {

	messageList, listErr := gmailService.Users.Messages.List(userID).Do()
	if listErr != nil {
		fmt.Println("Retrieve message list error: ", listErr)
		return
	}

	var message *gmail.Message
	var messageErr error
	var mailContent MailContent

	for _, element := range messageList.Messages {

		currentMessageID = element.Id
		message, messageErr = gmailService.Users.Messages.Get(userID, currentMessageID).Do()
		if messageErr != nil {
			fmt.Println("get message error : ", messageErr)
		}

		_, found := Find(message.LabelIds, "INBOX")
		if found {

			msgHeader := message.Payload.Headers

			for _, i := range msgHeader {
				if i.Name == "From" {
					mailContent.From = i.Value
				} else if i.Name == "To" {
					mailContent.To = i.Value
				} else if i.Name == "Subject" {
					mailContent.Subject = i.Value
				}
			}
			mailContent.Body = message.Snippet
			break
		}
	}

	contentType := "application/json"

	transport, err := cloudevents.NewHTTPTransport(cloudevents.WithTarget(sub.Endpoint), cloudevents.WithStructuredEncoding())
	if err != nil {
		fmt.Println("failed to create transport : ", err)
		return
	}

	client, err := cloudevents.NewClient(transport, cloudevents.WithTimeNow())
	if err != nil {
		fmt.Println("failed to create client : ", err)
		return
	}

	source, err := url.Parse(sub.Endpoint)
	event := cloudevents.Event{
		Context: cloudevents.EventContextV01{
			EventID:     sub.ID,
			EventType:   "mail",
			Source:      cloudevents.URLRef{URL: *source},
			ContentType: &contentType,
		}.AsV01(),
		Data: mailContent,
	}

	if oldMessageID != "" && currentMessageID != oldMessageID {
		oldMessageID = currentMessageID
		_, resp, err := client.Send(context.Background(), event)
		if err != nil {
			log.Printf("failed to send: %v", err)
		}
		fmt.Printf("Response: \n%s\n", resp)
	} else if oldMessageID == "" {
		oldMessageID = currentMessageID
	}
}

//Find string
func Find(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}
