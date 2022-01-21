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
  encryption_key = ""
}