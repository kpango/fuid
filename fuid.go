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
	machineID []byte
	t         *fastime.Fastime
}

const (
	encodedLen = 20 // string encoded len
	rawLen     = 12 // binary raw len
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

	objectIDCounter = func() uint32 {
		b := make([]byte, 3)
		if _, err := rand.Reader.Read(b); err != nil {
			panic(fmt.Errorf("cannot generate random number: %v;", err))
		}
		return uint32(b[0])<<16 | uint32(b[1])<<8 | uint32(b[2])
	}()
)

func New() *FUID {
	return &FUID{
		t: fastime.New().StartTimerD(context.Background(), time.Nanosecond*10),
		machineID: func() []byte {
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
					panic(fmt.Errorf("xid: cannot get hostname nor generate a random number: %v; %v", err, randErr))
				}
			}
			return id
		}(),
	}
}

func String() string {
	return instance.String()
}

func (f *FUID) String() string {
	var id [rawLen]byte
	binary.BigEndian.PutUint32(id[:], uint32(f.t.UnixNanoNow()))
	id[4] = f.machineID[0]
	id[5] = f.machineID[1]
	id[6] = f.machineID[2]
	id[7] = byte(pid >> 8)
	id[8] = byte(pid)
	i := atomic.AddUint32(&objectIDCounter, 1)
	id[9] = byte(i >> 16)
	id[10] = byte(i >> 8)
	id[11] = byte(i)
	dst := make([]byte, encodedLen)
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
	return *(*string)(unsafe.Pointer(&dst))
}
