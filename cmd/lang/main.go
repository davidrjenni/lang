// Copyright (c) 2023 David Jenni. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"davidrjenni.io/lang/compiler"
	"davidrjenni.io/lang/ir"
	"davidrjenni.io/lang/parser"
)

func main() {
	arg := ""
	if len(os.Args) > 1 {
		arg = os.Args[1]
	}

	switch arg {
	case "help", "":
		printHelp()

	case "run":
		if len(os.Args) < 3 {
			die("lang: no lang files listed")
		}

		filename := os.Args[2]
		b, err := parser.ParseFile(filename)
		if err != nil {
			die("%v\n", err)
		}

		asmFile, err := ioutil.TempFile("", "lang_build*.s")
		if err != nil {
			die("%v\n", err)
		}
		defer os.Remove(asmFile.Name())

		n := ir.Translate(b, ir.Loads)
		compiler.Compile(asmFile, n)
		if err := asmFile.Close(); err != nil {
			die("%v\n", err)
		}

		for _, a := range os.Args {
			if a == "-S" {
				sz, err := os.ReadFile(asmFile.Name())
				if err != nil {
					die("%v\n", err)
				}
				if err = ioutil.WriteFile(filename+".S", sz, 0644); err != nil {
					die("%v\n", err)
				}
				return
			}
		}

		exeFile, err := ioutil.TempFile("", "lang_build*.out")
		if err != nil {
			die("%v\n", err)
		}
		if err := exeFile.Close(); err != nil {
			die("%v\n", err)
		}
		defer os.Remove(exeFile.Name())

		asm := exec.Command("gcc", "-no-pie", asmFile.Name(), "-o", exeFile.Name())
		if err := asm.Run(); err != nil {
			die("%v\n", err)
		}

		run := exec.Command(exeFile.Name())
		run.Stdout = os.Stdout
		run.Stderr = os.Stderr
		if err := run.Run(); err != nil {
			os.Remove(asmFile.Name())
			os.Remove(exeFile.Name())
			die("%v\n", err)
		}

	default:
		dieUnknown()
	}
}

func die(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format, args...)
	os.Exit(1)
}

func dieUnknown() {
	die(`lang: unknown command
Run 'lang help' for usage.
`)
}

func printHelp() {
	fmt.Print(`usage: lang <cmd> [arguments]

Commands:
    run run a lang file
`)
}
