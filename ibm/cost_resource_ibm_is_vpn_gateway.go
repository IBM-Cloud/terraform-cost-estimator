package costcalculator

import (
	"fmt"

	"github.com/IBM-Cloud/terraform-cost-estimator/ibm/pricing"
	rest "github.com/IBM-Cloud/terraform-cost-estimator/ibm/rest"
	"go.uber.org/zap"
)

const (
	vpnID        = "is.vpn"
	floatingIPID = "is.floating-ip"
)

//Parse vpn details, configure the body, call restapi and return the estimated cost
func getVpnCost(resdata Resource, token string, generation int) (*BillOfMaterial, float64, error) {

	// var planID string

	// if generation == 1 {
	// 	planID = "a725be39-2286-4d2b-8e8f-52978940036c"
	// } else {
	// 	planID = "5d7ef689-7a6b-4a03-9282-5de2be997291"
	// }

	vpnPricingClient := pricing.NewPlanService(generation, vpnID)

	planID, err := vpnPricingClient.GetVPCPlan()
	if err != nil {
		return nil, 0, err
	}
	vpncost, err := vpnCost(planID, token)
	if err != nil {
		return nil, 0, err
	}

	//Additional cost for floating IP
	// var planIDip string

	//check the generation and set the planID for the floating-ip
	// if generation == 1 {
	// 	planIDip = "0ecd991d-3c0a-4fb1-b358-327322a882bf"
	// } else {
	// 	planIDip = "74a80442-144f-488f-85ee-f81150ff0169"
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
	vpncost += ipcost
	//configure BOM
	billdata := BillOfMaterial{}
	billdata.AddLineItemData(resdata, planID, vpncost)
	billdata.AddDependencyData(1, floatingIPID, ipcost)

	return &billdata, vpncost, nil

}

//New function IncCostFuncMap
func getVpnCost2(logger *zap.Logger, changeData ResourceConf, token string) (float64, error) {
	logger.Info("Entry:getVPNCost")

	vpnPricingClient := pricing.NewPlanService(generation, vpnID)

	planID, err := vpnPricingClient.GetVPCPlan()
	if err != nil {
		return 0, err
	}
	vpncost, err := vpnCost(planID, token)
	if err != nil {
		return 0, err
	}
	ipPricingClient := pricing.NewPlanService(generation, floatingIPID)

	planIDIP, err := ipPricingClient.GetIPPlan()
	if err != nil {
		return 0, err
	}
	ipcost, err := iPCost(planIDIP, token)
	if err != nil {
		return 0, err
	}
	vpncost += ipcost

	return vpncost, nil

}

func vpnCost(planID, token string) (float64, error) {

	//configure request payload
	//INSTANCE_HOUR =720 , CONNECTION_HOUR =720
	body := fmt.Sprintf(`{ 
		"service_id":"is.vpn",
		"plan_id":"%s",
		"currency":"USD",
		"country":"USA",
		"region": "us-south",
		"measures":[
			{
				"measure":"INSTANCE_HOUR",
				"quantity":720
			  },
			  {
				"measure":"CONNECTION_HOUR",
				"quantity":720
			  }
		]
	}`, planID)

	//restapi call
	newResp, err := rest.PostRequest(body, token)
	if err != nil {
		return 0, err
	}

	return newResp.Cost, nil

}
