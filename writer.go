package fsplit

import (
	"io"
	"os"
	"strconv"
)

type Writer interface {
	io.WriteCloser

	File() *os.File

	Sync() error

	Files() ([]string, error)
}

type writer struct {
	name string
	opts WriterOptions

	f  *os.File
	fi int
	n  int64
	fs []string
}

func (w *writer) rotate() (err error) {
	if w.f != nil {
		if err = w.f.Close(); err != nil {
			return
		}
		w.f = nil
	}

	name := w.name

	if w.fi > 0 {
		// on first rotation, append .1 to the file name
		if w.fi == 1 {
			if err = os.Rename(w.name, w.name+".1"); err != nil {
				return
			}
			w.fi += 1
			w.fs[0] = w.name + ".1"
		}

		name = w.name + "." + strconv.Itoa(w.fi)
	}

	if w.f, err = os.OpenFile(name, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, w.opts.Perm); err != nil {
		return
	}

	w.fi += 1
	w.fs = append(w.fs, name)

	return err
}

func (w *writer) Write(p []byte) (n int, err error) {
	if w.f == nil {
		err = os.ErrClosed
		return
	}

	if int64(len(p))+w.n > w.opts.SplitSize {
		buf1, buf2 := p[:w.opts.SplitSize-w.n], p[w.opts.SplitSize-w.n:]

		n, err = w.f.Write(buf1)
		w.n += int64(n)

		if err != nil {
			return
		}

		if err = w.rotate(); err != nil {
			return
		}

		var n2 int
		n2, err = w.f.Write(buf2)
		n += n2
		w.n = int64(n2)
	} else {
		n, err = w.f.Write(p)
		w.n += int64(n)
	}
	return
}

func (w *writer) Sync() error {
	if w.f == nil {
		return nil
	}
	return w.f.Sync()
}

func (w *writer) Close() error {
	if w.f == nil {
		return nil
	}
	f := w.f
	w.f = nil
	w.fi = 0
	return f.Close()
}

func (w *writer) File() *os.File {
	return w.f
}

func (w *writer) Files() ([]string, error) {
	if w.fi == 0 {
		return nil, os.ErrClosed
	}
	return w.fs, nil
}

type WriterOptions struct {
	Perm      os.FileMode
	SplitSize int64
}

func NewWriter(name string, opts WriterOptions) (Writer, error) {
	w := &writer{
		name: name,
		opts: opts,
	}
	return w, w.rotate()
}
