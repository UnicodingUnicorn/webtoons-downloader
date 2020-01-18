package main

import (
	"log"
	"os"

	"golang.org/x/net/html"
)

func GetAttr(n *html.Node, key string) string {
	for _, attr := range n.Attr {
		if attr.Key == key {
			return attr.Val
		}
	}
	return ""
}

func DirExists(dir string, verbose bool) error {
	_, err := os.Stat(dir)
	if os.IsNotExist(err) {
		if verbose {
			log.Printf("directory %s doesn't exist, creating it...", dir)
		}
		err = os.Mkdir(dir, 0755)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}
	return nil
}
