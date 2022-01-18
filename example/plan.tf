terraform {
  required_providers {
    ibm = {
      source = "IBM-Cloud/ibm"
      version = "1.38.0"
    }
  }
}

provider "ibm" {
  # Configuration options
}

resource "ibm_is_volume" "testacc_volume" {
  name     = "test-volume"
  profile  = "custom"
  zone     = "us-south-1"
  iops     = 10
  capacity = 200
  encryption_key = "crn:v1:bluemix:public:kms:us-south:a/dffc98a0f1f0f95f6613b3b752286b87:e4a29d1a-2ef0-42a6-8fd2-350deb1c647e:key:5437653b-c4b1-447f-9646-b2a2a4cd6179"
}