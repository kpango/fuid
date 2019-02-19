package fuid

import (
	"context"
	"sync"
	"testing"
)

func TestString(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "pass",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := String(); got == "" {
				t.Errorf("String() = %v", got)
			}
		})
	}
}

func TestFUID_String(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "pass",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New().String(); got == "" {
				t.Errorf("FUID.String() = %v", got)
			}
		})
	}
}

func TestFUIDCollisionRate(t *testing.T) {
	tests := []struct {
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
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ch := make(chan string, tt.buf)
			wg := new(sync.WaitGroup)
			wg.Add(1)
			go func() {
				for i := 0; i < tt.routines; i++ {
					wg.Add(1)
					go func() {
						ch <- String()
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
					t.Log(counter)
				case str := <-ch:
					if check[str] == true {
						counter++
					} else {
						check[str] = true
					}
				}
			}
		})
	}
}
