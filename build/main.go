package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	distRoot = "dist"
)

var (
	gooses   = &[]string{"windows", "linux"}
	goarches = &[]string{"amd64", "386"}
	name     = currentFolderName()
	version  = gitVersion()
)

func main() {
	goGenerate()

	for _, goos := range *gooses {
		for _, goarch := range *goarches {
			buildName := name + "-" + version + "-" + goos + "-" + goarch

			fmt.Println("Building " + buildName + "...")
			goBuild(goos, goarch, buildName)
		}
	}
}

// goGenerate executes the command "go generate"
func goGenerate() {
	cmd("go", "generate").Run()
}

// goBuild executes the command "go build" for the desired
// target OS and architecture, and writes the generated
// executable to the 'outDir' directory.
func goBuild(goos string, goarch string, name string) {
	os.Setenv("goos", goos)
	os.Setenv("goarch", goarch)

	out := filepath.Join(distRoot, name+exeSuffix())
	cmd("go", "build", "-o", out).Run()
}

// gitVersion returns the tag of the HEAD, if one exists,
// or else the commit hash.
func gitVersion() string {
	out := cmd("git", "tag", "--contains", "HEAD").Output()
	tag := strings.Split(string(out), "\n")[0]

	if tag != "" {
		return tag
	}

	out = cmd("git", "rev-parse", "--short", "HEAD").Output()
	return strings.Split(string(out), "\n")[0]
}

// currentFolderName returns the folder name
// of the current working directory
func currentFolderName() string {
	dir, err := os.Getwd()
	panicIf(err)
	return filepath.Base(dir)
}

// exeSuffix returns ".exe" if the GOOS
// environment variable is set to
// "windows".
func exeSuffix() string {
	if os.Getenv("GOOS") == "windows" {
		return ".exe"
	}
	return ""
}

func panicIf(err error) {
	if err != nil {
		panic(err.Error())
	}
}

////////////////////////////////////////////////////////////////////////////
// cmd
// Extends Go's exec.Cmd struct to enable automatic panicing
// when errors occur.
////////////////////////////////////////////////////////////////////////////

// cmd returns a Cmd struct
func cmd(name string, args ...string) *Cmd {
	cmd := exec.Command(name, args...)
	c := &Cmd{}
	c.Cmd = cmd
	return c
}

// Cmd extends exec Cmd
// from Go's standard library.
type Cmd struct {
	*exec.Cmd
}

// Run calls Go's cmd.Run() and panics
// is an error occurs.
func (c *Cmd) Run() {
	err := c.Cmd.Run()
	panicIf(err)
}

// Output calls Go's cmd.Output() and panics
// is an error occurs.
func (c *Cmd) Output() []byte {
	out, err := c.Cmd.Output()
	panicIf(err)
	return out
}
