package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"

	"gopkg.in/src-d/go-git.v4"
)

var (
	home   string
	pwd    string
	logger *log.Logger
	cfg    *config
)

func (cc *colorConfig) colorize(s string) string {
	return fmt.Sprintf("\u001b[38;5;%dm\u001b[48;5;%dm %s \u001b[0m", cc.fg, cc.bg, s)
}

type gitStatus struct {
	branch string
}
type colorConfig struct {
	fg int16
	bg int16
}
type config struct {
	gitColor colorConfig
}

func fetchGitStatus() (*gitStatus, error) {
	repo, err := git.PlainOpen(".")
	if err != nil {
		return nil, errors.Wrap(err, "cannot open git")
	}
	head, err := repo.Head()
	if err != nil {
		return nil, errors.Wrap(err, "cannot fetch HEAD")
	}
	refName := head.Name()
	ret := &gitStatus{
		branch: refName.Short(),
	}
	return ret, nil
}

func shorthandPwd() (string, error) {
	s := pwd
	if strings.Index(s, home) == 0 { // start with home dir
		s = strings.Replace(s, home, "~", 1)
	}
	ss := strings.Split(s, string(filepath.Separator))
	for i := 0; i < len(ss)-1; i++ {
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
	gs, err := fetchGitStatus()
	if err != nil {
		return "", err
	}
	return "PROMPT='%n@%m " + sp + " " + cfg.gitColor.colorize(gs.branch) + " %% '", nil
}

func init() {
	var err error
	home, err = os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	logpath := filepath.Join(home, ".mzp", "mzp.log")
	err = os.MkdirAll(filepath.Dir(logpath), 0777)
	if err != nil {
		panic(err)
	}
	f, err := os.OpenFile(logpath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}
	logger = log.New(f, "", log.LstdFlags)

	pwd, err = os.Getwd()
	if err != nil {
		panic(err)
	}

	cfg = &config{
		gitColor: colorConfig{
			bg: 42,
			fg: 0,
		},
	}
}

func main() {
	s, err := run()
	if err != nil {
		logger.Println(err)
		fmt.Print("PROMPT='failed > '")
		os.Exit(1)
	}
	fmt.Print(s)
}
