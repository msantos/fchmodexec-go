// Fchmodexec does an fchmod(2) on inherited file descriptors before
// exec(3)'ing a command.
//
// fchmodexec runs as part of an exec chain to change the permissions of
// any file descriptors inherited from the parent process before executing
// a program.
package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strconv"

	"golang.org/x/sys/unix"
)

const (
	version = "0.1.0"
)

func usage() {
	fmt.Fprintf(os.Stderr, `%s %s
Usage: <MODE> <FD> <...> -- <COMMAND> <...>
`, path.Base(os.Args[0]), version)
}

func fdget(argv []string) ([]int, error) {
	fds := make([]int, 0)

	for _, v := range argv {
		fd, err := strconv.Atoi(v)
		if err != nil {
			return fds, err
		}
		if _, err := unix.FcntlInt(uintptr(fd), unix.F_GETFD, 0); err != nil {
			continue
		}
		fds = append(fds, fd)
	}

	return fds, nil
}

func fdset(fds []int, mode uint32) error {
	for _, fd := range fds {
		if err := unix.Fchmod(fd, uint32(mode)); err != nil {
			return fmt.Errorf("%s: %w", fd, err)
		}
	}
	return nil
}

func at(a []string, s string) int {
	for n := 0; n < len(a); n++ {
		if s == a[n] {
			return n
		}
	}
	return -1
}

func main() {
	// 0: progname
	// 1: mode
	// 2..n: fd <...>
	// --: end of options
	// argv
	if len(os.Args) < 3 {
		usage()
		os.Exit(2)
	}

	mode, err := strconv.ParseInt(os.Args[1], 8, 32)
	if err != nil {
		fmt.Fprintln(os.Stderr, os.Args[1], err)
		os.Exit(2)
	}

	sep := at(os.Args[2:], "--") + 2
	argn := sep + 1

	if sep <= 2 || argn >= len(os.Args) {
		usage()
		os.Exit(2)
	}

	fds, err := fdget(os.Args[2:sep])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		usage()
		os.Exit(2)
	}

	if err := fdset(fds, uint32(mode)); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}

	arg0, err := exec.LookPath(os.Args[argn])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(127)
	}

	if err := unix.Exec(arg0, os.Args[argn:], os.Environ()); err != nil {
		fmt.Fprintln(os.Stderr, arg0, err)
	}

	os.Exit(126)
}