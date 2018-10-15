# Ignore

Download .gitignore files from the GitHub [gitignore](https://github.com/github/gitignore) repository!

[![Build Status](https://travis-ci.org/jonatasbaldin/ignore.svg?branch=master)](https://travis-ci.org/jonatasbaldin/ignore)

## Using it

Get the binary:

```
$ go get github.com/jonatasbaldin/ignore
```

Listing files:

```
$ ignore -list
Actionscript
Ada
Agda
Android
AppEngine
...
```

Downloading a file:

```
$ ignore Python
```

## For developer

```
cd /path/to/this/rep
dep ensure -v
go build
```

### License

[MIT](https://github.com/jonatasbaldin/ignore/blob/master/LICENSE).
