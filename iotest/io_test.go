// package io一些注意点的测试
//
// Read()阻塞与否由Read()实现者决定
// e.g.
// net.Conn 会阻塞
// os.File 不会阻塞
//
// note that
// ReadFull(), ReadAtLeast(), ReadAll()调用Read()，其行为由Read()的实现者决定

package iotest

import (
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Read()每次都是从头覆盖[]byte
// 复用时[]byte时，需注意[]byte内容
// 因此出现了ReadFull()和ReadAtLeast()
func TestRead(t *testing.T) {
	// abcde
	r := strings.NewReader("abcde")
	b := make([]byte, 3)

	// first read 3 bytes
	n, err := r.Read(b)
	if err != nil {
		panic(err)
	}
	assert.Equal(t, n, 3)

	// second read last 2 bytes
	m, err := r.Read(b)
	if err != nil {
		panic(err)
	}
	assert.Equal(t, m, 2)

	// assert b = []byte{'d', 'e', 'c'}
	assert.EqualValues(t, b, []byte{'d', 'e', 'c'})
}

// 每次把[]byte填满后返回，或遇到EOF
// 未读取到任何数据，返回err=io.EOF
// 读取到一些数据后遇到EOF，返回err=io.ErrUnexpectedEOF
// 把[]byte读满后，返回err=nil
func TestReadFull(t *testing.T) {
	// abced
	r := strings.NewReader("abcde")
	b := make([]byte, 3)

	// first read 3 bytes
	n, err := io.ReadFull(r, b)
	if err != nil {
		panic(err)
	}
	assert.Equal(t, n, 3)

	// second read last 2 bytes
	m, err := io.ReadFull(r, b)
	assert.EqualError(t, err, io.ErrUnexpectedEOF.Error())
	assert.Equal(t, m, 2)

	// assert b = []byte{'d', 'e', 'c'}
	assert.EqualValues(t, b, []byte{'d', 'e', 'c'})
}

// 每次至少读取len([]byte)个字节后返回，或遇到EOF
// 未读取到任何数据，返回err=io.EOF
// 读取到一些数据后遇到EOF，返回err=io.ErrUnexpectedEOF
// 读取到不少于len([]byte)个字节后，返回err=nil
func TestReadAtLeast(t *testing.T) {
	// abced
	r := strings.NewReader("abcde")
	b := make([]byte, 4)

	// first read at least 3 bytes, actually, 4 bytes
	n, err := io.ReadAtLeast(r, b, 3)
	if err != nil {
		panic(err)
	}
	assert.Equal(t, n, 4)

	// second read at least 3 bytes, actually, 1 byte
	m, err := io.ReadAtLeast(r, b, 3)
	assert.EqualError(t, err, io.ErrUnexpectedEOF.Error())
	assert.Equal(t, m, 1)

	// assert b = []byte{'e', 'b', 'c', 'd'}
	assert.EqualValues(t, b, []byte{'e', 'b', 'c', 'd'})
}

// 初次容量为512，如果数据量过大，会导致多次扩容，影响性能
// 可以尝试io.Copy() + bytes.Buffer [+ sync.Pool] 优化，可见httpdemo
func TestReadAll(t *testing.T) {
	// abced
	r := strings.NewReader("abcde")

	// read all abced
	b, err := io.ReadAll(r)
	if err != nil {
		panic(err)
	}

	// assert b = []byte{'a', 'b', 'c', 'd', 'e'}
	assert.EqualValues(t, b, []byte{'a', 'b', 'c', 'd', 'e'})
}
