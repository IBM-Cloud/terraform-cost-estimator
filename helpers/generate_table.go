package helpers

import (
	costcalculator "github.com/IBM-Cloud/terraform-cost-estimator/ibm"
)

//DisplayTable ...
type DisplayTable struct {
	Resource  string `header:"resource"`
	LocalName string `header:"local name"`
	Title     string `header:"title"`
	//EstimatedCost         float64 `header:"estimated cost"`
	CurrentEstimatedCost  float64 `header:"Current cost"`
	PreviousEstimatedCost float64 `header:"Previous cost"`
	ChangedEstimatedCost  float64 `header:"Changed cost"`
}

//GetTable ...
func GetTable(bom []costcalculator.BillOfMaterial) []DisplayTable {
	resourceMap := make([]DisplayTable, 0, len(bom))

	for _, item := range bom {
		var instance DisplayTable
		instance.Resource = item.TerraformItemID
		instance.LocalName = item.ID
		instance.Title = item.Title
		//instance.EstimatedCost = item.LineItemTotal
		instance.CurrentEstimatedCost = item.CurrLineItemTotal
		instance.PreviousEstimatedCost = item.PrevLineItemTotal
		instance.ChangedEstimatedCost = item.ChangeLineItemTotal
		resourceMap = append(resourceMap, instance)

	}
	return resourceMap
}
