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
	defaultVolumeProfile = "3iops-tier"
)

var HostCount float64

//Parse instance details, configure the body, call restapi and return the estimated cost
func getInstanceCost(logger *zap.Logger, changeData ResourceConf, token string) (float64, error) {
	//parse vcpu and memory from the profile
	logger.Info("Entry:getInstanceCost")
	profile := strings.Split(strings.Split(changeData.ISInstance.Profile, "-")[1], "x")
	vcpu := profile[0]
	memory := profile[1]
	imageid := changeData.ISInstance.Image
	HostCount, _ = strconv.ParseFloat(vcpu, 64)

	//check for image id and map with the metric value
	var image string
	if imageid == "51af68c9-5558-4425-825a-f9243a3b2c6c" || imageid == "624cde4a-b4fe-4426-8f60-150a019a67f9" || imageid == "a7a0626c-f97e-4180-afbe-0331ec62f32a" {
		image = "is.image.windows"
	} else if imageid == "54c1ba68-6d29-42e5-9ca7-e5f4a62c1503" {
		image = "is.image.rhel"
	} else {
		image = "is.image.centos"
	}

	pricingClient := pricing.NewPlanService(generation, "is.instance")

	planIDInstance, err := pricingClient.GetInstancePlan(changeData.ISInstance.Profile)
	if err != nil {
		return 0, err
	}
	//instanceCost will generate the body and call restapi calls
	inscost, err := instanceCost(vcpu, memory, image, token, planIDInstance)
	if err != nil {
		return 0, err
	}

	var volcost float64
	//Additional cost for floating IP
	if len(changeData.ISInstance.BootVolume) == 0 {
		var planIDVolume string
		// if generation == 1 {
		// 	planIDVolume = "e65b2a9c-039c-4ed5-8c30-6fb7bb285b92"

		// } else {
		// 	planIDVolume = "4579dbd3-80bb-43dc-a645-ba3b03335cc9"
		// }
		volPricingClient := pricing.NewPlanService(generation, volumeID)

		planIDVolume, err := volPricingClient.GetVolumePlan(defaultVolumeProfile)
		if err != nil {
			return 0, err
		}
		//volumeCost is called as a default volume is initialised with instance
		//100G volume and 3iops is the default plan for colume
		volcost, err = volumeCost(100, 3, planIDVolume, token)
		if err != nil {
			return 0, err
		}

		inscost += volcost

	}

	logger.Info("Exit:getInstanceCost")
	return inscost, nil
}

//configure the request body payload, call restapis and get the responce
func instanceCost(vcpu, memory, image, token, planID string) (float64, error) {

	//configure payload body as json string
	body := fmt.Sprintf(`{
		"service_id":"is.instance",
		"plan_id":"%s",
		"currency":"USD",
		"country":"USA",
		"region": "us-south",
		"measures":[
			{
				"measure":"MEMORY",
				"quantity":%s
			},
			{
				"measure":"VCPU",
				"quantity":%s
			},
			{
				"measure":"IMAGE",
				"quantity":"%s"
			}
		]
	}`, planID, memory, vcpu, image)

	//call rest-apis
	newResp, err := rest.PostRequest(body, token)
	if err != nil {
		return 0, err
	}
	// //configure bill of material
	// billdata := BillOfMaterial{
	// 	Quantity:      1,
	// 	ID:            "is.instance",
	// 	LineItemTotal: newResp.Cost,
	// }

	return newResp.Cost, nil

}
