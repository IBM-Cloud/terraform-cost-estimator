package costcalculator

import (
	rest "github.com/IBM-Cloud/terraform-cost-estimator/ibm/rest"
	"go.uber.org/zap"
)

func getCloudantCost(logger *zap.Logger, changeData ResourceConf, token string) (float64, error) {
	logger.Info("Entry:getCloudantCost")
	plan := changeData.Plan
	var objectID, serviceID string
	serviceID = "cloudant-" + plan
	planResp, err := rest.GetGlobalCatalogPlan(serviceID)
	if err != nil {
		logger.Error("Error occured while getting plan", zap.Error(err))
		return 0, err
	}
	// fmt.Printf("%+v\n", planResp)
	objectID = planResp.Resources[0].ID
	InstanceCostResp, err := rest.GetGlobalCatalogCost(objectID, "")
	if err != nil {
		logger.Error("Error occured while generating Cost plan", zap.Error(err))
		return 0, err
	}
	if InstanceCostResp.Type == "free" {
		logger.Warn("Cost of Instance is free, Resource have been created with Lite plan")
		return 0, nil
	}
	InstanceCost := InstanceCostResp.Metrics[0].Amounts[0].Prices[0].Price
	monthlyCost := getMonthlyCost(InstanceCost, InstanceCostResp.Metrics[0].ChargeUnitName, InstanceCostResp.Metrics[0].ChargeUnitQuantity)
	//configure BOM
	logger.Info("Exit:getCloudantCost")
	return monthlyCost, nil
}
