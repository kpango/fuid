package fuid_test

import (
	"testing"

	"github.com/kpango/fuid"
	"github.com/rs/xid"
)

func BenchmarkXID(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			xid.New().String()
		}
	})
}

func BenchmarkFUID(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			fuid.String()
		}
	})
}
