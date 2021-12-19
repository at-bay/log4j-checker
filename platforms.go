package main

type openJdk struct {
	url         string
	tgzFileName string
	sha256      string
	jpsPath     string
}

var openJdkPlatforms = map[string]openJdk{
	"linux": {
		url:         `https://github.com/adoptium/temurin17-binaries/releases/download/jdk-17.0.1%2B12/OpenJDK17U-jdk_x64_linux_hotspot_17.0.1_12.tar.gz`,
		tgzFileName: `OpenJDK17U-jdk_x64_linux_hotspot_17.0.1_12.tar.gz`,
		sha256:      `6ea18c276dcbb8522feeebcfc3a4b5cb7c7e7368ba8590d3326c6c3efc5448b6`,
		jpsPath:     `openjdk/jdk-17.0.1+12/bin/jps`,
	},
	"darwin": {
		url:         `https://github.com/adoptium/temurin17-binaries/releases/download/jdk-17.0.1%2B12/OpenJDK17U-jdk_x64_mac_hotspot_17.0.1_12.tar.gz`,
		tgzFileName: `OpenJDK17U-jdk_x64_mac_hotspot_17.0.1_12.tar.gz`,
		sha256:      `98a759944a256dbdd4d1113459c7638501f4599a73d06549ac309e1982e2fa70`,
		jpsPath:     `openjdk/jdk-17.0.1+12/Contents/Home/bin/jps`,
	},
	"windows": {
		url:         `https://github.com/adoptium/temurin17-binaries/releases/download/jdk-17.0.1%2B12/OpenJDK17U-jdk_x64_windows_hotspot_17.0.1_12.zip`,
		tgzFileName: `OpenJDK17U-jdk_x64_windows_hotspot_17.0.1_12.zip`,
		sha256:      `e5419773052ac6479ff211d5945f8625e0cdb036e69c0f71affaf02d5dc9aa0b`,
		jpsPath:     `openjdk\jdk-17.0.1+12\bin`,
	},
}
