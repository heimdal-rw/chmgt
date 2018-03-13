package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/heimdal-rw/chmgt/handling"
	"github.com/heimdal-rw/chmgt/models"
)

// App is the object to build the application routes and datastore configuration
type App struct {
	Router http.Handler
	DB     *models.Datasource
}

var a App

// executeRequest actually creates a ServeHTTP instance and then executes the request passed to it. Request body optional.
func executeRequest(method string, path string, authToken string, body io.Reader) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, body)
	if authToken != "" {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", authToken))
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)

	return rr
}

// checkResponseCode checks that the response code is what we expect.
func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

// getAuthToken authenticates to the API and returns a token to use for the rest of the requests.
func getAuthToken(user string, pw string) (string, error) {

	authString := fmt.Sprintf(`{"username": "%s", "password": "%s"}`, user, pw)
	data := []byte(authString)
	response := executeRequest("POST", "/api/authenticate", "", bytes.NewBuffer(data))

	resp1, err := formatResponse(response)
	if err != nil {
		return "", err
	}

	return resp1.Data.(string), nil
}

// formatResponse unmarshals the json response into an APIResponseJSON struct
func formatResponse(r *httptest.ResponseRecorder) (handling.APIResponseJSON, error) {
	response := handling.APIResponseJSON{}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return response, err
	}

	err = json.Unmarshal(body, &response)
	return response, err
}

// formatData returns the Data section of the response in a use-able map
func formatData(i interface{}) []map[string]interface{} {
	r := make([]map[string]interface{}, 0)
	ilist := i.([]interface{})
	for _, v := range ilist {
		n := v.(map[string]interface{})
		r = append(r, n)
	}
	return r
}

// insertUsers adds some users to the database so we have something to auth with and check for with API calls
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

// clearCollections wipes out the mongo collections in our test db.
func clearCollections() {
	sess := a.DB.Session
	sess.DB("chmgt_test").C("Users").RemoveAll(nil)
}
