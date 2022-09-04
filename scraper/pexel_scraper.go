package scraper

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type pexelsResult struct {
	Entries []pexelsEntry `json:"photos"`
}

type pexelsEntry struct {
	Urls pexelsImgUrl `json:"src"`
}

type pexelsImgUrl struct {
	Thumb string `json:"tiny"`
	Full  string `json:"original"`
}

type PexelsImgSearch struct {
	hc     http.Client
	apikey string
}

func NewPexelsImgSearch(timeout time.Duration, apikey string) *PexelsImgSearch {
	c := http.Client{
		Timeout: timeout,
	}
	s := PexelsImgSearch{
		hc:     c,
		apikey: apikey,
	}
	return &s
}

func (s *PexelsImgSearch) Search(searchTerm string) ([]Image, error) {
	url.QueryEscape(searchTerm)
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "https://api.pexels.com/v1/search?query="+searchTerm, nil)
	if err != nil {
		return nil, fmt.Errorf("Creating request failed: %w\n", err)
	}
	req.Header.Add("Authorization", s.apikey)
	res, err := s.hc.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Request failed: %w\n", err)
	}
	var r pexelsResult
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("Reading response body failed: %w\n", err)
	}
	defer res.Body.Close()
	err = json.Unmarshal(body, &r)
	if err != nil {
		return nil, fmt.Errorf("Unmarschalling json  failed: %w\n", err)
	}
	return s.convertToImage(&r), nil
}

func (s *PexelsImgSearch) convertToImage(r *pexelsResult) []Image {
	imgs := make([]Image, 0)
	if len(r.Entries) == 0 {
		return imgs
	}
	for _, v := range r.Entries {
		imgs = append(imgs, Image{
			Source: Pexels,
			URL:    v.Urls.Full,
			Thumb:  v.Urls.Thumb,
		})
	}
	return imgs
}
