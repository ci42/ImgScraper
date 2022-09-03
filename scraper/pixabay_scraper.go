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

type pixabayResult struct {
	Entries []pixabayEntry `json:"hits"`
}

type pixabayEntry struct {
	Thumb string `json:"previewURL"`
	Full  string `json:"largeImageURL"`
}

type PixabayImgSearch struct {
	hc     http.Client
	apikey string
}

func NewPixabayImgSearch(timeoutInSeconds int64, apikey string) *PixabayImgSearch {
	c := http.Client{
		Timeout: time.Second * time.Duration(timeoutInSeconds),
	}
	s := PixabayImgSearch{
		hc:     c,
		apikey: apikey,
	}
	return &s
}

func (s *PixabayImgSearch) Search(searchTerm string) ([]Image, error) {
	url.QueryEscape(searchTerm)
	u := "https://pixabay.com/api/?key=" + s.apikey + "&q=" + searchTerm
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, u, nil)
	if err != nil {
		return nil, fmt.Errorf("Creating request failed: %w\n", err)
	}
	res, err := s.hc.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Request failed: %w\n", err)
	}
	var r pixabayResult
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

func (s *PixabayImgSearch) convertToImage(r *pixabayResult) []Image {
	imgs := make([]Image, 0, len(r.Entries))
	if len(r.Entries) == 0 {
		return imgs
	}
	for _, v := range r.Entries {
		imgs = append(imgs, Image{
			Source: Pixabay,
			URL:    v.Full,
			Thumb:  v.Thumb,
		})
	}
	return imgs
}
