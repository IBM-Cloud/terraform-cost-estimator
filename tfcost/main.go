package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/IBM-Cloud/terraform-cost-estimator/authentication"
	"github.com/IBM-Cloud/terraform-cost-estimator/helpers"
	costcalculator "github.com/IBM-Cloud/terraform-cost-estimator/ibm"
	"github.com/fatih/color"
	"github.com/kataras/tablewriter"
	"github.com/landoop/tableprinter"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "tfcost"
	app.Usage = "A command line tool to calaculate cost of Terraform resources"
	app.Action = func(c *cli.Context) error {
		fmt.Println("Estimate cost of the resources provisioned")
		return nil
	}

	app.Commands = []cli.Command{
		{
			Name:    "plan",
			Aliases: []string{"p"},
			Usage:   "path to plan json file",
			Flags: []cli.Flag{
				cli.BoolFlag{Name: "json, j", Usage: "print BOM in json format"},
			},
			Action: func(c *cli.Context) error {
				tokenResp, err := authentication.AuthenticateAPIKey()
				if err != nil {
					//handle postform error
					if strings.Contains(err.Error(), "Insufficient credentials") {
						err = fmt.Errorf("%s (please export IC_API_KEY )", err.Error())
						log.Println("Error in IAM Authentication", err)

					}
					return nil

				}
				token := "Bearer " + tokenResp.AccessToken

				costClient, err := costcalculator.NewTFCostClient(&costcalculator.Config{
					IAMAccessToken: token,
				})
				if err != nil {
					log.Println("Cannot initialise config")
					return nil
				}
				planfile := c.Args().First()
				if planfile == "" {
					log.Fatal("Error invalid argument plan.json. Please enter a valid path to plan.json file or check usage")
					return nil
				}
				bom, err := costClient.GetCost(planfile)
				if err != nil {
					log.Println(err)
				}

				// fmt.Println("\nBOMnew= ", bom.Lineitem)
				notice := color.New(color.Bold, color.FgGreen).PrintlnFunc()
				notice("\nNOTE:\n cost displayed here is estimated on monthly basis in USD, it is not the actual cost")
				fmt.Print("\nRate Card Version: V1.0")
				notice("\n* represents cost acquired from the rate card\n")

				// if json flag enabeled
				if c.Bool("json") || c.Bool("j") {
					outputByte, err := json.Marshal(bom)
					if err != nil {
						log.Println(err)
						return nil
					}

					outputJSON := string(outputByte[:])
					fmt.Println("\njson:", outputJSON)
					err = ioutil.WriteFile("cost.json", outputByte, 0644)
					return nil
				}

				// Default o/p type is tabular format
				printer := tableprinter.New(os.Stdout)

				// Optionally, customize the table, import of the underline 'tablewriter' package is required for that.
				printer.BorderTop, printer.BorderBottom, printer.BorderLeft, printer.BorderRight = true, true, true, true
				printer.CenterSeparator = "│"
				printer.ColumnSeparator = "│"
				printer.RowSeparator = "─"
				printer.HeaderBgColor = tablewriter.BgBlackColor
				printer.HeaderFgColor = tablewriter.FgGreenColor

				// Print the slice of structs as table, as shown above.
				printer.Print(helpers.GetTable(bom.Lineitem))
				bom.TotalCost, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", bom.TotalCost), 64)
				notice("\nTotal Estimated Cost: $", bom.TotalCost)

				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
