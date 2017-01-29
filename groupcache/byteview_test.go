package groupcache

import (
	"fmt"
	"io"
	"io/ioutil"
	"testing"
)

func TestByteView(t *testing.T) {
	for _, s := range []string{"", "x", "yy"} {
		for _, v := range []ByteView{of([]byte(s)), of(s)} {
			name := fmt.Sprintf("string %q, view %+v", s, v)
			fmt.Println(name)
			if v.Len() != len(s) {
				t.Errorf("%s: Len = %d; want %d", name, v.Len(), len(s))
			}
			if v.String() != s {
				t.Errorf("%s: String = %q; want %q", name, v.String(), s)
			}
			var longDest [3]byte
			if n := v.Copy(longDest[:]); n != len(s) {
				t.Errorf("%s: long Copy = %d; want %d", name, n, len(s))
			}
			fmt.Println(longDest)
			var shortDest [1]byte
			if n := v.Copy(shortDest[:]); n != min(len(s), 1) {
				t.Errorf("%s: long Copy = %d; want %d", name, n, min(len(s), 1))
			}
			fmt.Println(shortDest)
			got, err := ioutil.ReadAll(v.Reader())
			fmt.Println(got)
			if err != nil || string(got) != s {
				t.Errorf("%s: Reader = %q, %v; want %q", name, got, err, s)
			}
			got, err = ioutil.ReadAll(io.NewSectionReader(v, 0, int64(len(s))))
			if err != nil || string(got) != s {
				t.Errorf("%s: SectionReader of ReaderAt = %q, %v; want %q", name, got, err, s)
			}
			fmt.Println(got)
		}
	}
}

func of(x interface{}) ByteView {
	if bytes, ok := x.([]byte); ok {
		return ByteView{b: bytes}
	}
	return ByteView{s: x.(string)}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func TestByteViewEqual(t *testing.T) {
	tests := []struct {
		a    interface{} //string or []byte
		b    interface{}
		want bool
	}{
		{"x", "x", true},
		{"x", "y", false},
		{"x", "yy", false},
		{[]byte("x"), []byte("x"), true},
		{[]byte("x"), []byte("y"), false},
		{[]byte("x"), []byte("yy"), false},
		{[]byte("x"), "x", true},
		{[]byte("x"), "y", false},
		{[]byte("x"), "yy", false},
		{"x", []byte("x"), true},
		{"x", []byte("y"), false},
		{"x", []byte("yy"), false},
	}
	for i, tt := range tests {
		va := of(tt.a)
		if bytes, ok := tt.b.([]byte); ok {
			if got := va.EqualBytes(bytes); got != tt.want {
				t.Errorf("%d, EqualBytes = %v; want %v", i, got, tt.want)
			}
		} else {
			if got := va.EqualString(tt.b.(string)); got != tt.want {
				t.Errorf("%d. EqualString = %v; want %v", i, got, tt.want)
			}
		}
		if got := va.Equal(of(tt.b)); got != tt.want {
			t.Errorf("%d. Equal = %v; want %v", i, got, tt.want)
		}
	}
}

func TestByteViewSlice(t *testing.T) {
	tests := []struct {
		in   string
		from int
		to   interface{} //nil表示末尾
		want string
	}{
		{
			in:   "abc",
			from: 1,
			to:   2,
			want: "b",
		},
		{
			in:   "abc",
			from: 1,
			want: "bc",
		},
		{
			in:   "abc",
			to:   2,
			want: "ab",
		},
	}
	for i, tt := range tests {
		for _, v := range []ByteView{of([]byte(tt.in)), of(tt.in)} {
			name := fmt.Sprintf("test %d, view %+v", i, v)
			fmt.Println(name)
			if tt.to != nil {
				v = v.Slice(tt.from, tt.to.(int))
			} else {
				v = v.SliceFrom(tt.from)
			}
			fmt.Println(v)
			if v.String() != tt.want {
				t.Errorf("%s: got %q; want %q", name, v.String(), tt.want)
			}
		}
	}
}
