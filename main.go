package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"text/template"

	"github.com/avamsi/ergo"
	"github.com/avamsi/ergo/assert"
	"github.com/google/shlex"
)

func defineFlags(args []string) (*flag.FlagSet, map[string]any, []string) {
	var (
		fset   = flag.NewFlagSet("", flag.ExitOnError)
		values = map[string]any{}
	)
	for i, arg := range args {
		if !strings.HasPrefix(arg, "-") {
			return fset, values, args[i:]
		}
		var (
			parts        = assert.Ok(shlex.Split(arg))
			name         = strings.TrimLeft(parts[0], "-")
			value, usage string
		)
		switch len(parts) {
		case 1:
		case 2:
			value = parts[1]
		case 3:
			value, usage = parts[1], parts[2]
		default:
			ergo.Panicf("not --name [value [usage]]: %v", arg)
		}
		values[name] = fset.String(name, value, usage)
	}
	panic(fmt.Sprintf("no command: %v", args))
}

func render(tmpl string, values map[string]any) (string, error) {
	var (
		t = template.Must(
			template.New("").
				Funcs(template.FuncMap{"join": strings.Join}).
				Parse(tmpl),
		)
		b   strings.Builder
		err = t.Execute(&b, values)
	)
	return b.String(), err
}

func run(noop, quiet bool, cmd *exec.Cmd) error {
	cmd.Env = os.Environ()
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	if !quiet {
		fmt.Println("$", cmd.String())
	}
	if noop {
		return nil
	}
	return cmd.Run()
}

func ash(noop, quiet bool, args []string) error {
	fset, values, args := defineFlags(args)
	assert.Nil(fset.Parse(args[1:]))
	values["args"] = fset.Args()
	cmd, err := render(args[0], values)
	if err != nil {
		return err
	}
	return run(noop, quiet, exec.Command("sh", "-c", cmd))
}

const help = `ash is a hybrid between getopts and 'sh -c'.

It lets you define CLI flags dynamically and run shell commands rendered from
Go templates. Flags and positional arguments can be referenced within the
command template -- flags by name, and the remaining arguments via {{.args}}.

Usage:
  ash [-n, --noop] [-q, --quiet] [-<name> [<value> [usage]]]... <template>

Example:
  alias rdiff='ash "-b main branch" "--remote origin" \
  	"git fetch {{.remote}} && git diff {{.remote}}/{{.b}}"'

  $ rdiff --help

  Usage:
    -b string
          branch (default "main")
    -remote string
           (default "origin")

See https://pkg.go.dev/text/template for more information on Go templates.`

func main() {
	var (
		noop, quiet bool
		args        = os.Args[1:]
	)
loop:
	for i, arg := range args {
		switch arg {
		case "help", "-h", "--help":
			fmt.Println(help)
			return
		case "-n", "--noop":
			noop = true
		case "-q", "--quiet":
			quiet = true
		default:
			args = args[i:]
			break loop
		}
	}
	if err := ash(noop, quiet, args); err != nil {
		fmt.Fprintf(os.Stderr, "ash: %v\n", err)
		if eerr := new(exec.ExitError); errors.As(err, &eerr) {
			os.Exit(eerr.ExitCode())
		}
		os.Exit(1)
	}
}
