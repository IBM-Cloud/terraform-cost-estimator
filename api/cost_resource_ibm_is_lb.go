package costcalculator

import (
	"fmt"

	"github.com/IBM-Cloud/terraform-cost-estimator/pricing"
	rest "github.com/IBM-Cloud/terraform-cost-estimator/rest"
	"go.uber.org/zap"
)

const (
	loadBalancerID = "is.load-balancer"
)

//Parse load-balancer details, configure the body, call restapi and return the estimated cost
func getLbCost(resdata Resource, token string, generation int) (*BillOfMaterial, float64, error) {
	// var planID string
	//check the generation and set the planID for the instance
	// if generation == 1 {
	// 	planIDLB = "6092eed9-b0a4-4933-86dd-5532b22274dc"
	// } else {
	// 	planIDLB = "f6975b83-c82b-41b7-bc1d-753f9989f144"
	// }

	pricingClient := pricing.NewPlanService(generation, loadBalancerID)

	planIDLB, err := pricingClient.GetLBPlan()
	if err != nil {
		return nil, 0, err
	}

	lbasscost, err := lBCost(planIDLB, token)
	if err != nil {
		return nil, 0, err
	}

	//Additional cost for floating IP

	//check the generation and set the planID for the floating IP
	// if generation == 1 {
	// 	planIDIp = "0ecd991d-3c0a-4fb1-b358-327322a882bf"
	// } else {
	// 	planIDIp = "74a80442-144f-488f-85ee-f81150ff0169"
	// }
	ipPricingClient := pricing.NewPlanService(generation, floatingIPID)

	planIDIP, err := ipPricingClient.GetIPPlan()
	if err != nil {
		return nil, 0, err
	}
	ipcost, err := iPCost(planIDIP, token)
	if err != nil {
		return nil, 0, err
	}
	lbasscost += ipcost
	//configure bill of material

	billdata := BillOfMaterial{}
	billdata.AddLineItemData(resdata, planIDLB, lbasscost+ipcost)
	billdata.AddDependencyData(1, floatingIPID, ipcost)

	return &billdata, lbasscost, nil

}

func lBCost(planID, token string) (float64, error) {

	//configure request payload
	//INSTANCE_HOUR =720 , GIGABYTE_MONTH =1
	body := fmt.Sprintf(`{ 
		"service_id":"is.load-balancer",
		"plan_id":"%s",
		"currency":"USD",
		"country":"USA",
		"region": "us-south",
		"measures":[
			{
				"measure":"INSTANCE_HOUR",
				"quantity":%d
			  },
			  {
				"measure":"GIGABYTE_MONTH",
				"quantity":%d
			  }
		]
	}`, planID, 720, 1)

	//restapicall
	newResp, err := rest.PostRequest(body, token)
	if err != nil {
		return 0, err
	}

	return newResp.Cost, nil

}

func getLoadBalancerCost(logger *zap.Logger, changeData ResourceConf, token string) (float64, error) {
	// var planID string
	//check the generation and set the planID for the instance
	// if generation == 1 {
	// 	planIDLB = "6092eed9-b0a4-4933-86dd-5532b22274dc"
	// } else {
	// 	planIDLB = "f6975b83-c82b-41b7-bc1d-753f9989f144"
	// }

	logger.Info("Entry:getLoadBalancerCost")
	pricingClient := pricing.NewPlanService(generation, loadBalancerID)

	planIDLB, err := pricingClient.GetLBPlan()
	if err != nil {
		return 0, err
	}

	lbasscost, err := lBCost(planIDLB, token)
	if err != nil {
		return 0, err
	}

	//Additional cost for floating IP

	//check the generation and set the planID for the floating IP
	// if generation == 1 {
	// 	planIDIp = "0ecd991d-3c0a-4fb1-b358-327322a882bf"
	// } else {
	// 	planIDIp = "74a80442-144f-488f-85ee-f81150ff0169"
	// }
	ipPricingClient := pricing.NewPlanService(generation, floatingIPID)

	planIDIP, err := ipPricingClient.GetIPPlan()
	if err != nil {
		return 0, err
	}
	ipcost, err := iPCost(planIDIP, token)
	if err != nil {
		return 0, err
	}
	lbasscost += ipcost

	logger.Info("Entry:getLoadBalancerCost")

	return lbasscost, nil

}
