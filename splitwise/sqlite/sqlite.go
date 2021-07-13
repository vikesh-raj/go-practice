package sqlite

import (
	"database/sql"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/vikesh-raj/go-practice/splitwise/models"
	"github.com/vikesh-raj/go-practice/splitwise/providers"
)

type sqliteProvider struct {
	sb sq.StatementBuilderType
	db *sql.DB
}

var createLedger string = `
CREATE TABLE IF NOT EXISTS ledger (
	id TEXT PRIMARY KEY,
	user TEXT,
	to TEXT,
	amount NUMERIC,
	remarks TEXT,
	time NUMERIC,
);`

// NewSqliteProvider creates a new sql user manager for the given path.
func NewSqliteProvider(dbPath string) (providers.DBProvider, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("unable to open database file : %s : %v", dbPath, err)
	}

	placeholder := sq.StatementBuilder.PlaceholderFormat(sq.Question)
	s := sqliteProvider{db: db, sb: placeholder}
	err = s.init()
	if err != nil {
		return nil, fmt.Errorf("unable to initialize database : %v", err)
	}
	return &s, nil
}

func (s *sqliteProvider) init() error {
	return nil
}

func (s *sqliteProvider) Search(user string, params providers.FilterParams) ([]models.LedgerEntry, error) {
	return nil, nil
}

func (s *sqliteProvider) AddTrasaction(transaction models.Transaction) error {
	return nil
}

func (s *sqliteProvider) FindSettlementAmount(user, otherUser string) (float64, error) {
	return 0.0, nil
}
