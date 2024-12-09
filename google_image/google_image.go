package google_image

import (
	"fmt"
	"github.com/abolfazlalz/scalesops/golang/database"
	"github.com/abolfazlalz/scalesops/golang/utils"
	"github.com/gocolly/colly/v2"
	"log"
	"strings"
	"sync"
)

type ImageDownloader struct {
	mu          sync.Mutex
	maxImage    int
	querySearch string
	requestID   uint
}

// NewImageDownloader create new instance from ImageDownloader
func NewImageDownloader(querySearch string, maxImage int) *ImageDownloader {
	db := database.Database()
	// save request in database
	if tx := db.Create(&database.ImageRequest{
		Query: querySearch,
	}); tx.Error != nil {
		log.Fatalf("error during create request model: %v", tx.Error)
	}
	return &ImageDownloader{
		// to prevent concurrency injection
		mu:          sync.Mutex{},
		maxImage:    maxImage,
		querySearch: querySearch,
	}
}

// handleImage on selector find image
func (id *ImageDownloader) handleImage(e *colly.HTMLElement) {
	id.mu.Lock()
	defer id.mu.Unlock()
	if id.maxImage == 0 {
		return
	}
	src := e.Attr("src")
	if strings.HasPrefix(src, "https") {
		// decrease to keep memorize how many images left
		id.maxImage -= 1
		file, err := utils.DownloadFile(src)

		if err != nil {
			log.Fatalf("error downloading image: %v", err)
		}
		db := database.Database()
		img, err := utils.ImageDetails(file)
		if err != nil {
			panic(err)
		}
		// save request in database
		db.Create(&database.Image{
			URL:            src,
			Path:           file,
			ImageRequestID: id.requestID,
			Width:          img.Width,
			Height:         img.Height,
			Size:           img.Size,
		})
		log.Println(file)
	}

}

// Process Start process on Google images and save them on system and keep history on database
func (id *ImageDownloader) Process() error {
	c := colly.NewCollector()
	c.OnHTML("img[src]", id.handleImage)
	c.OnError(func(r *colly.Response, err error) {
		log.Fatal(err)
	})
	searchUrl := fmt.Sprintf("https://www.google.com/search?hl=en&tbm=isch&q=%s", id.querySearch)
	if err := c.Visit(searchUrl); err != nil {
		panic(err)
	}
	return nil
}
