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
