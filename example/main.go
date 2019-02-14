package main

import (
	"os"

	"github.com/kpango/fuid"
	"github.com/kpango/glg"
	"github.com/rs/xid"
)

func main() {
	pid := os.Getpid()
	for i := 0; i < 10; i++ {
		glg.Debug(fuid.String())
		glg.Info(xid.New().String())
		glg.Info(byte(pid >> 8))
		glg.Info(byte(pid))
	}

}
