provider "ibm" {
  generation =2
 
}
terraform {
  required_providers {
    ibm = {
      source = "IBM-Cloud/ibm"
      version = "~> 1.15.0"
    }
  }
}
resource "ibm_container_cluster" "testacc_cluster" {
  name            = "test"
  datacenter      = "dal10"
  machine_type    = "b3c.4x16"
  hardware        = "shared"
  public_vlan_id  = "vlan"
  private_vlan_id = "vlan"
  subnet_id       = ["7654643"]

  default_pool_size = 1

  webhook {
    level = "Normal"
    type  = "slack"
    url   = "https://hooks.abc.com/"
  }
}

resource "ibm_container_worker_pool" "testacc_workerpool" {
  worker_pool_name = "mypool"
  machine_type     = "u2c.2x4"
  cluster          = ibm_container_cluster.testacc_cluster.id
  size_per_zone    = 1
  hardware         = "shared"
  disk_encryption  = "true"
  region           = "eu-de"

  //User can increase timeouts 
    timeouts {
      update = "180m"
    }
}

resource "ibm_container_worker_pool_zone_attachment" "test_zone" {
  cluster         = "mycluster"
  worker_pool     = element(split("/",ibm_container_worker_pool.testacc_workerpool.id),1)
  zone            = "dal12"
  private_vlan_id = "2333267"
  public_vlan_id  = "2333265"

  //User can increase timeouts
  timeouts {
      create = "90m"
      update = "3h"
      delete = "30m"
    }
}

data "ibm_resource_group" "group" {
  name = "Default"
}

resource "ibm_resource_instance" "resource_instance" {
  name              = "test"
  service           = "cloud-object-storage"
  plan              = "lite"
  location          = "global"
  resource_group_id = data.ibm_resource_group.group.id
  tags              = ["tag1", "tag2"]

  //User can increase timeouts
  timeouts {
    create = "15m"
    update = "15m"
    delete = "15m"
  }
}

resource "ibm_is_lb" "lb" {
  name    = "loadbalancer1"
  subnets = ["04813493-15d6-4150-9948-6246342"]
}