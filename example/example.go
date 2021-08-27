package main

import (
	argvPkg "github.com/atercattus/argv"
)

var (
	argv struct {
		Name  string
		Count int

		Help bool
	}
)

func init() {
	argvPkg.BoolVar(&argv.Help, `help`, false, argvPkg.WithName(`h`), argvPkg.WithUsage(`Show this help`))

	argvPkg.StringVar(&argv.Name, `name`, `traveler`, argvPkg.WithEnv(`NAME`), argvPkg.WithUsage(`Your name`))
	argvPkg.IntVar(&argv.Count, `count`, 5, argvPkg.WithEnv(`COUNT`), argvPkg.WithUsage(`Times to repeat`))

	argvPkg.Parse()
}

func main() {
	switch {
	case argv.Help:
		argvPkg.Usage()
		return
	}

	for i := 0; i < argv.Count; i++ {
		println(`Hello, ` + argv.Name)
	}

	// Example:
	// go run . -h
	// NAME=kek go run . -count 2

	println(`Bye!`)
}
