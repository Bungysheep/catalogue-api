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

func TestCatalogue(t *testing.T) {
	t.Run("Get all catalogues", getAllCatalogues)

	t.Run("Get catalogue", getCatalogue)

	t.Run("Create catalogue", createCatalogue)

	t.Run("Create catalogue without custom field definitions", createCatalogueWithoutFieldDef)

	t.Run("Update catalogue with invalid version", updateCatalogueWithInvalidVersion)

	t.Run("Update catalogue", updateCatalogue)

	t.Run("Delete catalogue", deleteCatalogue)
}

func getAllCatalogues(t *testing.T) {
	req, err := http.NewRequest("GET", "http://localhost:50051/v1/catalogues", bytes.NewBuffer([]byte("")))
	assert.NilError(t, err, "Failed to create get all request.")

	req.Header.Add("Authorization", accessTokenTest)

	client := &http.Client{}

	resp, err := client.Do(req)
	assert.NilError(t, err, "Failed to retrieve all catalogues.")
	assert.Equal(t, resp.StatusCode, http.StatusOK)

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	assert.NilError(t, err, "Failed to read body response.")

	var respData map[string]interface{}
	err = json.Unmarshal(body, &respData)
	assert.NilError(t, err, "Failed to decode body response.")
	assert.Equal(t, respData["success"], true)

	data := respData["data"].([]interface{})
	assert.Equal(t, len(data), 3)

	dataOutput := data[0].(map[string]interface{})
	assert.Equal(t, dataOutput["code"], "CLG_TEST_1")
	assert.Equal(t, dataOutput["description"], "Catalogue Test 1")
	assert.Equal(t, dataOutput["details"], "Catalogue Test 1")
	assert.Equal(t, dataOutput["status"], "A")
	assert.Equal(t, dataOutput["created_by"], "TESTUSER")
	assert.Equal(t, dataOutput["modified_by"], "TESTUSER")
}

func getCatalogue(t *testing.T) {
	req, err := http.NewRequest("GET", "http://localhost:50051/v1/catalogues/CLG_TEST_1", bytes.NewBuffer([]byte("")))
	assert.NilError(t, err, "Failed to create get all request.")

	req.Header.Add("Authorization", accessTokenTest)

	client := &http.Client{}

	resp, err := client.Do(req)
	assert.NilError(t, err, "Failed to retrieve catalogue.")
	assert.Equal(t, resp.StatusCode, http.StatusOK)

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	assert.NilError(t, err, "Failed to read body response.")

	var respData map[string]interface{}
	err = json.Unmarshal(body, &respData)
	assert.NilError(t, err, "Failed to decode body response.")
	assert.Equal(t, respData["success"], true)

	dataOutput := respData["data"].(map[string]interface{})
	assert.Equal(t, dataOutput["code"], "CLG_TEST_1")
	assert.Equal(t, dataOutput["description"], "Catalogue Test 1")
	assert.Equal(t, dataOutput["details"], "Catalogue Test 1")
	assert.Equal(t, dataOutput["status"], "A")
	assert.Equal(t, dataOutput["created_by"], "TESTUSER")
	assert.Equal(t, dataOutput["modified_by"], "TESTUSER")

	dataFieldDefs := dataOutput["field_definitions"].([]interface{})
	assert.Equal(t, len(dataFieldDefs), 3)

	dataFieldDefOutput := dataFieldDefs[0].(map[string]interface{})
	assert.Equal(t, dataFieldDefOutput["clg_code"], "CLG_TEST_1")
	assert.Equal(t, dataFieldDefOutput["caption"], "Field-1")
	assert.Equal(t, dataFieldDefOutput["type"], "A")
	assert.Equal(t, dataFieldDefOutput["created_by"], "TESTUSER")
	assert.Equal(t, dataFieldDefOutput["modified_by"], "TESTUSER")
}

func createCatalogue(t *testing.T) {
	dataInput := map[string]interface{}{
		"code":        "CLG_TEST",
		"description": "Catalogue Test",
		"details":     "Catalogue Test",
		"status":      "A",
		"created_at":  time.Now(),
		"modified_at": time.Now(),
		"vers":        1,
		"field_definitions": []interface{}{
			map[string]interface{}{
				"caption": "Field-1",
				"type":    "A",
			},
			map[string]interface{}{
				"caption": "Field-2",
				"type":    "N",
			},
			map[string]interface{}{
				"caption": "Field-3",
				"type":    "D",
			},
		},
	}

	bodyReq, err := json.Marshal(dataInput)
	assert.NilError(t, err, "Failed to encode body request.")

	req, err := http.NewRequest("POST", "http://localhost:50051/v1/catalogues", bytes.NewBuffer(bodyReq))
	assert.NilError(t, err, "Failed to create request.")

	req.Header.Add("Authorization", accessTokenTest)

	client := &http.Client{}

	resp, err := client.Do(req)
	assert.NilError(t, err, "Failed to create catalogue.")
	assert.Equal(t, resp.StatusCode, http.StatusAccepted)

	defer resp.Body.Close()

	bodyResp, err := ioutil.ReadAll(resp.Body)
	assert.NilError(t, err, "Failed to read body response.")

	var respData map[string]interface{}
	err = json.Unmarshal(bodyResp, &respData)
	assert.NilError(t, err, "Failed to decode body response.")
	assert.Equal(t, respData["success"], true)

	dataOutput := respData["data"].(map[string]interface{})
	assert.Equal(t, dataOutput["code"], "CLG_TEST")
	assert.Equal(t, dataOutput["description"], "Catalogue Test")
	assert.Equal(t, dataOutput["details"], "Catalogue Test")
	assert.Equal(t, dataOutput["status"], "A")
	assert.Equal(t, dataOutput["created_by"], "TESTUSER")
	assert.Equal(t, dataOutput["modified_by"], "TESTUSER")

	createdAt, _ := time.Parse(time.RFC3339, dataOutput["created_at"].(string))
	modifiedAt, _ := time.Parse(time.RFC3339, dataOutput["modified_at"].(string))

	assert.Equal(t, createdAt.Local().Format(configs.DATEFORMAT), dataInput["created_at"].(time.Time).Format(configs.DATEFORMAT))
	assert.Equal(t, modifiedAt.Local().Format(configs.DATEFORMAT), dataInput["modified_at"].(time.Time).Format(configs.DATEFORMAT))

	dataFieldDefs := dataOutput["field_definitions"].([]interface{})
	assert.Equal(t, len(dataFieldDefs), 3)

	dataFieldDefOutput := dataFieldDefs[0].(map[string]interface{})
	assert.Equal(t, dataFieldDefOutput["clg_code"], "CLG_TEST")
	assert.Equal(t, dataFieldDefOutput["caption"], "Field-1")
	assert.Equal(t, dataFieldDefOutput["type"], "A")
	assert.Equal(t, dataFieldDefOutput["created_by"], "TESTUSER")
	assert.Equal(t, dataFieldDefOutput["modified_by"], "TESTUSER")
}

func createCatalogueWithoutFieldDef(t *testing.T) {
	dataInput := map[string]interface{}{
		"code":              "CLG_TEST_4",
		"description":       "Catalogue Test",
		"details":           "Catalogue Test",
		"status":            "A",
		"created_at":        time.Now(),
		"modified_at":       time.Now(),
		"vers":              1,
		"field_definitions": []interface{}{},
	}

	bodyReq, err := json.Marshal(dataInput)
	assert.NilError(t, err, "Failed to encode body request.")

	req, err := http.NewRequest("POST", "http://localhost:50051/v1/catalogues", bytes.NewBuffer(bodyReq))
	assert.NilError(t, err, "Failed to create request.")

	req.Header.Add("Authorization", accessTokenTest)

	client := &http.Client{}

	resp, err := client.Do(req)
	assert.NilError(t, err, "Failed to create catalogue.")
	assert.Equal(t, resp.StatusCode, http.StatusAccepted)

	defer resp.Body.Close()

	bodyResp, err := ioutil.ReadAll(resp.Body)
	assert.NilError(t, err, "Failed to read body response.")

	var respData map[string]interface{}
	err = json.Unmarshal(bodyResp, &respData)
	assert.NilError(t, err, "Failed to decode body response.")
	assert.Equal(t, respData["success"], true)

	dataOutput := respData["data"].(map[string]interface{})
	assert.Equal(t, dataOutput["code"], "CLG_TEST_4")
	assert.Equal(t, dataOutput["description"], "Catalogue Test")
	assert.Equal(t, dataOutput["details"], "Catalogue Test")
	assert.Equal(t, dataOutput["status"], "A")
	assert.Equal(t, dataOutput["created_by"], "TESTUSER")
	assert.Equal(t, dataOutput["modified_by"], "TESTUSER")

	createdAt, _ := time.Parse(time.RFC3339, dataOutput["created_at"].(string))
	modifiedAt, _ := time.Parse(time.RFC3339, dataOutput["modified_at"].(string))

	assert.Equal(t, createdAt.Local().Format(configs.DATEFORMAT), dataInput["created_at"].(time.Time).Format(configs.DATEFORMAT))
	assert.Equal(t, modifiedAt.Local().Format(configs.DATEFORMAT), dataInput["modified_at"].(time.Time).Format(configs.DATEFORMAT))

	dataFieldDefs := dataOutput["field_definitions"].([]interface{})
	assert.Equal(t, len(dataFieldDefs), 0)
}

func updateCatalogueWithInvalidVersion(t *testing.T) {
	dataInput := map[string]interface{}{
		"code":        "CLG_TEST",
		"description": "Catalogue Test - Updated",
		"details":     "Catalogue Test - Updated",
		"status":      "I",
		"modified_at": time.Now(),
		"vers":        2,
	}

	bodyReq, err := json.Marshal(dataInput)
	assert.NilError(t, err, "Failed to encode body request.")

	req, err := http.NewRequest("PUT", "http://localhost:50051/v1/catalogues/CLG_TEST", bytes.NewBuffer(bodyReq))
	assert.NilError(t, err, "Failed to create update request.")

	req.Header.Add("Authorization", accessTokenTest)

	client := &http.Client{}

	resp, err := client.Do(req)
	assert.NilError(t, err, "Failed to update catalogue.")
	assert.Equal(t, resp.StatusCode, http.StatusBadRequest)

	defer resp.Body.Close()

	bodyResp, err := ioutil.ReadAll(resp.Body)
	assert.NilError(t, err, "Failed to read body response.")

	var respData map[string]interface{}
	err = json.Unmarshal(bodyResp, &respData)
	assert.NilError(t, err, "Failed to decode body response.")
	assert.Equal(t, respData["success"], false)
}

func updateCatalogue(t *testing.T) {
	dataInput := map[string]interface{}{
		"code":        "CLG_TEST",
		"description": "Catalogue Test - Updated",
		"details":     "Catalogue Test - Updated",
		"status":      "I",
		"modified_at": time.Now(),
		"vers":        1,
	}

	bodyReq, err := json.Marshal(dataInput)
	assert.NilError(t, err, "Failed to encode body request.")

	req, err := http.NewRequest("PUT", "http://localhost:50051/v1/catalogues/CLG_TEST", bytes.NewBuffer(bodyReq))
	assert.NilError(t, err, "Failed to create update request.")

	req.Header.Add("Authorization", accessTokenTest)

	client := &http.Client{}

	resp, err := client.Do(req)
	assert.NilError(t, err, "Failed to update catalogue.")
	assert.Equal(t, resp.StatusCode, http.StatusAccepted)

	defer resp.Body.Close()

	bodyResp, err := ioutil.ReadAll(resp.Body)
	assert.NilError(t, err, "Failed to read body response.")

	var respData map[string]interface{}
	err = json.Unmarshal(bodyResp, &respData)
	assert.NilError(t, err, "Failed to decode body response.")
	assert.Equal(t, respData["success"], true)

	dataOutput := respData["data"].(map[string]interface{})
	assert.Equal(t, dataOutput["code"], "CLG_TEST")
	assert.Equal(t, dataOutput["description"], "Catalogue Test - Updated")
	assert.Equal(t, dataOutput["details"], "Catalogue Test - Updated")
	assert.Equal(t, dataOutput["status"], "I")
	assert.Equal(t, dataOutput["created_by"], "TESTUSER")
	assert.Equal(t, dataOutput["modified_by"], "TESTUSER")

	modifiedAt, _ := time.Parse(time.RFC3339, dataOutput["modified_at"].(string))

	assert.Equal(t, modifiedAt.Local().Format(configs.DATEFORMAT), dataInput["modified_at"].(time.Time).Format(configs.DATEFORMAT))
}

func deleteCatalogue(t *testing.T) {
	req, err := http.NewRequest("DELETE", "http://localhost:50051/v1/catalogues/CLG_TEST", bytes.NewBuffer([]byte("")))
	assert.NilError(t, err, "Failed to create delete request.")

	req.Header.Add("Authorization", accessTokenTest)

	client := &http.Client{}

	resp, err := client.Do(req)
	assert.NilError(t, err, "Failed to create catalogue.")
	assert.Equal(t, resp.StatusCode, http.StatusOK)

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	assert.NilError(t, err, "Failed to read body response.")

	var respData map[string]interface{}
	err = json.Unmarshal(body, &respData)
	assert.NilError(t, err, "Failed to decode body response.")
	assert.Equal(t, respData["success"], true)
}
