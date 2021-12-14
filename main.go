package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

func findDirs(lines []string) []string {
	unix := `((?:[a-zA-Z]\:){0,1}(?:[\\/][\w.]+){1,})`
	re := regexp.MustCompile(unix)
	var found []string
	for _, line := range lines {
		matches := re.FindAllString(line, -1)
		for _, v := range matches {
			fmt.Println(v)
			found = append(found, v)
		}
	}
	return found
}

func findJars(lines []string) []string {
	re := regexp.MustCompile(`(/.*?.jar)`)

	var jars []string
	for _, line := range lines {
		matches := re.FindAllString(line, -1)
		for _, v := range matches {
			jars = append(jars, v)
		}
	}
	return jars
}

func verify() string {
	path, err := exec.LookPath("jps")
	if err != nil {
		log.Fatal("missing 'jps' command. please install the latest Oracle JDK")
	}

	fmt.Printf("found 'jps' command at %s\n", path)
	return path
}

func runJps(path string) ([]string, error) {
	cmd := exec.Command(path, "-l", "-v")
	r, err := cmd.StdoutPipe()

	if err != nil {
		log.Fatal(err)
	}

	var lines []string

	err = cmd.Start()
	if err != nil {
		log.Fatalf("failed running %s. error is: %v", path, err)
		return lines, err
	}
	cmd.Stderr = cmd.Stdout

	// Make a new channel which will be used to ensure we get all output
	done := make(chan struct{})

	// Create a scanner which scans r in a line-by-line fashion
	scanner := bufio.NewScanner(r)

	// Use the scanner to scan the output line by line and log it
	// It's running in a goroutine so that it doesn't block
	go func() {

		// Read line by line and process it
		for scanner.Scan() {
			line := scanner.Text()
			lines = append(lines, line)
		}

		// We're all done, unblock the channel
		done <- struct{}{}

	}()

	// Start the command and check for errors
	err = cmd.Start()

	// Wait for all output to be processed
	<-done

	// Wait for the command to finish
	err = cmd.Wait()
	return lines, nil
}

func findLog4j(root string) {
	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Fprintf(errFile, "%s: %s\n", path, err)
			return nil
		}
		if excludes.Has(path) {
			return filepath.SkipDir
		}
		if info.IsDir() {
			return nil
		}
		switch ext := strings.ToLower(filepath.Ext(path)); ext {
		case ".jar", ".war", ".ear":
			f, err := os.Open(path)
			if err != nil {
				fmt.Fprintf(errFile, "can't open %s: %v", path, err)
				return nil
			}
			defer f.Close()
			sz, err := f.Seek(0, os.SEEK_END)
			if err != nil {
				fmt.Fprintf(errFile, "can't seek in %s: %v", path, err)
				return nil
			}
			if _, err := f.Seek(0, os.SEEK_END); err != nil {
				fmt.Fprintf(errFile, "can't seek in %s: %v", path, err)
				return nil
			}
			handleJar(path, f, sz)
		default:
			return nil
		}
		return nil
	})
}

func main() {
	flag.Var(&excludes, "exclude", "paths to exclude")
	flag.BoolVar(&verbose, "verbose", false, "log every archive file considered")
	flag.StringVar(&logFileName, "log", "", "log file to write output to")
	flag.BoolVar(&quiet, "quiet", false, "no ouput unless vulnerable")
	flag.BoolVar(&ignore_v1, "ignore-v1", false, "ignore log4j 1.x versions")
	flag.Parse()

	if !quiet {
		fmt.Printf("%s - a simple local log4j vulnerability scanner\n\n", filepath.Base(os.Args[0]))
	}

	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s [--verbose] [--quiet] [--ignore-v1] [--exclude path] [ paths ... ]\n", os.Args[0])
		os.Exit(1)
	}

	if logFileName != "" {
		f, err := os.Create(logFileName)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Could not create log file")
			os.Exit(2)
		}
		logFile = f
		errFile = f
		defer f.Close()
	}

	if !quiet {
		fmt.Println("\nScan finished")
	}

	path := verify()
	lines, _ := runJps(path)
	jars := findJars(lines)
	for _, jar := range jars {
		fmt.Println(jar)
		findLog4j(jar)
	}

	dirs := findDirs(lines)
	for _, dir := range dirs {
		fmt.Println(dir)
		findLog4j(dir)
	}
}
