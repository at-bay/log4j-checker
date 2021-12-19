package main

import "fmt"

var (
	excludes        argsList
	includes        argsList
	logFileName     string
	verbose         bool
	ignoreV1        bool
	skipJpsDownload bool
	FoundVln        bool
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
