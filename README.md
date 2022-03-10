# IBM Cloud Terraform Cost estimator 

This project provides the Go implementation for Plan based terraform cost estimation for terraform-provider-ibm by parsing the input terraform plan data, and provides an estimated cost of the template. It can be used as an SDK as well as a CLI
# tfcost CLI
The following are the instructions and steps to use cli for IBM-Cloud terraform cost calculator

## Prerequisites

>Must have installed- terraform 0.12+, terraform-provider-ibm (https://github.com/IBM-Cloud/terraform-provider-ibm/tree/v1.36.0)
>You must have a working Terraform template with IBM-Cloud resources

## Using the CLI

1. Download the latest release of the tfcost binary

2. Put the tfcost binary in your PATH

3. Run ```tfcost``` command inside your terraform template directory over a tfplan.json file-


Steps to generate the tfplan.json planFile

1. Inside the directory which contains the tf files do
```bash
terraform plan --out tfplan.binary
```
2. After generating the binary file generate the respective plan json file using
```bash
terraform show -json tfplan.binary > tfplan.json
```

Now To generate the cost output from planJSON input run the tfcost command 

```bash
export IC_API_KEY=your_ibmcloud_api_key

tfcost plan tfplan.json
```

## Sample output
![Estimated cost](/image.png)



## Terraform Cost Estimator SDK

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

    costClient, err := costcalculator.NewTFCostClient(&costcalculator.Config{
	    IAMAccessToken: token,
    })

    bom, err := costClient.GetCost(planFile.json)
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


The following are the instructions and steps to use cli for IBM-Cloud terraform cost calculator

## Prerequisites

>Must have installed- terraform 0.12+, terraform-provider-ibm (https://github.com/IBM-Cloud/terraform-provider-ibm/tree/v1.36.0)
>You must have a working Terraform template with IBM-Cloud resources


1. Download the latest release of the tfcost binary

2. Put the tfcost binary in your PATH

3. Run inside your terraform template directory-


The Repository has ```examples``` folder which cites few examples plan.json that you can use as input.

```bash
go mod vendor
```
If you get issue while vendoring go mod then export
```bash
export GO111MODULE=on
```

Check Supported Resources docs [here](/supportedResources.md) 

## Sample output
![Estimated cost](/image.png)

## Developing the tool

update the latest go.mod dependencies inside ./tfcost directory if required

```
go get github.com/IBM-Cloud/terraform-cost-estimator
go mod vendor
```

If you wish to work on the tool, you'll first need [Go](http://www.golang.org) installed on your machine (version 1.8+ is *required*). You'll also need to correctly setup a [GOPATH](http://golang.org/doc/code.html#GOPATH), as well as adding `$GOPATH/bin` to your `$PATH`.

To compile the tool, run `make build`. This will build the tool and put the tfcost binary in the `$GOPATH/bin` directory.

```sh
go build
...
$GOPATH/bin/terraform-cost-estimator
...
```


## Release notes

Please refer to [here](https://github.com/IBM-Cloud/terraform-cost-estimator/releases) for details.


# Issues, defects and feature requests

Any issues/defects, or feature requests, please [file an issue](https://github.com/IBM-Cloud/terraform-cost-estimator/issues) if not raised before.