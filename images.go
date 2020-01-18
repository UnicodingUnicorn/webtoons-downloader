package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"mime"
	"path"
	"strconv"
	"time"
)

func (d *Downloader) GetChaptersImages(chapterUrls []string, titles []string, imageUrls [][]string, update func(int)) error {
	// Should probably use a struct for this
	if !(len(chapterUrls) == len(titles) && len(titles) == len(imageUrls)) {
		return errors.New("parameter length mismatch")
	}

	for i := 0; i < len(titles); i++ {
		err := d.GetChapterImages(chapterUrls[i], titles[i], imageUrls[i], update)
		if err != nil {
			return err
		}
	}

	return nil
}

func (d *Downloader) GetChapterImages(chapterUrl string, title string, imageUrls []string, update func(int)) error {
	// Headers
	headers := make(map[string] string)
	headers["Referer"] = chapterUrl

	// Create padded format string for image names
	nameFmtStr := fmt.Sprintf("%%0%dd%%s", len(strconv.Itoa(len(imageUrls))))

	// Get images
	for i, imageUrl := range imageUrls {
		res, err := d.Client.GetPageRaw(imageUrl, headers)
		if err != nil {
			return err
		}

		// Get extension (if extension corresponding to content-type exists
		extension := ""
		if content_type, exists := res.Header["Content-Type"]; exists {
			if len(content_type) > 0 {
				extensions, err := mime.ExtensionsByType(content_type[0])
				if err != nil {
					return err
				}

				if extensions != nil && len(extensions) > 0 {
					extension = extensions[0]
				}
			}
		}

		// I'm lazy sue me
		data, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return err
		}

		// Make sure subfolder exists
		err = DirExists(path.Join(d.OutputDir, title), false)
		if err != nil {
			return err
		}

		// Write file
		err = ioutil.WriteFile(path.Join(d.OutputDir, title, fmt.Sprintf(nameFmtStr, i + 1, extension)), data, 0644)
		if err != nil {
			return err
		}

		// Don't trigger DDOS protection or something
		time.Sleep(d.Delay)

		update(0)
	}

	return nil
}
