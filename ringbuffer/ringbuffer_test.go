package ringbuffer

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/docker/docker/daemon/logger"
)

var (
	COUNT = 1
	MSG   = &logger.Message{
		Line: []byte(`1234567890abcdefghijklmopqrstuvwxyzABCDEFGHIJKLMOPQRSTUVWXYZ!@#$%^&*()`),
	}
)

func TestSendReadAll10000(t *testing.T) {
	COUNT = 10000
	rb := New(COUNT)
	defer rb.Close()
	for i := 0; i < COUNT+100; i++ {
		rb.Write(MSG)
	}
	log.Println(`test msg len:`, len(MSG.Line)*COUNT)
	time.Sleep(time.Second * 2)
	bf := rb.Tail(0, time.Time{})
	length := 0
	for _, b := range bf {
		length += len(b.Line)
	}
	if len(MSG.Line)*COUNT != length {
		t.Errorf(`recv incorrected buffer len %v want %v`, length, len(MSG.Line)*COUNT)
	}
}

func TestSendTail50(t *testing.T) {
	COUNT = 10000
	n := 50
	rb := New(COUNT)
	defer rb.Close()
	testbuffer := ""
	for i := 0; i < COUNT; i++ {
		msg := &logger.Message{
			Line: []byte(fmt.Sprintf("%v", i)),
		}

		rb.Write(msg)
		if i >= COUNT-n {
			testbuffer += fmt.Sprintf("%v", i)
		}
	}
	log.Println(`test msg len:`, len(MSG.Line)*COUNT)
	time.Sleep(time.Second * 2)
	bfs := rb.Tail(n, time.Time{})

	resbuffer := ""
	for i, bf := range bfs {
		if bf == nil {
			log.Println(`err in`, i)
			continue
		}
		resbuffer += string(bf.Line)
	}
	if len(resbuffer) != len(testbuffer) {
		t.Errorf(`recv incorrected buffer len %v want %v`, len(resbuffer), len(testbuffer))
	}
	if resbuffer != testbuffer {
		t.Errorf(`recv incorrected buffer %v want %v`, resbuffer, testbuffer)
	}
}

func TestSendTailSince(t *testing.T) {
	COUNT = 10000
	n := 50
	rb := New(COUNT)
	defer rb.Close()
	testbuffer := ""
	for i := 0; i < COUNT; i++ {
		msg := &logger.Message{
			Line:      []byte(fmt.Sprintf("%v", i)),
			Timestamp: time.Unix(int64(i), 0),
		}

		rb.Write(msg)
		if i >= COUNT-n {
			testbuffer += fmt.Sprintf("%v", i)
		}
	}
	log.Println(`test msg len:`, len(MSG.Line)*COUNT)
	time.Sleep(time.Second * 2)
	bfs := rb.Tail(0, time.Unix(int64(COUNT-n), 0))

	resbuffer := ""
	for i, bf := range bfs {
		if bf == nil {
			log.Println(`err in`, i)
			continue
		}
		resbuffer += string(bf.Line)
	}
	if len(resbuffer) != len(testbuffer) {
		t.Errorf(`recv incorrected buffer len %v want %v`, len(resbuffer), len(testbuffer))
	}
	if resbuffer != testbuffer {
		t.Errorf(`recv incorrected buffer %v want %v`, resbuffer, testbuffer)
	}
}
