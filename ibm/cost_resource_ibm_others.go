package costcalculator

import "go.uber.org/zap"

//Parse subnet details, configure the body, call restapi and return the estimated cost
func getSubnetCost(logger *zap.Logger, changeData ResourceConf, token string) (float64, error) {

	logger.Info("Entry: getSubnetCost")

	return 0, nil

}
