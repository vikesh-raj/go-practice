package sqlite

import (
	"database/sql"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/vikesh-raj/go-practice/splitwise/models"
	"github.com/vikesh-raj/go-practice/splitwise/providers"
)

type sqliteProvider struct {
	sb         sq.StatementBuilderType
	db         *sql.DB
	insertStmt *sql.Stmt
}

const dbVersion int = 1

var dropLedgerTable string = "DROP TABLE IF EXISTS ledger;"
var createLedgerTable string = `
CREATE TABLE IF NOT EXISTS ledger (
	id TEXT PRIMARY KEY,
	user TEXT,
	to TEXT,
	amount NUMERIC,
	remarks TEXT,
	time NUMERIC,
);`
var indexLedgerTable string = `
CREATE INDEX idx_ledger_user ON ledger (user);
CREATE INDEX idx_ledger_to ON ledger (user_id, session_id);
CREATE INDEX idx_ledger_time ON ledger (time DESC);`

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

func (s *sqliteProvider) init() (err error) {

	err = s.createVersionTable()
	if err != nil {
		return err
	}

	err = s.migrate()
	if err != nil {
		return err
	}

	s.insertStmt, err = s.db.Prepare(
		"INSERT INTO ledger(id, user, to, amount, remarks, time) VALUES(?, ?, ?, ?, ?, ?)")
	if err != nil {
		return fmt.Errorf("unable to prepare the insert stmt : %v", err)
	}
	return nil
}

func (s *sqliteProvider) Search(user string, params providers.FilterParams) ([]models.LedgerEntry, error) {
	stmt := s.sb.Select("id", "user", "to", "amount", "remarks", "time").From("ledger").
		Where(sq.Eq{"user": user}).
		OrderByClause("time DESC")
	rows, err := stmt.RunWith(s.db).Query()
	if err != nil {
		return nil, fmt.Errorf("unable to select from dialogue history : %v", err)
	}

	records := make([]models.LedgerEntry, 0)
	for rows.Next() {
		var le models.LedgerEntry
		var timestamp int64
		err := rows.Scan(&le.ID, &le.User, &le.To, &le.Amount, &le.Remarks, &timestamp)
		if err != nil {
			return nil, fmt.Errorf("row scan error : %v", err)
		}
		le.Time = time.Unix(timestamp, 0)
		records = append(records, le)
	}

	rerr := rows.Close()
	if rerr != nil {
		return nil, fmt.Errorf("error while closing rows : %v", rerr)
	}
	return records, nil
}

func (s *sqliteProvider) AddTrasaction(transaction models.Transaction) error {
	providers.FillTransaction(&transaction)

	_, err := s.insertStmt.Exec(
		providers.NewUUID(), transaction.From,
		transaction.To, transaction.Amount,
		transaction.Remarks,
		toUnix(transaction.Time))
	if err != nil {
		return fmt.Errorf("unable to insert into ledger table : %v", err)
	}

	_, err = s.insertStmt.Exec(
		providers.NewUUID(), transaction.To,
		transaction.From, -transaction.Amount,
		transaction.Remarks,
		toUnix(transaction.Time))
	if err != nil {
		return fmt.Errorf("unable to insert into ledger table : %v", err)
	}

	return nil
}

func (s *sqliteProvider) FindSettlementAmount(user, otherUser string) (float64, error) {
	stmt := s.sb.Select("SUM(amount)").From("ledger").
		Where(sq.Eq{"user": user}).
		Where(sq.Eq{"to": otherUser})

	var amount float64
	err := stmt.RunWith(s.db).QueryRow().Scan(&amount)
	return amount, err
}

func (s *sqliteProvider) migrate() error {

	currentVersion, err := s.getVersion()
	if err != nil {
		return fmt.Errorf("unable to get version : %v", err)
	}

	if currentVersion == dbVersion {
		return nil
	}

	switch currentVersion {
	case 0:
		err := s.createLedgerTable()
		if err != nil {
			return fmt.Errorf("unable to create history table : %v", err)
		}

		err = s.setVersion(1)
		if err != nil {
			return fmt.Errorf("unable to set version : %v", err)
		}
	}
	return nil
}

func (s *sqliteProvider) createLedgerTable() error {
	return s.execMulti([]string{dropLedgerTable, createLedgerTable, indexLedgerTable})
}

func (s *sqliteProvider) createVersionTable() error {
	statement := `CREATE TABLE IF NOT EXISTS version (key TEXT PRIMARY KEY, value TEXT);`
	err := s.exec(statement)
	if err != nil {
		return fmt.Errorf("unable to create version table : %v", err)
	}

	var version int
	err = s.db.QueryRow("SELECT value FROM version WHERE key='version'").Scan(&version)
	if err == sql.ErrNoRows {
		statement := `INSERT INTO version(key, value) VALUES("version", 0);`
		return s.exec(statement)
	}
	return err
}

func (s *sqliteProvider) getVersion() (version int, err error) {
	err = s.db.QueryRow("SELECT value FROM version WHERE key='version'").Scan(&version)
	if err == sql.ErrNoRows {
		return 0, nil
	}
	return
}

func (s *sqliteProvider) setVersion(version int) error {
	statement := fmt.Sprint(`UPDATE version SET value=`, version, ` WHERE key='version'`)
	return s.exec(statement)
}

func (s *sqliteProvider) exec(statement string) error {
	_, err := s.db.Exec(statement)
	return err
}

func (s *sqliteProvider) execMulti(statements []string) error {
	for _, statement := range statements {
		err := s.exec(statement)
		if err != nil {
			return err
		}
	}
	return nil
}

func toUnix(time time.Time) int64 {
	return time.UTC().Unix()
}
