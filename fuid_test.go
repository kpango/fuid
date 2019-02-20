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

func TestCollisionRate(t *testing.T) {
	tests := []struct {
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
			strf:     String,
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
						for j := 0; j < tt.cycle; j++ {
							ch <- tt.strf()
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
			check := make(map[string]bool)
			for {
				select {
				case <-ctx.Done():
					return
				case str := <-ch:
					if check[str] == true {
						t.Error("collision")
					} else {
						check[str] = true
					}
				}
			}
		})
	}
}
