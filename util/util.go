package util

import (
	"io"
	"os"
)

func EqualBytes(x, y []byte) bool {
	if len(x) != len(y) {
		return false
	}
	for i := range x {
		if x[i] != y[i] {
			return false
		}
	}
	return true
}

func GetFileBytes(f *os.File) ([]byte, error) {
	stat, err := f.Stat()
	if err != nil {
		return nil, err
	}
	return ReadByte(f, stat.Size())
}

func ReadByte(f io.Reader, size int64) ([]byte, error) {
	var b = make([]byte, size)
	_, err := f.Read(b)
	if err != nil {
		return nil, err
	}
	// if n == 0 {
	// 	return nil, errors.New("ファイル読み込みに失敗")
	// }

	return b, nil
}
