package main

import (
	"fmt"
	"github.com/gocolly/colly"
	"io"
	"log"
	"net/http"
	"os"
	filepath2 "path/filepath"
	"strings"
)

func main() {
	c := colly.NewCollector()

	imageUrls := make([]string, 0)
	c.OnHTML("img[src]", func(e *colly.HTMLElement) {
		imageUrl := e.Attr("src")
		imageUrls = append(imageUrls, imageUrl)
	})
	c.OnError(func(r *colly.Response, err error) {
		log.Println("request url:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	//weburl is website url paste on that string
	weburl := ""
	err := c.Visit(weburl)
	if err != nil {
		log.Fatalf("Cannot visit %s", weburl)
	}

	for _, imageUrl := range imageUrls {
		err := downloadImage(imageUrl)
		if err != nil {
			log.Println("Error downloading image:", err)
		}
	}
	fmt.Println("Done!")
}

func downloadImage(imageUrl string) error {
	if !strings.HasPrefix(imageUrl, "http://") && !strings.HasPrefix(imageUrl, "https://") {
		fmt.Println("Skipping", imageUrl)
	}
	//download directory the files will be on there
	dirdownload := "./download"
	err := os.MkdirAll(dirdownload, os.ModePerm)
	if err != nil {
		fmt.Println("Failed Create Directory", err)
	}
	tokens := strings.Split(imageUrl, "/")
	filename := tokens[len(tokens)-1]

	filepath := filepath2.Join(dirdownload, filename)

	resp, err := http.Get(imageUrl)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return err
	}

	fmt.Println("Downloaded", imageUrl)
	return nil
}
