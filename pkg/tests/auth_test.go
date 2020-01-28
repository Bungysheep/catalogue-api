package tests

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/bungysheep/catalogue-api/pkg/configs"
	"gotest.tools/assert"
)

func TestAuth(t *testing.T) {
	t.Run("Get all users", getAllUsers)

	t.Run("Register user", registerUser)

	t.Run("Sign in user", signInUser)

	t.Run("Sign in user with inactive user", signInUserWithInactiveUser)

	t.Run("Sign in user with invalid username", signInUserWithInvalidUsername)

	t.Run("Sign in user with invalid password", signInUserWithInvalidPassword)
}

func getAllUsers(t *testing.T) {
	req, err := http.NewRequest("GET", configs.TESTDOMAIN+"/v1/users", bytes.NewBuffer([]byte("")))
	assert.NilError(t, err, "Failed to create get all request.")

	client := &http.Client{}

	resp, err := client.Do(req)
	assert.NilError(t, err, "Failed to retrieve all users.")
	assert.Equal(t, resp.StatusCode, http.StatusOK)

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	assert.NilError(t, err, "Failed to read body response.")

	var respData map[string]interface{}
	err = json.Unmarshal(body, &respData)
	assert.NilError(t, err, "Failed to decode body response.")
	assert.Equal(t, respData["success"], true)

	data := respData["data"].([]interface{})
	assert.Equal(t, len(data), 1)

	dataOutput := data[0].(map[string]interface{})
	assert.Equal(t, dataOutput["username"], "TESTUSER")
	assert.Equal(t, dataOutput["name"], "Test User")
	assert.Equal(t, dataOutput["email"], "Test.User@testmail.com")
	assert.Equal(t, dataOutput["status"], "I")
	assert.Equal(t, dataOutput["created_by"], "TESTUSER")
	assert.Equal(t, dataOutput["modified_by"], "TESTUSER")
}

func registerUser(t *testing.T) {
	dataInput := map[string]interface{}{
		"username":    "TESTACCOUNT",
		"name":        "Test Account",
		"email":       "test.account@testmail.com",
		"password":    "asdf1234",
		"status":      "A",
		"created_by":  "TESTACCOUNT",
		"created_at":  time.Now(),
		"modified_by": "TESTACCOUNT",
		"modified_at": time.Now(),
		"vers":        1,
	}
	bodyReq, err := json.Marshal(dataInput)
	assert.NilError(t, err, "Failed to encode body request.")

	resp, err := http.Post(configs.TESTDOMAIN+"/v1/auth/register", "application/json", bytes.NewBuffer(bodyReq))
	assert.NilError(t, err, "Failed to register user.")
	assert.Equal(t, resp.StatusCode, http.StatusAccepted)

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	assert.NilError(t, err, "Failed to read body response.")

	var respData map[string]interface{}
	err = json.Unmarshal(body, &respData)
	assert.NilError(t, err, "Failed to decode body response.")
	assert.Equal(t, respData["success"], true)

	dataOutput := respData["data"].(map[string]interface{})
	assert.Equal(t, dataOutput["username"], dataInput["username"])
	assert.Equal(t, dataOutput["name"], dataInput["name"])
	assert.Equal(t, dataOutput["email"], dataInput["email"])
	assert.Equal(t, dataOutput["status"], dataInput["status"])
	assert.Equal(t, dataOutput["created_by"], dataInput["created_by"])
	assert.Equal(t, dataOutput["modified_by"], dataInput["modified_by"])

	createdAt, _ := time.Parse(time.RFC3339, dataOutput["created_at"].(string))
	modifiedAt, _ := time.Parse(time.RFC3339, dataOutput["modified_at"].(string))

	assert.Equal(t, createdAt.Local().Format(configs.DATEFORMAT), dataInput["created_at"].(time.Time).Format(configs.DATEFORMAT))
	assert.Equal(t, modifiedAt.Local().Format(configs.DATEFORMAT), dataInput["modified_at"].(time.Time).Format(configs.DATEFORMAT))
}

func signInUser(t *testing.T) {
	requestBody, err := json.Marshal(map[string]interface{}{
		"username": "TESTACCOUNT",
		"password": "asdf1234",
	})
	assert.NilError(t, err, "Failed to encode body request.")

	resp, err := http.Post(configs.TESTDOMAIN+"/v1/auth/signin", "application/json", bytes.NewBuffer(requestBody))
	assert.NilError(t, err, "Failed to sign in user.")
	assert.Equal(t, resp.StatusCode, http.StatusOK)

	body, err := ioutil.ReadAll(resp.Body)
	assert.NilError(t, err, "Failed to read body response.")

	var respData map[string]interface{}
	err = json.Unmarshal(body, &respData)
	assert.NilError(t, err, "Failed to decode body response.")
	assert.Equal(t, respData["success"], true)

	dataOutput := respData["data"].(map[string]interface{})
	assert.Equal(t, dataOutput["username"], "TESTACCOUNT")
	assert.Equal(t, dataOutput["name"], "Test Account")
	assert.Equal(t, dataOutput["email"], "test.account@testmail.com")
	assert.Equal(t, dataOutput["status"], "A")
}

func signInUserWithInactiveUser(t *testing.T) {
	requestBody, err := json.Marshal(map[string]interface{}{
		"username": "TESTUSER",
		"password": "asdf1234",
	})
	assert.NilError(t, err, "Failed to encode body request.")

	resp, err := http.Post(configs.TESTDOMAIN+"/v1/auth/signin", "application/json", bytes.NewBuffer(requestBody))
	assert.NilError(t, err, "Failed to sign in user.")
	assert.Equal(t, resp.StatusCode, http.StatusBadRequest)

	body, err := ioutil.ReadAll(resp.Body)
	assert.NilError(t, err, "Failed to read body response.")

	var respData map[string]interface{}
	err = json.Unmarshal(body, &respData)
	assert.NilError(t, err, "Failed to decode body response.")
	assert.Equal(t, respData["success"], false)
}

func signInUserWithInvalidUsername(t *testing.T) {
	requestBody, err := json.Marshal(map[string]interface{}{
		"username": "TESTACCOUNT",
		"password": "asdf@1234",
	})
	assert.NilError(t, err, "Failed to encode body request.")

	resp, err := http.Post(configs.TESTDOMAIN+"/v1/auth/signin", "application/json", bytes.NewBuffer(requestBody))
	assert.NilError(t, err, "Failed to sign in user.")
	assert.Equal(t, resp.StatusCode, http.StatusBadRequest)

	body, err := ioutil.ReadAll(resp.Body)
	assert.NilError(t, err, "Failed to read body response.")

	var respData map[string]interface{}
	err = json.Unmarshal(body, &respData)
	assert.NilError(t, err, "Failed to decode body response.")
	assert.Equal(t, respData["success"], false)
}

func signInUserWithInvalidPassword(t *testing.T) {
	requestBody, err := json.Marshal(map[string]interface{}{
		"username": "TESTACCOUNT",
		"password": "asdf@1234",
	})
	assert.NilError(t, err, "Failed to encode body request.")

	resp, err := http.Post(configs.TESTDOMAIN+"/v1/auth/signin", "application/json", bytes.NewBuffer(requestBody))
	assert.NilError(t, err, "Failed to sign in user.")
	assert.Equal(t, resp.StatusCode, http.StatusBadRequest)

	body, err := ioutil.ReadAll(resp.Body)
	assert.NilError(t, err, "Failed to read body response.")

	var respData map[string]interface{}
	err = json.Unmarshal(body, &respData)
	assert.NilError(t, err, "Failed to decode body response.")
	assert.Equal(t, respData["success"], false)
}
