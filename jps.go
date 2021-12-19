package main

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func untar(tarball, target string) error {
	reader, err := os.Open(tarball)
	if err != nil {
		return err
	}
	defer reader.Close()
	tarReader := tar.NewReader(reader)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		path := filepath.Join(target, header.Name)
		info := header.FileInfo()
		if info.IsDir() {
			if err = os.MkdirAll(path, info.Mode()); err != nil {
				return err
			}
			continue
		}

		file, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, info.Mode())
		if err != nil {
			return err
		}
		defer file.Close()
		_, err = io.Copy(file, tarReader)
		if err != nil {
			return err
		}
	}
	return nil
}

func unGzip(source, target string) error {
	reader, err := os.Open(source)
	if err != nil {
		return err
	}
	defer reader.Close()

	archive, err := gzip.NewReader(reader)
	if err != nil {
		return err
	}
	defer archive.Close()

	target = filepath.Join(target, archive.Name)
	writer, err := os.Create(target)
	if err != nil {
		return err
	}
	defer writer.Close()

	_, err = io.Copy(writer, archive)
	return err
}

func download(downloadFolder string, url string) {
	sep := `\`
	if strings.Index(url, "/") > -1 {
		sep = "/"
	}

	if url[len(url)-1:] == sep {
		return
	}

	urlParts := strings.Split(url, sep)
	filename := urlParts[len(urlParts)-1]

	// Just a simple GET request to the image URL
	// We get back a *Response, and an error
	res, err := http.Get(url) //nolint:gosec
	if err != nil {
		log.Printf("failed downloading from url: %s. error: %s", url, err)
		return
	}

	// We read all the bytes of the image
	// Types: data []byte
	data, err := ioutil.ReadAll(res.Body)

	if err != nil {
		log.Printf("error reading download content: %v\n", err)
		return
	}

	defer res.Body.Close()

	// You can now save it to disk or whatever...
	if err = ioutil.WriteFile(downloadFolder+string(filepath.Separator)+filename, data, 0400); err != nil {
		log.Println("error Saving:", filename, err)
	} else {
		log.Println("saved:", filename)
	}
}
