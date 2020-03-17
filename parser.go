package quickfix

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/daominah/gomicrokit/log"
	"github.com/quickfixgo/quickfix/internal"
)

const (
	defaultBufSize = 4096
)

var bufferPool internal.BufferPool

type parser struct {
	//buffer is a slice of bigBuffer
	bigBuffer, buffer []byte
	reader            io.Reader
	lastRead          time.Time
}

func newParser(reader io.Reader) *parser {
	return &parser{reader: reader}
}

const errFromIOReader = "error from IOReader"

func (p *parser) readMore() (int, error) {
	if len(p.buffer) == cap(p.buffer) {
		var newBuffer []byte
		switch {
		//initialize the parser
		case len(p.bigBuffer) == 0:
			p.bigBuffer = make([]byte, defaultBufSize)
			newBuffer = p.bigBuffer[0:0]

			//shift buffer back to the start of bigBuffer
		case 2*len(p.buffer) <= len(p.bigBuffer):
			newBuffer = p.bigBuffer[0:len(p.buffer)]

			//reallocate big buffer with enough space to shift buffer
		default:
			p.bigBuffer = make([]byte, 2*len(p.buffer))
			newBuffer = p.bigBuffer[0:len(p.buffer)]
		}

		copy(newBuffer, p.buffer)
		p.buffer = newBuffer
	}

	n, e := p.reader.Read(p.buffer[len(p.buffer):cap(p.buffer)])
	p.lastRead = time.Now()
	p.buffer = p.buffer[:len(p.buffer)+n]
	if e != nil {
		return n, fmt.Errorf("%v: %v", errFromIOReader, e)
	}
	return n, nil
}

func (p *parser) findIndex(delim []byte) (int, error) {
	return p.findIndexAfterOffset(0, delim)
}

// findIndexAfterOffset finds from the offset character
func (p *parser) findIndexAfterOffset(offset int, delim []byte) (int, error) {
	for {
		if offset > len(p.buffer) {
			if n, err := p.readMore(); n == 0 && err != nil {
				return -1, err
			}

			continue
		}

		if index := bytes.Index(p.buffer[offset:], delim); index != -1 {
			return index + offset, nil
		}

		n, err := p.readMore()

		if n == 0 && err != nil {
			return -1, err
		}
	}
}

// findStart returns index of FIXTagBeginString (8=)
func (p *parser) findStart() (int, error) {
	return p.findIndex([]byte("8="))
}

func (p *parser) findEndAfterOffset(offset int) (int, error) {
	index, err := p.findIndexAfterOffset(offset, []byte("\00110="))
	if err != nil {
		return index, err
	}

	index, err = p.findIndexAfterOffset(index+1, []byte("\001"))
	if err != nil {
		return index, err
	}

	return index + 1, nil
}

// jumpLength returns an index base on value of FIXTagBodyLength (9=).
// In case of HNX InfoGate msg (without FIXTagCheckSum), this func returns
// index of the last '\001'
func (p *parser) jumpLength() (int, error) {
	lengthIndex, err := p.findIndex([]byte("9="))
	if err != nil {
		return 0, err
	}

	lengthIndex += 3

	offset, err := p.findIndexAfterOffset(lengthIndex, []byte("\001"))
	if err != nil {
		return 0, err
	}

	if offset == lengthIndex {
		return 0, errors.New("No length given")
	}

	length, err := atoi(p.buffer[lengthIndex:offset])
	if err != nil {
		return length, err
	}

	if length <= 0 {
		return length, errors.New("Invalid length")
	}

	return offset + length, nil
}

func (p *parser) ReadMessage() (msgBytes *bytes.Buffer, err error) {
	start, err := p.findStart()
	if err != nil {
		return nil, fmt.Errorf("error when findStart: %v", err)
	}
	p.buffer = p.buffer[start:]

	index, err := p.jumpLength()
	if err != nil {
		return nil, fmt.Errorf("error when jumpLength: %v", err)
	}

	if IsHNXInfoGateProtocol { // HNX removed FIXTagCheckSum
		index += 1
	} else {
		index, err = p.findEndAfterOffset(index)
		if err != nil {
			return nil, fmt.Errorf("error when findEndAfterOffset: %v", err)
		}
	}

	if index > len(p.buffer) {
		_ = log.Debugf
		currentBuff := make([]byte, len(p.buffer))
		copy(currentBuff, p.buffer)
		//log.Debugf("error not received full FIXMsg: %s", currentBuff)
		_, ioErr := p.readMore()
		if ioErr != nil {
			return nil, ioErr
		}
		return nil, errors.New("error not received full FIXMsg")
	}

	msgBytes = bufferPool.Get()
	msgBytes.Reset()
	msgBytes.Write(p.buffer[:index])
	p.buffer = p.buffer[index:]

	return
}
