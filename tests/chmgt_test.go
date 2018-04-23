package tests

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/heimdal-rw/chmgt/config"
	"github.com/heimdal-rw/chmgt/handling"
	"github.com/heimdal-rw/chmgt/models"
	"github.com/stretchr/testify/assert"
	"gopkg.in/mgo.v2/bson"
)

func TestMain(m *testing.M) {
	conf, err := config.ReadConfig("test_config.toml")
	if err != nil {
		log.Fatal(err)
	}
	handler, err = handling.NewHandler(conf)
	if err != nil {
		log.Fatal(err)
	}
	handler.Datasource, err = models.NewDatasource(conf)
	if err != nil {
		log.Fatal(err)
	}

	code := m.Run()

	clearCollections()

	os.Exit(code)
}

func TestGetSingleUser(t *testing.T) {
	clearCollections()
	insertUsers(handler.Datasource, []string{"admin"})

	response, code, err := executeRequest("GET", "/api/users", nil, "admin", "password_admin")
	if err != nil {
		t.Errorf("Got %s when GETing '/api/users'. Error: ", err)
	}

	checkResponseCode(t, http.StatusOK, code)
	body := formatGetData(response.Data)
	assert.Exactly(t, bool(true), response.Success, fmt.Sprintf("Expected 'success' object to be true. Got: %v", response.Success))
	assert.Exactly(t, string("success"), response.Message, fmt.Sprintf("Expected response message to be 'success'. Got: %v", response.Message))
	assert.Exactly(t, int(1), len(body), fmt.Sprintf("Expected the body to only have 1 hit. Got: %v", len(body)))
	assert.Exactly(t, string("********"), body[0]["password"], fmt.Sprintf("Expected masked password to be returned. Got %v", body[0]["password"]))
	assert.Exactly(t, string("admin"), body[0]["username"], fmt.Sprintf("Expected the username to be 'admin'. Got %v", body[0]["username"]))
	assert.Exactly(t, string("admin@example.com"), body[0]["email"], fmt.Sprintf("Expected the username to be 'admin@example.com'. Got %v", body[0]["email"]))
	assert.Exactly(t, string("admin"), body[0]["firstname"], fmt.Sprintf("Expected the firstname to be 'Admin'. Got %v", body[0]["firstname"]))
	assert.Exactly(t, string("User"), body[0]["lastname"], fmt.Sprintf("Expected the lastname to be 'User'. Got %v", body[0]["lastname"]))
}

func TestGetMultipleUsers(t *testing.T) {
	clearCollections()
	insertUsers(handler.Datasource, []string{"admin", "tester1", "tester2"})

	response, code, err := executeRequest("GET", "/api/users", nil, "admin", "password_admin")
	if err != nil {
		t.Errorf("Got %s when GETing '/api/users'. Error: ", err)
	}

	checkResponseCode(t, http.StatusOK, code)
	body := formatGetData(response.Data)
	assert.Exactly(t, bool(true), response.Success, fmt.Sprintf("Expected 'success' object to be true. Got: %v", response.Success))
	assert.Exactly(t, string("success"), response.Message, fmt.Sprintf("Expected response message to be 'success'. Got: %v", response.Message))
	assert.Exactly(t, int(3), len(body), fmt.Sprintf("Expected the body to only 3 hits. Got: %v", len(body)))
	assert.Exactly(t, string("********"), body[0]["password"], fmt.Sprintf("Expected masked password to be returned. Got %v", body[0]["password"]))
	assert.Exactly(t, string("********"), body[1]["password"], fmt.Sprintf("Expected masked password to be returned. Got %v", body[1]["password"]))
	assert.Exactly(t, string("********"), body[2]["password"], fmt.Sprintf("Expected masked password to be returned. Got %v", body[2]["password"]))
	assert.Exactly(t, string("admin"), body[0]["username"], fmt.Sprintf("Expected the username to be 'admin'. Got %v", body[0]["username"]))
	assert.Exactly(t, string("tester1"), body[1]["username"], fmt.Sprintf("Expected the username to be 'tester1'. Got %v", body[1]["username"]))
	assert.Exactly(t, string("tester2"), body[2]["username"], fmt.Sprintf("Expected the username to be 'tester3'. Got %v", body[2]["username"]))
}

func TestGetNonExistingUser(t *testing.T) {
	clearCollections()
	insertUsers(handler.Datasource, []string{"admin"})

	// Generate a ObjectId to use. This should not match any user since part of it being generated is using the system clock.
	id := bson.NewObjectId().Hex()

	_, code, err := executeRequest("GET", fmt.Sprintf("/api/users/%s", id), nil, "admin", "password_admin")
	if err != nil {
		t.Errorf("Got %s when GETing '%s'. Error: ", err, id)
	}

	checkResponseCode(t, http.StatusNotFound, code)
}

func TestAddUser(t *testing.T) {
	clearCollections()
	insertUsers(handler.Datasource, []string{"admin"})

	userData := `
{"username": "harry.potter",
"firstname": "Harry",
"lastname": "Potter",
"email": "hpotter@gmail.com",
"password": "supersecret"}
`

	dataBytes := []byte(userData)

	response, code, err := executeRequest("POST", "/api/users", bytes.NewBuffer(dataBytes), "admin", "password_admin")
	if err != nil {
		t.Errorf("Got %s when POSTing '/api/users'. Error: ", err)
	}

	checkResponseCode(t, http.StatusOK, code)

	assert.Exactly(t, bool(true), response.Success, fmt.Sprintf("Expected 'success' object to be true. Got: %v", response.Success))
	assert.Exactly(t, bool(true), bson.IsObjectIdHex(response.Data.(string)), fmt.Sprintf("Expected return data to be ObjectHexID, got: %v", response.Data))
}

func TestUpdateUser(t *testing.T) {
	clearCollections()
	insertUsers(handler.Datasource, []string{"admin", "hpotter"})

	// Figure out the user id
	var testUser map[string]interface{}
	response, code, err := executeRequest("GET", "/api/users", nil, "admin", "password_admin")
	if err != nil {
		t.Errorf("Got %s when GETing '/api/users'. Error: ", err)
	}
	checkResponseCode(t, http.StatusOK, code)

	data := formatGetData(response.Data)
	for _, v := range data {
		if v["username"] == "hpotter" {
			testUser = v
		}
	}

	// Set data to send with PUT
	updateData := fmt.Sprintf(`
{"firstname": "Harry",
"lastname": "Potter",
"username": "%s",
"email": "%s"}
`, testUser["username"], testUser["email"])
	dataBytes := []byte(updateData)

	// Make the PUT update
	response, code, err = executeRequest("PUT", fmt.Sprintf("/api/users/%s", testUser["id"]), bytes.NewBuffer(dataBytes), "admin", "password_admin")
	if err != nil {
		t.Errorf("Got %s when PUTing '/api/users/%s'. Error: ", err, testUser["id"])
	}
	checkResponseCode(t, http.StatusOK, code)

	// Get the user again to verify changes are in place
	response, code, err = executeRequest("GET", fmt.Sprintf("/api/users/%s", testUser["id"]), nil, "admin", "password_admin")
	if err != nil {
		t.Errorf("Got %s when GETing '/api/users/%s'. Error: ", err, testUser["id"])
	}
	checkResponseCode(t, http.StatusOK, code)
	body := formatGetData(response.Data)
	assert.Exactly(t, bool(true), response.Success, fmt.Sprintf("Expected 'success' object to be true. Got: %v", response.Success))
	assert.Exactly(t, string("success"), response.Message, fmt.Sprintf("Expected response message to be 'success'. Got: %v", response.Message))
	assert.Exactly(t, int(1), len(body), fmt.Sprintf("Expected the body to only have 1 hit. Got: %v", len(body)))
	assert.Exactly(t, string("********"), body[0]["password"], fmt.Sprintf("Expected masked password to be returned. Got %v", body[0]["password"]))
	assert.Exactly(t, string("hpotter"), body[0]["username"], fmt.Sprintf("Expected the username to be 'hpotter'. Got %v", body[0]["username"]))
	assert.Exactly(t, string("hpotter@example.com"), body[0]["email"], fmt.Sprintf("Expected the email to be 'hpotter@example.com'. Got %v", body[0]["email"]))
	assert.Exactly(t, string("Harry"), body[0]["firstname"], fmt.Sprintf("Expected the firstname to be 'Harry'. Got %v", body[0]["firstname"]))
	assert.Exactly(t, string("Potter"), body[0]["lastname"], fmt.Sprintf("Expected the lastname to be 'Potter'. Got %v", body[0]["lastname"]))

}

func TestDeleteUser(t *testing.T) {
	clearCollections()
	insertUsers(handler.Datasource, []string{"admin", "hpotter"})

	// Figure out the user id
	var testUser map[string]interface{}
	response, code, err := executeRequest("GET", "/api/users", nil, "admin", "password_admin")
	if err != nil {
		t.Errorf("Got %s when GETing '/api/users'. Error: ", err)
	}
	checkResponseCode(t, http.StatusOK, code)

	data := formatGetData(response.Data)
	for _, v := range data {
		if v["username"] == "hpotter" {
			testUser = v
		}
	}

	response, code, err = executeRequest("DELETE", fmt.Sprintf("/api/users/%s", testUser["id"]), nil, "admin", "password_admin")
	checkResponseCode(t, http.StatusOK, code)

	response, code, err = executeRequest("DELETE", fmt.Sprintf("/api/users/%s", testUser["id"]), nil, "admin", "password_admin")
	checkResponseCode(t, http.StatusNotFound, code)

}
