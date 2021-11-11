package costcalculator

import (
	"fmt"

	"github.com/IBM-Cloud/terraform-cost-estimator/pricing"
	rest "github.com/IBM-Cloud/terraform-cost-estimator/rest"
	"go.uber.org/zap"
)

//Parse floating-IP details, configure the body, call restapi and return the estimated cost
func getIPCost(logger *zap.Logger, changeData ResourceConf, token string) (float64, error) {

	// var planID string

	//check the generation and set the planID for the floating-ip
	// if generation == 1 {
	// 	planID = "0ecd991d-3c0a-4fb1-b358-327322a882bf"
	// } else {
	// 	planID = "74a80442-144f-488f-85ee-f81150ff0169"
	// }
	logger.Info("Entry: getIPCost")
	ipPricingClient := pricing.NewPlanService(generation, "is.floating-ip")

	planIDIP, err := ipPricingClient.GetIPPlan()
	if err != nil {
		return 0, err
	}
	fipcost, err := iPCost(planIDIP, token)
	if err != nil {
		return 0, err
	}

	return fipcost, nil

}

func iPCost(planID, token string) (float64, error) {
	//configure request payload
	//no of instance =1
	body := fmt.Sprintf(`{ 
		"service_id":"is.floating-ip",
		"plan_id":"%s",
		"currency":"USD",
		"country":"USA",
		"region": "us-south",
		"measures":[
			{
				"measure":"INSTANCE",
				"quantity":1
			}
		]
	}`, planID)

	newResp, err := rest.PostRequest(body, token)
	if err != nil {
		return 0, err
	}

	return newResp.Cost, nil

}
