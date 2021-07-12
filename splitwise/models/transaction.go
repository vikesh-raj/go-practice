package models

import "time"

type Transaction struct {
	ID      string
	From    string
	To      string
	Amount  string
	Remarks string
	Time    time.Time
}
