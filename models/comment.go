package models

import "time"

type Comment struct {
	ID      int
	Task    int
	Date    time.Time
	Comment string
}
