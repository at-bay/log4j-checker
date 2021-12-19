package main

import (
	"flag"
	"fmt"
	"github.com/common-nighthawk/go-figure"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

var (
	dirRe         = `((?:[a-zA-Z]\:){0,1}(?:[\\/][\w.\-]+){1,})`
	compiledDirRe = regexp.MustCompile(dirRe)

	javaParamRe         = `(-[A-Zֿֿֿֿֿֿֿ\:]+(\w+\.*)+[=\:]*)`
	compiledJavaParamRe = regexp.MustCompile(javaParamRe)

	javaagentRe         = `-javaagent\:.*?=\d+\:`
	compiledJavaAgentRe = regexp.MustCompile(javaagentRe)

	sep = "/"
)

func findDirs(lines []string) []string {
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

var (
	unixJarRe         = `((?:\.?/\w+.*?.[jwe]ar)|(?:\w+/\w+.*?.[jwe]ar))`
	compiledUnixJarRe = regexp.MustCompile(unixJarRe)

	unixJarRe1         = `(?:[a-zA-Z]\:){0,1}[\w\+\-\.\\^/]+.[jwe]ar`
	compiledUnixJarRe1 = regexp.MustCompile(unixJarRe1)
)

func findJars(lines []string) []string {
	found := map[string]interface{}{}
	for _, line := range lines {
		line = strings.ReplaceAll(line, "-javaagent:", "")
		matches := compiledUnixJarRe.FindAllString(line, -1)
		for _, v := range matches {
			stripped := strings.Trim(v, " ")
			if len(stripped) > 0 {
				found[stripped] = nil
			}
		}

		if runtime.GOOS == "windows" {
			continue
		}

		// brute force matching
		matches = compiledUnixJarRe1.FindAllString(line, -1)
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

func createLogFile() *os.File {
	f, err := os.Create(logFileName)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Could not create log file")
		os.Exit(2)
	}
	logFile = f
	errFile = f
	return f
}

func createTmpDir() string {
	temp, err := os.MkdirTemp("", "lo4j-checker-")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	return temp
}

func scan(useJps bool, jpsInstallPath string) {
	// true, if user did NOT specify -include
	var err error
	if useJps {
		if skipJpsDownload {
			jpsInstallPath, err = verifyJpsInstalled()
			if err != nil {
				fmt.Fprintln(os.Stderr, missingJps)
				os.Exit(1)
			}
		}

		jpsOutputLines, err := runJps(jpsInstallPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed running jps scanning. error: %s", err)
			os.Exit(1)
		}
		jars := findJars(jpsOutputLines)
		for _, jar := range jars {
			findLog4j(jar)
		}
		includes = findDirs(jpsOutputLines)
	}

	for _, path := range includes {
		findLog4j(path)
	}
}

func main() {
	flag.Var(&excludes, "exclude", "path to exclude. example: -exclude PATH [-exclude ANOTHER]")
	flag.Var(&includes, "include", "path to include. example -include PATH [-include ANOTHER]")
	flag.StringVar(&logFileName, "log", "", "log file to write output to")
	flag.BoolVar(&verbose, "verbose", false, "no output unless vulnerable")
	flag.BoolVar(&ignoreV1, "ignore-v1", false, "ignore log4j 1.x versions")
	flag.BoolVar(&skipJpsDownload, "no-jps-download", false, "skip downloading of jps")
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
		f := createLogFile()
		defer f.Close()
	}

	var (
		jpsInstallPath string
		temp           string
	)

	if !skipJpsDownload {
		if verbose {
			fmt.Fprintf(os.Stdout, "downloading OpenJDK17 from adoptium.net\nextracted file and created temporary folders will be deleted upon termination\n")
		}
		temp, jpsInstallPath = getJps()
		if len(jpsInstallPath) == 0 {
			fmt.Fprintf(os.Stderr, "error downloading jps")
			os.Exit(1)
		}
		defer func(name string) {
			err := os.RemoveAll(name)
			if err != nil {
				fmt.Fprintf(os.Stderr, "failed deleting temporary folder at path: %s, with error: %s", temp, err)
			}
			if verbose {
				fmt.Fprintf(os.Stdout, "deleted temporary folder at path: %s", temp)
			}
		}(temp)
	}

	scan(useJps, jpsInstallPath)

	if verbose {
		fmt.Println("\nscan finished.")
	}

	if FoundVln {
		fmt.Printf("\n%s\n%s\n", foundVlnMsg, furtherInfoMsg)
	} else {
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
