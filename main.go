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

const usage = `ash is a simple command runner (think sh -c) that supports defining flags
dynamically and executing (Go) templated commands (referring to said flags).

Usage:
  ash (--name [value [usage]])... <command ...{{.name}}> [options]`

func main() {
	if len(os.Args) == 2 && (os.Args[1] == "-h" || os.Args[1] == "--help") {
		fmt.Println(usage)
		return
	}
	var (
		cmdTmpl string
		fs      = flag.NewFlagSet("", flag.ExitOnError)
		opts    = map[string]*string{}
		i       int
		arg     string
	)
	for i, arg = range os.Args[1:] {
		if strings.HasPrefix(arg, "-") {
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
			opts[name] = &value
			fs.StringVar(&value, name, value, usage)
		} else {
			cmdTmpl = arg
			break
		}
	}
	assert.Nil(fs.Parse(os.Args[i+2:]))
	assert.Truef(len(fs.Args()) == 0, "not empty: %v", fs.Args())
	var (
		t = template.Must(template.New("").Parse(cmdTmpl))
		b strings.Builder
	)
	assert.Nil(t.Execute(&b, opts))
	cmd := exec.Command("sh", "-c", b.String())
	cmd.Env = os.Environ()
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	fmt.Println("$", cmd.String())
	if err := cmd.Run(); err != nil {
		if eerr := new(exec.ExitError); errors.As(err, &eerr) {
			os.Exit(eerr.ExitCode())
		}
		os.Exit(1)
	}
}
