<p align="center">
  <h3 align="center">pget</h3>

  <p align="center">
    Command-line utility to downloads a file in parallel
</p>

## About The Project

This project is a CLI (inspired by `wget`) that downloads a file in parallel by leveraging 
[HTTP range requests](https://developer.mozilla.org/en-US/docs/Web/HTTP/Range_requests).

### Installation

A binary can be built by running `go build`, with the Go compiler installed.

## Usage

The CLI can be used as follows:

```
pget [-c 10] [-v] https://url.to/file 
```

`c` is the number of chunks to split the file into

`v` activates verbose mode and outputs debug logs 
