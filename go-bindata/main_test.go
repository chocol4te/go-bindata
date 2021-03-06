// This work is subject to the CC0 1.0 Universal (CC0 1.0) Public Domain
// Dedication license. Its contents can be found at:
// http://creativecommons.org/publicdomain/zero/1.0/

package main

import (
	"flag"
	"os"
	"reflect"
	"runtime"
	"testing"

	"github.com/shuLhan/go-bindata"
)

var (
	_traces = make([]byte, 1024)
)

func printStack() {
	var lines, start, end int

	runtime.Stack(_traces, false)

	for x, b := range _traces {
		if b != '\n' {
			continue
		}

		lines++
		if lines == 5 {
			start = x
		} else if lines == 7 {
			end = x + 1
			break
		}
	}

	os.Stderr.Write(_traces[start:end])
}

func assert(t *testing.T, exp, got interface{}, equal bool) {
	if reflect.DeepEqual(exp, got) == equal {
		return
	}

	printStack()

	t.Fatalf("\n"+
		">>> Expecting '%+v'\n"+
		"          got '%+v'\n", exp, got)
	os.Exit(1)
}

func TestParseArgs(t *testing.T) {
	oldArgs := os.Args
	defer func() {
		os.Args = oldArgs
	}()

	var (
		defConfig    = bindata.NewConfig()
		argInputPath = "."
		argPkg       = "pkgnametest"
		argOutPkg    = "assets"
		argOutFile   = argOutPkg + "/template.go"
	)

	tests := []struct {
		desc      string
		args      []string
		expConfig *bindata.Config
	}{{
		desc: `With "-pkg ` + argPkg + `"`,
		args: []string{
			"noop",
			"-pkg", argPkg,
			argInputPath,
		},
		expConfig: &bindata.Config{
			Output:  defConfig.Output,
			Package: argPkg,
			Input: []bindata.InputConfig{{
				Path: argInputPath,
			}},
			Ignore: defConfig.Ignore,
		},
	}, {
		desc: `With "-o ` + argOutFile + `" (package name should be "` + argOutPkg + `")`,
		args: []string{
			"noop",
			"-o", argOutFile,
			argInputPath,
		},
		expConfig: &bindata.Config{
			Output:  argOutFile,
			Package: argOutPkg,
			Input: []bindata.InputConfig{{
				Path: argInputPath,
			}},
			Ignore: defConfig.Ignore,
		},
	}, {

		desc: `With "-pkg ` + argPkg + ` -o ` + argOutPkg + `" (package name should be ` + argPkg + `)`,
		args: []string{
			"noop",
			"-pkg", argPkg,
			"-o", argOutFile,
			argInputPath,
		},
		expConfig: &bindata.Config{
			Output:  argOutFile,
			Package: argPkg,
			Input: []bindata.InputConfig{{
				Path: argInputPath,
			}},
			Ignore: defConfig.Ignore,
		},
	}}

	for _, test := range tests {
		t.Log(test.desc)

		os.Args = test.args

		flag.CommandLine = flag.NewFlagSet(test.args[0],
			flag.ExitOnError)

		got := parseArgs()

		assert(t, test.expConfig, got, true)
	}
}
