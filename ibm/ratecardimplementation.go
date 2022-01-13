package costcalculator

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"go.uber.org/zap"
)

type RateCard struct {
	Service_name   string  `json:"service_name"`
	Plan           string  `json:"plan_id"`
	Estimated_rate float64 `json:"estimated_rate"`
	Units          string  `json:"unit"`
	Currency       string  `json:"currency"`
	Unit_quantity  string  `json:"unit_quantity"`
	Usage_based    bool    `json:"usage_based"`
}

func ratecard(logger *zap.Logger, resource string, planData Planstruct) (ResourceChanges, float64, error) {
	rateCardFilename := "../ibm/rate_card.json"
	if os.Getenv("RATECARD") != "" {
		rateCardFilename = os.Getenv("RATECARD")
	}
	rateCard, _ := ioutil.ReadFile(rateCardFilename)
	card := []RateCard{}
	err := json.Unmarshal([]byte(rateCard), &card)
	if err != nil {
		logger.Error("Error while Unmarshalling Plan Data", zap.Error(err))
		fmt.Print(err)
	}
	classic_rateCard, _ := ioutil.ReadFile("../ibm/classic_vm.json")
	classic_card := ClassicRateCard{}
	err = json.Unmarshal([]byte(classic_rateCard), &classic_card)
	if err != nil {
		//logger.Error("Error while Unmarshalling Plan Data", zap.Error(err))
		fmt.Print(err)
	}

	var profile string
	var rate_card_profile string
	var instance_type = resource

	for _, item := range planData.ResourceChanges {

		if instance_type == "ibm_compute_vm_instance" && item.Type == instance_type {
			var cost float64
			var temp float64
			logger.Info("Entry:getComputeVMinstanceCost")
			//when flavor is present
			if item.Change.After.FlavorKeyName != "" {

				for _, classic_item := range classic_card.Flavors {
					if item.Change.After.FlavorKeyName == classic_item.Flavor.KeyName {
						if classic_item.Flavor.TotalMinimumRecurringFee != "" {
							temp, _ = strconv.ParseFloat(classic_item.Flavor.TotalMinimumRecurringFee, 64)
							logger.Info("Exit:getComputeVMinstanceCost")
							return item, temp, nil
						}
					}
				}
			} else {
				//when flavor is absent
				for _, classic_item := range classic_card.Memory {
					if item.Change.After.Memory == classic_item.Template.MaxMemory {
						temp, _ = strconv.ParseFloat(classic_item.ItemPrice.RecurringFee, 64)
						cost += temp
					}
				}
				for _, classic_item := range classic_card.Processors {
					if item.Change.After.Cores == classic_item.Template.StartCpus {
						if !classic_item.Template.DedicatedHost {
							temp, _ = strconv.ParseFloat(classic_item.ItemPrice.RecurringFee, 64)
							cost += temp
						}
					}
				}
				for _, classic_item := range classic_card.OperatingSystems {
					if item.Change.After.OperatingSystemReferenceCode == classic_item.Template.OperatingSystem {
						temp, _ = strconv.ParseFloat(classic_item.ItemPrice.RecurringFee, 64)
						cost += temp
					}
				}
				for _, classic_item := range classic_card.NetworkComponents {
					if !item.Change.After.DedicatedHostFlag {
						if (item.Change.After.NetworkSpeed == 1000 && classic_item.Template.NetworkComponent[0].MaxSpeed == 1000) && (classic_item.Template.PrivateNetworkOnly == false && item.Change.After.PrivateNetworkOnly == false) {
							temp, _ = strconv.ParseFloat(classic_item.ItemPrice.RecurringFee, 64)
							cost += temp
						} else if (item.Change.After.NetworkSpeed == 1000 && classic_item.Template.NetworkComponent[0].MaxSpeed == 1000) && (classic_item.Template.PrivateNetworkOnly == true && item.Change.After.PrivateNetworkOnly == true) {
							temp, _ = strconv.ParseFloat(classic_item.ItemPrice.RecurringFee, 64)
							cost += temp
						} else {
							continue
						}
					}

				}
				logger.Info("Exit:getComputeVMinstanceCost")
				return item, cost, nil
			}
		}
		image := ""
		for _, card_item := range card {
			card_elements := strings.Split(card_item.Plan, ".")

			if item.Type == instance_type {

				//virtual server for vpc
				if instance_type == "ibm_is_instance" && card_elements[0] == "ibm_is_instance" {
					logger.Info("Entry:getInstanceCost")
					profile = item.Change.After.Profile
					rate_card_image := strings.Split(card_item.Plan, ".")[1]
					rate_card_profile = strings.Split(card_item.Plan, ".")[2]
					if item.Change.After.ImageID == "51af68c9-5558-4425-825a-f9243a3b2c6c" || item.Change.After.ImageID == "624cde4a-b4fe-4426-8f60-150a019a67f9" || item.Change.After.ImageID == "a7a0626c-f97e-4180-afbe-0331ec62f32a" {
						image = "windows"
					} else if item.Change.After.ImageID == "54c1ba68-6d29-42e5-9ca7-e5f4a62c1503" {
						image = "rhel"
					} else {
						image = "centos"
					}
					if profile == rate_card_profile && image == rate_card_image {
						logger.Info("Exit:getInstanceCost")
						return item, card_item.Estimated_rate, nil
					} else {
						return item, 0, errors.New("profile not found")
					}

					//kubernetes classic infra
				} else if instance_type == "ibm_container_cluster" && card_elements[0] == "ibm_container_cluster" || instance_type == "ibm_container_worker_pool" && card_elements[0] == "ibm_container_worker_pool" {
					logger.Info("Entry:getKubernetesCost")
					profile = item.Change.After.MachineType

					rate_card_profile = strings.SplitAfterN(card_item.Plan, ".", 4)[3]

					hardware := strings.ToLower(item.Change.After.Hardware)

					card_hardware := strings.SplitAfterN(card_item.Plan, ".", 4)[2]

					card_hardware = strings.Trim(card_hardware, ".")

					if profile == rate_card_profile {
						if hardware == card_hardware {
							if instance_type == "ibm_container_cluster" {
								logger.Info("Exit:getKubernetesCost")
								return item, card_item.Estimated_rate * float64(item.Change.After.DefaultPoolSize), nil
							} else if instance_type == "ibm_container_worker_pool" {
								return item, card_item.Estimated_rate * float64(item.Change.After.SizePerZone), nil
							} else {
								return item, 0, errors.New("incorrect kubernetes configuration")
							}
						} else {
							return item, 0, errors.New("hardware information not found")
						}
					}

					//kubernetes vpc cluster (rate changes with the Operating System)
				} else if instance_type == "ibm_container_vpc_cluster" && card_elements[0] == "ibm_container_vpc_cluster" {
					logger.Info("Entry:getKubernetesCost")
					profile = item.Change.After.Flavour

					rate_card_profile = strings.Split(card_item.Plan, ".")[3]

					if profile == rate_card_profile {
						logger.Info("Exit:getKubernetesCost")
						return item, card_item.Estimated_rate * float64(item.Change.After.WorkerCount), nil
					} else {
						return item, 0, errors.New("profile not found")
					}

					//vpc cluster worker pool
				} else if instance_type == "ibm_container_vpc_worker_pool" && card_elements[0] == "ibm_container_vpc_worker_pool" {
					logger.Info("Entry:getVPCclusterCost")
					profile = item.Change.After.Flavour

					rate_card_profile = strings.SplitAfterN(card_item.Plan, ".", 3)[2]

					if profile == rate_card_profile {
						logger.Info("Exit:getVPCclusterCost")
						return item, card_item.Estimated_rate * float64(item.Change.After.WorkerCount), nil
					} else {
						return item, 0, errors.New("profile not found")
					}

					//app config environment
				} else if instance_type == "ibm_app_config_environment" && card_elements[0] == "ibm_app_config_environment" {
					logger.Info("Entry:getAppConfigCost")
					if card_item.Usage_based {
						err = errors.New("it is a usage based resource")
						return item, 0, err
					} else {
						logger.Info("Exit:getAppConfigCost")
						return item, card_item.Estimated_rate, nil
					}

					//app config feature
				} else if instance_type == "ibm_app_config_feature" && card_elements[0] == "ibm_app_config_feature" {
					logger.Info("Entry:getAppConfigCost")
					if card_item.Usage_based {
						err = errors.New("it is a usage based resource")
						return item, 0, err
					} else {
						logger.Info("Exit:getAppConfigCost")
						return item, card_item.Estimated_rate, nil
					}

					//resource instance (event notifications and cloud object storage (cos instance))
				} else if instance_type == "ibm_resource_instance" && card_elements[0] == "ibm_resource_instance" {

					if item.Change.After.Service == "cloud-object-storage" && card_elements[1] == "cos_instance" {
						logger.Info("Entry:getCosInstaneCost")
						if item.Change.After.Plan == card_elements[2] {

							if card_item.Usage_based {
								err = errors.New("it is a usage based resource")
								return item, 0, err

							} else {
								logger.Info("Exit:getCosInstaneCost")
								return item, card_item.Estimated_rate, nil
							}
						} else {
							return item, 0, errors.New("invalid configuration for " + instance_type)
						}

					} else if item.Change.After.Service == "event-notifications" && card_elements[1] == "event-notifications" {
						logger.Info("Entry:getEventNotificationCost")
						if item.Change.After.Plan == card_elements[2] {

							if card_item.Usage_based {
								err = errors.New("it is a usage based resource")
								return item, 0, err

							} else {
								logger.Info("Exit:getEventNotificationCost")
								return item, card_item.Estimated_rate, nil
							}
						} else {
							return item, 0, errors.New("invalid configuration for " + instance_type)
						}
					}

					//secondary subnets
				} else if instance_type == "ibm_subnet" && card_elements[0] == "ibm_subnet" {

					if strings.ToLower(item.Change.After.Type) == card_elements[1] {
						capacity, _ := strconv.Atoi(card_elements[2])
						if item.Change.After.Capacity == capacity {
							logger.Info("Entry:getSecondarySubnetCost")
							logger.Info("Exit:getSecondarySubnetCost")
							return item, card_item.Estimated_rate, nil
						} else {
							return item, 0, errors.New("capacity not found")
						}
					}

					//every other resource
				} else if instance_type == "ibm_cis" && card_elements[0] == "ibm_cis" {

					if item.Change.After.Plan == card_elements[1] {
						logger.Info("Entry:getCISCost")
						if card_item.Usage_based {
							err = errors.New("it is a usage based resource")
							return item, 0, err

						} else {
							logger.Info("Exit:getCISCost")
							return item, card_item.Estimated_rate, nil
						}
					} else {
						return item, 0, errors.New("plan not found")
					}

					// this is for dedicated host for VPC
				} else if instance_type == "ibm_is_dedicated_host" && card_elements[0] == "ibm_is_dedicated_host" {
					logger.Info("Entry:getDedicatedHostCost")
					profile = item.Change.After.Profile

					rate_card_profile = strings.Split(card_item.Plan, ".")[2]
					//fmt.Println(profile, rate_card_profile)

					if profile == rate_card_profile {
						logger.Info("Exit:getDedicatedHostCost")
						return item, card_item.Estimated_rate, nil
					} else {
						return item, 0, errors.New("profile not found")
					}

					//this is for every other resource
				} else {
					//logger.Info("Entry:getRateCardCost")
					if instance_type == card_elements[0] {
						if card_item.Usage_based {
							err = errors.New("it is a usage based resource")
							return item, 0, err

						} else {
							//logger.Info("Exit:getRateCardCost")
							return item, card_item.Estimated_rate, nil
						}
					}
				}
			}
		}
	}
	return ResourceChanges{}, 0, errors.New("no cost found")
}
