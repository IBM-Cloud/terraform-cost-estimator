package pricing

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// const (
// 	//URL of Cost Estimator Endpoint...
// 	URL = "https://globalcatalog.cloud.ibm.com/api/v1/is.instance/plan"
// )

type PricingPlan struct {
	Offset        int    `json:"offset"`
	Limit         int    `json:"limit"`
	Count         int    `json:"count"`
	ResourceCount int    `json:"resource_count"`
	First         string `json:"first"`
	Resources     []struct {
		Active      bool      `json:"active"`
		CatalogCrn  string    `json:"catalog_crn"`
		ChildrenURL string    `json:"children_url"`
		Created     time.Time `json:"created"`
		Disabled    bool      `json:"disabled"`
		GeoTags     []string  `json:"geo_tags"`
		ID          string    `json:"id"`
		Images      struct {
			FeatureImage string `json:"feature_image"`
			Image        string `json:"image"`
			MediumImage  string `json:"medium_image"`
			SmallImage   string `json:"small_image"`
		} `json:"images"`
		Kind     string `json:"kind"`
		Metadata struct {
			RcCompatible bool `json:"rc_compatible"`
		} `json:"metadata"`
		Name       string `json:"name"`
		OverviewUI struct {
			En struct {
				Description     string `json:"description"`
				DisplayName     string `json:"display_name"`
				LongDescription string `json:"long_description"`
			} `json:"en"`
		} `json:"overview_ui"`
		ParentID    string   `json:"parent_id"`
		ParentURL   string   `json:"parent_url"`
		PricingTags []string `json:"pricing_tags"`
		Provider    struct {
			Email string `json:"email"`
			Name  string `json:"name"`
		} `json:"provider"`
		Tags       []string  `json:"tags"`
		Updated    time.Time `json:"updated"`
		URL        string    `json:"url"`
		Visibility struct {
			Restrictions string `json:"restrictions"`
		} `json:"visibility"`
	} `json:"resources"`
}

// GetPricing ...
func (plan *planService) GetPricing(instanceID string) (PricingPlan, error) {
	client := &http.Client{}
	URL := fmt.Sprintf("https://globalcatalog.cloud.ibm.com/api/v1/%s/plan", instanceID)

	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		return PricingPlan{}, err
	}
	req.Header.Add("Accept", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return PricingPlan{}, err
	}

	defer resp.Body.Close()
	resp_body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return PricingPlan{}, err
	}

	var newResp PricingPlan
	//decode to response body
	err1 := json.Unmarshal([]byte(resp_body), &newResp)
	if err1 != nil {
		return PricingPlan{}, fmt.Errorf(
			"Error Unmarshaling response body (%s) :", err1)
	}

	return newResp, nil

}
