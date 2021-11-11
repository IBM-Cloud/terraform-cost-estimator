package main

import (
	"fmt"
	"log"
	"os"

	costcalculator "github.com/IBM-Cloud/terraform-cost-estimator/api"
	"go.uber.org/zap"
)

func main() {

	var planFile string
	//Add your bearer token for testing
	token := ""
	planFile = "testplan.json"
	costClient, err := costcalculator.NewTFCostClient(&costcalculator.Config{
		IAMAccessToken: token,
	})

	if err != nil {
		log.Fatal("Error while generating Client for Cost Estimator", zap.Error(err))
		os.Exit(1)
	}
	bom, err := costClient.GetCost(planFile)
	if err != nil {
		log.Fatal("Error while generating Cost", zap.Error(err))
		os.Exit(1)
	}

	fmt.Println("\nTotal cost= ", bom.TotalCost)
	fmt.Println("\nBOMnew= ", bom.Lineitem)

}
