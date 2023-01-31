package stringutil

import (
	"bytes"
	"compress/zlib"
	"io"
)

//ZlibCompress 压缩
func ZlibCompress(src []byte) []byte {
	var in bytes.Buffer
	w := zlib.NewWriter(&in)
	w.Write(src)
	w.Close()
	return in.Bytes()
}

//ZlibUnCompress 解压缩
func ZlibUnCompress(compressSrc []byte) []byte {
	var out bytes.Buffer
	b := bytes.NewReader(compressSrc)
	r, _ := zlib.NewReader(b)
	io.Copy(&out, r)
	return out.Bytes()
}
