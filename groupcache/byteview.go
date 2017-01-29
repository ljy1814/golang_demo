package groupcache

import (
	"bytes"
	"errors"
	"io"
	"strings"
)

type ByteView struct {
	b []byte
	s string
}

func (v ByteView) Len() int {
	if v.b != nil {
		return len(v.b)
	}
	return len(v.s)
}

func (v ByteView) ByteSlice() []byte {
	if v.b != nil {
		return cloneBytes(v.b)
	}
	return []byte(v.s)
}

func (v ByteView) String() string {
	//优先返回[]byte
	if v.b != nil {
		return string(v.b)
	}
	return v.s
}

func (v ByteView) At(i int) byte {
	if v.b != nil {
		return v.b[i]
	}
	return v.s[i]
}

func (v ByteView) Slice(from, to int) ByteView {
	if v.b != nil {
		return ByteView{b: v.b[from:to]}
	}
	return ByteView{s: v.s[from:to]}
}
func (v ByteView) SliceFrom(from int) ByteView {
	if v.b != nil {
		return ByteView{b: v.b[from:]}
	}
	return ByteView{s: v.s[from:]}
}

func (v ByteView) Copy(dest []byte) int {
	if v.b != nil {
		return copy(dest, v.b)
	}
	return copy(dest, v.s)
}

func (v ByteView) Equal(bv ByteView) bool {
	if bv.b == nil {
		return v.EqualString(bv.s)
	}
	return v.EqualBytes(bv.b)
}

func (v ByteView) EqualString(s string) bool {
	if v.b == nil {
		return v.s == s
	}
	l := v.Len()
	if len(s) != l {
		return false
	}
	for i, bi := range v.b {
		if bi != s[i] {
			return false
		}
	}
	return true
}

func (v ByteView) EqualBytes(bv []byte) bool {
	if v.b != nil {
		return bytes.Equal(v.b, bv)
	}
	l := v.Len()
	if len(bv) != l {
		return false
	}
	for i, bi := range bv {
		if bi != v.s[i] {
			return false
		}
	}
	return true
}

func (v ByteView) Reader() io.ReadSeeker {
	if v.b != nil {
		return bytes.NewReader(v.b)
	}
	return strings.NewReader(v.s)
}

func (v ByteView) ReadAt(p []byte, off int64) (n int, err error) {
	if off < 0 {
		return 0, errors.New("view: invalid offset")
	}
	if off >= int64(v.Len()) {
		return 0, io.EOF
	}
	n = v.SliceFrom(int(off)).Copy(p)
	if n < len(p) {
		err = io.EOF
	}
	return
}
