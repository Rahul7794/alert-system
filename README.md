# Alert System

### Introduction

An "alerting" service which will consume a file of currency conversion rates and produce alerts for a number of situations.

- When the spot rate for a currency pair changes by more than 10% from the 5-minute average for that currency pair.
- When the spot rate has been rising/falling for 15 minutes. This alert should be
  throttled to only output once per minute and should report the length of time of the rise/fall in seconds.
  
### Project Layout

```
├── main.go                             # main function of the program.
├── README.md                           # README for the project.
├── makefile                            # Makefile (or makefile), defines set of tasks to be executed.
├── .gitignore                          # files/folders to be ignored pushing to git.
├── go.mod                              # go.mod is the module definition file.
├── go.sum                              # go.sum contains all the dependency checksums, and is managed by the go tools
├── alertprocessor                      # processing of currency rates and sending alerts are done here.
│   ├── processor.go
│   ├── alertprocessor.go
│   └── alertprocessor_test.go
├── cmd                                 # commands and base commands are placed here.
│   ├── alert.go                        # initialization of alert processor objects are done here.
│   ├── root.go
│   └── version.go
├── config                              # config files for the project.
│   └── config.go
├── examples                            # all the sample input and output test files.
├── file                                # reader and writer interface and functions are implemented here.
│   ├── io.go
│   ├── reader.go
│   └── writer.go
├── log                                 # loging interfaces are implemented here for the project.
│   └── log.go
├── model                               # structs and data models for the project.
│   └── model.go
├── movingmean                          # logic for calculating moving mean and trends for currency pairs is done here.
│   ├── movingmean.go
│   ├── movingmean_test.go
│   ├── queue.go                        # implementation of Queue data structure and basic operations of Queue.
│   └── queue_test.go
└── version                             # keeps track of version for a project. 
    └── version.go
```
  
### Steps to set up GO

- Install Go for different OS via [this](https://golang.org/doc/install)
- [Download](https://www.jetbrains.com/go/download/download-thanks.html) and Install Goland - (A cross-platform IDE built specially for Go developers)
- Once Setup import project folder/zip.
- Run `go get` to download all the dependencies

### Steps to build

```
# Build Go executable/artifact
make build
```

### Steps to run

```
./alert-system alert -i {input_path} -o {outpath}
```

### Steps to run the tests

`make test`

Note: **The application expects input path where the currency conversion rates is available
and output path where it wants to write the alerts.**

  
  

