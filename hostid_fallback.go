// +build !darwin,!linux,!freebsd,!windows

package fuid

import "errors"

func readPlatformMachineID() (string, error) {
	return "", errors.New("not implemented")
}
