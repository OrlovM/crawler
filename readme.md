# Crawler

[![CI](https://github.com/OrlovM/crawler/actions/workflows/build%20and%20test.yml/badge.svg)](https://github.com/OrlovM/crawler/actions/workflows/build%20and%20test.yml)
[![Hits-of-Code](https://hitsofcode.com/github/OrlovM/crawler)](https://hitsofcode.com/view/github/OrlovM/crawler)
[![codebeat badge](https://codebeat.co/badges/af3d23f7-38bc-4972-9864-bae84f9d39a0)](https://codebeat.co/projects/github-com-orlovm-crawler-master)

[![GitHub release](https://img.shields.io/github/release/OrlovM/crawler.svg?label=version)](https://github.com/OrlovM/crawler/releases/latest)
[![License](https://img.shields.io/github/license/OrlovM/crawler.svg?style=flat-square)](https://github.com/OrlovM/crawler/blob/master/LICENSE.md)

A simple web crawler written on Go. It recursively crawls urls from defined url until the desired depth of recursion will be achieved.  
The result will be placed in a `json` file.

## Install

Download the latest release binary file from: https://github.com/OrlovM/crawler/releases/

### Or 

Get the source file with: 
```
go get github.com/OrlovM/crawler
```
And build with:
```
make build
```


## How to use

Run it with or without command line flags.

Example:

```
./crawler --startURL https://google.com --depth 3 --concurrency 20 -v -p ./results/result.json
```

Command line flags:

`--StartURL` URL to start from (default: "https://ya.ru")

`--depth` Depth refers to how far down into a website's page hierarchy crawler crawls (default 2)

`--concurrency` A maximum number of goroutines work at the same time (default: 50)

`--verbose`, `-v` Prints details about crawling process (default: false)

`--filepath`, `-p` A path there the file with results of crawl will be created.
Could be specified with file name e.g. "/a/b.json" or only a directory "/a/".
If file name is not specified it will be set to default pattern "2006-01-02T15:04:05.json"`, (default: ./crawl_results/)




