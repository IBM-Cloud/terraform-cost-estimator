# IBM Cloud Cost estimator for Terraform

A Go implementation for cost estimation of IBM Cloud resources provisioned using [terraform-provider-ibm](https://github.com/IBM-Cloud/terraform-provider-ibm).  The project uses the plan json file (produced by the [Terraform plan command](https://www.terraform.io/internals/json-format)), to prepare the bill-of-material with estimated cost (monthly).  The project implements an SDK, and a CLI.

The full list of supported resources and common assumptions are described [here](https://github.com/IBM-Cloud/terraform-cost-estimator/blob/main/supportedResources.md). 

### Disclaimer

* This is `experimental beta` release.  
* The generated bill-of-material includea an estimated cost and not the actual cost.
* The estimated cost does not include usage-cost.

### Prerequisites

* Terraform 0.12 or higher
* IBM Cloud provider for Terraform - [terraform-provider-ibm](https://github.com/IBM-Cloud/terraform-provider-ibm)
* Generate the Terraform plan json file (tfplan.json) from a working Terraform template (with IBM-Cloud resources)
  1. Run the `terraform plan` command from the current working directory (with the .tf files) to generate the tfplan.binary file.
     ```bash
     terraform plan --out tfplan.binary
     ```
  2. Run the `terraform show` command using the tfplan.binary file to generate the Terraform plan json file.
     ```bash
     terraform show -json tfplan.binary > tfplan.json
     ```
* You can try the workings of the tool using the sample tfplan.json in the `examples` folder.

---

## CLI Usage
The instructions to use the `IBM Cloud Cost estimator for Terraform` CLI

### Using the CLI

1. Download the latest release of the `tfcost` CLI from [here](https://github.com/IBM-Cloud/terraform-cost-estimator/releases) and place it your PATH
2. Run `tfcost` CLI with tfplan.json file as input
   ```bash
   export IC_API_KEY=your_ibmcloud_api_key
   
   tfcost plan tfplan.json
   ```

### Sample output
![Estimated cost](/image.png)

---

## SDK Usage

### Installing the SDK

1. Install the Go SDK using the following command
   ```bash
   go get github.com/IBM-Cloud/terraform-cost-estimator
   ```
2. Update the SDK to the latest version using the following command
   ```bash
   go get -u github.com/IBM-Cloud/terraform-cost-estimator
   ```
3. Vendor the `terraform-cost-estimator` in your project.
   ```bash
   go mod vendor
   ```
   If you get issue while vendoring go mod then export
   ```bash
   export GO111MODULE=on
   ```

### Using the SDK

```go
import "github.com/IBM-Cloud/terraform-cost-estimator/api"

func main(){

    costClient, err := costcalculator.NewTFCostClient(&costcalculator.Config{
	    IAMAccessToken: token,
    })

    bom, err := costClient.GetCost(tfplan.json)
    .....
}
```

#### Sample output - bill-of-materials

```json
{
    "country": "USA",
    "currency": "USD",
    "total_cost": "82.24",
    "line_item": [
        {
            "id": "is.instance",
            "terraform_item_id": "ibm_is_instance",
            "quantity": "1",
            "title": "testacc_instance",
            "plan_id": "66380d42-d4a9-4627-88fa-7b6631e54443 ",
            "short_desciption": "",
            "features": "",
            "estimate_type": "custom",
            "line_item_total": "84.24"
        }
    ]
}
```
---

## Contribution

You can contribute to the `IBM Cloud Cost estimator for Terraform` tool, refer to the details [here](https://github.com/IBM-Cloud/terraform-cost-estimator/blob/main/CONTRIBUTING.md)

### Prerequisites

* [Go](http://www.golang.org) version 1.8+ ot higher
* Setup a [GOPATH](http://golang.org/doc/code.html#GOPATH), and add `$GOPATH/bin` to your `$PATH`.

### Build 'IBM Cloud Cost estimator for Terraform'

* Update the latest go.mod dependencies 
  ```
  go get github.com/IBM-Cloud/terraform-cost-estimator
  go mod vendor
  ```
* Run `make build` to generate tfcost binary in the `$GOPATH/bin` directory.
  ```sh
  go build
  ...
  $GOPATH/bin/terraform-cost-estimator
  ...
  ```
---

## Rate Card Support

The full list of supported resources for cost-estimation is described [here](https://github.com/IBM-Cloud/terraform-cost-estimator/blob/main/supportedResources.md).  The 'IBM Cloud Cost estimator for Terraform' uses IBM Cloud BSS SDKs and other IBM Cloud Platform APIs to dynamically fetch the service-plan for the IBM Cloud services - to compute the estimated cost. 

The 'IBM Cloud Cost estimator for Terraform' supports the ability to override the service-plan using a custom [rate card](https://github.com/IBM-Cloud/terraform-cost-estimator/ibm/rate_card.json).   A default rate card is bundled with the tool. 

You can customize the `tfcost` tool (CLI or SDK) by providing your own `rate card`. 

```sh
export RATECARD=path_to_your_rate_card.json
export IC_API_KEY=your_ibmcloud_api_key

tfcost plan tfplan.json
```

---
## Releases

Please refer to [here](https://github.com/IBM-Cloud/terraform-cost-estimator/releases) for details.

---
# Issues, defects and feature requests

Submit your issue, defects, or feature requests [here](https://github.com/IBM-Cloud/terraform-cost-estimator/issues).
---
