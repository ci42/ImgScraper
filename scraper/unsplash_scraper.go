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

type unsplashResult struct {
	Entries []unsplashEntry `json:"results"`
}

type unsplashEntry struct {
	Urls UnsplashImgUrl `json:"Urls"`
}

type UnsplashImgUrl struct {
	Thumb string `json:"Thumb"`
	Full  string `json:"Full"`
}

type UnsplashImgSearch struct {
	hc     http.Client
	apikey string
}

func NewUnsplashImgSearch(timeout time.Duration, apikey string) *UnsplashImgSearch {
	c := http.Client{
		Timeout: timeout,
	}
	s := UnsplashImgSearch{
		hc:     c,
		apikey: apikey,
	}
	return &s
}

func (s *UnsplashImgSearch) Search(searchTerm string) ([]Image, error) {
	url.QueryEscape(searchTerm)
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "https://api.unsplash.com/search/photos?query="+searchTerm, nil)
	if err != nil {
		return nil, fmt.Errorf("Creating request failed: %w\n", err)
	}
	req.Header.Add("Authorization", "Client-ID "+s.apikey)
	res, err := s.hc.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Request failed: %w\n", err)
	}
	var r unsplashResult
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

func (s *UnsplashImgSearch) convertToImage(r *unsplashResult) []Image {
	imgs := make([]Image, 0)
	if len(r.Entries) == 0 {
		return imgs
	}
	for _, v := range r.Entries {
		imgs = append(imgs, Image{
			Source: Unsplash,
			URL:    v.Urls.Full,
			Thumb:  v.Urls.Thumb,
		})
	}
	return imgs
}
