package helpers

import (
	"fmt"

	costcalculator "github.com/IBM-Cloud/terraform-cost-estimator/ibm"
)

//DisplayTable ...
type DisplayTable struct {
	Resource  string `header:"resource"`
	LocalName string `header:"local name"`
	Title     string `header:"title"`
	//EstimatedCost         float64 `header:"estimated cost"`
	CurrentEstimatedCost  string `header:"Current cost"`
	PreviousEstimatedCost string `header:"Previous cost"`
	ChangedEstimatedCost  string `header:"Changed cost"`
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
		instance.CurrentEstimatedCost = fmt.Sprintf("%.2f", item.CurrLineItemTotal)
		instance.PreviousEstimatedCost = fmt.Sprintf("%.2f", item.PrevLineItemTotal)
		instance.ChangedEstimatedCost = fmt.Sprintf("%.2f", item.ChangeLineItemTotal)
		if item.RateCardCost {
			instance.CurrentEstimatedCost = instance.CurrentEstimatedCost + " *"
			instance.PreviousEstimatedCost = instance.PreviousEstimatedCost + " *"
			instance.ChangedEstimatedCost = instance.ChangedEstimatedCost + " *"
		}
		resourceMap = append(resourceMap, instance)

	}
	return resourceMap
}
