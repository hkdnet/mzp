package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

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

type colorConfig struct {
	fg int16
	bg int16
}
type config struct {
	gitColor colorConfig
}

type builder interface {
	build() (string, error)
}
type exprBuilder interface {
	buildExpr() (string, error)
}

type promptBuilder struct {
	builders []builder
}

func (pb *promptBuilder) buildExpr() (string, error) {
	ss := make([]string, len(pb.builders))
	for idx, b := range pb.builders {
		s, err := b.build()
		if err != nil {
			return "", err
		}
		ss[idx] = s
	}
	expr := strings.Join(ss, " ")

	return "PROMPT='" + expr + " %% '", nil
}

type rpromptBuilder struct {
	builders []builder
}

func (rpb *rpromptBuilder) buildExpr() (string, error) {
	return "RPROMPT=''", nil
}

type hostAndUserBuilder struct{}

func (b *hostAndUserBuilder) build() (string, error) {
	return "%n@%m", nil
}

type shorthandPwdBuilder struct{}

func (b *shorthandPwdBuilder) build() (string, error) {
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

type gitBuilder struct{}

func (gb *gitBuilder) build() (string, error) {
	repo, err := git.PlainOpen(".")
	if err != nil {
		if err == git.ErrRepositoryNotExists {
			return "no git", nil
		}
		logger.Println("cannot open git:", err)
		return cfg.gitColor.colorize("?"), nil
	}
	head, err := repo.Head()
	if err != nil {
		logger.Println("cannot fetch HEAD:", err)
		return cfg.gitColor.colorize("?"), nil
	}
	refName := head.Name()

	return cfg.gitColor.colorize(refName.Short()), nil
}

func run() (string, error) {
	exprBuiders := []exprBuilder{
		&promptBuilder{
			builders: []builder{
				&hostAndUserBuilder{},
				&shorthandPwdBuilder{},
				&gitBuilder{},
			},
		},
		&rpromptBuilder{},
	}
	exprs := make([]string, len(exprBuiders))
	for idx, eb := range exprBuiders {
		expr, err := eb.buildExpr()
		if err != nil {
			return "", err
		}
		exprs[idx] = expr
	}
	return strings.Join(exprs, " "), nil
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
