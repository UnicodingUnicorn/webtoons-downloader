package main

import (
	"flag"
	"log"
	"time"

	"github.com/cheggaaa/pb/v3"
)

type Downloader struct {
	Client *CachedHttpClient
	OutputDir string
	Url string
	Delay time.Duration
}

func NewDownloader(url string, cacheDir string, outputDir string, delay int64) (*Downloader, error) {
	client, err := NewCachedHttpClient(cacheDir)
	if err != nil {
		return nil, err
	}

	err = DirExists(outputDir, true)
	if err != nil {
		return nil, err
	}

	d := &Downloader {
		Client: client,
		Url: url,
		OutputDir: outputDir,
		Delay: time.Duration(delay) * time.Millisecond,
	}

	return d, nil
}

func main() {
	// url := "https://www.webtoons.com/en/challenge/furi2play/list?title_no=172967"
	var url string
	flag.StringVar(&url, "u", "", "Specify chapter list page to download from")
	var cacheDir string
	flag.StringVar(&cacheDir, "c", "", "Specify directory to cache pages. Blank means no caching. You probably don't want to cache.")
	var outputDir string
	flag.StringVar(&outputDir, "o", "./output", "output directory")
	var delay int64
	flag.Int64Var(&delay, "d", 500, "delay between requests in ms")
	flag.Parse()

	if url == "" {
		log.Fatal("Please specify a url!")
	}

	d, err := NewDownloader(url, cacheDir, outputDir, delay)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("getting chapters...")
	chapterUrls, err := d.GetChapterUrls()
	if err != nil {
		log.Fatal(err)
	}
	if len(chapterUrls) == 0 {
		log.Fatal("number of retrieved chapters is 0, are you sure you input the correct url?")
	}

	log.Println("getting image urls...")
	total := len(chapterUrls)
	bar := pb.StartNew(total)
	titles, imageUrls, err := d.GetImageUrls(chapterUrls, func(_ int) {
		bar.Increment()
	})
	bar.Finish()
	if err != nil {
		log.Fatal(err)
	}
	for i, iu := range imageUrls {
		if len(iu) == 0 {
			log.Printf("number of images retrievable for chapter %s is 0", titles[i])
		}
	}

	origins := make([]string, 0)
	for _, cu := range chapterUrls {
		origins = append(origins, cu[1])
	}

	log.Println("getting images...")
	total = 0
	for _, iu := range imageUrls {
		total += len(iu)
	}
	bar = pb.StartNew(total)
	err = d.GetChaptersImages(origins, titles, imageUrls, func(_ int) {
		bar.Increment()
	})
	bar.Finish()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("done!")
}
