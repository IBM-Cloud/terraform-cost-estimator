package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	//URL of Cost Estimator Endpoint...
	URL = "https://billing.cloud.ibm.com/v4/calculator/meter"
)

type Respstruct struct {
	Cost     float64     `json:"cost"`
	Measures interface{} `json:"measures"`
	PlanID   string      `json:"plan_id"`
}

//getRequest ...
func GetGlobalCatalogCost(objectID, token string) (*GlobalCatalogResponse, error) {
	url := fmt.Sprintf("https://globalcatalog.cloud.ibm.com/api/v1/%s/pricing", objectID)
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("accept", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if !(res.StatusCode >= 200 && res.StatusCode < 300) && res.StatusCode != 404 && res.StatusCode != 410 {
		// If it is NotSuccess && NotFound && NotGone
		// catch all errors and log.
		return nil, fmt.Errorf("Unexpected status: %d - Response: %s", res.StatusCode, string(body))
	}
	var gResp GlobalCatalogResponse
	//decode to response body
	err1 := json.Unmarshal([]byte(body), &gResp)
	if err1 != nil {
		return nil, fmt.Errorf(
			"Error Unmarshaling response body (%s) :", err1)
	}

	return &gResp, nil
}

//PostRequest ...
func PostRequest(requestBody, token string) (Respstruct, error) {

	jsonb := []byte(requestBody)

	req, err := http.NewRequest("POST", URL, bytes.NewBuffer(jsonb))
	if err != nil {
		return Respstruct{}, err
	}

	//set auth header token
	req.Header.Set("Authorization", token)
	req.Header.Set("Content-Type", "application/json")
	//rest call
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return Respstruct{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Respstruct{}, fmt.Errorf(
			"Error reading response body (%s) :", err)
	}

	var newResp Respstruct
	//decode to response body
	err1 := json.Unmarshal([]byte(body), &newResp)
	if err1 != nil {
		return Respstruct{}, fmt.Errorf(
			"Error Unmarshaling response body (%s) :", err1)
	}

	return newResp, nil
}

//getRequest ...
func GetGlobalCatalogPlan(serviceID string) (*GlobalCatalogPlanResponse, error) {
	url := fmt.Sprintf("https://globalcatalog.cloud.ibm.com/api/v1/%s/*", serviceID)
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("accept", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if !(res.StatusCode >= 200 && res.StatusCode < 300) && res.StatusCode != 404 && res.StatusCode != 410 {
		// If it is NotSuccess && NotFound && NotGone
		// catch all errors and log.
		return nil, fmt.Errorf("Unexpected status: %d - Response: %s", res.StatusCode, string(body))
	}
	var planResp GlobalCatalogPlanResponse
	//decode to response body
	err1 := json.Unmarshal([]byte(body), &planResp)
	if err1 != nil {
		return nil, fmt.Errorf("Error Unmarshaling response body (%s) :", err1)
	}

	return &planResp, nil
}
