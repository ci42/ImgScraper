package scraper

type ImgScraper interface {
	Search(query string) ([]Image, error)
}
