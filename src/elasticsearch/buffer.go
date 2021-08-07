package elasticsearch

import (
	"errors"
	"io"
)

//https://golangtc.com/t/5a49e2104ce40d740bbbc515

type buffer struct {
	reader io.Reader
	buf    []byte
	start  int
	end    int
}

func newBuffer(reader io.Reader, len int) buffer {
	buf := make([]byte, len)
	return buffer{reader, buf, 0, 0}
}

func (b *buffer) Len() int {
	return b.end - b.start
}

//将有用的字节前移
func (b *buffer) grow() {
	if b.start == 0 {
		return
	}
	copy(b.buf, b.buf[b.start:b.end])
	b.end -= b.start
	b.start = 0
}

//从reader里面读取数据，如果reader阻塞，会发生阻塞
func (b *buffer) readFromReader() (int, error) {
	b.grow()
	n, err := b.reader.Read(b.buf[b.end:])
	if err != nil {
		return n, err
	}
	b.end += n
	return n, nil
}

//返回n个字节，而不产生移位
func (b *buffer) seek(n int) ([]byte, error) {
	if b.end-b.start >= n {
		buf := b.buf[b.start : b.start+n]
		return buf, nil
	}
	return nil, errors.New("not enough")
}

//舍弃offset个字段，读取n个字段
func (b *buffer) read(offset, n int) []byte {
	b.start += offset
	buf := b.buf[b.start : b.start+n]
	b.start += n
	return buf
}
