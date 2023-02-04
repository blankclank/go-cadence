package customer

import "time"

type Customer struct {
	Name         string    `json:"name"`
	LastVisit    time.Time `json:"lastVisit"`
	TimesVisited int       `json:"timesVisited"`
}
