package costcalculator

import (
	"fmt"

	"github.com/IBM-Cloud/terraform-cost-estimator/pricing"
	rest "github.com/IBM-Cloud/terraform-cost-estimator/rest"
)

const (
	volumeID = "is.volume"
)

//Parse volume details, configure the body, call restapi and return the estimated cost
func getVolumeCost(resdata Resource, token string, generation int) (*BillOfMaterial, float64, error) {

	var planID string
	var iops int
	//parse configuration from properties
	// profile := resdata.Values.Profile
	volume := resdata.Values.Capacity

	pricingClient := pricing.NewPlanService(generation, volumeID)

	planID, err := pricingClient.GetVolumePlan(resdata.Values.Profile)
	if err != nil {
		return nil, 0, err
	}
	volcost, err := volumeCost(volume, iops, planID, token)
	if err != nil {
		return nil, 0, err
	}
	//configure BOM
	billdata := BillOfMaterial{}
	billdata.AddLineItemData(resdata, planID, volcost)

	return &billdata, volcost, nil

}

//configure the request body payload, call restapis and get the responce
func volumeCost(volume, iops int, planid, token string) (float64, error) {

	//configure request payload
	body := fmt.Sprintf(`{ 
		"service_id":"is.volume",
		"plan_id":"%s",
		"currency":"USD",
		"country":"USA",
		"region": "us-south",
		"measures":[
			{
				"measure":"VOLUME",
				"quantity":%d
			},
			{
				"measure":"IOPS",
				"quantity": %d
			}
		]
	}`, planid, volume, iops)

	//restapi call
	newResp, err := rest.PostRequest(body, token)
	if err != nil {
		return 0, err
	}

	return newResp.Cost, nil

}
