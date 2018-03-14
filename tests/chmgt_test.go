package tests

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/heimdal-rw/chmgt/config"
	"github.com/heimdal-rw/chmgt/handling"
	"github.com/heimdal-rw/chmgt/models"
	"github.com/stretchr/testify/assert"
)

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
	handler.Datasource, err = models.NewDatasource(conf)

	a.Router = handler.Router
	a.DB = handler.Datasource

	code := m.Run()

	clearCollections()

	os.Exit(code)
}

func TestGetSingleUser(t *testing.T) {
	clearCollections()
	insertUsers(a.DB, []string{"admin"})

	response, err := executeRequest("GET", "/api/users", nil, "admin", "password_admin")
	if err != nil {
		t.Errorf("Got %s when GETing '/api/users'. Error: ", err)
	}
	fmtresp, err := formatResponse(response)
	if err != nil {
		t.Errorf("Error in response format. Got %s", err)
	}

	// Testing the response object for this particular function
	checkResponseCode(t, http.StatusOK, response.Code)
	body := formatData(fmtresp.Data)
	assert.Exactly(t, bool(true), fmtresp.Success, fmt.Sprintf("Expected 'success' object to be true. Got: %v", fmtresp.Success))
	assert.Exactly(t, string("success"), fmtresp.Message, fmt.Sprintf("Expected response message to be 'success'. Got: %v", fmtresp.Message))
	assert.Exactly(t, int(1), len(body), fmt.Sprintf("Expected the body to only have 1 hit. Got: %v", len(body)))
	assert.Exactly(t, nil, body[0]["password"], fmt.Sprintf("Expected no password to be returned. Got %v", body[0]["password"]))
	assert.Exactly(t, string("admin"), body[0]["username"], fmt.Sprintf("Expected the username to be 'admin'. Got %v", body[0]["username"]))
	assert.Exactly(t, string("admin@example.com"), body[0]["email"], fmt.Sprintf("Expected the username to be 'admin@example.com'. Got %v", body[0]["email"]))
	assert.Exactly(t, string("admin"), body[0]["firstname"], fmt.Sprintf("Expected the firstname to be 'Admin'. Got %v", body[0]["firstname"]))
	assert.Exactly(t, string("User"), body[0]["lastname"], fmt.Sprintf("Expected the lastname to be 'User'. Got %v", body[0]["lastname"]))
}

// TODO: Populate tests below

func TestGetMultipleUsers(t *testing.T) {
	clearCollections()
	insertUsers(a.DB, []string{"admin", "tester1", "tester2"})

	response, err := executeRequest("GET", "/api/users", nil, "admin", "password_admin")
	if err != nil {
		t.Errorf("Got %s when GETing '/api/users'. Error: ", err)
	}
	fmtresp, err := formatResponse(response)
	if err != nil {
		t.Errorf("Error in response format. Got %s", err)
	}

	// Testing the response object for this particular function
	checkResponseCode(t, http.StatusOK, response.Code)
	body := formatData(fmtresp.Data)
	assert.Exactly(t, int(3), len(body), fmt.Sprintf("Expected the body to only 3 hits. Got: %v", len(body)))
}

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
