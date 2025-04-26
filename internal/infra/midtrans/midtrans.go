package midtrans

import (
	"batikin-be/config"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
)

func NewMidtransSnapClient() snap.Client {
	snapClient := snap.Client{}
	conf := config.Load()
	snapClient.New(conf.MidtransServerKey, midtrans.Sandbox)

	return snapClient
}
