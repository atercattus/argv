package argv

import (
	"flag"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

type (
	Flag struct {
		Names []string
		Usage string

		typ             string
		defaultValue    interface{}
		envName         string
		envDefault      string
		envDefaultIsSet bool
	}

	Enrichmenter func(fl *Flag)
)

var (
	Usage = func() {
		_, _ = fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
		PrintDefaults()
	}

	allFlags []*Flag
)

func Parse() {
	flag.Parse()
}

// Args returns the non-flag command-line arguments.
func Args() []string {
	return flag.Args()
}

func PrintDefaults() {
	sort.Slice(allFlags, func(i, j int) bool {
		return allFlags[i].Names[0] < allFlags[j].Names[0]
	})

	for _, fl := range allFlags {
		sort.Strings(fl.Names)

		names := `-` + strings.Join(fl.Names, `, -`)
		if len(fl.envName) > 0 {
			names += `, env ` + fl.envName
		}

		usage := fl.Usage
		if len(usage) > 0 {
			usage = "\n  \t" + usage
		}

		defaultStr := ``
		if fl.defaultValue != nil {
			defaultStr = fmt.Sprintf(` (default: %v)`, fl.defaultValue)
		}

		_, _ = fmt.Fprintf(flag.CommandLine.Output(), "  %s %s%s%s\n", names, fl.typ, usage, defaultStr)
	}
}

func BoolVar(p *bool, name string, defaultValue bool, opts ...Enrichmenter) {
	var fl Flag
	fl.typ = `bool`
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

	fl.defaultValue = defaultValue

	allFlags = append(allFlags, &fl)

	for _, name := range fl.Names {
		flag.BoolVar(p, name, defaultValue, fl.Usage)
	}
}

func StringVar(p *string, name string, defaultValue string, opts ...Enrichmenter) {
	var fl Flag
	fl.typ = `string`
	fl.Names = []string{name}
	for _, opt := range opts {
		opt(&fl)
	}

	if fl.envDefaultIsSet {
		defaultValue = fl.envDefault
	}

	fl.defaultValue = defaultValue

	allFlags = append(allFlags, &fl)

	for _, name := range fl.Names {
		flag.StringVar(p, name, defaultValue, fl.Usage)
	}
}

func IntVar(p *int, name string, defaultValue int, opts ...Enrichmenter) {
	var fl Flag
	fl.typ = `int`
	fl.Names = []string{name}
	for _, opt := range opts {
		opt(&fl)
	}

	if fl.envDefaultIsSet {
		defaultValue, _ = strconv.Atoi(fl.envDefault)
	}

	fl.defaultValue = defaultValue

	allFlags = append(allFlags, &fl)

	for _, name := range fl.Names {
		flag.IntVar(p, name, defaultValue, fl.Usage)
	}
}

func WithEnv(name string) Enrichmenter {
	return func(fl *Flag) {
		fl.envName = name
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
