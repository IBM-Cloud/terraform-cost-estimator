package costcalculator

import "time"

var resourceMap = map[string]string{
	"ibm_is_instance":           "Virtual Server",
	"ibm_is_lb":                 "Load Balancer",
	"ibm_is_floating_ip":        "Floating IP",
	"ibm_is_vpn_gateway":        "VPN Gateway",
	"ibm_is_volume":             "Storage Volume",
	"ibm_is_image":              "Custom image",
	"ibm_is_vpc":                "VPC",
	"ibm_container_vpc_cluster": "IKS on VPC",
	"ibm_is_subnet":             "Subnet",
	"ibm_container_cluster":     "IKS",
}

var planIDMap = map[string]string{
	"instance_1": "66380d42-d4a9-4627-88fa-7b6631e5bd63",
	"instance_2": "a736a57f-0584-474f-8411-55dc7d9dc811",
	"instance_p": "69a9646a-8f46-42bf-9834-1a92558bb618",
}

type Resource struct {
	Address       string `json:"address"`
	Mode          string `json:"mode"`
	Type          string `json:"type"`
	Name          string `json:"name"`
	ProviderName  string `json:"provider_name"`
	SchemaVersion int    `json:"schema_version"`
	Values        struct {
		Container
		Generation int    `json:"generation"`
		Name       string `json:"name"`
		// Subnets    []string      `json:"subnets"`
		Timeouts   interface{}   `json:"timeouts"`
		Type       string        `json:"type"`
		Profile    string        `json:"profile"`
		Image      string        `json:"image"`
		Iops       int           `json:"iops"`
		Capacity   int           `json:"capacity"`
		BootVolume []interface{} `json:"boot_volume"`
	} `json:"values"`
}

type Container struct {
	Datacenter       string        `json:"datacenter"`
	DefaultPoolSize  int           `json:"default_pool_size"`
	Hardware         string        `json:"hardware"`
	MachineType      string        `json:"machine_type"`
	Flavour          string        `json:"flavor"`
	Name             string        `json:"name"`
	UpdateAllWorkers bool          `json:"update_all_workers"`
	WorkerNum        int           `json:"worker_num"`
	SizePerZone      int           `json:"size_per_zone"`
	WorkerCount      int           `json:"worker_count"`
	Zone             []interface{} `json:"zones"`
	Plan             string        `json:"plan"`
	Service          string        `json:"service"`
}

//Planstruct consists of the entire json plandata
type Planstruct struct {
	FormatVersion    string `json:"format_version"`
	TerraformVersion string `json:"terraform_version"`
	PlannedValues    struct {
		RootModule struct {
			Resources []Resource `json:"resources"`
		} `json:"root_module"`
	} `json:"planned_values"`
	ResourceChanges []ResourceChanges `json:"resource_changes"`
	Configuration   struct {
		ProviderConfig struct {
			Ibm struct {
				Name        string `json:"name"`
				Expressions struct {
					Generation struct {
						ConstantValue int `json:"constant_value"`
					} `json:"generation"`
					IbmcloudAPIKey struct {
						ConstantValue string `json:"constant_value"`
					} `json:"ibmcloud_api_key"`
				} `json:"expressions"`
			} `json:"ibm"`
		} `json:"provider_config"`
		RootModule struct {
			Resources []struct {
				Address           string `json:"address"`
				Mode              string `json:"mode"`
				Type              string `json:"type"`
				Name              string `json:"name"`
				ProviderConfigKey string `json:"provider_config_key"`
				Expressions       struct {
					Name struct {
						ConstantValue string `json:"constant_value"`
					} `json:"name"`
					Subnets struct {
						ConstantValue []string `json:"constant_value"`
					} `json:"subnets"`
				} `json:"expressions"`
				SchemaVersion int `json:"schema_version"`
			} `json:"resources"`
		} `json:"root_module"`
	} `json:"configuration"`
}

type ResourceChanges struct {
	Address      string `json:"address"`
	Mode         string `json:"mode"`
	Type         string `json:"type"`
	Name         string `json:"name"`
	ProviderName string `json:"provider_name"`
	Change       struct {
		Actions      []string     `json:"actions"`
		Before       ResourceConf `json:"before"`
		After        ResourceConf `json:"after"`
		AfterUnknown struct {
			Hostname              bool `json:"hostname"`
			ID                    bool `json:"id"`
			OperatingStatus       bool `json:"operating_status"`
			PrivateIps            bool `json:"private_ips"`
			PublicIps             bool `json:"public_ips"`
			ResourceControllerURL bool `json:"resource_controller_url"`
			ResourceGroup         bool `json:"resource_group"`
			ResourceGroupName     bool `json:"resource_group_name"`
			ResourceName          bool `json:"resource_name"`
			Status                bool `json:"status"`
		} `json:"after_unknown"`
	} `json:"change"`
}
type ResourceConf struct {
	Container
	ISInstance
	Type                         string `json:"type"`
	Capacity                     int    `json:"capacity"`
	Iops                         int    `json:"iops"`
	FlavorKeyName                string `json:"flavor_key_name"`
	Memory                       int    `json:"memory"`
	Cores                        int    `json:"cores"`
	OperatingSystemReferenceCode string `json:"os_reference_code"`
	DedicatedHostFlag            bool   `json:"dedicated_acct_host_only"`
	NetworkSpeed                 int    `json:"network_speed"`
	PrivateNetworkOnly           bool   `json:"private_network_only"`
	MembersMemoryAllocationMB    int    `json:"members_memory_allocation_mb"`
	MembersDiskAllocationMB      int    `json:"members_disk_allocation_mb"`
	MembersCPUAllocationCount    int    `json:"members_cpu_allocation_count"`
	NodeCount                    int    `json:"node_count"`
	NodeMemoryAllocationMB       int    `json:"node_memory_allocation_mb"`
	NodeDiskAllocationMB         int    `json:"node_disk_allocation_mb"`
	NodeCPUAllocationCount       int    `json:"node_cpu_allocation_count"`
	ImageID                      string `json:"image"`
}

type ISInstance struct {
	Profile    string        `json:"profile"`
	Image      string        `json:"image"`
	BootVolume []interface{} `json:"boot_volume"`
}

//Respstruct contains restapi response
type Respstruct struct {
	Cost     float64     `json:"cost"`
	Measures interface{} `json:"measures"`
	PlanID   string      `json:"plan_id"`
}

//BillOfMaterial - output cost data for a resource
type BillOfMaterial struct {
	Quantity            int      `json:"quantity" header:"quantity"`
	TerraformItemID     string   `json:"terraformItemId" header:"terraformItemId"`
	ID                  string   `json:"id" header:"id"`
	RateCardCost        bool     `json:"rateCardCost" header:"rateCardCost"`
	Title               string   `json:"title" header:"title"`
	PlanID              string   `json:"planID" header:"planID"`
	ShortDescription    string   `json:"shortDescription" header:"shortDescription"`
	Features            []string `json:"features" header:"features"`
	EstimateType        string   `json:"estimateType" header:"estimateType"`
	LineItemTotal       float64  `json:"lineitemtotal" header:"lineitemtotal"`
	CurrLineItemTotal   float64  `json:"currlineitemtotal" header:"currlineitemtotal"`
	PrevLineItemTotal   float64  `json:"prevlineitemtotal" header:"prevlineitemtotal"`
	ChangeLineItemTotal float64  `json:"changelineitemtotal" header:"changelineitemtotal"`

	ResourceBreakdown []DependencyResource `json:"depends" header:"depends"`
}

type DependencyResource struct {
	Quantity int     `json:"quantity" header:"quantity"`
	Cost     float64 `json:"cost" header:"cost"`
	Title    string  `json:"title" header:"title"`
}

type BOM struct {
	Country   string  `json:"country"`
	Currency  string  `json:"currency"`
	TotalCost float64 `json:"totalcost"`
	Lineitem  []BillOfMaterial
}

//Request payload struct
type Request struct {
	ServiceID string    `json:"service_id"`
	PlanID    string    `json:"plan_id"`
	Currency  string    `json:"currency"`
	Country   string    `json:"country"`
	Region    string    `json:"region"`
	Measures  []Metrics `json:"measures"`
}

type Metrics struct {
	Measure  string `json:"measure"`
	Quantity string `json:"quantity"`
}

type GlobalCatalog struct {
	Origin string `json:"origin"`
	I18N   struct {
	} `json:"i18n"`
	StartingPrice struct {
	} `json:"starting_price"`
	EffectiveFrom  time.Time `json:"effective_from"`
	EffectiveUntil time.Time `json:"effective_until"`
	Metrics        []struct {
		PartRef               string `json:"part_ref"`
		MetricID              string `json:"metric_id"`
		TierModel             string `json:"tier_model"`
		ResourceDisplayName   string `json:"resource_display_name"`
		ChargeUnitDisplayName string `json:"charge_unit_display_name"`
		ChargeUnitName        string `json:"charge_unit_name"`
		ChargeUnit            string `json:"charge_unit"`
		ChargeUnitQuantity    int    `json:"charge_unit_quantity"`
		Amounts               []struct {
			Country  string `json:"country"`
			Currency string `json:"currency"`
			Prices   []struct {
				QuantityTier int     `json:"quantity_tier"`
				Price        float64 `json:"price"`
			} `json:"prices"`
		} `json:"amounts"`
		UsageCapQty    int    `json:"usage_cap_qty"`
		DisplayCap     int    `json:"display_cap"`
		EffectiveFrom  string `json:"effective_from"`
		EffectiveUntil string `json:"effective_until"`
	} `json:"metrics"`
}

type ClassicRateCard struct {
	BlockDevices interface{} `json:"blockDevices"`
	DataCenters  interface{} `json:"datacenters"`
	Flavors      []struct {
		Flavor struct {
			KeyName                  string `json:"keyName"`
			TotalMinimumRecurringFee string `json:"totalMinimumRecurringFee"`
		} `json:"flavor"`
	}
	Memory []struct {
		ItemPrice struct {
			RecurringFee string `json:"recurringFee"`
		} `json:"itemPrice"`
		Template struct {
			MaxMemory int `json:"maxMemory"`
		} `json:"template"`
	}
	NetworkComponents []struct {
		ItemPrice struct {
			RecurringFee              string `json:"recurringFee"`
			DedicatedHostInstanceFlag bool   `json:"dedicatedHostInstanceFlag"`
		} `json:"itemPrice"`
		Template struct {
			NetworkComponent []struct {
				MaxSpeed int `json:"maxSpeed"`
			} `json:"networkComponents"`
			PrivateNetworkOnly bool `json:"privateNetworkOnlyFlag"`
		} `json:"template"`
	} `json:"networkComponents"`
	OperatingSystems []struct {
		ItemPrice struct {
			RecurringFee string `json:"recurringFee"`
		} `json:"itemPrice"`
		Template struct {
			OperatingSystem string `json:"operatingSystemreferenceCode"`
		} `json:"template"`
	} `json:"operatingSystems"`
	Processors []struct {
		ItemPrice struct {
			RecurringFee string `json:"recurringFee"`
		} `json:"itemPrice"`
		Template struct {
			StartCpus     int  `json:"startCpus"`
			DedicatedHost bool `json:"dedicatedAccountHostOnlyFlag"`
		} `json:"template"`
	}
}

type RateCard struct {
	Service_name   string  `json:"service_name"`
	Plan           string  `json:"plan_id"`
	Estimated_rate float64 `json:"estimated_rate"`
	Units          string  `json:"unit"`
	Currency       string  `json:"currency"`
	Unit_quantity  string  `json:"unit_quantity"`
	Usage_based    bool    `json:"usage_based"`
}
