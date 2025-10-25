package status

import "time"

type URLStatus struct {
	URL      string
	Status   string
	Time     time.Time
	Duration time.Duration
	Error    string
}
