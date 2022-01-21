module github.com/IBM-Cloud/terraform-cost-estimator/tfcost

go 1.15

require (
	github.com/IBM-Cloud/terraform-cost-estimator v0.0.0-20220121043039-1804168cb96b
	github.com/dustin/go-humanize v1.0.0 // indirect
	github.com/fatih/color v1.13.0
	github.com/kataras/tablewriter v0.0.0-20180708051242-e063d29b7c23
	github.com/landoop/tableprinter v0.0.0-20201125135848-89e81fc956e7
	github.com/mattn/go-runewidth v0.0.13 // indirect
	github.com/urfave/cli v1.22.5
	go.uber.org/zap v1.19.1
)

// replace github.com/IBM-Cloud/terraform-cost-estimator => ./terraform-cost-estimator
