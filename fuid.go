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
	s1 []byte
	c1 byte
	c2 byte
	c3 byte
	c4 byte
	c5 byte
	c6 byte
	c7 byte
	c8 byte
	c9 byte
	t  *fastime.Fastime
}

const (
	encodedLen = 20
	rawLen     = 7
	encoding   = "0123456789abcdefghijklmnopqrstuv"
)

var (
	instance = New()

	pid = func() int {
		p := os.Getpid()
		b, err := ioutil.ReadFile("/proc/self/cpuset")
		if err == nil && len(b) > 1 {
			p ^= int(crc32.ChecksumIEEE(b))
		}
		return p
	}()

	p1 = byte(pid >> 8)
	p2 = byte(pid)

	objectIDCounter = func() uint32 {
		b := make([]byte, 3)
		if _, err := rand.Reader.Read(b); err != nil {
			panic(fmt.Errorf("cannot generate random number: %v;", err))
		}
		return uint32(b[0])<<16 | uint32(b[1])<<8 | uint32(b[2])
	}()
)

func New() *FUID {
	machineID := func() []byte {
		id := make([]byte, 3)
		hid, err := readPlatformMachineID()
		if err != nil || len(hid) == 0 {
			hid, err = os.Hostname()
		}
		if err == nil && len(hid) != 0 {
			hw := md5.New()
			hw.Write([]byte(hid))
			copy(id, hw.Sum(nil))
		} else {
			if _, randErr := rand.Reader.Read(id); randErr != nil {
				panic(fmt.Errorf("fuid: cannot get hostname nor generate a random number: %v; %v", err, randErr))
			}
		}
		return id
	}()
	return &FUID{
		t:  fastime.New().StartTimerD(context.Background(), time.Nanosecond*10),
		c1: machineID[0] >> 5,
		c2: encoding[machineID[0]&0x1F],
		c3: encoding[machineID[1]>>3],
		c4: encoding[(machineID[2]>>6)&0x1F|(machineID[1]<<2)&0x1F],
		c5: encoding[(machineID[2]>>1)&0x1F],
		c6: encoding[(p1>>4)&0x1F|(machineID[2]<<4)&0x1F],
		c7: encoding[p2>>7|(p1<<1)&0x1F],
		c8: encoding[(p2>>2)&0x1F],
		c9: (p2 << 3) & 0x1F,
	}
}

func String() string {
	return instance.String()
}

func (f *FUID) String() string {
	var id [rawLen]byte
	binary.BigEndian.PutUint32(id[:], f.t.UnixUNow())
	i := atomic.AddUint32(&objectIDCounter, 1)
	id[4] = byte(i >> 16)
	id[5] = byte(i >> 8)
	id[6] = byte(i)
	dst := make([]byte, 0, encodedLen)
	dst = append(append(append(append(append(append(append(append(append(append(append(append(append(append(append(append(append(append(append(append(
		dst, encoding[id[0]>>3]),
		encoding[(id[1]>>6)&0x1F|(id[0]<<2)&0x1F]),
		encoding[(id[1]>>1)&0x1F]),
		encoding[(id[2]>>4)&0x1F|(id[1]<<4)&0x1F]),
		encoding[id[3]>>7|(id[2]<<1)&0x1F]),
		encoding[(id[3]>>2)&0x1F]),
		encoding[f.c1|(id[3]<<3)&0x1F]),
		f.c2),
		f.c3),
		f.c4),
		f.c5),
		f.c6),
		f.c7),
		f.c8),
		encoding[(id[4]>>5)|f.c9]),
		encoding[id[4]&0x1F]),
		encoding[id[5]>>3]),
		encoding[(id[6]>>6)&0x1F|(id[5]<<2)&0x1F]),
		encoding[(id[6]>>1)&0x1F]),
		encoding[(id[6]<<4)&0x1F])
	return *(*string)(unsafe.Pointer(&dst))
}
