package goddgimagesapi

import (
	"encoding/json"
	"errors"
	"io"
	"log"
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

type Client struct {
	httpClient *http.Client
	debug      bool
}

func NewClient(httpClient *http.Client) Client {
	return Client{httpClient: httpClient, debug: false}
}

func (c *Client) EnableDebug() {
	c.debug = true
}

var re = regexp.MustCompile(`vqd=([\d-]+)\&`)

func (c *Client) token(keywords string) (string, error) {
	r, _ := http.NewRequest("GET", URL, nil)
	addParams(r, map[string]string{
		"q": keywords,
	})

	res, err := c.httpClient.Do(r)
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
		if c.debug {
			log.Println(string(body))
		}

		return "", errors.New("token parsing failed")
	}

	return string(token)[4:], nil
}

func (c *Client) Do(query Query) (*Result, error) {
	if len(query.Keywords) == 0 {
		return nil, errors.New("not enough keywords")
	}

	tok, err := c.token(query.Keywords)
	if err != nil {
		return nil, err
	}

	r, _ := http.NewRequest("GET", URL+"i.js", nil)
	addParams(r, map[string]string{
		"vqd": tok,
		"l":   "en-us",
		"o":   "json",
		"q":   query.Keywords,
		"f":   ",,,",
		"p": func(b bool) string {
			if b {
				return "1"
			}
			return "-1"
		}(query.Moderate),

		"v7exp": "a",
	})
	headers(r)

	res, err := c.httpClient.Do(r)
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
