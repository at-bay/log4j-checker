package main

import "testing"

func TestFindingJars(t *testing.T) {
	lines := []string{"/amit.jar /moshe.jar", "-Dthing.this=/about/me.jar"}
	matches := findJars(lines)
	if len(matches) != 3 {
		t.Error("failed parsing jars")
	}
}

func TestLog4J(t *testing.T) {
	path := "/Users/amitmor/.m2/repository/org/apache/logging/log4j"
	findLog4j(path)
}
