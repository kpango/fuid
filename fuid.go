package fuid

import (
	"context"
	"crypto/md5"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"io/ioutil"
	"os"
	"sync/atomic"
	"time"
	"unsafe"

	"github.com/kpango/fastime"
)

type FUID struct {
	ou  [rawLen]byte
	c1  byte
	c2  byte
	oic uint32
	t   *fastime.Fastime
}

const (
	encodedLen = 20
	rawLen     = 7
	encoding   = "0123456789abcdefghijklmnopqrstuv"
)

var (
	instance = New()
)

func New() *FUID {
	mid := func() []byte {
		id := make([]byte, 3)
		hid, err := readPlatformMachineID()
		if err != nil || len(hid) == 0 {
			hid, err = os.Hostname()
		}
		if err == nil && len(hid) != 0 {
			hw := md5.New()
			hw.Write(*(*[]byte)(unsafe.Pointer(&hid)))
			copy(id, hw.Sum(nil))
		} else {
			if _, randErr := rand.Reader.Read(id); randErr != nil {
				panic(fmt.Errorf("fuid: cannot get hostname nor generate a random number: %v; %v", err, randErr))
			}
		}
		return id
	}()

	pid := os.Getpid()
	b, err := ioutil.ReadFile("/proc/self/cpuset")
	if err == nil && len(b) > 1 {
		pid ^= int(crc32.ChecksumIEEE(b))
	}

	p1 := byte(pid >> 8)
	p2 := byte(pid)

	return &FUID{
		ou: [...]byte{
			encoding[mid[0]&0x1F],
			encoding[mid[1]>>3],
			encoding[(mid[2]>>6)&0x1F|(mid[1]<<2)&0x1F],
			encoding[(mid[2]>>1)&0x1F],
			encoding[(p1>>4)&0x1F|(mid[2]<<4)&0x1F],
			encoding[p2>>7|(p1<<1)&0x1F],
			encoding[(p2>>2)&0x1F],
		},
		c1: mid[0] >> 5,
		c2: (p2 << 3) & 0x1F,
		oic: func() uint32 {
			b := make([]byte, 3)
			if _, err := rand.Reader.Read(b); err != nil {
				panic(fmt.Errorf("cannot generate random number: %v;", err))
			}
			return uint32(b[0])<<16 | uint32(b[1])<<8 | uint32(b[2])
		}(),
		t: fastime.New().StartTimerD(context.Background(), time.Millisecond),
	}
}

func String() string {
	return instance.String()
}

func (f *FUID) String() string {
	var id [rawLen]byte
	binary.BigEndian.PutUint32(id[:], f.t.UnixUNow())
	i := atomic.AddUint32(&f.oic, 1)
	id[4] = byte(i >> 16)
	id[5] = byte(i >> 8)
	id[6] = byte(i)
	return *(*string)(unsafe.Pointer(&struct {
		array unsafe.Pointer
		len   int
		cap   int
	}{
		array: unsafe.Pointer(&[...]byte{
			f.ou[0],
			f.ou[1],
			f.ou[2],
			f.ou[3],
			f.ou[4],
			f.ou[5],
			f.ou[6],
			encoding[id[0]>>3],
			encoding[(id[1]>>6)&0x1F|(id[0]<<2)&0x1F],
			encoding[(id[1]>>1)&0x1F],
			encoding[(id[2]>>4)&0x1F|(id[1]<<4)&0x1F],
			encoding[id[3]>>7|(id[2]<<1)&0x1F],
			encoding[(id[3]>>2)&0x1F],
			encoding[f.c1|(id[3]<<3)&0x1F],
			encoding[(id[4]>>5)|f.c2],
			encoding[id[4]&0x1F],
			encoding[id[5]>>3],
			encoding[(id[6]>>6)&0x1F|(id[5]<<2)&0x1F],
			encoding[(id[6]>>1)&0x1F],
			encoding[(id[6]<<4)&0x1F],
		}),
		len: encodedLen,
		cap: encodedLen,
	}))
}
