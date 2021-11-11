package costcalculator

import (
	"fmt"
	"strings"

	rest "github.com/IBM-Cloud/terraform-cost-estimator/rest"
	"go.uber.org/zap"
)

//Parse container cost details, configure the body, call restapi and return the estimated cost
func getContainerCost(resdata Resource, token string, generation int) (*BillOfMaterial, float64, error) {
	hardware := resdata.Values.Container.Hardware
	machine := resdata.Values.Container.MachineType
	workerpoolSize := resdata.Values.Container.DefaultPoolSize
	// in case of workerpool resou=rse workerpool size is size_per_zone
	if resdata.Type == "ibm_container_worker_pool" {
		workerpoolSize = resdata.Values.SizePerZone
	}
	var objectID string

	//if machine name starts with b3c, then convert it to b1c,respectively as its the standard to global catalog api updated inputs
	//if machine name starts with c3c, m3c,me4c,mb4c then remove this prefix, as we dont need it in ObjectID in this case, its the standard to global catalog api input

	if strings.HasPrefix(machine, "b3c") {
		machine = strings.ReplaceAll(machine, "b3c", "b1c")

	}
	if strings.HasPrefix(machine, "u2c") {
		machine = strings.ReplaceAll(machine, "u2c", "u1c")

	}

	if strings.HasPrefix(machine, "c3c") || strings.HasPrefix(machine, "m3c") || strings.HasPrefix(machine, "me4c") || strings.HasPrefix(machine, "mb4c") {
		splitMachine := strings.Split(machine, ".")
		machine = strings.Join(splitMachine[1:], ".")
	}

	if hardware == "shared" {
		objectID = "public.containers.kubernetes." + machine + ":us-south"

	} else {
		if strings.HasPrefix(machine, "me4c") {
			bmachine := strings.ReplaceAll(resdata.Values.Profile, "me4c", "containers.kubernetes.baremetal")
			objectID = bmachine + ":us-south"

		} else if strings.HasPrefix(machine, "mb4c") {
			bmachine := strings.ReplaceAll(resdata.Values.Profile, "mb4c", "containers.kubernetes.baremetal")
			objectID = bmachine + ":us-south"
		} else {
			objectID = "private.containers.kubernetes." + machine + ":us-south"

		}

	}
	workerCostResp, err := rest.GetGlobalCatalogCost(objectID, "")
	if err != nil || workerCostResp.Origin == nil {
		return nil, 0, err
	}

	workerCost := workerCostResp.Metrics[0].Amounts[0].Prices[0].Price
	clusterCost := workerCost * float64(workerpoolSize)

	monthlyCost := getMonthlyCost(clusterCost, workerCostResp.Metrics[0].ChargeUnitName, workerCostResp.Metrics[0].ChargeUnitQuantity)
	//configure BOM
	billdata := BillOfMaterial{}
	billdata.AddLineItemData(resdata, objectID, monthlyCost)

	return &billdata, monthlyCost, nil

}

//Parse container cost details, configure the body, call restapi and return the estimated cost
func containerCost(logger *zap.Logger, changeData ResourceConf, token string) (float64, error) {
	logger.Info("Entry:containerCost")
	hardware := changeData.Hardware
	machine := changeData.MachineType
	workerpoolSize := changeData.DefaultPoolSize
	// in case of workerpool resou=rse workerpool size is size_per_zone
	// if resdata.Type == "ibm_container_worker_pool" {
	// 	workerpoolSize = resdata.Values.SizePerZone
	// }
	//Todo- Add Workerpool support
	objectID := getContainerObjectID(hardware, machine)
	workerCostResp, err := rest.GetGlobalCatalogCost(objectID, "")
	if err != nil || workerCostResp.Origin == nil {
		return 0, fmt.Errorf("Error getting workercost")
	}

	workerCost := workerCostResp.Metrics[0].Amounts[0].Prices[0].Price
	clusterCost := workerCost * float64(workerpoolSize)

	monthlyCost := getMonthlyCost(clusterCost, workerCostResp.Metrics[0].ChargeUnitName, workerCostResp.Metrics[0].ChargeUnitQuantity)
	//configure BOM

	logger.Info("Exit:containerCost")
	return monthlyCost, nil

}

//Calculate total monthly cost based on response config.. see response model for different types of clusters
func getMonthlyCost(instanceMetricCost float64, chargeUnitName string, chargeUnitQuantity int) float64 {
	if chargeUnitName == "INSTANCES" {
		return instanceMetricCost
	} else if chargeUnitName == "instance_hours" || chargeUnitName == "VIRTUAL_PROCESSOR_CORE_HOURS" {
		if chargeUnitQuantity == 1 {
			instanceMetricCost = instanceMetricCost * 720
		} else if chargeUnitQuantity == 10 {
			instanceMetricCost = instanceMetricCost * 72
		}
	}
	return instanceMetricCost

}

func getContainerObjectID(hardware, machine string) string {
	var objectID string

	//synchronizing ui vpc machine type with global catalog api machine type(currently doc is in syncup with ui)
	if strings.HasPrefix(hardware, "vpc") {
		if machine == "bx2.16x64" || machine == "bx2.32x128" || machine == "bx2.4x16" {
			machine = strings.ReplaceAll(machine, "bx2", "b1c")
		} else if machine == "bx2.8x32" || machine == "cx2.16x32" || machine == "cx2.32x64" {
			machine = strings.Split(machine, ".")[1]
		} else if machine == "cx2.2x4" {
			machine = strings.ReplaceAll(machine, "cx2", "u1c")
		}
		//updating hardware type with 'shared' type
		hardware = strings.Trim(hardware, "vpc_")
	}

	//if machine name starts with b3c, then convert it to b1c,respectively as its the standard to global catalog api updated inputs
	//if machine name starts with c3c, m3c,me4c,mb4c then remove this prefix, as we dont need it in ObjectID in this case, its the standard to global catalog api input
	if strings.HasPrefix(machine, "b3c") {
		machine = strings.ReplaceAll(machine, "b3c", "b1c")

	}
	if strings.HasPrefix(machine, "u2c") {
		machine = strings.ReplaceAll(machine, "u2c", "u1c")

	}

	if strings.HasPrefix(machine, "c3c") || strings.HasPrefix(machine, "m3c") || strings.HasPrefix(machine, "me4c") || strings.HasPrefix(machine, "mb4c") {
		splitMachine := strings.Split(machine, ".")
		machine = strings.Join(splitMachine[1:], ".")
	}

	if hardware == "shared" {
		objectID = "public.containers.kubernetes." + machine + ":us-south"

	} else {
		if strings.HasPrefix(machine, "me4c") {
			bmachine := strings.ReplaceAll(machine, "me4c", "containers.kubernetes.baremetal")
			objectID = bmachine + ":us-south"

		} else if strings.HasPrefix(machine, "mb4c") {
			bmachine := strings.ReplaceAll(machine, "mb4c", "containers.kubernetes.baremetal")
			objectID = bmachine + ":us-south"
		} else {
			objectID = "private.containers.kubernetes." + machine + ":us-south"

		}

	}

	return objectID
}
