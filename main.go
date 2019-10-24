package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var (
	home string
)

func shorthandPwd() (string, error) {
	s, err := os.Getwd()
	if err != nil {
		return "", err
	}
	if strings.Index(s, home) == 0 { // start with home dir
		s = strings.Replace(s, home, "~", 1)
	}
	ss := strings.Split(s, string(filepath.Separator))
	for i := 0; i < len(ss)-2; i++ {
		s := ss[i]
		if len(s) == 0 {
			continue
		}
		if s[0] == '.' {
			ss[i] = s[0:2]
		} else {
			ss[i] = s[0:1]
		}
	}
	return filepath.Join(ss...), nil
}

func run() (string, error) {
	sp, err := shorthandPwd()
	if err != nil {
		return "", err
	}
	return "PROMPT='%n@%m " + sp + " %# '", nil
}

func init() {
	var err error
	home, err = os.UserHomeDir()
	if err != nil {
		panic(err)
	}
}

func main() {
	s, err := run()
	if err != nil {
		fmt.Print("PROMPT='failed > '")
		os.Exit(1)
	}
	fmt.Print(s)
}
