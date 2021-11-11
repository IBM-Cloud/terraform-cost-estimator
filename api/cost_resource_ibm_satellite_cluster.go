package costcalculator

import (
	"fmt"
	"strings"

	rest "github.com/IBM-Cloud/terraform-cost-estimator/rest"
	"go.uber.org/zap"
)

//Parse container cost details, configure the body, call restapi and return the estimated cost
func getSatelliteCluster(logger *zap.Logger, changeData ResourceConf, token string) (float64, error) {
	// for vpc the hardware is always shared
	workerCount := 1
	logger.Info("Entry:getSatelliteCluster")
	Zone := len(changeData.Zone)
	if changeData.WorkerCount != 0 {
		logger.Info("worker count in the data is ", zap.Any("workerCount", changeData.WorkerCount))
		workerCount = changeData.WorkerCount
	}
	logger.Info("Total number of zones available is", zap.Any("Zone", Zone))
	serviceID := "containers.kubernetes.satellite.roks"
	planResp, err := rest.GetGlobalCatalogPlan(serviceID)
	if err != nil {
		logger.Error("Error occured while getting plan", zap.Error(err))
		return 0, err
	}

	objectID := planResp.Resources[0].ID
	logger.Info("Object id is", zap.Any("ObjectID", objectID))

	workerCostResp, err := rest.GetGlobalCatalogCost(objectID, "")
	if err != nil || workerCostResp.Origin == nil {
		return 0, fmt.Errorf("Error getting workercost")
	}

	workerCost := workerCostResp.Metrics[0].Amounts[0].Prices[0].Price
	licenseCost := workerCostResp.Metrics[1].Amounts[0].Prices[0].Price

	if changeData.Flavour != "" {
		// example b3c.4x16
		machine_type := strings.Trim(changeData.Flavour, ".")[1]
		HostCount = float64(strings.Trim(string(machine_type), "x")[0])
	}

	logger.Info("total host count is", zap.Any("hostCount", HostCount))
	logger.Info("total worker count is", zap.Any("workerCount", workerCount))

	totalCores := HostCount * float64(workerCount)
	logger.Info("totalCores is ", zap.Any("totalCores", totalCores))

	clusterCost := totalCores * (workerCost*float64(Zone) + licenseCost)

	monthlyCost := getMonthlyCost(clusterCost, workerCostResp.Metrics[0].ChargeUnitName, workerCostResp.Metrics[0].ChargeUnitQuantity)
	//configure BOM

	logger.Info("Exit:getSatelliteCluster")
	return monthlyCost, nil

}
