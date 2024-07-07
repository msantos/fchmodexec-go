// Fchmodexec sets permissions on a list of file descriptors.
package fchmodexec

import (
	"fmt"
	"strconv"

	"golang.org/x/sys/unix"
)

// Get returns a list of open file descriptors.
func Get(s []string) ([]int, error) {
	fds := make([]int, 0)

	for _, v := range s {
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

// Set changes the permissions of a list of file descriptors.
func Set(fds []int, mode uint32) error {
	for _, fd := range fds {
		if err := unix.Fchmod(fd, uint32(mode)); err != nil {
			return fmt.Errorf("%s: %w", fd, err)
		}
	}
	return nil
}
