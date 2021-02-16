# Crawler

A simple web crawler written on Go. It recursively crawls urls from defined url until the desired depth of recursion will be achieved.

## Install

Download the latest release binary file from: https://github.com/OrlovM/crawler/releases/

###Or 
Get the source file with: 
```
go get github.com/OrlovM/crawler
```
And build with:
```
make build
```


## How to use

Download and build project.

Run it with or without command line flags.

Example:

```
./crawler --startURL https://google.com --depth 3 --concurrency 20 -v
```

Command line flags:

--StartURL URL to start from URL to start from (default: "https://ya.ru")

--depth Depth refers to how far down into a website's page hierarchy crawler crawls (default 2)

--concurrency A maximum number of goroutines work at the same time (default: 50)

--verbose, -v Prints details about crawling process (default: false)




