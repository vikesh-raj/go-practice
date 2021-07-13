package server

import (
	"time"

	"github.com/vikesh-raj/go-practice/splitwise/models"
)

func (a *Application) getLedger(user string) ([]models.LedgerEntry, error) {
	return []models.LedgerEntry{
		{ID: "11", User: "sri", To: "malli", Amount: 200.00, Remarks: "Coffee", Time: time.Now()},
	}, nil
	// return nil, nil
}

func (a *Application) addTrasaction(transaction models.Transaction) error {
	return nil
}

func (a *Application) findSettlementAmount(from, to string) (float64, error) {
	return 0.0, nil
}
