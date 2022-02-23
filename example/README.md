How to use terraform cost estimator.You can take the template script provided in examples directory.

before running script export you cloud api key
```bash
export IC_API_KEY=your_ibmcloud_api_key
```
follow below steps for estimate cost of template to be provision
1. Inside the terraform directory which contains the tf files run
```bash
terraform init 
```
and then 
```bash
terraform plan --out tfplan.binary
```
2. After generating the binary file generate the respective plan json file using
```bash
terraform show -json tfplan.binary > tfplan.json
```

3. run cost estimator cli using 
```bash
tfcost plan tfplan.json
```