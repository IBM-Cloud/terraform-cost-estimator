package costcalculator

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/IBM-Cloud/terraform-cost-estimator/ibm/pricing"
	rest "github.com/IBM-Cloud/terraform-cost-estimator/ibm/rest"
	"go.uber.org/zap"
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

//New function IncCostFuncMap
func getVolumeCost2(logger *zap.Logger, changeData ResourceConf, token string) (float64, error) {
	logger.Info("Entry:getVolumeCost2")
	var planID string
	iops, _ := strconv.Atoi(strings.Replace(changeData.ISInstance.Profile, "iops-tier", "", 1))
	if iops == 0 {
		iops = changeData.Iops
	}
	pricingClient := pricing.NewPlanService(generation, volumeID)

	planID, err := pricingClient.GetVolumePlan(changeData.ISInstance.Profile)
	if err != nil {
		return 0, err
	}

	volcost, err := volumeCost(changeData.Capacity, iops, planID, token)
	if err != nil {
		return 0, err
	}

	return volcost, nil

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
