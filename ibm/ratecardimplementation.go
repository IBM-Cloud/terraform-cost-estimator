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

func ratecard(logger *zap.Logger, resource string, planData ResourceConf) (float64, error) {
	rateCardFilename := "../ibm/rate_card.json"
	if os.Getenv("RATECARD") != "" {
		rateCardFilename = os.Getenv("RATECARD")
	}
	rateCard, _ := ioutil.ReadFile(rateCardFilename)
	card := []RateCard{}
	err := json.Unmarshal([]byte(rateCard), &card)
	if err != nil {
		//logger.Error("Error while Unmarshalling Plan Data", zap.Error(err))
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

	//Virtual server on classic infrastructure: check for flavour, if absent, check for parameters such as processors, memory, OS and network component.
	//Then fetch the price

	if instance_type == "ibm_compute_vm_instance" {
		var cost float64
		var temp float64
		var flavor_cost float64
		logger.Info("Entry:getComputeVMinstanceCost")
		//when flavor is present
		if planData.FlavorKeyName != "" {

			for _, classic_item := range classic_card.Flavors {
				if planData.FlavorKeyName == classic_item.Flavor.KeyName {
					if classic_item.Flavor.TotalMinimumRecurringFee != "" {
						flavor_cost, _ = strconv.ParseFloat(classic_item.Flavor.TotalMinimumRecurringFee, 64)
						logger.Info("Exit:getComputeVMinstanceCost")
						return flavor_cost, nil
					}
				}
			}
		} else {
			//when flavor is absent
			if planData.Memory != 0 && planData.Cores != 0 && planData.OperatingSystemReferenceCode != "" && planData.NetworkSpeed != 0 {
				for _, classic_item := range classic_card.Memory {
					if planData.Memory == classic_item.Template.MaxMemory {
						temp, _ = strconv.ParseFloat(classic_item.ItemPrice.RecurringFee, 64)
						cost += temp
					}
				}
				for _, classic_item := range classic_card.Processors {
					if planData.Cores == classic_item.Template.StartCpus {
						if !classic_item.Template.DedicatedHost {
							temp, _ = strconv.ParseFloat(classic_item.ItemPrice.RecurringFee, 64)
							cost += temp
						}
					}
				}
				for _, classic_item := range classic_card.OperatingSystems {
					if planData.OperatingSystemReferenceCode == classic_item.Template.OperatingSystem {
						temp, _ = strconv.ParseFloat(classic_item.ItemPrice.RecurringFee, 64)
						cost += temp
					}
				}
				for _, classic_item := range classic_card.NetworkComponents {
					if !planData.DedicatedHostFlag {
						if (planData.NetworkSpeed == 1000 && classic_item.Template.NetworkComponent[0].MaxSpeed == 1000) && (classic_item.Template.PrivateNetworkOnly == false && planData.PrivateNetworkOnly == false) {
							temp, _ = strconv.ParseFloat(classic_item.ItemPrice.RecurringFee, 64)
							cost += temp
						} else if (planData.NetworkSpeed == 1000 && classic_item.Template.NetworkComponent[0].MaxSpeed == 1000) && (classic_item.Template.PrivateNetworkOnly == true && planData.PrivateNetworkOnly == true) {
							temp, _ = strconv.ParseFloat(classic_item.ItemPrice.RecurringFee, 64)
							cost += temp
						} else {
							continue
						}
					}

				}
			} else {
				logger.Error("insufficient parameters")
			}
			logger.Info("Exit:getComputeVMinstanceCost")
			return cost, nil
		}
	}
	image := ""
	for _, card_item := range card {
		card_elements := strings.Split(card_item.Plan, ".")

		//virtual server for vpc: checks for the OS, profile and fetches the cost
		if instance_type == "ibm_is_instance" && card_elements[0] == "ibm_is_instance" {
			logger.Info("Entry:getInstanceCost")
			profile = planData.Profile
			rate_card_image := strings.Split(card_item.Plan, ".")[1]
			rate_card_profile = strings.Split(card_item.Plan, ".")[2]
			if planData.ImageID == "51af68c9-5558-4425-825a-f9243a3b2c6c" || planData.ImageID == "624cde4a-b4fe-4426-8f60-150a019a67f9" || planData.ImageID == "a7a0626c-f97e-4180-afbe-0331ec62f32a" {
				image = "windows"
			} else if planData.ImageID == "54c1ba68-6d29-42e5-9ca7-e5f4a62c1503" {
				image = "rhel"
			} else {
				image = "centos"
			}
			if profile == rate_card_profile && image == rate_card_image {
				logger.Info("Exit:getInstanceCost")
				return card_item.Estimated_rate, nil
			}

			//kubernetes cluster and worker pool, primarily we see if its VPC or classic, then,
			//we check the profile,then the hardware whether it is dedicated or shared after this, price is fetched
		} else if instance_type == "ibm_container_cluster" && card_elements[0] == "ibm_container_cluster" || instance_type == "ibm_container_worker_pool" && card_elements[0] == "ibm_container_worker_pool" {
			logger.Info("Entry:getKubernetesCost")
			profile = planData.MachineType

			rate_card_profile = strings.SplitAfterN(card_item.Plan, ".", 4)[3]

			hardware := strings.ToLower(planData.Hardware)

			card_hardware := strings.SplitAfterN(card_item.Plan, ".", 4)[2]

			card_hardware = strings.Trim(card_hardware, ".")

			if profile == rate_card_profile {
				if hardware == card_hardware {
					if instance_type == "ibm_container_cluster" {
						logger.Info("Exit:getKubernetesCost")
						return card_item.Estimated_rate * float64(planData.DefaultPoolSize), nil
					} else if instance_type == "ibm_container_worker_pool" {
						return card_item.Estimated_rate * float64(planData.SizePerZone), nil
					}
				}
			}

			//kubernetes vpc cluster
		} else if instance_type == "ibm_container_vpc_cluster" && card_elements[0] == "ibm_container_vpc_cluster" {
			logger.Info("Entry:getKubernetesCost")
			profile = planData.Flavour

			rate_card_profile = strings.SplitAfterN(card_item.Plan, ".", 3)[2]

			if profile == rate_card_profile {
				logger.Info("Exit:getKubernetesCost")
				return card_item.Estimated_rate * float64(planData.WorkerCount), nil
			}

			//vpc cluster worker pool
		} else if instance_type == "ibm_container_vpc_worker_pool" && card_elements[0] == "ibm_container_vpc_worker_pool" {
			logger.Info("Entry:getVPCclusterCost")
			profile = planData.Flavour

			rate_card_profile = strings.SplitAfterN(card_item.Plan, ".", 3)[2]

			if profile == rate_card_profile {
				logger.Info("Exit:getVPCclusterCost")
				return card_item.Estimated_rate * float64(planData.WorkerCount), nil
			}

			//app config environment
		} else if instance_type == "ibm_app_config_environment" && card_elements[0] == "ibm_app_config_environment" {
			logger.Info("Entry:getAppConfigCost")
			if card_item.Usage_based {
				err = errors.New("it is a usage based resource")
				return 0, err
			} else {
				logger.Info("Exit:getAppConfigCost")
				return card_item.Estimated_rate, nil
			}

			//app config feature
		} else if instance_type == "ibm_app_config_feature" && card_elements[0] == "ibm_app_config_feature" {
			logger.Info("Entry:getAppConfigCost")
			if card_item.Usage_based {
				err = errors.New("it is a usage based resource")
				return 0, err
			} else {
				logger.Info("Exit:getAppConfigCost")
				return card_item.Estimated_rate, nil
			}

			//resource instance (event notifications and cloud object storage (cos instance))
			//2 resources are handled, plan is checked and matched in order to get the price
		} else if instance_type == "ibm_resource_instance" && card_elements[0] == "ibm_resource_instance" {

			if planData.Service == "cloud-object-storage" && card_elements[1] == "cos_instance" {
				logger.Info("Entry:getCosInstaneCost")
				if planData.Plan == card_elements[2] {

					if card_item.Usage_based {
						err = errors.New("it is a usage based resource")
						return 0, err

					} else {
						logger.Info("Exit:getCosInstaneCost")
						return card_item.Estimated_rate, nil
					}
				}

			} else if planData.Service == "event-notifications" && card_elements[1] == "event-notifications" {
				logger.Info("Entry:getEventNotificationCost")
				if planData.Plan == card_elements[2] {

					if card_item.Usage_based {
						err = errors.New("it is a usage based resource")
						return 0, err

					} else {
						logger.Info("Exit:getEventNotificationCost")
						return card_item.Estimated_rate, nil
					}
				} else {
					return 0, errors.New("invalid configuration for " + instance_type)
				}
			}

			//secondary subnets: capacity is checked in order to get the price
		} else if instance_type == "ibm_subnet" && card_elements[0] == "ibm_subnet" {

			if strings.ToLower(planData.Type) == card_elements[1] {
				capacity, _ := strconv.Atoi(card_elements[2])
				if planData.Capacity == capacity {
					logger.Info("Entry:getSecondarySubnetCost")
					logger.Info("Exit:getSecondarySubnetCost")
					return card_item.Estimated_rate, nil
				}
			}

			//cis: plan is checked in order to fetch the price
		} else if instance_type == "ibm_cis" && card_elements[0] == "ibm_cis" {

			if planData.Plan == card_elements[1] {
				logger.Info("Entry:getCISCost")
				if card_item.Usage_based {
					err = errors.New("it is a usage based resource")
					logger.Info("Exit:getCISCost")
					return 0, err

				} else {
					logger.Info("Exit:getCISCost")
					return card_item.Estimated_rate, nil
				}

			}

			// this is for dedicated host for VPC, profile is fetched and matched to get the price
		} else if instance_type == "ibm_is_dedicated_host" && card_elements[0] == "ibm_is_dedicated_host" {
			logger.Info("Entry:getDedicatedHostCost")
			profile = planData.Profile

			rate_card_profile = strings.Split(card_item.Plan, ".")[2]

			if profile == rate_card_profile {
				logger.Info("Exit:getDedicatedHostCost")
				return card_item.Estimated_rate, nil
			}

			//this is for every other resource
		} else {
			if instance_type == card_elements[0] {
				if card_item.Usage_based {
					err = errors.New("it is a usage based resource")
					return 0, err

				} else {
					logger.Info("Exit:getRateCardCost")
					return card_item.Estimated_rate, nil
				}
			}
		}
	}
	return 0, errors.New("no cost found")
}
