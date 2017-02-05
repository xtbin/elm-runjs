package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/robertkrimen/otto"
)

func sysfail(format string, a ...interface{}) {
	logf(format, a...)
	os.Exit(1)
}

func usrfail(format string, a ...interface{}) {
	flag.Usage()
	logf("\n"+format, a...)
	os.Exit(2)
}

func logf(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", a...)
}

func main() {
	flag.Usage = func() {
		logf("usage: %s file.js", os.Args[0])
	}
	flag.Parse()
	if flag.NArg() != 1 {
		usrfail("error: exactly one file parameter required")
	}
	f, err := os.Open(flag.Arg(0))
	if err != nil {
		sysfail("error opening %q: %v", flag.Arg(0), err)
	}
	defer f.Close()
	vm := otto.New()
	err = stubvm(vm)
	if err != nil {
		sysfail("js init error: %v", err)
	}
	_, err = vm.Run(f)
	if err != nil {
		sysfail("js runtime error: %v", err)
	}
}

func stubvm(vm *otto.Otto) error {
	obj, err := vm.Object(`({on:function(){}})`)
	if err != nil {
		return err
	}
	return vm.Set("process", obj)
}
