package status

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
