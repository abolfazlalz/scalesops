package database

import "gorm.io/gorm"

type ImageRequest struct {
	gorm.Model
	Query  string
	Images []Image
}

type Image struct {
	gorm.Model
	URL            string
	Path           string
	ImageRequestID uint
	Width          int
	Height         int
	Size           int64
}
