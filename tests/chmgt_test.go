package chmgt_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/heimdal-rw/chmgt/config"
	"github.com/heimdal-rw/chmgt/handling"
	"github.com/heimdal-rw/chmgt/models"
)

type App struct {
	Router http.Handler
	DB     *models.Datasource
}

var a App

func TestMain(m *testing.M) {
	a = App{}

	conf, err := config.ReadConfig("test_config.toml")
	if err != nil {
		log.Fatal(err)
	}
	handler, err := handling.NewHandler(conf)
	if err != nil {
		log.Fatal(err)
	}
	handler.Datasource, err = models.NewDatasource(
		conf.DatabaseConnection(),
		conf.Database.Name,
	)

	a.Router = handler.Router
	a.DB = handler.Datasource

	code := m.Run()

	clearCollections()

	os.Exit(code)
}

func TestGetSingleUser(t *testing.T) {
	clearCollections()
	insertUsers(a.DB, []string{"admin"})

	token, err := getAuthToken("admin", "password_admin")
	if err != nil {
		t.Errorf("Got %s when attempting to auth", err)
	}

	req, _ := http.NewRequest("GET", "/api/users", nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	fmtresp, err := formatResponse(response)
	if err != nil {
		t.Errorf("Error in response format. Got %s", err)
	}

	if fmtresp.Success != true {
		t.Errorf("Expected success to be 'true'. Got %v", fmtresp.Success)
	}

	if fmtresp.Message != "success" {
		t.Errorf("Expected message to be 'success'. Got %v", fmtresp.Message)
	}

	body := formatData(fmtresp.Data)

	if len(body) > 1 {
		t.Errorf("Expected the body to only have 1 hit. Got %v", len(body))
	}

	if body[0]["username"] != "admin" {
		t.Errorf("Expected the username to be 'admin'. Got %v", body[0]["username"])
	}

	if body[0]["email"] != "admin@example.com" {
		t.Errorf("Expected the username to be 'admin@example.com'. Got %v", body[0]["email"])
	}

	if body[0]["firstname"] != "admin" {
		t.Errorf("Expected the firstname to be 'Admin'. Got %v", body[0]["firstname"])
	}

	if body[0]["lastname"] != "User" {
		t.Errorf("Expected the lastname to be 'User'. Got %v", body[0]["firstname"])
	}
}

// TODO: Populate tests below

// func TestGetMultipleUsers(t *testing.T) {

// }

// func TestGetNonExistingUser(t *testing.T) {

// }

// func TestAddUser(t *testing.T) {

// }

// func TestUpdateUser(t *testing.T) {

// }

// func TestDeleteUser(t *testing.T) {

// }

// func TestGetSingleChangeRequest(t *testing.T) {

// }

// func TestGetMultipleChangeRequest(t *testing.T) {

// }

// func TestGetNonExistingChangeRequest(t *testing.T) {

// }

// func TestAddChangeRequest(t *testing.T) {

// }

// func TestUpdateChangeRequest(t *testing.T) {

// }

// func TestDeleteChangeRequest(t *testing.T) {

// }

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func getAuthToken(user string, pw string) (string, error) {

	authString := fmt.Sprintf(`{"username": "%s", "password": "%s"}`, user, pw)
	data := []byte(authString)

	req, err := http.NewRequest("POST", "/api/authenticate", bytes.NewBuffer(data))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	response := executeRequest(req)

	resp1, err := formatResponse(response)
	if err != nil {
		return "", err
	}

	return resp1.Data.(string), nil
}

func formatResponse(r *httptest.ResponseRecorder) (handling.APIResponseJSON, error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return handling.APIResponseJSON{}, err
	}
	response := handling.APIResponseJSON{}
	jsonErr := json.Unmarshal(body, &response)

	if jsonErr != nil {
		return handling.APIResponseJSON{}, jsonErr
	}

	return response, nil
}

func formatData(i interface{}) []map[string]interface{} {
	r := make([]map[string]interface{}, 0)
	ilist := i.([]interface{})
	for _, v := range ilist {
		n := v.(map[string]interface{})
		r = append(r, n)
	}
	return r
}

func insertUsers(d *models.Datasource, s []string) {
	for _, v := range s {
		u := models.Item{
			"username":  v,
			"password":  fmt.Sprintf("password_%s", v),
			"email":     fmt.Sprintf("%s@example.com", v),
			"firstname": v,
			"lastname":  "User"}
		d.InsertUser(u)
	}
}

func clearCollections() {
	sess := a.DB.Session
	sess.DB("chmgt_test").C("Users").RemoveAll(nil)
}
