[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

### Description
`log4j-checker` tool helps identify whether a certain system is running a vulnerable version of the log4j library. Download and run the tool on each suspected system in your organization. Please refer to the usage section for more details.
This cli tool will try to locate java processes that are using a vulnerable version of log4j library and provide an indication if found.

### About At-Bay
At-Bay is an insurance provider founded by security experts specifically to address cyber risk.

### License
The project is licensed under MIT License.

### Supported Operating Systems
* [Linux](#DownloadLinux) (amd64)
* [Freebsd/darwin (MacOS)](#DownloadLinux) (amd64)
* [Windows](#DownloadWindows) (*experimental*) (amd64)

### Legal Disclaimer
This project is made for non-commercial and ethical testing purposes only, and is offered as-is, without warranty. 

Use of `log4j-checker` for attacking targets is illegal. It is the end user's responsibility to obey all applicable local, state and federal laws. At-Bay assumes no liability and is not responsible for any misuse or damage caused by this program.


### Prerequisites
#### Java
The log4j vulnerabilities are only relevant, in general, to machines running Java processes.
To know if Java runtime is installed on the machine at question, open terminal and hit (as root): `java -version`.
If no Java runtime is present, you can proceed to another machine.

### Installation

#### <a id="DownloadLinux"></a> Linux/Mac - Download Latest
* Download the latest release:
    ```shell
    # Linux
    wget -L https://github.com/at-bay/log4j-checker/releases/download/v1.1.4.alpha/log4j-checker-linux-amd64-v1.1.4.alpha.bin -O log4j-checker-linux-amd64-v1.1.4.alpha.bin
    # Mac
    wget -L https://github.com/at-bay/log4j-checker/releases/download/v1.1.4.alpha/log4j-checker-darwin-amd64-v1.1.4.alpha.bin -O log4j-checker-darwin-amd64-v1.1.4.alpha.bin
    ```
* (**Important**) Verify downloaded binary integrity (compare command results with the releases page SHA256):
    ```shell
    # Linux
    sha256sum log4j-checker-linux-amd64-v1.1.4.alpha.bin
    # Mac
    shasum -a 256 log4j-checker-darwin-amd64-v1.1.4.alpha.bin
    ```
* Make runnable:
    ```shell
    # Linux
    chmod +x log4j-checker-linux-amd64-v1.1.4.alpha.bin
    # Mac
    chmod +x log4j-checker-darwin-amd64-v1.1.4.alpha.bin
    ```



#### <a id="DownloadWindows"></a> Windows - Download Latest
* In your browser, open the [releases](https://github.com/at-bay/log4j-checker/releases) page 
* Click on the `log4j-checker-windows-amd64-v1.1.4.alpha.exe` file
* Your browser would most like notify you of the dangers of downloading exe files. Select to download anyway
* In your browser, click on the downloaded file dropdown menu and select "Show in folder"
* The file explorer would then open up. Click the address bar in it so that the download location would appear as a highlighted text such as: `C:\Users\Administrator\Downloads` and copy that address
* Open Windows `Command Prompt` application as an Administrator
  
* (**Important**) Verify the integrity of the downloaded exe file by pasting to the command prompt:
  
    ```
    CertUtil -hashfile C:\Users\Administrator\Downloads\log4j-checker-windows-amd64-v1.1.4.alpha.exe SHA256
    ```

* The result would print the SHA256 of the downloaded exe file and should match the published SHA256 of the release you've used
* Head to the [Windows Usage](#UsageWindows) section below

### Usage

#### Linux/Mac
To scan *all* running Java processes on the current machine, we recommend running the tool with root permissions:
```
Usage of sudo ./log4j-checker-darwin-amd64-v1.1.4.alpha.bin:
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

#### <a id="UsageWindows"></a> Windows
To scan *all* running Java processes on the current machine, we recommend running the tool with Administrator permissions:
* Open `Command Prompt` as `Adminstrator`
* Go to the tool download directory, for example: `chdir C:\Users\Administrator\Downloads`
* Run the tool (same options as for Linux/Mac apply): `log4j-checker-windows-amd64-v1.1.4.alpha.exe` 

### Issues
Please report issues in the project [issues page](https://github.com/at-bay/log4j-checker/issues).

### Notes

#### OpenJDK
To automatically scan the system' running Java processes, `log4j-checker` downloads and extracts OpenJDK17.
If you don't want this behaviour, either:
* Run `log4j-checker` with the arg `--no-jps-download` to avoid the download and extract. In that case `log4j-checker` expects to have `jps` (a tool that is part of OpenJDK) available in your `$PATH`
* Run `log4j-checker` with the arg(s) `-include PATH` to not automatically detect the running Java process and provide a manual path to scan.

#### Performance Hit
The scan can take up to a couple of minutes, depending mostly on the number of running Java processes and the amount of libraries those are dependent on. 

The scan might hinder your system performance during its run so be extra mindful to the possibility that the system would incur a performance hit during that time 