package main

import (
	"fmt"
	"github.com/abolfazlalz/scalesops/golang/database"
	"github.com/abolfazlalz/scalesops/golang/google_image"
	"log"
)

func main() {
	database.Database()
	fmt.Println("Database connected successfully !")
	var query string
	fmt.Print("Enter search query: ")
	_, err := fmt.Scanln(&query)

	if err != nil {
		log.Fatalf("error during give query from input: %v", err)
	}

	var maxImage int
	fmt.Print("Enter max number of images to download: ")
	_, err = fmt.Scanf("%d", &maxImage)
	if err != nil {
		log.Fatalf("error during input for max number of images: %v", err)
	}

	id := google_image.NewImageDownloader(query, maxImage)
	if err := id.Process(); err != nil {
		log.Fatalf("error during process to download images: %v", err)
	}
}
