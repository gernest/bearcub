[![Build Status](https://travis-ci.org/gernest/bearcub.svg?branch=master)](https://travis-ci.org/gernest/bearcub) [![Coverage Status](https://coveralls.io/repos/github/gernest/bearcub/badge.svg?branch=master)](https://coveralls.io/github/gernest/bearcub?branch=master) [![GoDoc](https://godoc.org/github.com/gernest/bearcub?status.svg)](https://godoc.org/github.com/gernest/bearcub)


# Installation

```shell
go get -u github.com/gernest/bearcub
```

# Testing

Make sure the package is installed

```shell
cd $GOPATH/src/github.com/gernest/bearcub

go test -v
```

# Benchmarks

Make sure the package is installed

```shell
cd $GOPATH/src/github.com/gernest/bearcub

go test -run none -v  -bench=.
```