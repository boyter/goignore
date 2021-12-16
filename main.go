package main

import (
	"fmt"
	"github.com/gobwas/glob"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

func main() {
	readDir("./", nil) // account for the fact that ./ is not counted
}

// recursively walks all directories starting from this one
func readDir(name string, gitIgnoreGlob []Globs) {
	d, err := os.Open(name)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer d.Close()

	files, err := d.ReadDir(-1)
	if err != nil {
		fmt.Println(err.Error())
	}

	for _, f := range files {
		if f.Name() == ".gitignore" || f.Name() == ".ignore" {
			content, err := ioutil.ReadFile(path.Join(name, f.Name()))
			if err != nil {
				fmt.Println(err.Error())
			}
			fmt.Println("name>>>", name, f.Name())
			ignore := parseIgnore(string(content), name)
			gitIgnoreGlob = append(gitIgnoreGlob, ignore...)
		}
	}

	for _, f := range files {
		process := true
		for _, g := range gitIgnoreGlob {
			fmt.Println(path.Join(strings.TrimLeft(name, "./"), f.Name()), g.Pattern,  g.Glob.Match(path.Join(strings.TrimLeft(name, "./"), f.Name())))
			if g.Glob.Match(path.Join(strings.TrimLeft(name, "./"), f.Name())) {
				process = false
				break
			}
		}

		if !process {
			continue
		}

		if f.IsDir() {
			if f.Name() == ".git" {
				continue
			}

			readDir(path.Join(name, f.Name()), gitIgnoreGlob)
		} else {
			content, _ := ioutil.ReadFile(path.Join(name, f.Name()))

			if strings.HasPrefix(string(content), "foo:") {
				fmt.Println(len(gitIgnoreGlob), path.Join(name, f.Name()), strings.TrimSpace(string(content)))
			}
		}
	}
}

type Globs struct {
	Glob glob.Glob
	Pattern string
}

func parseIgnore(content string, directory string) []Globs {
	globPatterns := GlobifyGitIgnore(content, directory)

	var compiledGlob []Globs
	for _, s := range globPatterns {
		compiled, err := glob.Compile(s)
		if err == nil {
			compiledGlob = append(compiledGlob, Globs{
				Glob:    compiled,
				Pattern: s,
			})
		} else {
			fmt.Println(err.Error())
		}
	}

	return compiledGlob
}
