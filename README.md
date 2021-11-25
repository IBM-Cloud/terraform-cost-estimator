# IBM Cloud Terraform Cost estimator 

This project provides the Go implementation for parsing the input terraform plan data, and provides an estimated cost of the template. It can be used as an SDK as well as a CLI

## Installing the SDK

1. Install the SDK using the following command

```bash
go get github.com/IBM-Cloud/terraform-cost-estimator
```

2. Update the SDK to the latest version using the following command

```bash
go get -u github.com/IBM-Cloud/terraform-cost-estimator
```


## Using the SDK

>You must have a working terraform template with IBM cloud resources.

The SDK has ```examples``` folder which cites few examples on how to use the SDK.

```go
import "github.com/IBM-Cloud/terraform-cost-estimator/api"

func main(){

    bom, err := costcalculator.GetCost(planFile.json, iam-oauth-token)

    costClient, err := costcalculator.NewTFCostClient(&costcalculator.Config{
	    IAMAccessToken: token,
    })
    .....
}
```

### Steps to generate the json planFile
 Note: running terraform commands requires terraform and its plugins to be present as a prerequisite
1. Inside the directory which contains the tf files do
```bash
terraform plan --out tfplan.binary
```
2. After generating the binary file generate the respective plan json file using
```bash
terraform show -json tfplan.binary > tfplan.json
```

## Sample BOM output

```json{
    "country": "USA",
    "currency": "USD",
    "total_cost": "82.24",
    "line_item": [
        {
            "id": "is.instance",
            "terraform_item_id": "ibm_is_instance",
            "quantity": "1",
            "title": "testacc_instance",
            "plan_id": "66380d42-d4a9-4627-88fa-7b6631e5bd63 ",
            "short_desciption": "",
            "features": "",
            "estimate_type": "custom",
            "line_item_total": "84.24"
        }
    ]
}```

List of terraform Resources Supported by the SDK are

- ibm_is_instance
- ibm_is_volume
- ibm_is_lb
- ibm_is_floating_ip
- ibm_is_vpn_gateway
- ibm_is_image      
- ibm_is_vpc   
- ibm_is_subnet     
- ibm_container_cluster
- ibm_container_vpc_cluster
- ibm_container_worker_pool
- ibm_satellite_cluster
- ibm_service_instance
```
# Terraform cost CLI
The following are the instructions and steps to use cli for IBM-Cloud terraform cost calculator

## Using the CLI
### Requirements
>You must have a working Terraform template with IBM-Cloud resources


The CLI has ```examples``` folder which cites few examples plan.json that you can use as input.

```bash
go mod vendor
```
If you get issue while vendoring go mod then export
```bash

export GO111MODULE=on
export GOPRIVATE=*.ibm.com
```

Steps to generate the json planFile
1. Inside the directory which contains the tf files do
```bash
terraform plan --out tfplan.binary
```
2. After generating the binary file generate the respective plan json file using
```bash
terraform show -json tfplan.binary > tfplan.json
```

Now To generate the cost output from planJSON input do

```bash
export IC_API_KEY=your_ibmcloud_api_key
cd tfcost
go run main.go plan=../example/tfplan.json
```

## Sample output
![Estimated cost](/image.png)


