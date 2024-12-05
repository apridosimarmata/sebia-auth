package infrastructure

import (
	"context"
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"fmt"
	"mini-wallet/domain/inquiry"
	"strings"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
)

type Payment interface {
	CreatePaymentLink(ctx context.Context, inquiry inquiry.InquiryEntity) (url string, err error)
	VerifyCallback(ctx context.Context, signatureKey string, orderId string, total string, status string) (err error)
}

type paymentImplementation struct {
	snapClient *snap.Client
	serverKey  string
}

func NewPayment(snapClient *snap.Client) Payment {
	return &paymentImplementation{
		snapClient: snapClient,
		serverKey:  "Mid-server-9ZpQXdhK925cSTsjNV7bbcJX",
	}
}

func (payment *paymentImplementation) CreatePaymentLink(ctx context.Context, inquiry inquiry.InquiryEntity) (url string, err error) {
	names := strings.Split(inquiry.FullName, " ")
	fName := names[0]
	lName := names[1]
	// 2. Initiate Snap request param
	req := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  inquiry.ID,
			GrossAmt: 1,
			// GrossAmt: int64(inquiry.TotalPayment),
		},
		CreditCard: &snap.CreditCardDetails{
			Secure: true,
		},
		CustomerDetail: &midtrans.CustomerDetails{
			FName: fName,
			LName: lName,
			Email: inquiry.Email,
			Phone: inquiry.PhoneNumber,
		},
	}

	// 3. Execute request create Snap transaction to Midtrans Snap API
	snapResp, _ := payment.snapClient.CreateTransaction(req)

	return snapResp.RedirectURL, nil
}

func (payment *paymentImplementation) VerifyCallback(ctx context.Context, signatureKey string, orderId string, total string, status string) (err error) {

	comp := sha512.New()
	comp.Write([]byte(fmt.Sprintf("%s%s%s%s", orderId, status, total, payment.serverKey)))
	hash := hex.EncodeToString(comp.Sum(nil))
	if hash != signatureKey {
		return errors.New("invalid signature")
	}

	return nil
}
