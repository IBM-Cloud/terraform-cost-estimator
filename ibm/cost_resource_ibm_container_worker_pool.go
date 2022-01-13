package costcalculator

import (
	"fmt"

	rest "github.com/IBM-Cloud/terraform-cost-estimator/ibm/rest"
	"go.uber.org/zap"
)

//Parse container cost details, configure the body, call restapi and return the estimated cost
func workerPoolContainerCost(logger *zap.Logger, changeData ResourceConf, token string) (float64, error) {
	logger.Info("Entry:workerPoolContainerCost")
	hardware := changeData.Hardware
	machine := changeData.MachineType
	workerpoolSize := changeData.SizePerZone
	objectID := getContainerObjectID(hardware, machine)
	workerCostResp, err := rest.GetGlobalCatalogCost(objectID, "")
	if err != nil || workerCostResp.Origin == nil {
		return 0, fmt.Errorf("error getting workercost")
	}

	workerCost := workerCostResp.Metrics[0].Amounts[0].Prices[0].Price
	clusterCost := workerCost * float64(workerpoolSize)

	monthlyCost := getMonthlyCost(clusterCost, workerCostResp.Metrics[0].ChargeUnitName, workerCostResp.Metrics[0].ChargeUnitQuantity)
	//configure BOM

	logger.Info("Exit:workerPoolContainerCost")
	return monthlyCost, nil

}
