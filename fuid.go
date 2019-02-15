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
	"sync"
	"sync/atomic"
	"time"
	"unsafe"

	"github.com/kpango/fastime"
)

type FUID struct {
	m1 byte
	m2 byte
	m3 byte
	t  *fastime.Fastime
}

const (
	encodedLen = 20
	rawLen     = 12
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

	buf = sync.Pool{
		New: func() interface{} {
			return make([]byte, encodedLen)
		},
	}
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
		m1: machineID[0],
		m2: machineID[1],
		m3: machineID[2],
	}
}

func String() string {
	return instance.String()
}

func (f *FUID) String() string {
	var id [rawLen]byte
	binary.BigEndian.PutUint32(id[:], uint32(f.t.UnixNanoNow()))
	id[4] = f.m1
	id[5] = f.m2
	id[6] = f.m3
	id[7] = p1
	id[8] = p2
	i := atomic.AddUint32(&objectIDCounter, 1)
	id[9] = byte(i >> 16)
	id[10] = byte(i >> 8)
	id[11] = byte(i)
	dst := buf.Get().([]byte)
	dst[0] = encoding[id[0]>>3]
	dst[1] = encoding[(id[1]>>6)&0x1F|(id[0]<<2)&0x1F]
	dst[2] = encoding[(id[1]>>1)&0x1F]
	dst[3] = encoding[(id[2]>>4)&0x1F|(id[1]<<4)&0x1F]
	dst[4] = encoding[id[3]>>7|(id[2]<<1)&0x1F]
	dst[5] = encoding[(id[3]>>2)&0x1F]
	dst[6] = encoding[id[4]>>5|(id[3]<<3)&0x1F]
	dst[7] = encoding[id[4]&0x1F]
	dst[8] = encoding[id[5]>>3]
	dst[9] = encoding[(id[6]>>6)&0x1F|(id[5]<<2)&0x1F]
	dst[10] = encoding[(id[6]>>1)&0x1F]
	dst[11] = encoding[(id[7]>>4)&0x1F|(id[6]<<4)&0x1F]
	dst[12] = encoding[id[8]>>7|(id[7]<<1)&0x1F]
	dst[13] = encoding[(id[8]>>2)&0x1F]
	dst[14] = encoding[(id[9]>>5)|(id[8]<<3)&0x1F]
	dst[15] = encoding[id[9]&0x1F]
	dst[16] = encoding[id[10]>>3]
	dst[17] = encoding[(id[11]>>6)&0x1F|(id[10]<<2)&0x1F]
	dst[18] = encoding[(id[11]>>1)&0x1F]
	dst[19] = encoding[(id[11]<<4)&0x1F]
	buf.Put(dst)
	return *(*string)(unsafe.Pointer(&dst))
}
