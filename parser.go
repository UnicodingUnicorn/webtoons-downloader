package main

import (
	"golang.org/x/net/html"
)

type Step struct {
	Element, Id, Class string
}

func step(n *html.Node, i int, steps []*Step, endFn func(*html.Node)[][]string) [][]string {
	if i == len(steps) {
		return endFn(n)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if IsValid(c, steps[i]) {
				return step(c, i + 1, steps, endFn)
		}
	}

	return nil
}

func IsValid(n *html.Node, step *Step) bool {
	if n.Type == html.ElementNode && n.Data == step.Element {
		id := GetAttr(n, "id")
		class := GetAttr(n, "class")
		if (step.Id != "" && step.Class != "" && step.Id == id && step.Class == class) || (step.Id != "" && step.Id == id) || (step.Class != "" && step.Class == class) || (step.Id == "" && step.Class == "") {
			return true
		}
	}
	return false
}
