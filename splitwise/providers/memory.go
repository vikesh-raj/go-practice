package providers

import (
	"sync"

	"github.com/vikesh-raj/go-practice/splitwise/models"
)

type inMemoryDB struct {
	mutex  *sync.Mutex
	ledger []models.LedgerEntry
}

func NewInMemoryDB() DBProvider {
	return &inMemoryDB{
		mutex:  &sync.Mutex{},
		ledger: make([]models.LedgerEntry, 0),
	}
}

func (db *inMemoryDB) Search(user string, params FilterParams) ([]models.LedgerEntry, error) {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	output := make([]models.LedgerEntry, 0)

	for i := len(db.ledger) - 1; i >= 0; i-- {
		item := db.ledger[i]
		if item.User != user {
			continue
		}

		output = append(output, item)
	}
	return output, nil
}

func (db *inMemoryDB) AddTrasaction(transaction models.Transaction) error {
	FillTransaction(&transaction)

	db.mutex.Lock()
	defer db.mutex.Unlock()

	le := models.LedgerEntry{
		ID:      NewUUID(),
		User:    transaction.From,
		To:      transaction.To,
		Amount:  transaction.Amount,
		Remarks: transaction.Remarks,
		Time:    transaction.Time,
	}

	db.ledger = append(db.ledger, le)

	le = models.LedgerEntry{
		ID:      NewUUID(),
		User:    transaction.To,
		To:      transaction.From,
		Amount:  -transaction.Amount,
		Remarks: transaction.Remarks,
		Time:    transaction.Time,
	}

	db.ledger = append(db.ledger, le)
	return nil
}

func (db *inMemoryDB) FindSettlementAmount(user, otherUser string) (float64, error) {
	db.mutex.Lock()
	defer db.mutex.Unlock()
	amount := 0.0

	for i := len(db.ledger) - 1; i >= 0; i-- {
		item := db.ledger[i]
		if item.User != user {
			continue
		}
		if item.To != otherUser {
			continue
		}
		amount += item.Amount
	}

	return amount, nil
}
