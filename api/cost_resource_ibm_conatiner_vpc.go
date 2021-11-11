package costcalculator

import (
	"fmt"

	rest "github.com/IBM-Cloud/terraform-cost-estimator/rest"
	"go.uber.org/zap"
)

///Parse container cost details, configure the body, call restapi and return the estimated cost
func vpcContainerCost(logger *zap.Logger, changeData ResourceConf, token string) (float64, error) {
	// for vpc the hardware is always shared
	logger.Info("Entry:vpcContainerCost")
	hardware := "vpc_shared"
	machine := changeData.Flavour
	WorkerCount := changeData.WorkerCount
	Zone := len(changeData.Zone)
	objectID := getContainerObjectID(hardware, machine)
	logger.Info("Object id is", zap.Any("ObjectID", objectID))
	workerCostResp, err := rest.GetGlobalCatalogCost(objectID, "")
	if err != nil || workerCostResp.Origin == nil {
		return 0, fmt.Errorf("Error getting workercost")
	}

	workerCost := workerCostResp.Metrics[0].Amounts[0].Prices[0].Price
	clusterCost := workerCost * float64(WorkerCount) * float64(Zone)

	monthlyCost := getMonthlyCost(clusterCost, workerCostResp.Metrics[0].ChargeUnitName, workerCostResp.Metrics[0].ChargeUnitQuantity)
	//configure BOM

	logger.Info("Exit:vpcContainerCost")
	return monthlyCost, nil

}
