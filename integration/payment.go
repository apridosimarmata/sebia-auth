package integration

import (
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
)

func NewSnapClient() *snap.Client {
	var s = snap.Client{}
	s.New("Mid-server-9ZpQXdhK925cSTsjNV7bbcJX", midtrans.Production)

	return &s
}
