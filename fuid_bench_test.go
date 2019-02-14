package fuid_test

import (
	"testing"

	gu "github.com/google/uuid"
	"github.com/kpango/fuid"
	"github.com/rs/xid"
	su "github.com/satori/go.uuid"
)

func BenchmarkFUID(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			fuid.String()
		}
	})
}

func BenchmarkXID(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			xid.New().String()
		}
	})
}

func BenchmarkSatoriUUID(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			su.NewV4()
		}
	})
}

func BenchmarkGoogleUUID(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			gu.New()
		}
	})
}
