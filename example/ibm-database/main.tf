########################################################
# Create VM configured to access ICD database
########################################################
terraform {
  required_providers {
    ibm = {
      source = "IBM-Cloud/ibm"
      version = "~> 1.15.0"
    }
  }
}

data "ibm_resource_group" "group" {
  name = "Default"
}

resource "ibm_database" "test_acc" {
  resource_group_id = data.ibm_resource_group.group.id
  name              = "demo-postgres"
  service           = "databases-for-postgresql"
  plan              = "standard"
  location          = "eu-gb"
  adminpassword     = "adminpassword"

  whitelist {
    address     = "${ibm_compute_vm_instance.webapp1[0].ipv4_address}/32"
    description = ibm_compute_vm_instance.webapp1[0].hostname
  }

  tags = ["tag1", "tag2"]

  // adminpassword                = "password12"
  members_memory_allocation_mb = 3072
  members_disk_allocation_mb   = 20480

  users {
    name     = "user123"
    password = "password12"
  }
}

