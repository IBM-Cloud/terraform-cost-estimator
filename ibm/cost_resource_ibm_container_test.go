package costcalculator

import (
	"os"
	"testing"
)

func TestContainerCost(t *testing.T) {

	token := os.Getenv("IC_IAM_TOKEN")
	costClient, _ := NewTFCostClient(&Config{
		IAMAccessToken: token,
	})
	bom, err := costClient.GetCost("testplan.json")
	if err != nil {
		t.Errorf("failed to get container cost. recieved error from getCost %v", err.Error())

	}
	for _, lineItem := range bom.Lineitem {
		if lineItem.TerraformItemID == "ibm_container_cluster" {
			if lineItem.CurrLineItemTotal != 0.0 {
				t.Log("success")
			} else {
				t.Errorf("failed to get container cost. expected %v recieved %v", "non zero", lineItem.CurrLineItemTotal)
			}
		}
	}

}
