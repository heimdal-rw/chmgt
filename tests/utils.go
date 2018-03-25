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
func executeRequest(method string, path string, body io.Reader, user string, pw string) (handling.APIResponse, int, error) {
	ret := handling.APIResponse{}

	authString := fmt.Sprintf(`{"username": "%s", "password": "%s"}`, user, pw)
	authBytes := []byte(authString)
	authRequest, _ := http.NewRequest("POST", "/api/authenticate", bytes.NewBuffer(authBytes))
	authRequest.Header.Set("Content-Type", "application/json")
	authRec := httptest.NewRecorder()
	a.Router.ServeHTTP(authRec, authRequest)

	resp, err := formatResponse(authRec)
	if err != nil {
		return ret, 0, err
	}
	authToken := resp.Data.(string)

	req, _ := http.NewRequest(method, path, body)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", authToken))
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)

	ret, err = formatResponse(rr)
	if err != nil {
		return handling.APIResponse{}, 0, err
	}

	return ret, rr.Code, nil
}

// checkResponseCode checks that the response code is what we expect.
func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

// formatResponse unmarshals the json response into an APIResponse struct
func formatResponse(r *httptest.ResponseRecorder) (handling.APIResponse, error) {
	response := handling.APIResponse{}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return response, err
	}

	err = json.Unmarshal(body, &response)
	return response, err
}

// formatData returns the Data section of the response in a use-able map
func formatGetData(i interface{}) []map[string]interface{} {
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
		d.InsertItem(u, "Users")
	}
}

// clearCollections wipes out the mongo collections in our test db.
func clearCollections() {
	sess := a.DB.Session
	sess.DB("chmgt_test").C("Users").RemoveAll(nil)
}
