package goddgimagesapi

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"regexp"
)

const URL = "https://duckduckgo.com/"

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

type Query struct {
	Keywords string
	Moderate bool
}

var re = regexp.MustCompile(`vqd=([\d-]+)\&`)

func token(keywords string) (string, error) {
	r, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		return "", err
	}

	addParams(r, map[string]string{
		"q": keywords,
	})

	res, err := http.DefaultClient.Do(r)
	if err != nil {
		return "", err
	}

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

func mod(b bool) string {
	if b {
		return "1"
	}

	return "-1"
}

func Do(query Query) (*Result, error) {
	if len(query.Keywords) == 0 {
		return nil, errors.New("not enough keywords")
	}

	tok, err := token(query.Keywords)
	if err != nil {
		return nil, err
	}

	r, err := http.NewRequest("GET", URL+"i.js", nil)
	if err != nil {
		return nil, err
	}
	addParams(r, map[string]string{
		"vqd":   tok,
		"l":     "en-us",
		"o":     "json",
		"q":     query.Keywords,
		"f":     ",,,",
		"p":     mod(query.Moderate),
		"v7exp": "a",
	})
	headers(r)

	res, err := http.DefaultClient.Do(r)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	result := Result{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
