package main

import (
	"github.com/kpango/fuid"
	"github.com/kpango/glg"
	"github.com/rs/xid"
)

func main() {
	for i := 0; i < 10; i++ {
		glg.Debug(fuid.String())
		glg.Info(xid.New().String())
	}

}
