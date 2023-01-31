package stringutil

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/speps/go-hashids/v2"
	"github.com/spf13/cast"
	"hash/crc32"
	"math/rand"
	"time"
)

func MD5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

// GenUniqueID 唯一id 最短为 13位
func GenUniqueID(len int) string {
	hd := hashids.NewData()
	hd.Salt = RandStr(16)

	hd.MinLength = len
	h, _ := hashids.NewWithData(hd)

	e, _ := h.EncodeInt64([]int64{time.Now().UnixNano()})
	return e
}

func RandStr(n int) string {
	result := make([]byte, n/2)
	rand.Read(result)
	return hex.EncodeToString(result)
}

func CRC32(message []byte) int {
	return cast.ToInt(crc32.ChecksumIEEE(message))
}
