package goddgimagesapi

import "net/http"

func addParams(r *http.Request, p map[string]string) {
	q := r.URL.Query()

	for k, v := range p {
		q.Add(k, v)
	}

	r.URL.RawQuery = q.Encode()
}

func headers(r *http.Request) {
	r.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:93.0) Gecko/20100101 Firefox/93.0")
	r.Header.Add("accept", "application/json, text/javascript, */* q=0.01")
	r.Header.Add("authority", "duckduckgo.com")
	r.Header.Add("sec-fetch-dest", "Empty")
	r.Header.Add("x-requested-with", "XMLHttpRequest")
	r.Header.Add("sec-fetch-site", "same-origin")
	r.Header.Add("sec-fetch-mode", "cors")
	r.Header.Add("referer", "https://duckduckgo.com/")
	r.Header.Add("accept-language", "en-US,enq=0.9")
}
