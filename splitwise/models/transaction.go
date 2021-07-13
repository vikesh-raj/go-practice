package models

import "time"

type Transaction struct {
	ID      string
	From    string
	To      string
	Amount  float64
	Remarks string
	Time    time.Time
}
