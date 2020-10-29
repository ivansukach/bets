package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"mime/multipart"
	"net/http/httptest"
	"strings"
	"sync"
	"time"

	"github.com/ivansukach/bets/internal/handlers/blocking"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"testing"
)

type ErrorResponse struct {
	Status string `json:"status"`
	Error  string `json:"error"`
}
type MessageResponse struct {
	Message string `json:"message"`
}

func TestBlockUsersAndBlockedReport(t *testing.T) {
	srv := httptest.NewServer(GetRouter())
	defer srv.Close()

	//BlockUsers
	rand.Seed(time.Now().UnixNano())
	ids := []int64{int64(rand.Intn(10000)), int64(rand.Intn(17657)), 123}
	requestBody, err := json.Marshal(blocking.BlockUsersReqModel{
		Ids: ids,
	})
	if err != nil {
		log.Error(err)
		return
	}
	log.Debug("Body of request", string(requestBody))
	responseBlockingUsers, err := http.Post(srv.URL+"/block-users",
		"application/json",
		bytes.NewBuffer(requestBody))
	defer func() {
		if responseBlockingUsers != nil {
			responseBlockingUsers.Body.Close()
		}
	}()
	if err != nil {
		log.Error(err)
		return
	}
	bodyBlockingUsers, err := ioutil.ReadAll(responseBlockingUsers.Body)
	if err != nil {
		log.Error(err)
		return
	}
	log.Println(string(bodyBlockingUsers))

	//BlockedReport
	responseBlockedReport, err := http.Get(srv.URL + "/blocked-report")
	defer func() {
		if responseBlockedReport != nil {
			responseBlockedReport.Body.Close()
		}
	}()
	if err != nil {
		log.Error(err)
		return
	}
	boundary := strings.Split(responseBlockedReport.Header.Get("Content-Type"), "boundary=")[1]
	bodyBlockedReport, err := ioutil.ReadAll(responseBlockedReport.Body)
	if err != nil {
		log.Error(err)
		return
	}
	idsOfBlockedUsers := strings.Split(strings.Split(string(bodyBlockedReport), "--"+boundary)[1], "--"+boundary+"--")[0]
	for i := range ids {
		if !strings.Contains(idsOfBlockedUsers, fmt.Sprintf("%d", ids[i])) {
			log.Error("Report does not contain ids, that have just been sent to the server")
			return
		}

	}
	log.Println(string(bodyBlockedReport))
}
func TestReport(t *testing.T) {
	srv := httptest.NewServer(GetRouter())
	defer srv.Close()

	//Report
	responseReport, err := http.Get(srv.URL + "/report")
	defer responseReport.Body.Close()
	if err != nil {
		log.Error(err)
		return
	}
	bodyReport, err := ioutil.ReadAll(responseReport.Body)
	if err != nil {
		log.Error(err)
		return
	}
	log.Println(string(bodyReport))
}

func TestProcessing(t *testing.T) {
	srv := httptest.NewServer(GetRouter())
	defer srv.Close()

	//Processing
	var requestBodyProcessing bytes.Buffer
	multiPartWriter := multipart.NewWriter(&requestBodyProcessing)
	rand.Seed(time.Now().UnixNano())
	content := GenerateRandomContentOfFile(rand.Intn(30))

	requestBodyProcessing.Reset()
	part, err := multiPartWriter.CreateFormFile("file", "example.csv")
	if err != nil {
		log.Error(err)
		return
	}
	_, err = part.Write([]byte(content))
	if err != nil {
		log.Error(err)
		return
	}
	contentType := multiPartWriter.FormDataContentType()
	log.Debug("FormDataContentType: ", multiPartWriter.FormDataContentType())
	log.Debug("Body of request: ", requestBodyProcessing.String())
	multiPartWriter.Close()
	responseProcessing, err := http.Post(srv.URL+"/process",
		contentType, &requestBodyProcessing)
	defer responseProcessing.Body.Close()
	if err != nil {
		log.Error(err)
		return
	}
	bodyProcessing, err := ioutil.ReadAll(responseProcessing.Body)
	if err != nil {
		log.Error(err)
		return
	}
	log.Println(string(bodyProcessing))
}
func TestProcessingAndReport(t *testing.T) {
	srv := httptest.NewServer(GetRouter())
	defer srv.Close()

	//Processing
	var requestBodyProcessing bytes.Buffer
	multiPartWriter := multipart.NewWriter(&requestBodyProcessing)
	rand.Seed(time.Now().UnixNano())
	content := GenerateRandomContentOfFile(rand.Intn(30))

	requestBodyProcessing.Reset()
	part, err := multiPartWriter.CreateFormFile("file", "example.csv")
	if err != nil {
		log.Error(err)
		return
	}
	_, err = part.Write([]byte(content))
	if err != nil {
		log.Error(err)
		return
	}
	contentType := multiPartWriter.FormDataContentType()
	log.Debug("FormDataContentType: ", multiPartWriter.FormDataContentType())
	log.Debug("Body of request: ", requestBodyProcessing.String())
	multiPartWriter.Close()
	responseProcessing, err := http.Post(srv.URL+"/process",
		contentType, &requestBodyProcessing)
	defer responseProcessing.Body.Close()
	if err != nil {
		log.Error(err)
		return
	}
	bodyProcessing, err := ioutil.ReadAll(responseProcessing.Body)
	if err != nil {
		log.Error(err)
		return
	}
	log.Println(string(bodyProcessing))

	//Report
	responseReport, err := http.Get(srv.URL + "/report")
	defer responseReport.Body.Close()
	if err != nil {
		log.Error(err)
		return
	}
	bodyReport, err := ioutil.ReadAll(responseReport.Body)
	if err != nil {
		log.Error(err)
		return
	}
	requestContent := strings.Split(content, "\n")
	for i := range requestContent {
		values := strings.Split(requestContent[i], ", ")
		id := ""
		name := ""
		fmt.Sscanf(values[0], "%d", &id)
		fmt.Sscanf(values[1], "%s", &name)
		if !strings.Contains(string(bodyReport), id+", "+name) {
			log.Error("We should get in our report data, that we have entered to processing")
			return
		}
	}
	log.Println(string(bodyReport))
}
func GenerateRandomContentOfFile(numOfBets int) string {
	idsCollection := []int64{123, 126, 128, 131, 10012, 1, 6, 118, 8567, 8732}
	namesCollection := []string{"Petr", "Vasiliy", "Andrey", "Ivan", "Egor",
		"Valentin", "Denis", "Vadim", "Alena", "Kate"}
	coefficientsCollection := []float64{1.35, 1.42, 1.57, 5.8, 3.11, 1.14, 4.77, 2.53, 1.99, 1.87}
	amountCollection := []float64{200.2, 106.6, 20.43, 51.8, 91.1, 1.12, 47.77, 25.3, 16.99, 10.87}
	content := ""
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < numOfBets; i++ {
		idAndNameIndex := rand.Intn(10)
		content += fmt.Sprintf("%d, %s, %f, %f\n", idsCollection[idAndNameIndex], namesCollection[idAndNameIndex],
			amountCollection[rand.Intn(10)], coefficientsCollection[rand.Intn(10)])
	}
	return content
}
func TestProcessingInvalidBody(t *testing.T) {
	srv := httptest.NewServer(GetRouter())
	defer srv.Close()

	//Processing
	var requestBodyProcessing bytes.Buffer
	multiPartWriter := multipart.NewWriter(&requestBodyProcessing)
	contentType := multiPartWriter.FormDataContentType()
	log.Debug("FormDataContentType: ", multiPartWriter.FormDataContentType())
	log.Debug("Body of request: ", requestBodyProcessing.String())
	multiPartWriter.Close()
	responseProcessing, err := http.Post(srv.URL+"/process",
		contentType, &requestBodyProcessing)
	defer responseProcessing.Body.Close()
	if err != nil {
		log.Error(err)
		return
	}
	bodyProcessing, err := ioutil.ReadAll(responseProcessing.Body)
	if err != nil {
		log.Error(err)
		return
	}
	errorResp := ErrorResponse{}
	json.Unmarshal(bodyProcessing, &errorResp)
	if len(errorResp.Error) == 0 {
		log.Error("We send request with empty body. We should get error")
		return
	}
	log.Println(string(bodyProcessing))
}
func TestProcessingEmptyRequest(t *testing.T) {
	srv := httptest.NewServer(GetRouter())
	defer srv.Close()

	//Processing
	var requestBodyProcessing bytes.Buffer
	responseProcessing, err := http.Post(srv.URL+"/process",
		"application/json", &requestBodyProcessing)
	defer responseProcessing.Body.Close()
	if err != nil {
		log.Error(err)
		return
	}
	bodyProcessing, err := ioutil.ReadAll(responseProcessing.Body)
	if err != nil {
		log.Error(err)
		return
	}
	errorResp := ErrorResponse{}
	json.Unmarshal(bodyProcessing, &errorResp)
	if len(errorResp.Error) == 0 {
		log.Error("We send request with empty body. We should get error")
		return
	}
	log.Println(string(bodyProcessing))
}
func TestProcessingAndReportSimultaneously(t *testing.T) {
	srv := httptest.NewServer(GetRouter())
	defer srv.Close()

	//Processing
	var requestBodyProcessing bytes.Buffer
	multiPartWriter := multipart.NewWriter(&requestBodyProcessing)
	rand.Seed(time.Now().UnixNano())
	content := GenerateRandomContentOfFile(rand.Intn(300))

	requestBodyProcessing.Reset()
	part, err := multiPartWriter.CreateFormFile("file", "example.csv")
	if err != nil {
		log.Error(err)
		return
	}
	_, err = part.Write([]byte(content))
	if err != nil {
		log.Error(err)
		return
	}
	contentType := multiPartWriter.FormDataContentType()
	//log.Debug("FormDataContentType: ", multiPartWriter.FormDataContentType())
	//log.Debug("Body of request: ", requestBodyProcessing.String())
	multiPartWriter.Close()
	var responseProcessing *http.Response
	var responseReport *http.Response
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		responseProcessing, err = http.Post(srv.URL+"/process",
			contentType, &requestBodyProcessing)

	}()
	responseReport, err = http.Get(srv.URL + "/report")
	wg.Wait()
	defer responseProcessing.Body.Close()
	defer responseReport.Body.Close()
	if err != nil {
		log.Error(err)
		return
	}
	bodyProcessing, err := ioutil.ReadAll(responseProcessing.Body)
	if err != nil {
		log.Error(err)
		return
	}
	bodyReport, err := ioutil.ReadAll(responseReport.Body)
	if err != nil {
		log.Error(err)
		return
	}
	message := MessageResponse{}
	json.Unmarshal(bodyProcessing, &message)
	if len(message.Message) == 0 {
		log.Error("We have executed both requests simultaneously. So we should get an error: Processing running")
	}
	log.Println("Response of processing: ", string(bodyProcessing))
	log.Println("Response of report: ", string(bodyReport))
}
func TestBlockUsersEmptyIdsSlice(t *testing.T) {
	srv := httptest.NewServer(GetRouter())
	defer srv.Close()

	//BlockUsers
	rand.Seed(time.Now().UnixNano())
	requestBody, err := json.Marshal(blocking.BlockUsersReqModel{
		Ids: []int64{},
	})
	if err != nil {
		log.Error(err)
		return
	}
	log.Debug("Body of request", string(requestBody))
	responseBlockingUsers, err := http.Post(srv.URL+"/block-users",
		"application/json",
		bytes.NewBuffer(requestBody))
	defer func() {
		if responseBlockingUsers != nil {
			responseBlockingUsers.Body.Close()
		}
	}()
	if err != nil {
		log.Error(err)
		return
	}
	bodyBlockingUsers, err := ioutil.ReadAll(responseBlockingUsers.Body)
	if err != nil {
		log.Error(err)
		return
	}
	errorResp := ErrorResponse{}
	json.Unmarshal(bodyBlockingUsers, &errorResp)
	if len(errorResp.Error) == 0 {
		log.Error("We send request with empty body. We should get error")
		return
	}
	log.Println(string(bodyBlockingUsers))
}
