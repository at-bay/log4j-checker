package main

import (
	"archive/tar"
	"archive/zip"
	"bufio"
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// unTarUnGzip takes a destination path and a reader; a tar reader loops over the tarfile
// creating the file structure at 'dst' along the way, and writing any files
func unTarUnGzip(srcReader io.Reader, dst string) error {
	gzr, err := gzip.NewReader(srcReader)
	if err != nil {
		return err
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)
	for {
		header, err := tr.Next()
		switch {
		// if no more files are found return
		case err == io.EOF:
			return nil
		// return any other error
		case err != nil:
			return err
		// if the header is nil, just skip it (not sure how this happens)
		case header == nil:
			continue
		}
		// the target location where the dir/file should be created
		target := filepath.Join(dst, header.Name)
		// the following switch could also be done using fi.Mode(), not sure if there
		// a benefit of using one vs. the other.
		// fi := header.FileInfo()

		// check the file type
		switch header.Typeflag {
		// if it's a dir and it doesn't exist create it
		case tar.TypeDir:
			if _, err := os.Stat(target); err != nil {
				if err := os.MkdirAll(target, 0755); err != nil {
					return err
				}
			}

		// if it's a file create it
		case tar.TypeReg:
			f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return err
			}

			// copy over contents
			if _, err := io.Copy(f, tr); err != nil {
				return err
			}

			// manually close here after each file operation; defering would cause each file close
			// to wait until all operations have completed.
			f.Close()
		}
	}
}

func unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer func() {
		if err := r.Close(); err != nil {
			panic(err)
		}
	}()

	os.MkdirAll(dest, 0400)

	// Closure to address file descriptors issue with all the deferred .Close() methods
	extractAndWriteFile := func(f *zip.File) error {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer func() {
			if err := rc.Close(); err != nil {
				panic(err)
			}
		}()

		path := filepath.Join(dest, f.Name)

		// Check for ZipSlip (Directory traversal)
		if !strings.HasPrefix(path, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("illegal file path: %s", path)
		}

		if f.FileInfo().IsDir() {
			os.MkdirAll(path, f.Mode())
		} else {
			os.MkdirAll(filepath.Dir(path), f.Mode())
			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer func() {
				if err := f.Close(); err != nil {
					panic(err)
				}
			}()

			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}
		return nil
	}

	for _, f := range r.File {
		err := extractAndWriteFile(f)
		if err != nil {
			return err
		}
	}

	return nil
}

func download(downloadFolder string, jdk openJdk) error {
	sep := `\`
	if strings.Index(jdk.url, "/") > -1 {
		sep = "/"
	}

	urlParts := strings.Split(jdk.url, sep)
	filename := urlParts[len(urlParts)-1]

	// Just a simple GET request to the image URL
	// We get back a *Response, and an error
	res, err := http.Get(jdk.url) //nolint:gosec
	if err != nil {
		return err
	}

	// We read all the bytes of the image
	// Types: data []byte
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if err := verifyDownload(data, jdk.sha256); err != nil {
		return err
	}
	// You can now save it to disk or whatever...
	if err = ioutil.WriteFile(downloadFolder+string(filepath.Separator)+filename, data, 0400); err != nil {
		return err
	} else if verbose {
		fmt.Fprintf(os.Stdout, "saved: %s\n", filename)
		return nil
	}
	return nil
}

func verifyDownload(b []byte, hash string) error {
	hasher := sha256.New()
	io.Copy(hasher, bytes.NewBuffer(b))

	sum := hex.EncodeToString(hasher.Sum(nil))
	if sum != hash {
		return fmt.Errorf("sha256 of downloaded file: %s does not match expected: %s", sum, hash)
	}

	return nil
}

func verifyJpsInstalled() (string, error) {
	path, err := exec.LookPath("jps")
	if err != nil {
		return "", err
	}
	if verbose {
		fmt.Printf("found 'jps' command at %s\n", path)
	}
	return path, nil
}

func runJps(path string) ([]string, error) {
	cmd := exec.Command(path, "-l", "-v")
	r, err := cmd.StdoutPipe()
	if err != nil {
		wrappedErr := fmt.Errorf("failed execing command: %s. error is: %w", path, err)
		return nil, wrappedErr
	}

	var lines []string

	err = cmd.Start()
	if err != nil {
		wrappedErr := fmt.Errorf("failed running %s. error is: %w", path, err)
		return lines, wrappedErr
	}
	cmd.Stderr = cmd.Stdout
	// Make a new channel which will be used to ensure we get all output
	done := make(chan struct{})
	// Create a scanner which scans r in a line-by-line fashion
	scanner := bufio.NewScanner(r)
	// Use the scanner to scan the output line by line and log if it's running in a goroutine so that it doesn't block
	go func() { // Read line by line and process it
		for scanner.Scan() {
			line := scanner.Text()
			lines = append(lines, line)
		}
		// We're all done, unblock the channel
		done <- struct{}{}
	}()
	// Start the command and check for errors
	cmd.Start()
	// Wait for all output to be processed
	<-done
	// Wait for the command to finish
	if err := cmd.Wait(); err != nil {
		return nil, fmt.Errorf("failed reading 'jps' output with error: %w", err)
	}

	return lines, nil
}

func getJps(jpsInstallPath string) (string, string) {
	temp := createTmpDir()

	if verbose {
		fmt.Fprintf(os.Stdout, "created tmp dir to download openjdk: %s\n", temp)
	}

	currentOS := runtime.GOOS
	jdk := openJdkPlatforms[currentOS]

	if err := download(temp, jdk); err != nil {
		fmt.Fprintf(os.Stderr, "failed downloading from url: %s. error: %s\n", jdk.url, err)
	}

	tmpFileSrc := temp + sep + jdk.tgzFileName
	tmpFileDst := temp + sep + "openjdk"
	if currentOS != "windows" {
		file, err := os.Open(tmpFileSrc)
		if err != nil {
			fmt.Printf("error opening file: %s. error is: %v\n", tmpFileSrc, err)
			return "", ""
		}
		defer file.Close()
		err = unTarUnGzip(file, tmpFileDst)
		if err != nil {
			return "", ""
		}
	} else {
		sep = `\`
		if err := unzip(tmpFileSrc, tmpFileDst); err != nil {
			fmt.Printf("error unzipping file: %s. error is: %v\n", tmpFileSrc, err)
			return "", ""
		}
	}

	jpsInstallPath = temp + sep + jdk.jpsPath
	return temp, jpsInstallPath
}
