package costcalculator

import (
	"fmt"

	"github.com/IBM-Cloud/terraform-cost-estimator/pricing"
	rest "github.com/IBM-Cloud/terraform-cost-estimator/rest"
)

const (
	isImageID = "is.image"
)

//Parse image details, configure the body, call restapi and return the estimated cost
func getImageCost(resdata Resource, token string, generation int) (BillOfMaterial, float64, error) {
	// var planID string

	//check the generation and set the planID for the instance
	// if generation == 1 {
	// 	planID = "07a7cd01-53f7-4165-a917-03deb0103695"
	// } else {
	// 	planID = "2a8d400a-4ea9-494a-aabc-2096e6477155"
	// }

	PricingClient := pricing.NewPlanService(generation, isImageID)

	planID, err := PricingClient.GetVPCPlan()
	if err != nil {
		return BillOfMaterial{}, 0, err
	}
	imagecost, err := imageCost(planID, token)
	if err != nil {
		return BillOfMaterial{}, 0, err
	}
	//configure BOM
	billdata := BillOfMaterial{}
	billdata.AddLineItemData(resdata, planID, imagecost)

	return billdata, imagecost, nil

}

func imageCost(planID, token string) (float64, error) {

	//configure request payload
	body := fmt.Sprintf(`{ 
		"service_id":"is.image",
		"plan_id":"%s",
		"currency":"USD",
		"country":"USA",
		"region": "us-south",
		"measures":[
			{
				"measure":"VOLUME",
				"quantity":%d
			  }
		]
	}`, planID, 5)

	newResp, err := rest.PostRequest(body, token)
	if err != nil {
		return 0, err
	}

	return newResp.Cost, nil

}
