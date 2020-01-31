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

func TestProduct(t *testing.T) {
	t.Run("Get product by catalogue", getProductByCatalogue)

	t.Run("Get product", getProduct)

	t.Run("Create product", createProduct)

	t.Run("Update product with invalid version", updateProductWithInvalidVersion)

	t.Run("Update product", updateProduct)

	t.Run("Delete product", deleteProduct)
}

func getProductByCatalogue(t *testing.T) {
	req, err := http.NewRequest("GET", "http://localhost:50051/v1/products/bycatalogue/CLG_TEST_1", bytes.NewBuffer([]byte("")))
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
	assert.Equal(t, dataOutput["code"], "P-0001")
	assert.Equal(t, dataOutput["description"], "Book")
	assert.Equal(t, dataOutput["details"], "Book")
	assert.Equal(t, dataOutput["status"], "A")
	assert.Equal(t, dataOutput["created_by"], "TESTUSER")
	assert.Equal(t, dataOutput["modified_by"], "TESTUSER")
}

func getProduct(t *testing.T) {
	req, err := http.NewRequest("GET", "http://localhost:50051/v1/products/2", bytes.NewBuffer([]byte("")))
	assert.NilError(t, err, "Failed to create get all request.")

	req.Header.Add("Authorization", accessTokenTest)

	client := &http.Client{}

	resp, err := client.Do(req)
	assert.NilError(t, err, "Failed to retrieve product.")
	assert.Equal(t, resp.StatusCode, http.StatusOK)

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	assert.NilError(t, err, "Failed to read body response.")

	var respData map[string]interface{}
	err = json.Unmarshal(body, &respData)
	assert.NilError(t, err, "Failed to decode body response.")
	assert.Equal(t, respData["success"], true)

	dataOutput := respData["data"].(map[string]interface{})
	assert.Equal(t, dataOutput["code"], "P-0002")
	assert.Equal(t, dataOutput["description"], "Pen")
	assert.Equal(t, dataOutput["details"], "Pen")
	assert.Equal(t, dataOutput["status"], "A")
	assert.Equal(t, dataOutput["created_by"], "TESTUSER")
	assert.Equal(t, dataOutput["modified_by"], "TESTUSER")

	dataUoms := dataOutput["uoms"].([]interface{})
	assert.Equal(t, len(dataUoms), 2)

	dataUomOutput := dataUoms[0].(map[string]interface{})
	assert.Equal(t, dataUomOutput["prod_id"], float64(2))
	assert.Equal(t, dataUomOutput["code"], "EACH")
	assert.Equal(t, dataUomOutput["descr"], "Each")
	assert.Equal(t, dataUomOutput["is_default"], true)
	assert.Equal(t, dataUomOutput["ratio"], float64(1))
}

func createProduct(t *testing.T) {
	dataInput := map[string]interface{}{
		"clg_code":    "CLG_TEST_2",
		"code":        "Q-0001",
		"description": "Hardisk",
		"details":     "Hardisk",
		"status":      "A",
		"created_at":  time.Now(),
		"modified_at": time.Now(),
		"vers":        1,
	}

	bodyReq, err := json.Marshal(dataInput)
	assert.NilError(t, err, "Failed to encode body request.")

	req, err := http.NewRequest("POST", "http://localhost:50051/v1/products", bytes.NewBuffer(bodyReq))
	assert.NilError(t, err, "Failed to create request.")

	req.Header.Add("Authorization", accessTokenTest)

	client := &http.Client{}

	resp, err := client.Do(req)
	assert.NilError(t, err, "Failed to create product.")
	assert.Equal(t, resp.StatusCode, http.StatusAccepted)

	defer resp.Body.Close()

	bodyResp, err := ioutil.ReadAll(resp.Body)
	assert.NilError(t, err, "Failed to read body response.")

	var respData map[string]interface{}
	err = json.Unmarshal(bodyResp, &respData)
	assert.NilError(t, err, "Failed to decode body response.")
	assert.Equal(t, respData["success"], true)

	dataOutput := respData["data"].(map[string]interface{})
	assert.Equal(t, dataOutput["clg_code"], "CLG_TEST_2")
	assert.Equal(t, dataOutput["code"], "Q-0001")
	assert.Equal(t, dataOutput["description"], "Hardisk")
	assert.Equal(t, dataOutput["details"], "Hardisk")
	assert.Equal(t, dataOutput["status"], "A")
	assert.Equal(t, dataOutput["created_by"], "TESTUSER")
	assert.Equal(t, dataOutput["modified_by"], "TESTUSER")

	createdAt, _ := time.Parse(time.RFC3339, dataOutput["created_at"].(string))
	modifiedAt, _ := time.Parse(time.RFC3339, dataOutput["modified_at"].(string))

	assert.Equal(t, createdAt.Local().Format(configs.DATEFORMAT), dataInput["created_at"].(time.Time).Format(configs.DATEFORMAT))
	assert.Equal(t, modifiedAt.Local().Format(configs.DATEFORMAT), dataInput["modified_at"].(time.Time).Format(configs.DATEFORMAT))
}

func updateProductWithInvalidVersion(t *testing.T) {
	dataInput := map[string]interface{}{
		"clg_code":    "CLG_TEST_2",
		"code":        "Q-0001",
		"description": "Hardisk - Updated",
		"details":     "Hardisk - Updated",
		"status":      "A",
		"created_at":  time.Now(),
		"modified_at": time.Now(),
		"vers":        2,
	}

	bodyReq, err := json.Marshal(dataInput)
	assert.NilError(t, err, "Failed to encode body request.")

	req, err := http.NewRequest("PUT", "http://localhost:50051/v1/products/4", bytes.NewBuffer(bodyReq))
	assert.NilError(t, err, "Failed to create update request.")

	req.Header.Add("Authorization", accessTokenTest)

	client := &http.Client{}

	resp, err := client.Do(req)
	assert.NilError(t, err, "Failed to update product.")
	assert.Equal(t, resp.StatusCode, http.StatusBadRequest)

	defer resp.Body.Close()

	bodyResp, err := ioutil.ReadAll(resp.Body)
	assert.NilError(t, err, "Failed to read body response.")

	var respData map[string]interface{}
	err = json.Unmarshal(bodyResp, &respData)
	assert.NilError(t, err, "Failed to decode body response.")
	assert.Equal(t, respData["success"], false)
}

func updateProduct(t *testing.T) {
	dataInput := map[string]interface{}{
		"clg_code":    "CLG_TEST_2",
		"code":        "Q-0001",
		"description": "Hardisk - Updated",
		"details":     "Hardisk - Updated",
		"status":      "A",
		"created_at":  time.Now(),
		"modified_at": time.Now(),
		"vers":        1,
	}

	bodyReq, err := json.Marshal(dataInput)
	assert.NilError(t, err, "Failed to encode body request.")

	req, err := http.NewRequest("PUT", "http://localhost:50051/v1/products/4", bytes.NewBuffer(bodyReq))
	assert.NilError(t, err, "Failed to create update request.")

	req.Header.Add("Authorization", accessTokenTest)

	client := &http.Client{}

	resp, err := client.Do(req)
	assert.NilError(t, err, "Failed to update product.")
	assert.Equal(t, resp.StatusCode, http.StatusAccepted)

	defer resp.Body.Close()

	bodyResp, err := ioutil.ReadAll(resp.Body)
	assert.NilError(t, err, "Failed to read body response.")

	var respData map[string]interface{}
	err = json.Unmarshal(bodyResp, &respData)
	assert.NilError(t, err, "Failed to decode body response.")
	assert.Equal(t, respData["success"], true)

	dataOutput := respData["data"].(map[string]interface{})
	assert.Equal(t, dataOutput["clg_code"], "CLG_TEST_2")
	assert.Equal(t, dataOutput["code"], "Q-0001")
	assert.Equal(t, dataOutput["description"], "Hardisk - Updated")
	assert.Equal(t, dataOutput["details"], "Hardisk - Updated")
	assert.Equal(t, dataOutput["status"], "A")
	assert.Equal(t, dataOutput["created_by"], "TESTUSER")
	assert.Equal(t, dataOutput["modified_by"], "TESTUSER")

	modifiedAt, _ := time.Parse(time.RFC3339, dataOutput["modified_at"].(string))

	assert.Equal(t, modifiedAt.Local().Format(configs.DATEFORMAT), dataInput["modified_at"].(time.Time).Format(configs.DATEFORMAT))
}

func deleteProduct(t *testing.T) {
	req, err := http.NewRequest("DELETE", "http://localhost:50051/v1/products/4", bytes.NewBuffer([]byte("")))
	assert.NilError(t, err, "Failed to create delete request.")

	req.Header.Add("Authorization", accessTokenTest)

	client := &http.Client{}

	resp, err := client.Do(req)
	assert.NilError(t, err, "Failed to create product.")
	assert.Equal(t, resp.StatusCode, http.StatusOK)

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	assert.NilError(t, err, "Failed to read body response.")

	var respData map[string]interface{}
	err = json.Unmarshal(body, &respData)
	assert.NilError(t, err, "Failed to decode body response.")
	assert.Equal(t, respData["success"], true)
}
