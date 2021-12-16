package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/common-nighthawk/go-figure"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

type argsList []string

func (flags *argsList) String() string {
	return fmt.Sprint(*flags)
}

func (flags *argsList) Set(value string) error {
	*flags = append(*flags, value)
	return nil
}

func (flags argsList) Has(path string) bool {
	for _, exclude := range flags {
		if path == exclude {
			return true
		}
	}
	return false
}

var (
	excludes    argsList
	includes    argsList
	logFileName string
	verbose     bool
	ignoreV1    bool
	foundVln    bool
)

func findDirs(lines []string) []string {
	dirRe := `((?:[a-zA-Z]\:){0,1}(?:[\\/][\w.\-]+){1,})`
	compiledDirRe := regexp.MustCompile(dirRe)

	javaParamRe := `(-[A-Zֿֿֿֿֿֿֿ\:]+(\w+\.*)+[=\:]*)`
	compiledJavaParamRe := regexp.MustCompile(javaParamRe)

	javaagentRe := `-javaagent\:.*?=\d+\:`
	compiledJavaAgentRe := regexp.MustCompile(javaagentRe)

	found := map[string]interface{}{}
	for _, line := range lines {
		matches := compiledJavaParamRe.Split(line, -1)
		for _, v := range matches {
			// handle javaagent param which might point to a lib/bin
			if compiledJavaAgentRe.MatchString(v) {
				javaAgentSplitted := compiledJavaAgentRe.Split(v, -1)
				if len(javaAgentSplitted) > 1 {
					// take the last element
					v = javaAgentSplitted[len(javaAgentSplitted)-1]
				}
			}
			if compiledDirRe.MatchString(v) {
				stripped := strings.Trim(v, " ")
				if len(stripped) > 0 {
					found[stripped] = nil
				}
			}
		}
	}
	distinct := mapKeysToSlice(found)

	return distinct
}

func mapKeysToSlice(m map[string]interface{}) []string {
	var keys []string //nolint:prealloc
	for key := range m {
		keys = append(keys, key)
	}

	return keys
}

func findJars(lines []string) []string {
	unixJarRe := `((?:/\w+.*?.[jwe]ar)|(?:\w+/\w+.*?.[jwe]ar))`
	re := regexp.MustCompile(unixJarRe)

	found := map[string]interface{}{}
	for _, line := range lines {
		line = strings.ReplaceAll(line, "javaagent:", "")
		matches := re.FindAllString(line, -1)
		for _, v := range matches {
			stripped := strings.Trim(v, " ")
			if len(stripped) > 0 {
				found[stripped] = nil
			}
		}
	}
	distinct := mapKeysToSlice(found)

	return distinct
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
	err = cmd.Start()

	// Wait for all output to be processed
	<-done

	// Wait for the command to finish
	if err = cmd.Wait(); err != nil {
		log.Fatalf("failed reading 'jps' output with error: %v", err)
	}

	return lines, nil
}

func findLog4j(root string) {
	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			if verbose {
				fmt.Fprintf(errFile, "skipping %s. %s\n", path, err)
			}
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
	return
}

func main() {
	flag.Var(&excludes, "exclude", "path to exclude. example: -exclude PATH [-exclude ANOTHER]")
	flag.Var(&includes, "include", "path to include. example -include PATH [-include ANOTHER]")
	flag.StringVar(&logFileName, "log", "", "log file to write output to")
	flag.BoolVar(&verbose, "verbose", false, "no output unless vulnerable")
	flag.BoolVar(&ignoreV1, "ignore-v1", false, "ignore log4j 1.x versions")
	flag.Parse()

	useJps := true
	if len(includes) > 0 {
		useJps = false
	}

	if verbose {
		myFigure := figure.NewColorFigure("At-Bay, Inc.", "", "blue", true)
		myFigure.Print()
		fmt.Printf("%s - a light, local log4j vulnerability scanner\n\n", filepath.Base(os.Args[0]))
	}

	if len(os.Args) < 1 {
		fmt.Fprintf(os.Stderr, "Usage: %s [--verbose] [--ignore-v1] [--exclude path] [ include ... ]\n", os.Args[0])
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

	// true, if user did NOT specify -include
	if useJps {
		jpsInstallPath, err := verifyJpsInstalled()
		if err != nil {
			fmt.Fprintln(os.Stderr, missingJps)
			os.Exit(1)
		}
		jpsOutputLines, _ := runJps(jpsInstallPath)
		jars := findJars(jpsOutputLines)
		for _, jar := range jars {
			findLog4j(jar)
		}
		includes = findDirs(jpsOutputLines)
	}

	for _, path := range includes {
		findLog4j(path)
	}

	if verbose {
		fmt.Println("\nscan finished.")
	}

	if foundVln {
		fmt.Printf("\n%s\n%s\n", foundVlnMsg, furtherInfoMsg)
	} else if verbose {
		fmt.Printf("\n%s\n%s\n", noVlnMsg, furtherInfoMsg)
	}
}

var missingJps = `
The 'jps' binary is not installed on your system. You need to either:
* install a version of Oracle JDK or OpenJDK depending on your current installation (use java -version to find what is your installation): 
  * official OpenJDK site: https://openjdk.java.net/projects/jdk/
  * official Oracle site (JDK17): https://docs.oracle.com/en/java/javase/17/install/installation-jdk-linux-platforms.html
  * DigitalOcean tutorial for installing various OpenJDK versions: https://www.digitalocean.com/community/tags/java?subtype=tutorial&q=openjdk
* run this with specific directory/ies to scan using the '-include' argument`

var foundVlnMsg = `the system is vulnerable, please update immediately.`
var noVlnMsg = `the system might not be vulnerable, but we encourage you to verify further with the system vendor.`
var furtherInfoMsg = `for details refer to the blog (https://www.at-bay.com/articles/security-alert-log4j/).`
