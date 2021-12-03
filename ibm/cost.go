package costcalculator

import (
	"encoding/json"
	"io/ioutil"

	"github.com/IBM-Cloud/terraform-cost-estimator/utils"
	"go.uber.org/zap"
)

type Config struct {
	IAMAccessToken           string
	IamAPIKey                string
	GlobalCatalogAPIEndpoint string
	BSSEndpoint              string
}
type CostV1Config struct {
	*Config
}

var generation int

type costFunction func(Resource, string, int) (*BillOfMaterial, float64, error)
type pf map[string]costFunction

var ResourceFuncMap = pf{
	// "ibm_is_volume":      getVolumeCost,
	// "ibm_is_vpn_gateway": getVpnCost,
	// "ibm_is_image":       getVolumeCost,
}

type incCostFunction func(*zap.Logger, ResourceConf, string) (float64, error)
type cf map[string]incCostFunction

var IncCostFuncMap = cf{
	"ibm_is_volume":                     getVolumeCost2,
	"ibm_is_vpn_gateway":                getVpnCost2,
	"ibm_is_image":                      getImageCost2,
	"ibm_container_cluster":             containerCost,
	"ibm_container_worker_pool":         workerPoolContainerCost,
	"ibm_is_lb":                         getLoadBalancerCost,
	"ibm_container_vpc_cluster":         vpcContainerCost,
	"ibm_container_vpc_worker_pool":     vpcContainerCost,
	"ibm_is_vpc":                        getVpcCost,
	"ibm_is_instance":                   getInstanceCost,
	"ibm_service_instance":              serviceInstanceCost,
	"ibm_resource_instance":             serviceInstanceCost,
	"ibm_is_subnet":                     getSubnetCost,
	"ibm_is_floating_ip":                getIPCost,
	"ibm_satellite_cluster":             getSatelliteCluster,
	"ibm_satellite_cluster_worker_pool": getSatelliteCluster,
}

func NewTFCostClient(config *Config) (*CostV1Config, error) {
	return &CostV1Config{config}, nil
}

//GetCost will take plan file of a terraform template and returns the properties of estimated cost as  BOM...
func (costConfig *CostV1Config) GetCost(planFile string) (BOM, error) {

	logger, err := utils.NewLogger()
	if err != nil {
		panic(err)
	}
	logger.Info("Entry:GetCost")

	planData, _ := ioutil.ReadFile(planFile)
	data := Planstruct{}
	bom := BOM{}
	err = json.Unmarshal([]byte(planData), &data)
	if err != nil {
		logger.Error("Error while Unmarshalling Plan Data", zap.Error(err))
		return BOM{}, err
	}

	token := costConfig.IAMAccessToken
	cost, bom, err := calculateCost(data, token, logger)
	if err != nil {
		return bom, err
	}

	// bom.Country = "USA"
	// bom.Currency = "USD"
	bom.TotalCost = cost
	logger.Info("Exit:GetCost")

	return bom, nil

}

//calculateCost will take the plandata, parse the configuration and returns total cost with a BOM
func calculateCost(planData Planstruct, token string, logger *zap.Logger) (float64, BOM, error) {
	logger.Info("Entry:calculateCost")
	var cost float64
	if planData.Configuration.ProviderConfig.Ibm.Expressions.Generation.ConstantValue == 0 {
		generation = 2
	} else {
		generation = planData.Configuration.ProviderConfig.Ibm.Expressions.Generation.ConstantValue
	}
	bom := BOM{
		Country:  "USA",
		Currency: "USD",
	}
	cost = 0

	//Range through all resource in terraform template
	for _, resdata := range planData.PlannedValues.RootModule.Resources {

		resourceType := resdata.Type

		//check for terraform resource type in and call respective functions
		for k, _ := range ResourceFuncMap {
			if k == resourceType {

				billdata, resourceCost, err := ResourceFuncMap[resourceType](resdata, token, generation)
				if err != nil || billdata == nil {
					logger.Error("Error while getting cost from mapped function", zap.Error(err))
					continue
				}
				//add the bill data to list of BOM
				bom.AddItem(*billdata)

				cost += resourceCost
			}

		}

	}

	//Use this for incremental Cost
	for _, reschanges := range planData.ResourceChanges {
		resourceType := reschanges.Type

		//check for terraform resource type in function map and call respective functions
		for k, _ := range IncCostFuncMap {
			if k == resourceType {
				actions := reschanges.Change.Actions
				var before, after float64
				var err error

				//if it is create new resource or no change
				if len(actions) == 1 && (actions[0] == "create" || actions[0] == "no-op") {
					after, err = IncCostFuncMap[resourceType](logger, reschanges.Change.After, token)

					if err != nil {
						logger.Error("Error while trying to get after activity for create actions", zap.Error(err))
						// return 0, bom, err
						continue
					}

				}

				//if it is delete existing resource
				if len(actions) == 1 && (actions[0] == "delete") {
					before, err = IncCostFuncMap[resourceType](logger, reschanges.Change.Before, token)
					if err != nil {
						logger.Error("Error while trying to get before activity for delete actions", zap.Error(err))
						// return 0, bom, err
						continue
					}
				}
				//if change in existing resource
				if len(actions) == 2 || stringInSlice("update", actions) {
					before, err = IncCostFuncMap[resourceType](logger, reschanges.Change.Before, token)
					if err != nil {
						logger.Error("Error while trying to get before activity for update actions", zap.Error(err))
						// return 0, bom, err
						continue
					}
					after, err = IncCostFuncMap[resourceType](logger, reschanges.Change.After, token)
					if err != nil {
						logger.Error("Error while trying to get after activity for update actions", zap.Error(err))
						// return 0, bom, err
						continue
					}
				}

				billdata := BillOfMaterial{}
				billdata.AddIncrementalCostData(reschanges, before, after)

				//add the bill data to list of BOM
				bom.AddItem(billdata)
				cost += billdata.CurrLineItemTotal

			}
		}

	}

	logger.Info("Exit:calculateCost")
	return cost, bom, nil
}

//AddItem to the BOM
func (bom *BOM) AddItem(item BillOfMaterial) []BillOfMaterial {
	bom.Lineitem = append(bom.Lineitem, item)
	return bom.Lineitem
}

//AddMetric to the Request
func (req *Request) AddMetric(item Metrics) []Metrics {
	req.Measures = append(req.Measures, item)
	return req.Measures
}

//AddLineItemData to the BOM
func (billdata *BillOfMaterial) AddLineItemData(resdata Resource, planID string, cost float64) BillOfMaterial {
	billdata.Quantity = 1
	billdata.ID = resdata.Name
	billdata.TerraformItemID = resdata.Type
	billdata.PlanID = planID
	billdata.LineItemTotal = cost
	billdata.Title = resourceMap[resdata.Type]
	return *billdata
}

//AddLineItemData to the BOM
func (billdata *BillOfMaterial) AddIncrementalCostData(resdata ResourceChanges, before, after float64) BillOfMaterial {
	billdata.Quantity = 1
	billdata.ID = resdata.Name
	billdata.TerraformItemID = resdata.Type
	billdata.CurrLineItemTotal = after
	billdata.PrevLineItemTotal = before
	billdata.ChangeLineItemTotal = after - before
	billdata.Title = resourceMap[resdata.Type]
	return *billdata
}

//AddDependencyData to the BOM
func (billdata *BillOfMaterial) AddDependencyData(quantity int, title string, cost float64) BillOfMaterial {

	dependsItem := DependencyResource{
		Quantity: quantity,
		Title:    title,
		Cost:     cost,
	}
	billdata.ResourceBreakdown = append(billdata.ResourceBreakdown, dependsItem)
	return *billdata
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
