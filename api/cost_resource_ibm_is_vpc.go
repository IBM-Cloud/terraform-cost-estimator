package costcalculator

import "go.uber.org/zap"

const (
	vpcID = "is.vpc"
)

//Parse floating-IP details, configure the body, call restapi and return the estimated cost
func getVpcCost(logger *zap.Logger, changeData ResourceConf, token string) (float64, error) {

	logger.Info("Entry:getVpcCost ")
	// var planID string

	//check the generation and set the planID for the floating-ip
	// if generation == 1 {
	// 	planID = "aeb480c6-11ae-4abc-929a-eaaefcdc5615"
	// } else {
	// 	planID = "5c6f17b8-1be9-4ae8-a328-afbfe41b02f8"
	// }

	// vpcPricingClient := pricing.NewPlanService(generation, vpcID)

	// planID, err := vpcPricingClient.GetVPCPlan()
	// if err != nil {
	// 	return 0, err
	// }
	//configure BOM
	// billdata := BillOfMaterial{}
	// billdata.AddLineItemData(resdata, planID, 0)

	logger.Info("Exit:getVpcCost ")
	return 0, nil

}
