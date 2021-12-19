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

### Legal Disclaimer
This project is made for non-commercial and ethical testing purposes only, and is offered as-is, without warranty. 

Use of `log4j-checker` for attacking targets is illegal. It is the end user's responsibility to obey all applicable local, state and federal laws. At-Bay assumes no liability and is not responsible for any misuse or damage caused by this program.


### Prerequisites
#### Java
The log4j vulnerabilities are only relevant, in general, to machines running Java processes.
To know if Java runtime is installed on the machine at question, open terminal and hit (as root): `java -version`.
If no Java runtime is present, you can proceed to another machine.

### Installation

#### OpenJDK
To automatically scan the system' running Java processes (using the packaged `jps` program), `log4j-checker` tool downloads and extracts OpenJDK17.
If you don't want this behaviour, either:
* Run `log4j-checker` with the arg `--no-jps-download` to avoid the download and extract and have `jps` available in your `$PATH` (manually, on your own) 
* Run `log4j-checker` with the arg(s) `-include PATH` to not automatically detect the running Java process and provide a manual path to scan.

#### <a id="MyHeading"></a> Download Latest `log4j-checker`
```shell
wget -L https://github.com/at-bay/log4j-checker/releases/download/v1.1.0/log4j-checker-linux-amd64-v1.1.0.bin -O log4j-checker-linux-amd64-v1.1.0.bin
chmod +x log4j-checker-linux-amd64-v1.1.0.bin
```

### Usage
To scan *all* running Java processes on the current machine, we recommend running the tool with root permissions:
```
Usage of sudo ./log4j-checker-darwin-amd64-v1.1.0.bin:
  -exclude value
        path to exclude. example: -exclude PATH [-exclude ANOTHER]
  -ignore-v1
        ignore log4j 1.x versions
  -include value
        path to include. example -include PATH [-include ANOTHER]
  -log string
        log file to write output to
  -no-jps-download
        skip downloading of jps
  -verbose
        no output unless vulnerable
```

### Issues
Please report issues in the project [issues page](https://github.com/at-bay/log4j-checker/issues).