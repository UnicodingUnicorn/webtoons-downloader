package main

import (
	"fmt"
	"strings"
	"time"

	"golang.org/x/net/html"
)

func (d *Downloader) GetChapterUrls() ([][]string, error) {
	pages := make([]string, 0)
	prevPage := ""
	i := 1
	for {
		p, err := d.Client.GetPage(fmt.Sprintf("%s&page=%d", d.Url, i), nil)
		if err != nil {
			return nil, err
		}

		page := string(p)

		if (prevPage != "" && prevPage == page) || i > 8 {
			break
		}

		i++
		prevPage = page
		pages = append(pages, page)
		time.Sleep(500 * time.Millisecond)
	}

	chapterUrls := make([][]string, 0)
	steps := GetChapterListSteps()
	for _, page := range pages {
		doc, err := html.Parse(strings.NewReader(page))
		if err != nil {
			return nil, err
		}

		urls := step(doc, 0, steps, GetChapterPageUrls)
		chapterUrls = append(chapterUrls, urls...)
	}

	reversedChapterUrls := make([][]string, len(chapterUrls))
	for i := len(chapterUrls) - 1; i >= 0; i-- {
		reversedChapterUrls[(len(chapterUrls) - i - 1)] = chapterUrls[i]
	}

	return reversedChapterUrls, nil
}

func GetAllText(n *html.Node) string {
	var content strings.Builder
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode {
			content.WriteString(GetAllText(c))
		} else if c.Type == html.TextNode {
			content.WriteString(c.Data)
		}
	}
	return content.String()
}

func GetChapterTitle(n *html.Node) [][]string {
	var content strings.Builder
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		content.WriteString(GetAllText(c))
	}

	val := make([][]string, 1)
	val[0] = make([]string, 1)
	val[0][0] = content.String()
	return val
}

func GetChapterLinkUrl(n *html.Node) [][]string {
	urls := make([][]string, 1)
	urls[0] = make([]string, 2)
	link := GetAttr(n, "href")

	titleSteps := GetChapterTitleSteps()
	t := step(n, 0, titleSteps, GetChapterTitle)
	var title string
	if len(t) > 0 {
		title = t[0][0]
	}

	urls[0][0] = title
	urls[0][1] = link

	return urls
}

func GetChapterPageUrls(n *html.Node) [][]string {
	chapterUrls := make([][]string, 0)
	listSteps := GetChapterListLinkSteps()
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if IsValid(c, listSteps[0]) {
			url := step(c, 1, listSteps, GetChapterLinkUrl)
			if len(url) > 0 {
				chapterUrls = append(chapterUrls, url[0])
			}
		}
	}
	return chapterUrls
}

func GetChapterListSteps() []*Step {
	steps := make([]*Step, 9)

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
		Id: "",
		Class: "cont_box",
	}

	steps[6] = &Step {
		Element: "div",
		Id: "",
		Class: "detail_body banner",
	}

	steps[7] = &Step {
		Element: "div",
		Id: "",
		Class: "detail_lst",
	}

	steps[8] = &Step {
		Element: "ul",
		Id: "_listUl",
		Class: "",
	}

	return steps
}

func GetChapterListLinkSteps() []*Step {
	steps := make([]*Step, 2)

	steps[0] = &Step {
		Element: "li",
		Id: "",
		Class: "",
	}

	steps[1] = &Step {
		Element: "a",
		Id: "",
		Class: "",
	}

	return steps
}

func GetChapterTitleSteps() []*Step {
	steps := make([]*Step, 1)

	steps[0] = &Step {
		Element: "span",
		Id: "",
		Class: "subj",
	}

	return steps
}
