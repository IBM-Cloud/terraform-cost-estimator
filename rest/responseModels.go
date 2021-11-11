package rest

import "time"

type GlobalCatalogResponse struct {
	Origin *string `json:"origin"`
	Type   string  `json:"type"`
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

type GlobalCatalogPlanResponse struct {
	Offset        int         `json:"offset" description:"The number of resources to skip over"`
	Limit         int         `json:"limit" description:"The number of resources to return."`
	Count         int         `json:"count" description:"The total number of resources available"`
	ResourceCount int         `json:"resource_count" description:"The total number of resources available"`
	First         string      `json:"first"`
	Resources     []Resources `json:"resources"`
}

type Resources struct {
	Active      bool        `json:"active"`
	CatalogCRN  string      `json:"catalog_crn"`
	ChildrenURL string      `json:"children_url"`
	Created     string      `json:"created"`
	Disabled    bool        `json:"disabled"`
	GeoTags     []string    `json:"geo_tags"`
	ID          string      `json:"id"`
	Images      images      `json:"images"`
	Kind        string      `json:"kind"`
	MetaData    interface{} `json:"metadata"`
	Name        string      `json:"name"`
	OverViewUI  interface{} `json:"overview_ui"`
	ParentID    string      `json:"parent_id"`
	ParentURL   string      `json:"parent_url"`
	PricingTags []string    `json:"pricing_tags"`
	Provider    interface{} `json:"provider"`
	Tags        []string    `json:"tags"`
	Updated     string      `json:"updated"`
	URL         string      `json:"url"`
	Visibility  interface{} `json:"visibility"`
}

type images struct {
	FeatureImage string `json:"feature_image"`
	Image        string `json:"image"`
	MediumImage  string `json:"medium_image"`
	SmallImage   string `json:"small_image"`
}
