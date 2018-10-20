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

## Contributing

```
go get -u github.com/golang/dep/cmd/dep # if dep is not installed

git clone git@github.com:jonatasbaldin/ignore.git
cd ignore
dep ensure -v
go build
```

### License

[MIT](https://github.com/jonatasbaldin/ignore/blob/master/LICENSE).
