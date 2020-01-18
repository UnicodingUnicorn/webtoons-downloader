package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/html"
)

func (d *Downloader) GetImageUrls(chapterUrls [][]string, update func(int)) ([]string, [][]string, error) {
	titles := make([]string, 0)
	dataUrls := make([][]string, 0)
	steps := GetDataUrlSteps()

	titleFmtStr := fmt.Sprintf("%%0%dd. %%s", len(strconv.Itoa(len(chapterUrls))))

	for i, chapterUrl := range chapterUrls {
		page, err := d.Client.GetPage(chapterUrl[1], nil)
		if err != nil {
			return nil, nil, err
		}

		doc, err := html.Parse(strings.NewReader(page))
		if err != nil {
			return nil, nil, err
		}

		pageDataUrls := step(doc, 0, steps, GetImageUrl)
		titles = append(titles, fmt.Sprintf(titleFmtStr, i + 1, chapterUrl[0]))
		if len(pageDataUrls) > 0 {
			dataUrls = append(dataUrls, pageDataUrls[0])
		}

		if update != nil {
			update(i + 1)
		}
		time.Sleep(d.Delay)
	}

	return titles, dataUrls, nil
}

func GetImageUrl(n *html.Node) [][]string {
	dataUrls := make([][]string, 1)
	dataUrls[0] = make([]string, 0)
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode && c.Data == "img" {
			dataUrl := GetAttr(c, "data-url")
			if dataUrl != "" {
				dataUrls[0] = append(dataUrls[0], dataUrl)
			}
		}
	}
	return dataUrls
}

func GetDataUrlSteps() []*Step {
	steps := make([]*Step, 8)

	steps[0] = &Step {
		Element: "html",
		Id: "",
		Class: "",
	}

	steps[1] = &Step {
		Element: "body",
		Id: "",
		Class: "",
	}

	steps[2] = &Step {
		Element: "div",
		Id: "wrap",
		Class: "",
	}

	steps[3] = &Step {
		Element: "div",
		Id: "container",
		Class: "",
	}

	steps[4] = &Step {
		Element: "div",
		Id: "content",
		Class: "",
	}

	steps[5] = &Step {
		Element: "div",
		Id: "_viewerBox",
		Class: "",
	}

	steps[6] = &Step {
		Element: "div",
		Id: "",
		Class: "viewer_lst",
	}

	steps[7] = &Step {
		Element: "div",
		Id: "_imageList",
		Class: "",
	}

	return steps
}
