package pricing

import (
	"strings"
)

//PlanRepository -
type PlanRepository interface {
	GetInstancePlan(profile string) (string, error)
	GetVolumePlan(profile string) (string, error)
	GetLBPlan() (string, error)
	GetIPPlan() (string, error)
	GetImagePlan() (string, error)
	GetVPCPlan() (string, error)
	GetVPNPlan() (string, error)

	GetPricing(intanceID string) (PricingPlan, error)
}

// The kmsSettingsService - KMS Service
type planService struct {
	generation int
	resource   string
}

// NewPlanService -
func NewPlanService(generation int, resource string) PlanRepository {
	return &planService{
		generation: generation,
		resource:   resource,
	}
}

// GetInstancePlan ...
func (plan *planService) GetInstancePlan(profile string) (string, error) {
	var planID string
	var planName string

	//check the generation and set the planName for the instance
	if plan.generation == 1 {
		planName = "standard-vsi"
		//power resource if string contains letter 'p'
	} else if strings.Contains(strings.Split(profile, "-")[0], "p") {
		planName = "gen2-instance-power"
	} else {
		planName = "advanced-vsi"
	}

	instancePlan, err := plan.GetPricing("is.instance")
	if err != nil {
		return "", err
	}

	resouces := instancePlan.Resources

	// set the plan id for instance
	for _, value := range resouces {
		if value.Name == planName {
			planID = value.ID
		}
	}

	return planID, nil

}

// GetPricing ...
func (plan *planService) GetVolumePlan(profile string) (string, error) {
	var planID string
	var planName string

	switch profile {
	case "10iops-tier":
		if plan.generation == 1 {
			planName = "10iops-tier"

		} else {
			planName = "gen2-volume-10iops-tier"
		}
	case "5iops-tier":
		if plan.generation == 1 {
			planName = "gen1-volume-5iops-tier"

		} else {
			planName = "gen2-volume-5iops-tier"

		}
	case "3iops-tier":
		if plan.generation == 1 {
			planName = "general-purpose"

		} else {
			planName = "gen2-volume-general-purpose"
		}
	case "custom":
		if plan.generation == 1 {
			planName = "custom-"

		} else {
			planName = "gen2-volume-custom"
		}
	}

	volumePlan, err := plan.GetPricing("is.volume")
	if err != nil {
		return "", err
	}
	resouces := volumePlan.Resources

	// set the plan id for volume
	for _, value := range resouces {
		if value.Name == planName {
			planID = value.ID
		}
	}

	return planID, nil
}

func (plan *planService) GetLBPlan() (string, error) {
	var planID string
	var planName string

	if plan.generation == 1 {
		planName = "gen1-load-balancer"
	} else {
		planName = "gen2-load-balancer"
	}

	lbPlan, err := plan.GetPricing("is.load-balancer")
	if err != nil {
		return "", err
	}
	resouces := lbPlan.Resources

	// set the plan id for volume
	for _, value := range resouces {
		if value.Name == planName {
			planID = value.ID
		}
	}

	return planID, nil
}
func (plan *planService) GetIPPlan() (string, error) {
	var planID string
	var planName string

	if plan.generation == 1 {
		planName = "gen1-floating-ip"
	} else {
		planName = "default-nextgen"
	}
	ipPlan, err := plan.GetPricing("is.floating-ip")
	if err != nil {
		return "", err
	}
	resouces := ipPlan.Resources

	// set the plan id for volume
	for _, value := range resouces {
		if value.Name == planName {
			planID = value.ID
		}
	}

	return planID, nil
}

func (plan *planService) GetImagePlan() (string, error) {
	var planID string
	var planName string

	if plan.generation == 1 {
		planName = "gen1-image"
	} else {
		planName = "gen2-image"
	}
	imagePlan, err := plan.GetPricing("is.image")
	if err != nil {
		return "", err
	}
	resouces := imagePlan.Resources

	// set the plan id for volume
	for _, value := range resouces {
		if value.Name == planName {
			planID = value.ID
		}
	}

	return planID, nil
}

func (plan *planService) GetVPCPlan() (string, error) {
	var planID string
	var planName string

	if plan.generation == 1 {
		planName = "-vpc-egress-data-transfer"
	} else {
		planName = "nextgen-egress"
	}
	vpcPlan, err := plan.GetPricing("is.vpc")
	if err != nil {
		return "", err
	}
	resouces := vpcPlan.Resources

	// set the plan id for volume
	for _, value := range resouces {
		if value.Name == planName {
			planID = value.ID
		}
	}

	return planID, nil
}

func (plan *planService) GetVPNPlan() (string, error) {
	var planID string
	var planName string

	if plan.generation == 1 {
		planName = "gen1-vpn"
	} else {
		planName = "gen2-vpn"
	}
	vpnPlan, err := plan.GetPricing("is.vpn")
	if err != nil {
		return "", err
	}
	resouces := vpnPlan.Resources

	// set the plan id for volume
	for _, value := range resouces {
		if value.Name == planName {
			planID = value.ID
		}
	}

	return planID, nil
}
