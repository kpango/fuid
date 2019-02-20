package main

import (
	"context"
	"sync"
	"time"

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
		strf     func() string
	}{
		{
			name:     "fuid",
			routines: 100,
			cycle:    100000,
			buf:      100000,
			strf:     fuid.String,
		},
		{
			name:     "kpango/xid",
			routines: 100,
			cycle:    100000,
			buf:      100000,
			strf: func() string {
				return xid.New().String()
			},
		},
	}
	for _, e := range env {
		func() {
			start := time.Now()
			ch := make(chan string, e.buf)
			wg := new(sync.WaitGroup)
			wg.Add(1)
			go func() {
				for i := 0; i < e.routines; i++ {
					wg.Add(1)
					go func() {
						for j := 0; j < e.cycle; j++ {
							ch <- e.strf()
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
					glg.Info(e.name)
					glg.Info(counter)
					glg.Info(time.Now().Sub(start))
					return
				case str := <-ch:
					if check[str] == true {
						counter++
					} else {
						check[str] = true
					}
				}
			}
		}()
	}
}
