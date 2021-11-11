package models

//Resource contains properties of plan for a terraform resource
type Resource struct {
	Address       string `json:"address"`
	Mode          string `json:"mode"`
	Type          string `json:"type"`
	Name          string `json:"name"`
	ProviderName  string `json:"provider_name"`
	SchemaVersion int    `json:"schema_version"`
	Values        struct {
		Generation int           `json:"generation"`
		Name       string        `json:"name"`
		Subnets    []string      `json:"subnets"`
		Timeouts   interface{}   `json:"timeouts"`
		Type       string        `json:"type"`
		Profile    string        `json:"profile"`
		Image      string        `json:"image"`
		Iops       int           `json:"iops"`
		Capacity   int           `json:"capacity"`
		BootVolume []interface{} `json:"boot_volume"`
	} `json:"values"`
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
	ResourceChanges []struct {
		Address      string `json:"address"`
		Mode         string `json:"mode"`
		Type         string `json:"type"`
		Name         string `json:"name"`
		ProviderName string `json:"provider_name"`
		Change       struct {
			Actions []string    `json:"actions"`
			Before  interface{} `json:"before"`
			After   struct {
				Name     string      `json:"name"`
				Subnets  []string    `json:"subnets"`
				Timeouts interface{} `json:"timeouts"`
				Type     string      `json:"type"`
			} `json:"after"`
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
	} `json:"resource_changes"`
	Configuration struct {
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
