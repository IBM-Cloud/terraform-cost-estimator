package costcalculator

import (
	"fmt"

	rest "github.com/IBM-Cloud/terraform-cost-estimator/ibm/rest"
	"go.uber.org/zap"
)

///Parse container cost details, configure the body, call restapi and return the estimated cost
func getDatabaseCost(logger *zap.Logger, changeData ResourceConf, token string) (float64, error) {
	logger.Info("Entry:getDatabaseCost")
	plan := changeData.Plan
	service := changeData.Service
	var objectID, serviceID string
	switch service {
	//cases for all database services
	case "databases-for-mongodb":
		serviceID = "databases-for-mongodb-" + plan
	case "databases-for-etcd":
		serviceID = "databases-for-etcd-" + plan
	case "databases-for-postgresql":
		serviceID = "databases-for-postgresql-" + plan
	case "databases-for-redis":
		serviceID = "databases-for-redis-" + plan
	case "databases-for-elasticsearch":
		serviceID = "databases-for-elasticsearch-" + plan
	case "messages-for-rabbitmq":
		serviceID = "messages-for-rabbitmq-" + plan
	case "databases-for-cassandra":
		serviceID = "databases-for-cassandra-" + plan
	case "databases-for-enterprisedb":
		serviceID = "databases-for-enterprisedb-" + plan
	default:
		err := fmt.Errorf("invalid Service Provided")
		logger.Error("Invalid Service Provided", zap.Any("Service", service))
		return 0, err
	}

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

	logger.Info("Exit:getDatabaseCost")

	return monthlyCost, nil

}
