package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	commentMarker  = "#"
	gogetsFilename = ".gogets"
)

var _stdout = log.New(os.Stdout, "", 0)
var _stderr = log.New(os.Stderr, "", 0)

func getHomePath() string {
	if usr, err := user.Current(); err != nil {
		_stderr.Fatal(err)
	} else {
		return usr.HomeDir
	}
	return ""
}

func loadPackages(filepath string) (packages []string, err error) {
	file, err := os.Open(filepath)

	if err == nil {
		defer file.Close()

		scanner := bufio.NewScanner(file)
		scanner.Split(bufio.ScanLines)

		var line string
		packages = []string{}
		re := regexp.MustCompile(`\s*#.*$`)

		for scanner.Scan() {
			line = strings.TrimSpace(scanner.Text())

			// skip comments
			if strings.HasPrefix(line, commentMarker) {
				continue
			}

			// remove trailing comments
			line = re.ReplaceAllString(line, "")

			// skip empty line
			if len(line) > 0 {
				packages = append(packages, line)
			}
		}

		return packages, nil
	}

	return []string{}, err
}

// do: go get -u packageName
func goGet(packageName string) (string, error) {
	b, err := exec.Command("go", "get", "-u", packageName).CombinedOutput()

	if err == nil {
		return string(b), nil
	}

	return string(b), err
}

func printUsage() {
	_stdout.Printf(`> usage:

# Show this help message

$ goget -h
$ goget --help


# Generate a sample .gogets file

$ goget -g
$ goget --generate


# Install/update all Go packages listed in ~/.gogets

$ goget
`)

	os.Exit(0)
}

func printSample() {
	_stdout.Printf(`# sample .gogets file
#
# $ go get -u github.com/meinside/goget
# $ goget
#
# then it will automatically 'go get -u' all packages listed in this file(~/.gogets)

# official
golang.org/x/tools/cmd/godoc

# useful packages
github.com/mailgun/godebug
github.com/motemen/gore
`)

	os.Exit(0)
}

// check if given params are given
func paramExists(params []string, shortParam string, longParam string) bool {
	for _, param := range params {
		if param == shortParam || param == longParam {
			return true
		}
	}
	return false
}

func main() {
	params := os.Args[1:]

	// check params
	if paramExists(params, "-h", "--help") {
		printUsage()
	} else if paramExists(params, "-g", "--generate") {
		printSample()
	}

	homeDir := getHomePath()
	goGetsFilepath := strings.Join([]string{homeDir, gogetsFilename}, string(filepath.Separator))

	_stdout.Printf(">>> loading packages from: %s\n", goGetsFilepath)

	// chdir to $HOME,
	if err := os.Chdir(homeDir); err == nil {
		if packages, err := loadPackages(goGetsFilepath); err == nil {
			_stdout.Println()

			for _, pkg := range packages {
				fmt.Printf("> go get -u %s ... ", pkg)

				if msg, err := goGet(pkg); err == nil {
					fmt.Printf("=> successful\n")
				} else {
					fmt.Printf("=> failed: %s\n%s----\n", err, msg)
				}
			}
		} else {
			_stderr.Println(err)
		}
	} else {
		_stderr.Println(err)
	}

}
