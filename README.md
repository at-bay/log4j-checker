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
To automatically scan the system running Java processes, you would need *Open Java Development Kit (OpenJDK)* that packs the `jps` command-line tool, available on your `$PATH`. If you don't want or can't install OpenJDK on your system, you can still run the checker tool in manual mode as [described below](#Usage).

To find out if `jps` is installed on your system simply issue the command `jps -h` on the target system terminal.
If `jps` is installed you would be returned with something similar to:
```shell
usage: jps [--help]
       jps [-q] [-mlvV] [<hostid>]
...
```
This means you can proceed to the [installation](#Installation) section below and skip the OpenJDK installation step straight to [downloading](#MyHeading) the `log4j-checker`.

### Installation
#### Download OpenJDK and extract
If `jps` is not installed that means that OpenJDK is not available on your `$PATH`, and you would need to download OpenJDK.

The below instructions *do not* install OpenJDK (as there is no need for actual installation), but download and extract a prepared OpenJDK binary and use it for the sole use of the `log4j-checker` tool. Feel free to delete the downloaded and extracted OpenJDK directory when you're done.
```shell
# for other versions and more instructions checkout: https://adoptopenjdk.net/installation.html
wget -L https://github.com/adoptium/temurin16-binaries/releases/download/jdk-16.0.2%2B7/OpenJDK16U-jdk_x64_linux_hotspot_16.0.2_7.tar.gz -O OpenJDK16U-jdk_x64_linux_hotspot_16.0.2_7.tar.gz
tar xzf OpenJDK16U-jdk_x64_linux_hotspot_16.0.2_7.tar.gz
export PATH=$PWD/jdk-16.0.2+7/bin:$PATH
```
#### <a id="MyHeading"></a> Download the latest log4j-checker
```shell
wget -L https://github.com/at-bay/log4j-checker/releases/download/v1.0.8/log4j-checker-linux-amd64-v1.0.8.bin -O log4j-checker-linux-amd64-v1.0.8.bin
chmod +x log4j-checker-linux-amd64-vv1.0.8.bin
```

### Usage
To scan *all* running Java processes, we recommend running the tool with root permissions:
```
Usage of sudo ./log4j-scanner-amd64-darwin-vv1.0.8.bin:
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