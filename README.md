[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

### Description
`log4j-checker` tool helps identify whether a certain system is running a vulnerable version of the log4j library. Download and run the tool on each suspected system in your organization. Please refer to the usage section for more details.
This cli tool will try to locate java processes that are using a vulnerable version of log4j library and provide an indication if found.

### About At-Bay
At-Bay is an insurance provider founded by security experts specifically to address cyber risk.

### License
The project is licensed under MIT License.

### Supported Operating Systems
* linux (amd64)
* freebsd/darwin (MacOS) (amd64)

### Prerequisites
To automatically scan all running Java processes, you would need an installation of [jps](https://docs.oracle.com/javase/8/docs/technotes/tools/unix/jps.html).
First verify your Java version (from cmd: `java -version`) and install the appropriate JDK using:
```
# on Debian/Ubuntu systems
sudo apt install openjdk-VERSION-jdk-headless
```
Alternatively, specify (multiple pairs of) `--include PATH` argument to scan specific directories but not the currently running Java processes

### Legal Disclaimer
This project is made for non-commercial and ethical testing purposes only, and is offered as-is, without warranty. 

Use of `log4j-checker` for attacking targets is illegal. It is the end user's responsibility to obey all applicable local, state and federal laws. At-Bay assumes no liability and is not responsible for any misuse or damage caused by this program.

### Installation
Download the latest precompiled binary from the [releases page](https://github.com/at-bay/log4j-checker/releases)
or use Golang build tool: `GOOS=linux GOARCH=amd64 go build` from the root of this repository.

### Usage
To scan all running Java processes, we recommend running the tool as with root permissions:
```
Usage of sudo ./log4j-scanner-amd64-darwin-v1.0.5.bin:
  -exclude value
        path to exclude. example: -exclude PATH [-exclude ANOTHER]
  -ignore-v1
        ignore log4j 1.x versions checks
  -include value
        path to include. example: -include PATH [-include ANOTHER]
  -log string
        log file to write output to
  -verbose
        verbose output. without this flag, no output unless vulnerable
```

### Issues
Please report issues in the project [issues page](https://github.com/at-bay/log4j-checker/issues).