# GoGet

List Go package names in your **~/.gogets** file and install/update all of them at once.

## Install

```bash
# with go 1.16+
$ go install github.com/meinside/goget@latest

# or with older versions
$ go get -u github.com/meinside/goget
```

## Usage

```bash
# print a sample .gogets file
$ goget -g > ~/.gogets

# install/update packages in .gogets file
$ goget
```
