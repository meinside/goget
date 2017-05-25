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
	CommentMarker  = "#"
	GogetsFilename = ".gogets"
)

func getHomePath() string {
	if usr, err := user.Current(); err != nil {
		log.Fatal(err)
	} else {
		return usr.HomeDir
	}
	return ""
}

func loadPackages(filepath string) (packages []string, err error) {
	if file, err := os.Open(filepath); err == nil {
		defer file.Close()

		scanner := bufio.NewScanner(file)
		scanner.Split(bufio.ScanLines)

		var line string
		packages = []string{}
		re := regexp.MustCompile(`\s*#.*$`)

		for scanner.Scan() {
			line = strings.TrimSpace(scanner.Text())

			// skip comments
			if strings.HasPrefix(line, CommentMarker) {
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
	} else {
		return []string{}, err
	}
}

// do: go get -u packageName
func goGet(packageName string) (string, error) {
	if b, err := exec.Command("go", "get", "-u", packageName).CombinedOutput(); err == nil {
		return string(b), nil
	} else {
		return string(b), err
	}
}

func printUsage() {
	fmt.Printf(`> Usage

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
	fmt.Printf(`# Sample .gogets file
# $ go get github.com/meinside/goget
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

	goGetsFilepath := strings.Join([]string{getHomePath(), GogetsFilename}, string(filepath.Separator))

	if packages, err := loadPackages(goGetsFilepath); err == nil {
		for _, pkg := range packages {
			fmt.Printf("> go get -u %s", pkg)

			if msg, err := goGet(pkg); err == nil {
				fmt.Printf(" => successful\n")
			} else {
				fmt.Printf(" => failed: %s\n%s----\n", err, msg)
			}
		}
	} else {
		fmt.Println(err.Error())
	}
}