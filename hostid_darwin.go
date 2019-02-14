// +build darwin

package fuid

import "syscall"

func readPlatformMachineID() (string, error) {
	return syscall.Sysctl("kern.uuid")
}
