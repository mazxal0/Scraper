package status

import "time"

func SetStatus(url, status string) {
	mu.Lock()
	defer mu.Unlock()
	if _, ok := Statuses[url]; !ok {
		Statuses[url] = &URLStatus{}
	}
	Statuses[url].Status = status
	Statuses[url].URL = url
}

func GetStatus(url string) string {
	mu.Lock()
	defer mu.Unlock()
	return Statuses[url].Status
}

func GetAll() map[string]string {
	mu.Lock()
	defer mu.Unlock()

	result := make(map[string]string)
	for url, s := range Statuses {
		result[url] = s.Status
	}

	return result
}

func SetTimeToAction(url string, t time.Time) {
	mu.Lock()
	defer mu.Unlock()

	if _, ok := Statuses[url]; !ok {
		Statuses[url] = &URLStatus{}
	}

	Statuses[url].Time = t
}

func GetTimeToAction(url string) time.Time {
	mu.Lock()
	defer mu.Unlock()

	if _, ok := Statuses[url]; !ok {
		return time.Now()
	}

	return Statuses[url].Time
}

func SetDuration(url string, d time.Duration) {
	mu.Lock()
	defer mu.Unlock()

	if _, ok := Statuses[url]; !ok {
		Statuses[url] = &URLStatus{}
	}

	Statuses[url].Duration = d
}

func GetDuration(url string) time.Duration {
	mu.Lock()
	defer mu.Unlock()

	return Statuses[url].Duration
}
