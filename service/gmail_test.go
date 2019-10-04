package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
)

var (
	credentialJSON = os.Getenv("GMAIL_CREDENTIAL_JSON")
	emailAddress   = os.Getenv("GMAIL_EMAIL_ADDRESS")
)

var _ = Describe("Authorization with valid base64 CREDENTIAL_JSON", func() {

	gmail := GmailArgument{}
	requestBody := new(bytes.Buffer)
	jsonErr := json.NewEncoder(requestBody).Encode(gmail)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	request, err := http.NewRequest("POST", "/health", requestBody)
	if err != nil {
		log.Fatal(err)
	}
	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(HealthCheck)
	handler.ServeHTTP(recorder, request)

	Describe("Health Check", func() {
		Context("health check", func() {
			It("Should result http.StatusOK", func() {
				Expect(http.StatusOK).To(Equal(recorder.Code))
			})
		})
	})
})

var _ = Describe("Authorization with invalid base64 CREDENTIAL_JSON", func() {

	os.Setenv("CREDENTIAL_JSON", "mockJSON")

	gmail := GmailArgument{}
	requestBody := new(bytes.Buffer)
	jsonErr := json.NewEncoder(requestBody).Encode(gmail)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	request, err := http.NewRequest("POST", "/authorization", requestBody)
	if err != nil {
		log.Fatal(err)
	}
	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(Authorization)
	handler.ServeHTTP(recorder, request)

	Describe("Authorization", func() {
		Context("authorization", func() {
			It("Should result http.StatusBadRequest", func() {
				Expect(http.StatusBadRequest).To(Equal(recorder.Code))
			})
		})
	})
})

var _ = Describe("Authorization with invalid args", func() {

	os.Setenv("CREDENTIAL_JSON", "mockJSON")

	gmail := []byte(`{"status":false}`)
	requestBody := new(bytes.Buffer)
	jsonErr := json.NewEncoder(requestBody).Encode(gmail)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	request, err := http.NewRequest("POST", "/authorization", requestBody)
	if err != nil {
		log.Fatal(err)
	}
	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(Authorization)
	handler.ServeHTTP(recorder, request)

	Describe("Authorization", func() {
		Context("authorization", func() {
			It("Should result http.StatusBadRequest", func() {
				Expect(http.StatusBadRequest).To(Equal(recorder.Code))
			})
		})
	})
})

var _ = Describe("Authorization with valid base64 CREDENTIAL_JSON", func() {

	os.Setenv("CREDENTIAL_JSON", credentialJSON)

	gmail := GmailArgument{}
	requestBody := new(bytes.Buffer)
	jsonErr := json.NewEncoder(requestBody).Encode(gmail)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	request, err := http.NewRequest("POST", "/authorization", requestBody)
	if err != nil {
		log.Fatal(err)
	}
	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(Authorization)
	handler.ServeHTTP(recorder, request)

	Describe("Authorization", func() {
		Context("authorization", func() {
			It("Should result http.StatusOK", func() {
				Expect(http.StatusOK).To(Equal(recorder.Code))
			})
		})
	})
})

var _ = Describe("AccessToken with invalid base64 CREDENTIAL_JSON", func() {

	os.Setenv("CREDENTIAL_JSON", "mockJSON")

	gmail := GmailArgument{}
	requestBody := new(bytes.Buffer)
	jsonErr := json.NewEncoder(requestBody).Encode(gmail)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	request, err := http.NewRequest("POST", "/authorization", requestBody)
	if err != nil {
		log.Fatal(err)
	}
	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(AccessToken)
	handler.ServeHTTP(recorder, request)

	Describe("Authorization", func() {
		Context("authorization", func() {
			It("Should result http.StatusBadRequest", func() {
				Expect(http.StatusBadRequest).To(Equal(recorder.Code))
			})
		})
	})
})

var _ = Describe("AccessToken with invalid args", func() {

	os.Setenv("CREDENTIAL_JSON", "mockJSON")

	gmail := []byte(`{"status":false}`)
	requestBody := new(bytes.Buffer)
	jsonErr := json.NewEncoder(requestBody).Encode(gmail)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	request, err := http.NewRequest("POST", "/authorization", requestBody)
	if err != nil {
		log.Fatal(err)
	}
	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(AccessToken)
	handler.ServeHTTP(recorder, request)

	Describe("Authorization", func() {
		Context("authorization", func() {
			It("Should result http.StatusBadRequest", func() {
				Expect(http.StatusBadRequest).To(Equal(recorder.Code))
			})
		})
	})
})

var _ = Describe("AccessToken with invalid auth code", func() {

	os.Setenv("CREDENTIAL_JSON", credentialJSON)

	gmail := GmailArgument{AuthorizationCode: "mockAuthCode"}
	requestBody := new(bytes.Buffer)
	jsonErr := json.NewEncoder(requestBody).Encode(gmail)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	request, err := http.NewRequest("POST", "/accessToken", requestBody)
	if err != nil {
		log.Fatal(err)
	}
	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(AccessToken)
	handler.ServeHTTP(recorder, request)

	Describe("Access Token", func() {
		Context("access Token", func() {
			It("Should result http.StatusBadRequest", func() {
				Expect(http.StatusBadRequest).To(Equal(recorder.Code))
			})
		})
	})
})

var _ = Describe("Send Mail with invalid base64 CREDENTIAL_JSON", func() {

	os.Setenv("CREDENTIAL_JSON", "mockJSON")

	gmail := GmailArgument{}
	requestBody := new(bytes.Buffer)
	jsonErr := json.NewEncoder(requestBody).Encode(gmail)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	request, err := http.NewRequest("POST", "/sendMail", requestBody)
	if err != nil {
		log.Fatal(err)
	}
	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(SendMail)
	handler.ServeHTTP(recorder, request)

	Describe("Access Token", func() {
		Context("access Token", func() {
			It("Should result http.StatusBadRequest", func() {
				Expect(http.StatusBadRequest).To(Equal(recorder.Code))
			})
		})
	})
})

var _ = Describe("Access Token with invalid args", func() {

	os.Setenv("CREDENTIAL_JSON", "mockJSON")

	gmail := []byte(`{"status":false}`)
	requestBody := new(bytes.Buffer)
	jsonErr := json.NewEncoder(requestBody).Encode(gmail)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	request, err := http.NewRequest("POST", "/sendMail", requestBody)
	if err != nil {
		log.Fatal(err)
	}
	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(SendMail)
	handler.ServeHTTP(recorder, request)

	Describe("Access Token", func() {
		Context("access Token", func() {
			It("Should result http.StatusBadRequest", func() {
				Expect(http.StatusBadRequest).To(Equal(recorder.Code))
			})
		})
	})
})

var _ = Describe("Send mail with valid arg", func() {

	os.Setenv("CREDENTIAL_JSON", credentialJSON)

	toList := []string{emailAddress}
	gmail := GmailArgument{AccessToken: "mockAccessToken", UserID: emailAddress, To: toList, Subject: "Test Subject", Body: "Mail body goes here"}
	requestBody := new(bytes.Buffer)
	jsonErr := json.NewEncoder(requestBody).Encode(gmail)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	request, err := http.NewRequest("POST", "/sendMail", requestBody)
	if err != nil {
		log.Fatal(err)
	}
	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(SendMail)
	handler.ServeHTTP(recorder, request)

	Describe("Access Token", func() {
		Context("access Token", func() {
			It("Should result http.StatusBadRequest", func() {
				Expect(http.StatusBadRequest).To(Equal(recorder.Code))
			})
		})
	})
})

var _ = Describe("Subscribe gmail account for new incoming message", func() {

	os.Setenv("CREDENTIAL_JSON", "mockJSON")

	sub := []byte(`{"status":false}`)
	requestBody := new(bytes.Buffer)
	err := json.NewEncoder(requestBody).Encode(sub)
	if err != nil {
		fmt.Println(" request err :", err)
	}
	req, err := http.NewRequest("POST", "/receive", requestBody)
	if err != nil {
		log.Fatal(err)
	}
	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(ReceiveMail)
	handler.ServeHTTP(recorder, req)

	Describe("Subscribe", func() {
		Context("Subscribe", func() {
			It("Should result http.StatusBadRequest", func() {
				Expect(http.StatusBadRequest).To(Equal(recorder.Code))
			})
		})
	})
})

var _ = Describe("Subscribe gmail account for new incoming message", func() {

	os.Setenv("CREDENTIAL_JSON", credentialJSON)

	data := RequestParam{UserID: emailAddress, AccessToken: "mockAccessToken"}
	sub := Subscribe{Endpoint: "https://webhook.site/3cee781d-0a87-4966-bdec-9635436294e9",
		ID:        "1",
		IsTesting: true,
		Data:      data,
	}
	requestBody := new(bytes.Buffer)
	err := json.NewEncoder(requestBody).Encode(sub)
	if err != nil {
		fmt.Println(" request err :", err)
	}
	req, err := http.NewRequest("POST", "/receive", requestBody)
	if err != nil {
		log.Fatal(err)
	}
	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(ReceiveMail)
	handler.ServeHTTP(recorder, req)

	Describe("Subscribe", func() {
		Context("Subscribe", func() {
			It("Should result http.StatusOK", func() {
				Expect(http.StatusOK).To(Equal(recorder.Code))
			})
		})
	})
})
