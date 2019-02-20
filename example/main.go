package main

import (
	"context"
	"sync"

	"github.com/kpango/fuid"
	"github.com/kpango/glg"
	"github.com/rs/xid"
)

func main() {
	for i := 0; i < 10; i++ {
		glg.Debugf("FUID:\t%s", fuid.String())
		glg.Infof("XID:\t%s", xid.New().String())
	}

	checkCollisionRate()
}

func checkCollisionRate() {
	env := []struct {
		name     string
		routines int
		cycle    int
		buf      int
	}{
		{
			name:     "normal",
			routines: 100,
			cycle:    100000,
			buf:      100000,
		},
	}
	for _, e := range env {
		ch := make(chan string, e.buf)
		wg := new(sync.WaitGroup)
		wg.Add(1)
		go func() {
			for i := 0; i < e.routines; i++ {
				wg.Add(1)
				go func() {
					for j := 0; j < e.cycle; j++ {
						ch <- fuid.String()
					}
					wg.Done()
				}()
			}
			wg.Done()
		}()
		ctx, cancel := context.WithCancel(context.Background())
		go func() {
			wg.Wait()
			cancel()
		}()
		check := map[string]bool{}
		counter := 0
		for {
			select {
			case <-ctx.Done():
				glg.Info(counter)
			case str := <-ch:
				if check[str] == true {
					counter++
				} else {
					check[str] = true
				}
			}
		}
	}
}
