package providers

import (
	"time"

	"github.com/vikesh-raj/go-practice/splitwise/models"
)

type FilterParams struct {
	After      time.Time
	Before     time.Time
	CreditOnly bool
	DebitOnly  bool
}

type DBProvider interface {
	Search(user string, params FilterParams) ([]models.LedgerEntry, error)
	AddTrasaction(transaction models.Transaction) error
	FindSettlementAmount(user, otherUser string) (float64, error)
}

func FillTransaction(transaction *models.Transaction) {
	if transaction.ID == "" {
		transaction.ID = NewUUID()
	}
	if transaction.Time.IsZero() {
		transaction.Time = time.Now()
	}
}
