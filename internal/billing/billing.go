package billing

import (
	"io/ioutil"
	"log"
	"os"
	"statusPage/internal/entities"
	"strconv"
	"sync"
)

const (
	CreateCustomerMask int64 = 1 << iota
	PurchaseMask
	PayoutMask
	RecurringMask
	FraudControlMask
	CheckoutPageMask
)

func BillingDataReader(path string, wg *sync.WaitGroup) entities.BillingData {
	out := make(chan entities.BillingData)
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(out)
		file, err := os.Open(path)
		if err != nil {
			log.Fatal("Cannot open billing file:", err)
		}
		defer file.Close()
		reader, err := ioutil.ReadAll(file)
		if err != nil {
			log.Fatal("Cannot read billing file:", err)
		}

		mask, err := strconv.ParseInt(string(reader), 2, 0)

		result := entities.BillingData{
			CreateCustomer: mask&CreateCustomerMask != 0,
			Purchase:       mask&PurchaseMask != 0,
			Payout:         mask&PayoutMask != 0,
			Recurring:      mask&RecurringMask != 0,
			FraudControl:   mask&FraudControlMask != 0,
			CheckoutPage:   mask&CheckoutPageMask != 0,
		}
		out <- result
	}()
	var result = <-out
	return result
}
