package socketio

import (
	"io"
)

type (
	LogFunc func(direction string, content interface{})

	SocketBuffer interface {
		Content() string
	}

	writeWrapper struct {
		io.WriteCloser
		SocketBuffer
		writerDelegate io.WriteCloser
		buf            []byte
	}

	readWrapper struct {
		io.ReadCloser
		SocketBuffer
		readerDelegate io.ReadCloser
		buf            []byte
	}
)

func wrapReader(r io.ReadCloser) (io.ReadCloser, SocketBuffer) {

	reader := &readWrapper{
		readerDelegate: r,
		buf:            []byte{},
	}

	return reader, reader
}

func (r *readWrapper) Read(p []byte) (n int, err error) {
	n, err = r.readerDelegate.Read(p)

	r.buf = append(r.buf, p[0:n]...)

	return n, err
}

func (r *readWrapper) Close() error {
	return r.readerDelegate.Close()
}

func (r *readWrapper) Content() string {
	return string(r.buf)
}

func wrapWriter(w io.WriteCloser) (io.WriteCloser, SocketBuffer) {

	wrapper := &writeWrapper{
		writerDelegate: w,
		buf:            []byte{},
	}

	return wrapper, wrapper
}

func (w *writeWrapper) Write(p []byte) (n int, err error) {

	n, err = w.writerDelegate.Write(p)

	w.buf = append(w.buf, p[0:n]...)

	return n, err
}

func (w *writeWrapper) Close() error {
	return w.writerDelegate.Close()
}

func (w *writeWrapper) Content() string {
	return string(w.buf)
}
