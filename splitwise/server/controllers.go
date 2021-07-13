package server

import (
	"github.com/vikesh-raj/go-practice/splitwise/models"
	"github.com/vikesh-raj/go-practice/splitwise/providers"
)

func (a *Application) getLedger(user string) ([]models.LedgerEntry, error) {
	return a.provider.Search(user, providers.FilterParams{})
}

func (a *Application) addTrasaction(transaction models.Transaction) error {
	return a.provider.AddTrasaction(transaction)
}

func (a *Application) findSettlementAmount(from, to string) (float64, error) {
	return a.provider.FindSettlementAmount(from, to)
}
