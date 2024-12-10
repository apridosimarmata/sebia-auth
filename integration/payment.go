package integration

import (
	"mini-wallet/utils"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
)

func NewSnapClient(config *utils.AppConfig) *snap.Client {
	var s = snap.Client{}

	envType := midtrans.Production
	if config.AppEnvironment == "development" {
		envType = midtrans.Sandbox
	}

	s.New(config.MidtransServerKey, envType)

	return &s
}
