package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func TestFindingJars(t *testing.T) {
	lines := getJpsLines()
	found := findJars(lines)
	expected := []string{
		"/Applications/IntelliJ IDEA CE.app/Contents/lib/idea_rt.jar",
		"target/log4j-checkout-1.0-SNAPSHOT.jar", "log4j-checkout-1.0-SNAPSHOT.jar", `C:\Users\Administrator\Downloads\OpenJDK17U-jdk_x64_windows_hotspot_17.0.1_12\jdk-17.0.1+12\lib\jrt-fs.jar`,
	}
	for _, item := range expected {
		_, exists := Find(found, item)
		if !exists {
			t.Errorf("%s jar was not found in extracted list", item)
		}
	}
}

func TestLog4jJars(t *testing.T) {
	testData := []struct {
		version string
		vuln    bool
	}{
		{"2.14.0", true},
		{"2.17.0", false},
	}

	for _, testCase := range testData {
		cwd, _ := os.Getwd()
		path := fmt.Sprintf("%s/testdata/log4j-core-%s.jar", cwd, testCase.version)
		file, _ := os.Open(path)
		defer file.Close()
		buf, _ := ioutil.ReadAll(file)
		handleJar("", bytes.NewReader(buf), int64(len(buf)))
		if FoundVln != testCase.vuln {
			t.Errorf("version: %s should have been detected as vulnerabile: %v", testCase.version, testCase.vuln)
		}
		// reset
		FoundVln = false
	}

}

func TestParseDirs(t *testing.T) {
	lines := getJpsLines()
	expected := getExpected()
	found := findDirs(lines)
	for _, item := range expected {
		_, exists := Find(found, item)
		if !exists {
			t.Errorf("%s dir was not found in extracted list", item)
		}
	}
}

func getJpsLines() []string {
	line := `70087 org.jetbrains.jps.cmdline.Launcher -Xmx700m -Djava.awt.headless=true -Djdt.compiler.useSingleThread=true -Dpreload.project.path=/DirA/DirB/ExternalProjects/log4j-checkout -Dpreload.config.path=/DirA/DirB/Library/Application Support/JetBrains/IdeaIC2020.3/options`
	line1 := `70016  -Xms128m -Xmx2048m -XX:ReservedCodeCacheSize=240m -XX:+UseConcMarkSweepGC -XX:SoftRefLRUPolicyMSPerMB=50 -ea -XX:CICompilerCount=2 -Dsun.io.useCanonPrefixCache=false -Djdk.http.auth.tunneling.disabledSchemes="" -XX:+HeapDumpOnOutOfMemoryError -XX:-OmitStackTraceInFastThrow -Djdk.attach.allowAttachSelf=true -Dkotlinx.coroutines.debug=off -Djdk.module.illegalAccess.silent=true -XX:+UseCompressedOops -Dfile.encoding=UTF-8 -XX:ErrorFile=/DirA/DirB/java_error_in_idea_%p.log -XX:HeapDumpPath=/DirA/DirB/java_error_in_idea.hprof -Djb.vmOptionsFile=/DirA/DirB/Library/Application Support/JetBrains/IdeaIC2020.3/idea.vmoptions -Didea.paths.selector=IdeaIC2020.3 -Didea.executable=idea -Didea.platform.prefix=Idea -Didea.vendor.name=JetBrains -Didea.home.path=/Applications/IntelliJ IDEA CE.app/Contents`
	line2 := `17924 org.jetbrains.jps.cmdline.Launcher -Xmx700m -Djava.awt.headless=true -Djdt.compiler.useSingleThread=true -Dpreload.project.path=/DirA/DirB/ExternalProjects/log4j-checkout -Dpreload.config.path=/DirA/DirB/Library/Application Support/JetBrains/IdeaIC2020.3/options -Dexternal.project.config=/DirA/DirB/Library/Caches/JetBrains/IdeaIC2020.3/external_build_system/log4j-checkout.8cd078d7 -Dcompile.parallel=false -Drebuild.on.dependency.change=true -Dio.netty.initialSeedUniquifier=-7854164269029676952 -Dfile.encoding=UTF-8 -Duser.language=en -Duser.country=IL -Didea.paths.selector=IdeaIC2020.3 -Didea.home.path=/Applications/IntelliJ IDEA CE.app/Contents -Didea.config.path=/DirA/DirB/Library/Application Support/JetBrains/IdeaIC2020.3 -Didea.plugins.path=/DirA/DirB/Library/Application Support/JetBrains/IdeaIC2020.3/plugins -Djps.log.dir=/DirA/DirB/Library/Logs/JetBrains/IdeaIC2020.3/build-log -Djps.fallback.jdk.home=/Applications/IntelliJ IDEA CE.app/Contents/jbr/Contents/Home -Djps.fallback.jdk.version=11.0.10 -Dio.netty.noUnsaf`
	line3 := `17925 com.mykong.Main -Dlog4j.debug -javaagent:/Applications/IntelliJ IDEA CE.app/Contents/lib/idea_rt.jar=64745:/Applications/IntelliJ IDEA CE.app/Contents/bin -Dfile.encoding=UTF-8`
	line4 := `68532  -Xms128m -Xmx750m -XX:ReservedCodeCacheSize=512m -XX:+IgnoreUnrecognizedVMOptions -XX:+UseG1GC -XX:SoftRefLRUPolicyMSPerMB=50 -XX:CICompilerCount=2 -XX:+HeapDumpOnOutOfMemoryError -XX:-OmitStackTraceInFastThrow -ea -Dsun.io.useCanonCaches=false -Djdk.http.auth.tunneling.disabledSchemes="" -Djdk.attach.allowAttachSelf=true -Djdk.module.illegalAccess.silent=true -Dkotlinx.coroutines.debug=off -XX:ErrorFile=/DirA/DirB/java_error_in_goland_%p.log -XX:HeapDumpPath=/DirA/DirB/java_error_in_goland.hprof -Xmx1262m -Djb.vmOptionsFile=/DirA/DirB/Library/Application Support/JetBrains/GoLand2021.3/goland.vmoptions -Dsplash=true -Didea.home.path=/Applications/GoLand.app/Contents -Didea.executable=goland -Djava.system.class.loader=com.intellij.util.lang.PathClassLoader -Didea.platform.prefix=GoLand -Didea.paths.selector=GoLand2021.3 -Didea.vendor.name=JetBrains`
	line5 := `39142 jdk.jcmd/sun.tools.jps.Jps -Dapplication.home=/Library/Java/JavaVirtualMachines/adoptopenjdk-15.jdk/Contents/Home -Xms8m -Djdk.module.main=jdk.jcmd`
	line6 := `41946 target/log4j-checkout-1.0-SNAPSHOT.jar`
	line7 := `20635 org.apache.catalina.startup.Bootstrap --add-opens=java.base/java.lang=ALL-UNNAMED --add-opens=java.base/java.io=ALL-UNNAMED --add-opens=java.rmi/sun.rmi.transport=ALL-UNNAMED -Djava.util.logging.config.file=/var/lib/tomcat9/conf/logging.properties -Djava.util.logging.manager=org.apache.juli.ClassLoaderLogManager -Djava.awt.headless=true -XX:+UseG1GC -Djdk.tls.ephemeralDHKeySize=2048 -Djava.protocol.handler.pkgs=org.apache.catalina.webresources -Dorg.apache.catalina.security.SecurityListener.UMASK=0027 -Dignore.endorsed.dirs= -Dcatalina.base=/var/lib/tomcat9 -Dcatalina.home=/usr/share/tomcat9 -Djava.io.tmpdir=/tmp`
	line8 := `41946 log4j-checkout-1.0-SNAPSHOT.jar`
	line9 := `1300 jdk.jcmd/sun.tools.jps.Jps -Dapplication.home=C:\Users\Administrator\Downloads\OpenJDK17U-jdk_x64_windows_hotspot_17.0.1_12\jdk-17.0.1+12 -Xms8m -Djdk.module.main=jdk.jcmd`
	line10 := `3844  C:\Users\Administrator\Downloads\OpenJDK17U-jdk_x64_windows_hotspot_17.0.1_12\jdk-17.0.1+12\lib\jrt-fs.jar exit -Xms128m -Xmx750m -XX:ReservedCodeCacheSize=512m -XX:+IgnoreUnrecognizedVMOptions -XX:+UseG1GC -XX:SoftRefLRUPolicyMSPerMB=50 -XX:CICompilerCount=2 -XX:+HeapDumpOnOutOfMemoryError -XX:-OmitStackTraceInFastThrow -ea -Dsun.io.useCanonCaches=false -Djdk.http.auth.tunneling.disabledSchemes="" -Djdk.attach.allowAttachSelf=true -Djdk.module.illegalAccess.silent=true -Dkotlinx.coroutines.debug=off -Djb.vmOptionsFile=C:\Program Files\JetBrains\IntelliJ IDEA Community Edition 2021.3\bin\idea64.exe.vmoptions -Djava.system.class.loader=com.intellij.util.lang.PathClassLoader -Didea.vendor.name=JetBrains -Didea.paths.selector=IdeaIC2021.3 -Didea.platform.prefix=Idea -Didea.jre.check=true -Dsplash=true -Dide.native.launcher=true -XX:ErrorFile=C:\Users\Administrator\java_error_in_idea64_%p.log -XX:HeapDumpPath=C:\Users\Administrator\java_error_in_idea64.hprof`
	lines := []string{line, line1, line2, line3, line4, line5, line6, line7, line8, line9, line10}
	return lines
}

func getExpected() []string {
	lineExpected := []string{
		"/DirA/DirB/ExternalProjects/log4j-checkout",
		"/DirA/DirB/Library/Application Support/JetBrains/IdeaIC2020.3/options",
	}
	line1Expected := []string{
		"/DirA/DirB/Library/Application Support/JetBrains/IdeaIC2020.3/idea.vmoptions",
		"/DirA/DirB/java_error_in_idea.hprof", "/DirA/DirB/java_error_in_idea_%p.log",
	}
	line2Expected := []string{
		// "/DirA/DirB/ExternalProjects/log4j-checkout", not distinct
		"/DirA/DirB/Library/Application Support/JetBrains/IdeaIC2020.3/options",
		"/DirA/DirB/Library/Caches/JetBrains/IdeaIC2020.3/external_build_system/log4j-checkout.8cd078d7",
		"/DirA/DirB/Library/Application Support/JetBrains/IdeaIC2020.3",
		"/DirA/DirB/Library/Application Support/JetBrains/IdeaIC2020.3/plugins",
		"/DirA/DirB/Library/Logs/JetBrains/IdeaIC2020.3/build-log",
	}
	line3Expected := []string{"/Applications/IntelliJ IDEA CE.app/Contents/bin"}
	line4Expected := []string{
		"/DirA/DirB/java_error_in_goland_%p.log",
		"/DirA/DirB/Library/Application Support/JetBrains/GoLand2021.3/goland.vmoptions",
	}
	line5Expected := []string{"/Library/Java/JavaVirtualMachines/adoptopenjdk-15.jdk/Contents/Home"}
	line6Expected := []string{"/tmp", "/usr/share/tomcat9", "/var/lib/tomcat9/conf/logging.properties"}
	line9Expected := []string{`C:\Users\Administrator\Downloads\OpenJDK17U-jdk_x64_windows_hotspot_17.0.1_12\jdk-17.0.1+12`}
	line10Expected := []string{`C:\Users\Administrator\java_error_in_idea64_%p.log`, `C:\Users\Administrator\java_error_in_idea64.hprof`}
	var expected []string
	expected = append(expected, lineExpected...)
	expected = append(expected, line1Expected...)
	expected = append(expected, line2Expected...)
	expected = append(expected, line3Expected...)
	expected = append(expected, line4Expected...)
	expected = append(expected, line5Expected...)
	expected = append(expected, line6Expected...)
	expected = append(expected, line9Expected...)
	expected = append(expected, line10Expected...)
	return expected
}

// Find takes a slice and looks for an element in it. If found it will
// return it's key, otherwise it will return -1 and a bool of false.
func Find(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}
