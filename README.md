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

## Example of `packages` file

```
# $XDG_CONFIG_HOME/goget/packages

# official packages
golang.org/x/tools/cmd/godoc
github.com/tools/godep

# useful packages
github.com/google/gops

# others
github.com/go-delve/delve/cmd/dlv
golang.org/x/tools/gopls@latest
#github.com/golangci/golangci-lint/cmd/golangci-lint@v1.47.2
github.com/golangci/golangci-lint/cmd/golangci-lint@latest
github.com/rclone/rclone
```

## License

MIT

