package authentication

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/IBM-Cloud/terraform-cost-estimator/helpers"
)

type IAMTokenResponse struct {
	AccessToken     string `json:"access_token"`
	RefreshToken    string `json:"refresh_token"`
	UAAAccessToken  string `json:"uaa_token"`
	UAARefreshToken string `json:"uaa_refresh_token"`
	TokenType       string `json:"token_type"`
}

const (
	//ErrCodeInvalidToken  ...
	ErrCodeInvalidToken = "InvalidToken"
	URL                 = "https://iam.cloud.ibm.com/identity/token"
)

//Description ...
// func (e IAMError) Description() string {
// 	if e.ErrorDetails != "" {
// 		return e.ErrorDetails
// 	}
// 	return e.ErrorMessage
// }

//AuthenticateAPIKey ...
func AuthenticateAPIKey() (IAMTokenResponse, error) {
	icAPIKey := helpers.EnvFallBack([]string{"IC_API_KEY", "IBMCLOUD_API_KEY", "BM_API_KEY", "BLUEMIX_API_KEY"}, "")
	if icAPIKey == "" {
		return IAMTokenResponse{}, fmt.Errorf("Insufficient credentials Please check the documentation on how to configure the IBM Cloud credentials")
	}
	return getToken(icAPIKey)
}

func getToken(apiKey string) (IAMTokenResponse, error) {
	response, err := http.PostForm(URL, url.Values{
		"grant_type": {"urn:ibm:params:oauth:grant-type:apikey"},
		"apikey":     {apiKey}})
	//okay, moving on...
	if err != nil {
		//handle postform error
		return IAMTokenResponse{}, err

	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		//handle read response error
		return IAMTokenResponse{}, fmt.Errorf(
			"Error reading response body (%s) :", err)
	}

	var newResp IAMTokenResponse
	//decode to response body
	err1 := json.Unmarshal([]byte(body), &newResp)
	if err1 != nil {
		return IAMTokenResponse{}, fmt.Errorf(
			"Error Unmarshaling response body1 (%s) :", err)
	}
	return newResp, nil
}
