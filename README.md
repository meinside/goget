# GoGet

List Go package names in your **$XDG_CONFIG_HOME/goget/packages** file and install/update all of them at once.

## Install

```bash
# with go 1.16+
$ go install github.com/meinside/goget@latest

# or with older versions
$ go get -u github.com/meinside/goget
```

## Usage

```bash
# if $XDG_CONFIG_HOME is not set,
$ export XDG_CONFIG_HOME="$HOME/.config"

$ mkdir -p $XDG_CONFIG_HOME/goget/

# print a sample `packages` file
$ goget -g > $XDG_CONFIG_HOME/goget/packages

# install/update packages in `packages` file
$ goget
```
