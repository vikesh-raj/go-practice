package models

import "time"

type LedgerEntry struct {
	ID      string
	User    string
	To      string
	Amount  float64
	Remarks string
	Time    time.Time
}
