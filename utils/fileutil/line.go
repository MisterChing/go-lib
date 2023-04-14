package fileutil

import (
	"bufio"
	"io"
	"os"
)

type FileLineHelper struct {
	f         *os.File
	bufReader *bufio.Reader
	line      string
	err       error
}

func NewFileLineHelper(name string) (*FileLineHelper, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	return &FileLineHelper{
		f:         f,
		bufReader: bufio.NewReader(f),
	}, nil
}

func (lh *FileLineHelper) Next() bool {
	line, _, err := lh.bufReader.ReadLine()
	if err != nil && err != io.EOF {
		lh.err = err
		lh.line = ""
		return false
	}
	if err == io.EOF {
		lh.line = ""
		return false
	}
	lh.line = string(line)
	return true
}

func (lh *FileLineHelper) GetLine() string {
	return lh.line
}

func (lh *FileLineHelper) Error() error {
	return lh.err
}

func (lh *FileLineHelper) Close() {
	_ = lh.f.Close()
}
