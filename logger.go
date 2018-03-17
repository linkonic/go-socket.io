package socketio

import (
	"io"
	"bufio"
)

type (
	LogMessage func(message string, content interface{})

	SocketBuffer interface {
		Content() string
	}

	writeWrapper struct {
		io.WriteCloser
		SocketBuffer
		writerDelegate io.WriteCloser
		buf            []byte
		size           int
	}
	readWrapper struct {
		io.Reader
		SocketBuffer
		readerDelegate io.Reader
		buf            []byte
		size           int
	}
)

func wrapReader(rd io.Reader) (io.Reader, SocketBuffer) {

	reader := &readWrapper{
		readerDelegate: bufio.NewReader(rd),
		buf:            []byte{},
		size:           0,
	}

	return reader, reader
}

func wrapWriter(writer io.WriteCloser) (io.WriteCloser, SocketBuffer) {

	wrapper := &writeWrapper{
		writerDelegate: writer,
		buf:            []byte{},
		size:           0,
	}

	return wrapper, wrapper
}

func (ww *writeWrapper) Write(p []byte) (n int, err error) {

	n, err = ww.writerDelegate.Write(p)

	ww.buf = append(ww.buf, p[0:n]...)
	ww.size += n

	return n, err
}

func (l readWrapper) Read(p []byte) (n int, err error) {
	n, err = l.readerDelegate.Read(p)

	l.buf = append(l.buf, p[0:n]...)
	l.size += n

	return n, err
}

func (l writeWrapper) Content() string {
	return string(l.buf[0:l.size])
}

func (l readWrapper) Content() string {
	return string(l.buf[0:l.size])
}
