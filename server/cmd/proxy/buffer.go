package main

import (
	"bufio"
	"io"
	"net"
)

// ReadAll is a fork of the Go Standard Library's io.ReadAll which
// also explicitly returns the length of the buffer that was read.
func ReadAll(r io.Reader) ([]byte, int, error) {
	b := make([]byte, 0, 512)

	for {
		if len(b) == cap(b) {
			// Add more capacity (let append pick how much).
			b = append(b, 0)[:len(b)]
		}

		n, err := r.Read(b[len(b):cap(b)])
		b = b[:len(b)+n]

		if err != nil {
			if err == io.EOF {
				err = nil
			}

			return b, len(b) + n, err
		}
	}
}

// Buffer is a (maybe unnecessary) implementation of a safe buffer.
// The length of a byte slice will always be equal to the number of
// bytes read from a net.Conn.
type Buffer struct {
	buffer []byte
	length int
}

// NewBuffer creates a new Buffer object.
func NewBuffer() *Buffer {
	return &Buffer{
		buffer: make([]byte, ReadBufferSize), length: 0,
	}
}

// NewBufferFrom creates a new Buffer object from an existing buffer
// and length.
func NewBufferFrom(b []byte, n int) *Buffer {
	return &Buffer{buffer: b, length: n}
}

// Buffer returns a byte slice with a safe length.
func (b *Buffer) Buffer() []byte {
	return b.buffer[:b.length]
}

// Size returns the size of the internal buffer.
func (b *Buffer) Size() int {
	return b.length
}

// Recv reads from a connection and saves the buffer and bytes read.
func (b *Buffer) Recv(conn net.Conn) (err error) {
	b.length, err = conn.Read(b.buffer)
	return err
}

// Recvall reads from a connection until an EOF is read.
func (b *Buffer) Recvall(conn net.Conn) (err error) {
	b.buffer, b.length, err = ReadAll(bufio.NewReader(conn))
	return err
}

// Send writes a the internal buffer to a connection.
func (b *Buffer) Send(conn net.Conn) (err error) {
	// Proxy the request to its destination.
	_, err = conn.Write(b.buffer[:b.length])
	return err
}
