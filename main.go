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
	"runtime"
	"strconv"
	"strings"
)

const (
	commentMarker  = "#"
	gogetsFilename = ".gogets"

	latestTag = "latest"
)

var _stdout = log.New(os.Stdout, "", 0)
var _stderr = log.New(os.Stderr, "", 0)

var _commentRe = regexp.MustCompile(`\s*#.*$`)
var _verRe = regexp.MustCompile(`([^@]+)@(.*?)$`)

// get $HOME
func getHomePath() string {
	if usr, err := user.Current(); err != nil {
		_stderr.Fatal(err)
	} else {
		return usr.HomeDir
	}
	return ""
}

// load packages from given path (.gogets file)
func loadPackages(filepath string) (packages map[string]string, err error) {
	file, err := os.Open(filepath)

	if err == nil {
		defer file.Close()

		scanner := bufio.NewScanner(file)
		scanner.Split(bufio.ScanLines)

		var line string
		packages = map[string]string{}

		for scanner.Scan() {
			line = strings.TrimSpace(scanner.Text())

			// skip comments
			if strings.HasPrefix(line, commentMarker) {
				continue
			}

			// remove trailing comments
			line = _commentRe.ReplaceAllString(line, "")

			// skip empty line
			if len(line) > 0 {
				name, tag := lineToPackageNameAndTag(line)
				packages[name] = tag
			}
		}

		return packages, nil
	}

	return map[string]string{}, err
}

// parse given line to package name and tag
func lineToPackageNameAndTag(line string) (packageName, tag string) {
	packageName = line
	tag = latestTag

	if _verRe.Match([]byte(line)) {
		matches := _verRe.FindStringSubmatch(line)

		count := len(matches)
		if count > 2 {
			packageName = matches[count-2]
			tag = matches[count-1]
		}
	}

	return packageName, tag
}

// check if go module is on by default
func isGoModDefault() bool {
	// eg: "go1.15.8", "go1.16"
	segments := strings.Split(strings.Replace(runtime.Version(), "go", "", 1), ".")
	if len(segments) >= 2 {
		major, _ := strconv.ParseInt(segments[0], 10, 32)
		minor, _ := strconv.ParseInt(segments[1], 10, 32)

		// GO111MODULE=on since 1.16
		if major == 1 && minor < 16 {
			return false
		}

		return true
	}

	return false
}

// run command `go install` or `go get` with given package name and tag
func runGoInstallCommand(packageName, tag string) (output []byte, err error) {
	// for older versions with GOPATH mode
	if !isGoModDefault() {
		fmt.Printf("> go get -u %s ... ", packageName)
		return exec.Command("go", "get", "-u", packageName).CombinedOutput()
	}

	fmt.Printf("> go install %s@%s ... ", packageName, tag)
	return exec.Command("go", "install", packageName+"@"+tag).CombinedOutput()
}

// do install
func goInstall(packageName, tag string) (output string, err error) {
	var b []byte
	if b, err = runGoInstallCommand(packageName, tag); err == nil {
		fmt.Printf("=> successful\n")

		return string(b), nil
	}

	fmt.Printf("=> failed: %s\n%s----\n", err, string(b))

	return string(b), err
}

// print usage to stdout
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

// print sample config file to stdout
func printSample() {
	_stdout.Printf(`# sample .gogets file
#
# $ go install github.com/meinside/goget
# $ goget
#
# then it will automatically 'go install' all packages listed in this file(~/.gogets)

# without version (latest)
golang.org/x/tools/cmd/godoc

# with version tag
github.com/mailgun/godebug@latest
github.com/motemen/gore/cmd/gore@v0.5.2
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

func run() {
	homeDir := getHomePath()
	goGetsFilepath := strings.Join([]string{homeDir, gogetsFilename}, string(filepath.Separator))

	_stdout.Printf(">>> loading packages from: %s\n", goGetsFilepath)

	// chdir to $HOME,
	if err := os.Chdir(homeDir); err == nil {
		if packages, err := loadPackages(goGetsFilepath); err == nil {
			_stdout.Println()

			for pkg, tag := range packages {
				goInstall(pkg, tag)
			}
		} else {
			_stderr.Println(err)
		}
	} else {
		_stderr.Println(err)
	}
}

func main() {
	params := os.Args[1:]

	// check params
	if paramExists(params, "-h", "--help") {
		printUsage()
	} else if paramExists(params, "-g", "--generate") {
		printSample()
	}

	run()
}
