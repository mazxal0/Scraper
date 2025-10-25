package scraper

type Result struct {
	URL   string `json:"url"`
	Title string `json:"title"`
	Error string `json:"error"`
}
