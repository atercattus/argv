package argv

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type (
	Flag struct {
		Names []string
		Usage string

		envDefault      string
		envDefaultIsSet bool
	}

	Enrichmenter func(fl *Flag)
)

func Parse() {
	flag.Parse()
}

func BoolVar(p *bool, name string, defaultValue bool, opts ...Enrichmenter) {
	var fl Flag
	fl.Names = []string{name}
	for _, opt := range opts {
		opt(&fl)
	}

	if fl.envDefaultIsSet {
		switch strings.ToLower(fl.envDefault) {
		case `true`, `1`:
			defaultValue = true
		case `false`, `0`, ``:
			defaultValue = false
			// default: silent ignore?
		}
	}

	for _, name := range fl.Names {
		flag.BoolVar(p, name, defaultValue, fl.Usage)
	}
}

func StringVar(p *string, name string, defaultValue string, opts ...Enrichmenter) {
	var fl Flag
	fl.Names = []string{name}
	for _, opt := range opts {
		opt(&fl)
	}

	if fl.envDefaultIsSet {
		defaultValue = fl.envDefault
	}

	for _, name := range fl.Names {
		flag.StringVar(p, name, defaultValue, fl.Usage)
	}
}

func IntVar(p *int, name string, defaultValue int, opts ...Enrichmenter) {
	var fl Flag
	fl.Names = []string{name}
	for _, opt := range opts {
		opt(&fl)
	}

	if fl.envDefaultIsSet {
		defaultValue, _ = strconv.Atoi(fl.envDefault)
	}

	for _, name := range fl.Names {
		flag.IntVar(p, name, defaultValue, fl.Usage)
	}
}

func WithEnv(name string) Enrichmenter {
	return func(fl *Flag) {
		fl.envDefault, fl.envDefaultIsSet = os.LookupEnv(name)
	}
}

func WithName(name string) Enrichmenter {
	return func(fl *Flag) {
		fl.Names = append(fl.Names, name)
	}
}

func WithUsage(usage string) Enrichmenter {
	return func(fl *Flag) {
		fl.Usage = usage
	}
}

func ShowUsage() {
	_, _ = fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
	flag.PrintDefaults()
}
