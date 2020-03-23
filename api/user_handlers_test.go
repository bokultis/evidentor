package main

// import (
// 	"bytes"
// 	"encoding/json"
// 	"evidentor/api/db"
// 	router "evidentor/api/router"
// 	user "evidentor/api/user"
// 	"io"
// 	"math/rand"
// 	"net/http"
// 	"net/http/httptest"
// 	"reflect"
// 	"strconv"
// 	"testing"
// 	"time"
// )

// func init() {
// 	db.DB = db.SetupDB()
// }

// var bearer = "Bearer " + "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImJvcmlzQGhvcmlzZW4uY29tIiwiZXhwIjoxNTY4NjIyNDE3LCJmaXJzdF9uYW1lIjoiQm9yaXMiLCJsYXN0X25hbWUiOiJLdW5kYWNpbmEiLCJ1c2VyX2lkIjoxfQ.zUrDOArvEKyPTbWPerVloGrN51edEuy6BY4iTfL9_Ic"

// func TestUsersIndexHandler(t *testing.T) {

// 	var body io.Reader
// 	body = nil
// 	req := prepareRequest(t, "GET", "/users", body)

// 	rr := executeRequest(req)

// 	// Check the status code is what we expect.
// 	if status := rr.Code; status != http.StatusOK {
// 		t.Errorf("handler returned wrong status code: got %v want %v",
// 			status, http.StatusOK)
// 	}

// 	// Check the response body is what we expect.
// 	byt := rr.Body.Bytes()
// 	var dat []user.UserWO
// 	if err := json.Unmarshal(byt, &dat); err != nil {
// 		panic(err)
// 	}

// 	rt := reflect.TypeOf(dat)

// 	if rt.Kind() != reflect.Slice {
// 		t.Errorf("handler returned unexpected body: got %v want slice of UserWO",
// 			rr.Body.String())
// 	}

// }

// func TestUsersShowHandler(t *testing.T) {

// 	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
// 	// pass 'nil' as the third parameter.
// 	var body io.Reader
// 	body = nil
// 	req := prepareRequest(t, "GET", "/users/1", body)
// 	rr := executeRequest(req)

// 	// Check the status code is what we expect.
// 	if status := rr.Code; status != http.StatusOK {
// 		t.Errorf("handler returned wrong status code: got %v want %v",
// 			status, http.StatusOK)
// 	}

// 	// Check the response body is what we expect.
// 	byt := rr.Body.Bytes()
// 	var dat user.UserWO
// 	if err := json.Unmarshal(byt, &dat); err != nil {
// 		panic(err)
// 	}

// 	rt := reflect.TypeOf(dat)
// 	if rt.Kind() != reflect.Struct {
// 		t.Errorf("handler returned unexpected body: got %v want UserWO",
// 			rr.Body.String())
// 	}

// }

// func TestUsersCreateHandler(t *testing.T) {

// 	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
// 	// pass 'nil' as the third parameter.
// 	s1 := rand.NewSource(time.Now().UnixNano())
// 	r1 := rand.New(s1)
// 	payload := []byte(`{
// 		"firstName": "Boris",
// 		"lastName": "Kundacina",
// 		"gender": "male",
// 		"eMail": "djoka` + strconv.Itoa(r1.Intn(100)) + `@horisen.com",
// 		"address": "Viz. Bul. 78",
// 		"password":"xyz",
// 		"optin": "yes"
// 	}`)

// 	body := bytes.NewBuffer(payload)
// 	req := prepareRequest(t, "POST", "/users", body)
// 	rr := executeRequest(req)

// 	// Check the status code is what we expect.
// 	if status := rr.Code; status != http.StatusOK {
// 		t.Errorf("handler returned wrong status code: got %v want %v",
// 			status, http.StatusOK)
// 	}

// 	// Check the response body is what we expect.
// 	byt := rr.Body.Bytes()
// 	var dat user.UserWO
// 	if err := json.Unmarshal(byt, &dat); err != nil {
// 		panic(err)
// 	}

// 	rt := reflect.TypeOf(dat)
// 	if rt.Kind() != reflect.Struct {
// 		t.Errorf("handler returned unexpected body: got %v want UserWO",
// 			rr.Body.String())
// 	}

// }

// func TestUsersUpdateHandler(t *testing.T) {

// 	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
// 	// pass 'nil' as the third parameter.

// 	payload := []byte(`{
// 		"firstName": "Boris",
// 		"lastName": "Kundacina",
// 		"gender": "male",
// 		"address": "Viz. Bul. 7555",
// 		"birthdate":"1979-25-06",
// 		"optin": "no"
// 	}`)

// 	body := bytes.NewBuffer(payload)
// 	req := prepareRequest(t, "PUT", "/users/17", body)
// 	rr := executeRequest(req)

// 	// Check the status code is what we expect.
// 	if status := rr.Code; status != http.StatusOK {
// 		t.Errorf("handler returned wrong status code: got %v want %v",
// 			status, http.StatusOK)
// 	}

// 	// Check the response body is what we expect.
// 	byt := rr.Body.Bytes()
// 	var dat user.UserWO
// 	if err := json.Unmarshal(byt, &dat); err != nil {
// 		panic(err)
// 	}

// 	rt := reflect.TypeOf(dat)
// 	if rt.Kind() != reflect.Struct {
// 		t.Errorf("handler returned unexpected body: got %v want UserWO",
// 			rr.Body.String())
// 	}

// }

// func TestUsersDeleteHandler(t *testing.T) {

// 	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
// 	// pass 'nil' as the third parameter.

// 	req := prepareRequest(t, "DELETE", "/users/18", nil)
// 	rr := executeRequest(req)

// 	// Check the status code is what we expect.
// 	if status := rr.Code; status != http.StatusOK {
// 		t.Errorf("handler returned wrong status code: got %v want %v",
// 			status, http.StatusOK)
// 	}

// 	expected := "User deleted"
// 	if rr.Body.String() != expected {
// 		t.Errorf("handler returned unexpected body: got %v want %s",
// 			rr.Body.String(), expected)
// 	}

// }

// func prepareRequest(t *testing.T, method string, uri string, body io.Reader) *http.Request {

// 	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
// 	// pass 'nil' as the third parameter.
// 	req, err := http.NewRequest(method, uri, body)
// 	// add authorization header to the req
// 	req.Header.Add("Authorization", bearer)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	return req
// }

// func executeRequest(req *http.Request) *httptest.ResponseRecorder {
// 	rr := httptest.NewRecorder()
// 	var AppRoutes []router.RoutePrefix
// 	AppRoutes = append(AppRoutes, user.Routes)
// 	router := router.NewRouter(&AppRoutes)
// 	router.ServeHTTP(rr, req)
// 	return rr
// }
