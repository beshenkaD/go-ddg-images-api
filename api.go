package goddgimagesapi

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"regexp"
)

type Result struct {
	Ads          interface{} `json:"ads"`
	Next         string      `json:"next"`
	Query        string      `json:"query"`
	QueryEncoded string      `json:"queryEncoded"`
	ResponseType string      `json:"response_type"`
	Results      []struct {
		Height    int    `json:"height"`
		Image     string `json:"image"`
		Source    string `json:"source"`
		Thumbnail string `json:"thumbnail"`
		Title     string `json:"title"`
		URL       string `json:"url"`
		Width     int    `json:"width"`
	} `json:"results"`
}

const URL = "https://duckduckgo.com/"

var client = http.DefaultClient

func SetClient(c *http.Client) {
	client = c
}

func Do(keywords string, moderate bool) (*Result, error) {
	if len(keywords) == 0 {
		return nil, errors.New("not enough keywords")
	}

	tok, err := token(keywords)
	if err != nil {
		return nil, err
	}

	r, _ := http.NewRequest("GET", URL+"i.js", nil)
	addParams(r, map[string]string{
		"vqd": tok,
		"l":   "en-us",
		"o":   "json",
		"q":   keywords,
		"f":   ",,,",
		"p": func(b bool) string {
			if b {
				return "1"
			}
			return "-1"
		}(moderate),

		"v7exp": "a",
	})
	headers(r)

	res, err := client.Do(r)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	result := Result{}
	err = json.NewDecoder(res.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

var re = regexp.MustCompile(`vqd=([\d-]+)\&`)

func token(keywords string) (string, error) {
	r, _ := http.NewRequest("GET", URL, nil)
	addParams(r, map[string]string{
		"q": keywords,
	})

	res, err := client.Do(r)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	token := re.Find(body)

	if token == nil {
		return "", errors.New("token parsing failed")
	}

	return string(token)[4:], nil
}
