package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
)

var (
	workDir         = "."
	workDirBasepath = ""
)

func main() {

	fmt.Print("Checking your Ginkgo Installation...")

	cmd := exec.Command("ginkgo", "help")
	if err := cmd.Run(); err != nil {
		color.Red("ERROR")
		log.Fatal("Please install ginkgo \n`go get github.com/onsi/ginkgo/ginkgo && go get github.com/onsi/gomega`")
	} else {
		color.Green("OK!")
	}

	err := os.Chdir(".")
	if err != nil {
		log.Fatal(err)
	}

	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	workDirBasepath = filepath.Base(wd)

	fInfos, err := ioutil.ReadDir(workDir)
	if err != nil {
		log.Fatal(err)
	}

	suiteFound := false

	fmt.Println("Generating specs")

	for _, finfo := range fInfos {
		if !finfo.IsDir() && validFilename(finfo.Name()) {
			cmd := exec.Command("ginkgo", "generate", finfo.Name())

			err := cmd.Run()
			if err != nil {
				log.Println(err)
			}
		}

		if isSuite(finfo.Name()) {
			suiteFound = true
		}
	}

	if !suiteFound {
		fmt.Println("Generating suite")
		cmd := exec.Command("ginkgo", "bootstrap", workDirBasepath)

		err := cmd.Run()
		if err != nil {
			log.Println(err)
		}

	}
}

func validFilename(filename string) bool {
	return !strings.Contains(filename, "_test") && strings.Contains(filename, ".go")
}

func isSuite(currFile string) bool {
	return strings.Contains(currFile, fmt.Sprintf("%s_suite", workDirBasepath))
}
