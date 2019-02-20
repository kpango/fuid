package main

import (
	"github.com/kpango/fuid"
	"github.com/kpango/glg"
	"github.com/rs/xid"
)

func main() {
	for i := 0; i < 10; i++ {
		glg.Debugf("FUID:\t%s", fuid.String())
		glg.Infof("XID:\t%s", xid.New().String())
	}
}
