package display

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"reflect"
	"strconv"
	"text/scanner"
	"time"
)

func Any(value interface{}) string {
	return formatAtom(reflect.ValueOf(value))
}

func formatAtom(v reflect.Value) string {
	switch v.Kind() {
	case reflect.Invalid:
		return "invalid"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(v.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(v.Uint(), 10)
	case reflect.Bool:
		return strconv.FormatBool(v.Bool())
	case reflect.String:
		return strconv.Quote(v.String())
	case reflect.Chan, reflect.Func, reflect.Ptr, reflect.Slice, reflect.Map:
		return v.Type().String() + "0x" + strconv.FormatUint(uint64(v.Pointer()), 16)
	default:
		return v.Type().String() + " value"
	}
}

//func testAny(t *testing.T) {
func main2() {
	fmt.Println("vim-go")
	var x int64 = 1
	var d time.Duration = 1 * time.Nanosecond
	/*	fmt.Println(format.Any(x))
		fmt.Println(format.Any(d))
		fmt.Println(format.Any([]int64{x}))
		fmt.Println(format.Any([]time.Duration{d}))
	*/
	fmt.Println(Any(x))
	fmt.Println(Any(d))
	fmt.Println(Any([]int64{x}))
	fmt.Println(Any([]time.Duration{d}))
	fmt.Println(Any(make(map[string]int)))
	fmt.Println(Any(make(chan []string)))
	fmt.Println("-----------------------------\n")
	Display("os.Stderr", os.Stderr)
	strangelove := Movie{
		Title:    "Dr. Strangelove",
		Subtitle: "How I Learned to Stop Worrying and Love the Bomb",
		Year:     1964,
		Color:    false,
		Actor: map[string]string{
			"Dr. StrangeLove":            "Peter Sellers",
			"Grp. Capt. Lionel Mandrake": "Peter Sellers",
			"Pres. Buck Turgidson":       "Peter Sellers",
			`Maj. T.J. "king" Kong`:      "Slim Pickens",
		},
		Oscars: []string{
			"Best Actor (Nomin.)",
			"Best Adapted Scareenplay (Nomin.)",
			"Best Director (Nomin.)",
			"Best Picture (Nomin.)",
		},
	}
	Display("strangelove", strangelove)
	Display("chan", make(chan map[string]string))
	fmt.Println("-----------------------------\n")
	encodeStrange, err := Marshal(strangelove)
	if err != nil {
		fmt.Println("xxxxxxxxxxxx")
	}
	Display("encode", encodeStrange)
	fmt.Println(encodeStrange)
	fmt.Println("\n-----------------------------\n")
	x = 22
	a := reflect.ValueOf(2)
	fmt.Printf("%v , %T\n", a, a)
	Display("a", a)
	b := reflect.ValueOf(&x)
	Display("&x", b)
	//解引用,之后可以取地址
	Display("Emel", b.Elem())
	fmt.Println(a.CanAddr())        //false
	fmt.Println(b.CanAddr())        //false
	fmt.Println(b.Elem().CanAddr()) //true
	fmt.Println("\n-----------------------------\n")
	c := reflect.ValueOf(&x).Elem()
	px := c.Addr().Interface().(*int64)
	*px = 5
	fmt.Println(x)
	fmt.Println("\n-----------------------------\n")
	//要注意变量的类型,int64不能赋值给int类型,否则panic
	c.Set(reflect.ValueOf(int64(44)))
	fmt.Println(x)
	fmt.Println("\n-----------------------------\n")
	var movie1 Movie
	err = Unmarshal(encodeStrange, &movie1)
	Display("movie", movie1)
	fmt.Println("\n-----------------------------\n")
	showByte2ascii(encodeStrange)
	fmt.Println("\n-----------------------------\n")
	e := reflect.ValueOf(x)
	fmt.Println(e)
	fmt.Println(e.Kind())
	Display("e.Kind", e.Kind())
	Display("e", e)

	fmt.Println("\n-----------------------------\n")
	var v *int16
	Display("v", v)
	var i interface{}
	Display("i", i)
	i = v
	Display("i", i)
	var xx = 192
	Display("xx", xx)
	var xxx int
	Display("xxx", xxx)
	fmt.Println("\n-----------------------------\n")
	testError()
}

func showByte2ascii(bytes []byte) {
	if 0 == len(bytes) {
		return
	}
	for _, b := range bytes {
		fmt.Printf("%c", b)
	}
}

func display(path string, v reflect.Value) {
	switch v.Kind() {
	case reflect.Invalid:
		fmt.Printf("%s = invalid\n", path)
	case reflect.Slice, reflect.Array:
		for i := 0; i < v.Len(); i++ {
			display(fmt.Sprintf("%s[%d]", path, i), v.Index(i))
		}
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			fieldPath := fmt.Sprintf("%s.%s", path, v.Type().Field(i).Name)
			display(fieldPath, v.Field(i))
		}
	case reflect.Map:
		for _, key := range v.MapKeys() {
			display(fmt.Sprintf("%s[%s]", path, formatAtom(key)), v.MapIndex(key))
		}
	case reflect.Ptr:
		if v.IsNil() {
			fmt.Printf("%s = nil\n", path)
		} else {
			display(fmt.Sprintf("(*%s)", path), v.Elem())
		}
	case reflect.Interface:
		if v.IsNil() {
			fmt.Printf("%s = nil\n", path)
		} else {
			fmt.Printf("%s.type = %s\n", path, v.Elem().Type())
			display(path+".value", v.Elem())
		}
	default:
		fmt.Printf("%s = %s\n", path, formatAtom(v))
	}
}
func fdisplay(w io.Writer, path string, v reflect.Value) {
	switch v.Kind() {
	case reflect.Invalid:
		fmt.Fprintf(w, "%s = invalid\n", path)
	case reflect.Slice, reflect.Array:
		for i := 0; i < v.Len(); i++ {
			fdisplay(w, fmt.Sprintf("%s[%d]", path, i), v.Index(i))
		}
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			fieldPath := fmt.Sprintf("%s.%s", path, v.Type().Field(i).Name)
			fdisplay(w, fieldPath, v.Field(i))
		}
	case reflect.Map:
		for _, key := range v.MapKeys() {
			fdisplay(w, fmt.Sprintf("%s[%s]", path, formatAtom(key)), v.MapIndex(key))
		}
	case reflect.Ptr:
		if v.IsNil() {
			fmt.Fprintf(w, "%s = nil\n", path)
		} else {
			fdisplay(w, fmt.Sprintf("(*%s)", path), v.Elem())
		}
	case reflect.Interface:
		if v.IsNil() {
			fmt.Fprintf(w, "%s = nil\n", path)
		} else {
			fmt.Fprintf(w, "%s.type = %s\n", path, v.Elem().Type())
			fdisplay(w, path+".value", v.Elem())
		}
	default:
		fmt.Fprintf(w, "%s = %s\n", path, formatAtom(v))
	}
}

func Fdisplay(w io.Writer, name string, x interface{}) {
	fmt.Fprintf(w, "\nDisplay %s (%T) [%s]:\n", name, x, time.Now().Format("2006-01-02 15:03:02"))
	fdisplay(w, name, reflect.ValueOf(x))
}
func Display(name string, x interface{}) {
	fmt.Printf("Display %s (%T):\n", name, x)
	//	fmt.Println(reflect.ValueOf(x))
	//	fmt.Println(reflect.ValueOf(x).Type())
	display(name, reflect.ValueOf(x))
}

type Movie struct {
	Title, Subtitle string
	Year            int
	Color           bool
	Actor           map[string]string
	Oscars          []string
	Sequel          *string
}

func encode(buf *bytes.Buffer, v reflect.Value) error {
	//	Display("current", v.Kind())
	switch v.Kind() {
	case reflect.Invalid:
		buf.WriteString("nil")
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		fmt.Fprintf(buf, "%d", v.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		fmt.Fprintf(buf, "%d", v.Uint())
	case reflect.String:
		fmt.Fprintf(buf, "%q", v.String())
	case reflect.Ptr:
		return encode(buf, v.Elem())
	case reflect.Bool:
		/*			fmt.Println(v.Bool())
					fmt.Printf("%v , %s", v.Bool(), v.String())
					fmt.Println("\n--------")
					os.Exit(1)
		*/
		fmt.Fprintf(buf, "%v", v.Bool())
	case reflect.Array, reflect.Slice:
		buf.WriteByte('(')
		for i := 0; i < v.Len(); i++ {
			if i > 0 {
				buf.WriteByte(' ')
			}
			if err := encode(buf, v.Index(i)); err != nil {
				return err
			}
		}
		buf.WriteByte(')')
	case reflect.Struct:
		//			fmt.Printf("struct...%d fields \n", v.NumField())
		buf.WriteByte('(')
		for i := 0; i < v.NumField(); i++ {
			if i > 0 {
				buf.WriteByte(' ')
			}
			fmt.Fprintf(buf, "(%s ", v.Type().Field(i).Name)
			if err := encode(buf, v.Field(i)); err != nil {
				return err
			}
			buf.WriteByte(')')
		}
		buf.WriteByte(')')
	case reflect.Map:
		buf.WriteByte('(')
		for i, key := range v.MapKeys() {
			if i > 0 {
				buf.WriteByte(' ')
			}
			buf.WriteByte('(')
			if err := encode(buf, key); err != nil {
				return err
			}
			buf.WriteByte(' ')
			if err := encode(buf, v.MapIndex(key)); err != nil {
				return err
			}
			buf.WriteByte(')')
		}
		buf.WriteByte(')')
	default:
		return fmt.Errorf("unsupported type: %s", v.Type())
	}
	return nil
}

func Marshal(v interface{}) ([]byte, error) {
	var buf bytes.Buffer
	if err := encode(&buf, reflect.ValueOf(v)); err == nil {
		return nil, err
	}
	//	fmt.Println("-----------------------------\n")
	//	Display("buf", buf)
	//	fmt.Println("-----------------------------\n")
	return buf.Bytes(), nil
}

type lexer struct {
	scan  scanner.Scanner
	token rune
}

func (lex *lexer) next() {
	lex.token = lex.scan.Scan()
}

func (lex *lexer) text() string {
	return lex.scan.TokenText()
}

func (lex *lexer) consumer(want rune) {
	if lex.token != want {
		panic(fmt.Sprintf("got %q, want %q", lex.text(), want))
	}
	lex.next()
}

func Unmarshal(data []byte, out interface{}) (err error) {
	lex := &lexer{scan: scanner.Scanner{Mode: scanner.GoTokens}}
	lex.scan.Init(bytes.NewReader(data))
	lex.next()
	defer func() {
		if x := recover(); x != nil {
			err = fmt.Errorf("error at %s: %v", lex.scan.Position, x)
		}
	}()
	read(lex, reflect.ValueOf(out).Elem())
	return nil
}

func read(lex *lexer, v reflect.Value) {
	switch lex.token {
	case scanner.Ident:
		if "nil" == lex.text() {
			v.Set(reflect.Zero(v.Type()))
			lex.next()
			return
		}
	case scanner.String:
		s, _ := strconv.Unquote(lex.text())
		v.SetString(s)
		lex.next()
		return
	case scanner.Int:
		i, _ := strconv.Atoi(lex.text())
		v.SetInt(int64(i))
		lex.next()
		return
	case '(':
		lex.next()
		readList(lex, v)
		lex.next() //consume ')
		return
	}
	panic(fmt.Sprintf("unexcepted token %q", lex.text()))
}

func readList(lex *lexer, v reflect.Value) {
	switch v.Kind() {
	case reflect.Array: //(item ...)
		for i := 0; !endList(lex); i++ {
			read(lex, v.Index(i))
		}
	case reflect.Slice: //(item ...)
		for !endList(lex) {
			item := reflect.New(v.Type().Elem()).Elem()
			read(lex, item)
			v.Set(reflect.Append(v, item))
		}
	case reflect.Struct: // ((key value) ...)
		for !endList(lex) {
			lex.consumer('(')
			if lex.token != scanner.Ident {
				panic(fmt.Sprintf("got token %q, want field name", lex.text()))
			}
			name := lex.text()
			lex.next()
			read(lex, v.FieldByName(name))
			lex.consumer(')')
		}
	case reflect.Map:
		v.Set(reflect.MakeMap(v.Type()))
		for !endList(lex) {
			lex.consumer('(')
			key := reflect.New(v.Type().Key()).Elem()
			read(lex, key)
			value := reflect.New(v.Type().Elem()).Elem()
			read(lex, value)
			v.SetMapIndex(key, value)
			lex.consumer(')')
		}
	default:
		panic(fmt.Sprintf("cannot decode list into %v", v.Type()))
	}
}

func endList(lex *lexer) bool {
	switch lex.token {
	case scanner.EOF:
		panic("end of file")
	case ')':
		return true
	}
	return false
}

//Error test
type Error struct {
	errCode uint8
}

func (e *Error) Error() string {
	switch e.errCode {
	case 1:
		return "file not found"
	case 2:
		return "time out"
	case 3:
		return "permission denied"
	default:
		return "unknown error"
	}
}

func checkError(err error) {
	Display("checkError", err)
	//err的指针是nil,但是其并不是nil
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}

func testError() {
	//此处的err并不是真正的nil,其type为Error
	var e *Error
	Display("err", e)
	fmt.Printf("%p , %p\n", e, nil)
	fmt.Println(e)
	fmt.Println(nil)
	checkError(e)
}
