package json

import (
    "bufio"
    "github.com/valyala/fastjson"
    "os"
    "testing"
)

func TestNewJson(t *testing.T) {
    {
        j := NewJson(None)
        if !j.IsNull() {
            t.Fail()
            t.Log()
        }
    }
    {
        j := NewJson(1)
        if !j.IsInt() {
            t.Fail()
            t.Log()
        }
    }
    {
        var a int64 = 1
        j := NewJson(a)
        if !j.IsInt() {
            t.Fail()
            t.Log()
        }
    }
    {
        j := NewJson(1.1)
        if !j.IsFloat() {
            t.Fail()
            t.Log()
        }
    }
    {
        j := NewJson([]any{1, 2, 3})
        if !j.IsArray() {
            t.Fail()
            t.Log()
        }
    }
    {
        j := NewJson(false)
        if !j.IsBool() {
            t.Fail()
            t.Log()
        }
    }
    {
        j := NewJson(true)
        if !j.IsBool() {
            t.Fail()
            t.Log()
        }
    }
    {
        j := NewJson("string")
        if !j.IsString() {
            t.Fail()
            t.Log()
        }
    }
    {
        m := make(map[string]any, 0)
        m["a"] = 1
        m["b"] = "2"
        m["c"] = false
        m["d"] = true
        j := NewJson(m)
        if !j.IsObject() {
            t.Fail()
            t.Log()
        }
    }
    {
        a := NewJson(true)
        j := NewJson(a)
        if !j.IsBool() {
            t.Fail()
            t.Log()
        }
    }
}

func TestLoads(t *testing.T) {
    {
        s := `{"a":"A","b":1,"c":1.1,"d":null,"e":[],"f":{},"g":null}`
        j := Loads(s)
        if !j.IsObject() {
            t.Fail()
            t.Log()
        }
        if Dumps(j) != s {
            t.Fail()
            t.Log()
        }
    }
    {
        s := `-1.23456`
        j := Loads(s)
        if j == nil {
            t.Fail()
            t.Log()
        }
    }
    {
        s := `1.23456`
        j := Loads(s)
        if j == nil {
            t.Fail()
            t.Log()
        }
    }
    {
        s := `1.2.3`
        j := Loads(s)
        if j != nil {
            t.Fail()
            t.Log()
        }
    }
    {
        s := `1.2k`
        j := Loads(s)
        if j != nil {
            t.Fail()
            t.Log()
        }
    }
    {
        s := `{"a": 1.2.3}`
        j := Loads(s)
        if j != nil {
            t.Fail()
            t.Log()
        }
    }
}

func TestDumps(t *testing.T) {
    {
        s := `1234567`
        j := Loads(s)
        if Dumps(j) != s {
            t.Fail()
            t.Log()
        }
    }
    {
        s := `-1234567`
        j := Loads(s)
        if Dumps(j) != s {
            t.Fail()
            t.Log()
        }
    }
    {
        s := `1.0`
        j := Loads(s)
        if Dumps(j) != `1` {
            t.Fail()
            t.Log()
        }
    }
    {
        s := `1.23456`
        j := Loads(s)
        if Dumps(j) != s {
            t.Fail()
            t.Log()
        }
    }
    {
        s := `true`
        j := Loads(s)
        if Dumps(j) != s {
            t.Fail()
            t.Log()
        }
    }
    {
        s := `false`
        j := Loads(s)
        if Dumps(j) != s {
            t.Fail()
            t.Log()
        }
    }
    {
        s := `null`
        j := Loads(s)
        if Dumps(j) != s {
            t.Fail()
            t.Log()
        }
    }
    {
        s := `[1,2,3,4]`
        j := Loads(s)
        if Dumps(j) != s {
            t.Fail()
            t.Log()
        }
    }
}

func TestEqual(t *testing.T) {
    {
        j := NewJson(1)
        i := NewJson(1)
        if !j.Equal(i) {
            t.Fail()
            t.Log()
        }
    }
    {
        j := NewJson(1.0)
        i := NewJson(1.0)
        if !j.Equal(i) {
            t.Fail()
            t.Log()
        }
    }
    {
        j := NewJson("1.0")
        i := NewJson("1.0")
        if !j.Equal(i) {
            t.Fail()
            t.Log()
        }
    }
    {
        j := NewJson(1.0)
        i := NewJson(1)
        if j.Equal(i) {
            t.Fail()
        }
    }
    {
        j := NewJson(true)
        i := NewJson(true)
        if !j.Equal(i) {
            t.Fail()
        }
    }
    {
        j := NewJson(None)
        i := NewJson(1)
        if j.Equal(i) {
            t.Fail()
        }
    }
    {
        j := NewJson(None)
        i := NewJson(None)
        if !j.Equal(i) {
            t.Fail()
        }
    }
    {
        j := NewArray()
        j.Append(NewJson(1))
        j.Append(NewJson("1"))
        i := NewArray()
        i.Append(NewJson(1))
        i.Append(NewJson("1"))
        if !j.Equal(i) {
            t.Fail()
            t.Log()
        }
    }
    {
        j := NewArray()
        j.Append(NewJson(1))
        j.Append(NewJson("1"))
        i := NewArray()
        i.Append(NewJson(1))
        i.Append(NewJson(1))
        if j.Equal(i) {
            t.Fail()
            t.Log()
        }
    }
    {
        j := NewArray()
        j.Append(NewJson(1))
        i := NewArray()
        i.Append(NewJson(1))
        i.Append(NewJson("1"))
        if j.Equal(i) {
            t.Fail()
            t.Log()
        }
    }
    {
        j := NewObject()
        j.Set("a", NewJson("A"))
        j.Set("b", NewJson("B"))
        i := NewObject()
        i.Set("a", NewJson("A"))
        i.Set("b", NewJson("B"))
        if !j.Equal(i) {
            t.Fail()
            t.Log()
        }
    }
    {
        j := NewObject()
        j.Set("a", NewJson("A"))
        j.Set("b", NewJson("B"))
        i := NewObject()
        i.Set("a", NewJson("A"))
        i.Set("b", NewJson("B"))
        i.Set("c", NewJson("C"))
        if j.Equal(i) {
            t.Fail()
            t.Log()
        }
    }
    {
        j := NewObject()
        j.Set("a", NewJson("A"))
        j.Set("b", NewJson("B"))
        i := NewObject()
        i.Set("a", NewJson("A"))
        i.Set("c", NewJson("C"))
        if j.Equal(i) {
            t.Fail()
            t.Log()
        }
    }
    {
        j := NewObject()
        j.Set("a", NewJson("a"))
        j.Set("b", NewJson("a"))
        i := NewObject()
        i.Set("a", NewJson("A"))
        i.Set("b", NewJson("B"))
        if j.Equal(i) {
            t.Fail()
            t.Log()
        }
    }
}

func TestSetGet(t *testing.T) {
    {
        j := NewObject()
        j.Set("a", NewJson("a"))
        j.Set("b", NewJson("b"))
        j.Set("c", NewJson("c"))
        if Dumps(j.Get("a")) != `"a"` {
            t.Fail()
            t.Log()
        }
        if !j.Get("d").IsNull() {
            t.Fail()
            t.Log()
        }
    }
    {
        arr := NewArray()
        arr.Append(NewJson(1))
        if !arr.Get(0).Equal(NewJson(1)) {
            t.Fail()
            t.Log()
        }
        arr.Set(0, NewJson(0))
        if !arr.Get(0).Equal(NewJson(0)) {
            t.Fail()
            t.Log()
        }
    }
    {
        arr := NewArray()
        arr.Append(NewJson(1))
        if !arr.Get(0).Equal(NewJson(1)) {
            t.Fail()
            t.Log()
        }
        var index int8 = 0
        arr.Set(index, NewJson(0))
        if !arr.Get(0).Equal(NewJson(0)) {
            t.Fail()
            t.Log()
        }
    }
    {
        arr := NewArray()
        arr.Append(NewJson(1))
        if !arr.Get(0).Equal(NewJson(1)) {
            t.Fail()
            t.Log()
        }
        var index int16 = 0
        arr.Set(index, NewJson(0))
        if !arr.Get(0).Equal(NewJson(0)) {
            t.Fail()
            t.Log()
        }
    }
    {
        arr := NewArray()
        arr.Append(NewJson(1))
        if !arr.Get(0).Equal(NewJson(1)) {
            t.Fail()
            t.Log()
        }
        var index int32 = 0
        arr.Set(index, NewJson(0))
        if !arr.Get(0).Equal(NewJson(0)) {
            t.Fail()
            t.Log()
        }
    }
    {
        arr := NewArray()
        arr.Append(NewJson(1))
        if !arr.Get(0).Equal(NewJson(1)) {
            t.Fail()
            t.Log()
        }
        var index int64 = 0
        arr.Set(index, NewJson(0))
        if !arr.Get(0).Equal(NewJson(0)) {
            t.Fail()
            t.Log()
        }
    }
}

func TestString(t *testing.T) {
    {
        j := NewJson(`"123"`)
        if !j.IsString() {
            t.Fail()
            t.Log()
        }
        v := j.String()
        if v != `"123"` {
            t.Fail()
            t.Log()
        }
    }
}

func TestInt(t *testing.T) {
    {
        j := NewJson(123)
        if !j.IsInt() {
            t.Fail()
            t.Log()
        }
        v := j.Int()
        if v != 123 {
            t.Fail()
            t.Log()
        }
    }
}

func TestFloat(t *testing.T) {
    {
        j := NewJson(1.23)
        if !j.IsFloat() {
            t.Fail()
            t.Log()
        }
        v := j.Float()
        if v != 1.23 {
            t.Fail()
            t.Log()
        }
    }
}

func TestBool(t *testing.T) {
    {
        j := NewJson(true)
        if !j.IsBool() {
            t.Fail()
            t.Log()
        }
        v := j.Bool()
        if v != true {
            t.Fail()
            t.Log()
        }
    }
}

func TestArray(t *testing.T) {
    {
        j := NewArray()
        j.Append(NewJson(1))
        j.Append(NewJson(2))
        if !j.IsArray() {
            t.Fail()
            t.Log()
        }
        v := j.Array()
        if len(v) != 2 {
            t.Fail()
            t.Log()
        }
    }
}

func TestObject(t *testing.T) {
    {
        j := NewObject()
        j.Set("a", NewJson("a"))
        j.Set("b", NewJson("b"))
        if !j.IsObject() {
            t.Fail()
            t.Log()
        }
        v := j.Object()
        if len(v) != 2 {
            t.Fail()
            t.Log()
        }
    }
}

func TestSetDefault(t *testing.T) {
    {
        j := NewObject()
        j.SetDefault("a", NewJson("b"))
        j.SetDefault("a", NewJson("c"))
        if j.Get("a").String() != "b" {
            t.Fail()
            t.Log()
        }
    }
}

func TestKeys(t *testing.T) {
    {
        j := NewObject()
        j.SetDefault("a", NewJson("b"))
        j.SetDefault("b", NewJson("c"))
        keys := j.Keys()
        if len(keys) != 2 {
            t.Fail()
            t.Log()
        }
    }
}

func TestValues(t *testing.T) {
    {
        j := NewObject()
        j.SetDefault("a", NewJson("b"))
        j.SetDefault("b", NewJson("c"))
        values := j.Values()
        if len(values) != 2 {
            t.Fail()
            t.Log()
        }
    }
}

func TestContains(t *testing.T) {
    {
        j := NewObject()
        j.SetDefault("a", NewJson("b"))
        j.SetDefault("b", NewJson("c"))
        if j.Contains("a") == false {
            t.Fail()
            t.Log()
        }
        if j.Contains("c") != false {
            t.Fail()
            t.Log()
        }
    }
    {
        arr := NewArray()
        arr.Append(NewJson(1))
        if !arr.Contains(NewJson(1)) {
            t.Fail()
            t.Log()
        }
    }
}

func TestClear(t *testing.T) {
    {
        j := NewObject()
        j.SetDefault("a", NewJson("b"))
        j.SetDefault("b", NewJson("c"))
        j.Clear()
        if j.Contains("a") != false {
            t.Fail()
            t.Log()
        }
        if j.Contains("b") != false {
            t.Fail()
            t.Log()
        }
    }
    {
        arr := NewArray()
        arr.Append(NewJson(1))
        if len(arr.Array()) != 1 {
            t.Fail()
            t.Log()
        }
        arr.Clear()
        if len(arr.Array()) != 0 {
            t.Fail()
            t.Log()
        }
    }
}

func TestUpdate(t *testing.T) {
    {
        j := NewObject()
        i := NewObject()
        i.Set("a", NewJson(1))
        i.Set("b", NewJson(2))
        j.Update(i)
        if j.Contains("a") == false {
            t.Fail()
            t.Log()
        }
        if j.Contains("b") == false {
            t.Fail()
            t.Log()
        }
    }
}

func TestPop(t *testing.T) {
    {
        arr := NewArray()
        arr.Append(NewJson(1))
        arr.Append(NewJson(2))
        arr.Pop(0)
        if !arr.Get(0).Equal(NewJson(2)) {
            t.Fail()
            t.Log()
        }
    }
    {
        j := NewObject()
        j.Set("a", NewJson(1))
        j.Set("b", NewJson(2))
        j.Pop("a")
        if j.Contains("a") == true {
            t.Fail()
            t.Log()
        }
        if j.Contains("b") == false {
            t.Fail()
            t.Log()
        }
        j.Pop("b")
        if j.Contains("b") == true {
            t.Fail()
            t.Log()
        }
    }
}

func TestInsert(t *testing.T) {
    {
        j := NewArray()
        j.Append(NewJson(1))
        j.Append(NewJson(2))
        j.Insert(0, NewJson(0))
        if j.Get(0).Int() != 0 {
            t.Fail()
            t.Log()
        }
        if j.Get(1).Int() != 1 {
            t.Fail()
            t.Log()
        }
        if j.Get(2).Int() != 2 {
            t.Fail()
            t.Log()
        }
    }
}

func TestAppend(t *testing.T) {
    {
        j := NewArray()
        j.Append(NewJson(1))
        j.Append(NewJson(2))
        if j.Get(0).Int() != 1 {
            t.Fail()
            t.Log()
        }
        if j.Get(1).Int() != 2 {
            t.Fail()
            t.Log()
        }
    }
}

func TestExtend(t *testing.T) {
    {
        j := NewArray()
        i := NewArray()
        i.Append(NewJson(1))
        i.Append(NewJson(2))
        j.Extend(i)
        if j.Get(0).Int() != 1 {
            t.Fail()
            t.Log()
        }
        if j.Get(1).Int() != 2 {
            t.Fail()
            t.Log()
        }
    }
}

func TestRemove(t *testing.T) {
    {
        j := NewArray()
        j.Append(NewJson(1))
        j.Append(NewJson(2))
        j.Append(NewJson(3))
        j.Remove(NewJson(2))
        if j.Get(0).Int() != 1 {
            t.Fail()
            t.Log()
        }
        if j.Get(1).Int() != 3 {
            t.Fail()
            t.Log()
        }
    }
}

func TestIndex(t *testing.T) {
    {
        j := NewArray()
        j.Append(NewJson(1))
        j.Append(NewJson(2))
        j.Append(NewJson(3))
        if j.Index(NewJson(1)) != 0 {
            t.Fail()
            t.Log()
        }
        if j.Index(NewJson(2)) != 1 {
            t.Fail()
            t.Log()
        }
        if j.Index(NewJson(3)) != 2 {
            t.Fail()
            t.Log()
        }
    }
}

func benchmark1(t *testing.T, repo string) {
    filePaths := []string{}
    filePaths = append(filePaths, "testdata/canada.json")
    for _, filePath := range filePaths {
        if file, err := os.Open(filePath); err == nil {
            buf := bufio.NewReader(file)
            for {
                line, err := buf.ReadString('\n')
                if err != nil {
                    break
                }
                if repo == "valyala/fastjson" {
                    j, err := fastjson.Parse(line)
                    if j == nil || err != nil {
                        t.Fail()
                        t.Log(filePath)
                    }
                } else {
                    j := Loads(line)
                    if j == nil {
                        t.Fail()
                        t.Log(filePath)
                    }
                }
            }
        }
    }
}

func benchmark2(t *testing.T, repo string) {
    filePaths := []string{}
    filePaths = append(filePaths, "testdata/citm_catalog.json")
    for _, filePath := range filePaths {
        if file, err := os.Open(filePath); err == nil {
            buf := bufio.NewReader(file)
            for {
                line, err := buf.ReadString('\n')
                if err != nil {
                    break
                }
                if repo == "valyala/fastjson" {
                    j, err := fastjson.Parse(line)
                    if j == nil || err != nil {
                        t.Fail()
                        t.Log(filePath)
                    }
                } else {
                    j := Loads(line)
                    if j == nil {
                        t.Fail()
                        t.Log(filePath)
                    }
                }
            }
        }
    }
}

func benchmark3(t *testing.T, repo string) {
    filePaths := []string{}
    filePaths = append(filePaths, "testdata/large.json")
    for _, filePath := range filePaths {
        if file, err := os.Open(filePath); err == nil {
            buf := bufio.NewReader(file)
            for {
                line, err := buf.ReadString('\n')
                if err != nil {
                    break
                }
                if repo == "valyala/fastjson" {
                    j, err := fastjson.Parse(line)
                    if j == nil || err != nil {
                        t.Fail()
                        t.Log(filePath)
                    }
                } else {
                    j := Loads(line)
                    if j == nil {
                        t.Fail()
                        t.Log(filePath)
                    }
                }
            }
        }
    }
}

func benchmark4(t *testing.T, repo string) {
    filePaths := []string{}
    filePaths = append(filePaths, "testdata/medium.json")
    for _, filePath := range filePaths {
        if file, err := os.Open(filePath); err == nil {
            buf := bufio.NewReader(file)
            for {
                line, err := buf.ReadString('\n')
                if err != nil {
                    break
                }
                if repo == "valyala/fastjson" {
                    j, err := fastjson.Parse(line)
                    if j == nil || err != nil {
                        t.Fail()
                        t.Log(filePath)
                    }
                } else {
                    j := Loads(line)
                    if j == nil {
                        t.Fail()
                        t.Log(filePath)
                    }
                }
            }
        }
    }
}

func benchmark5(t *testing.T, repo string) {
    filePaths := []string{}
    filePaths = append(filePaths, "testdata/small.json")
    for _, filePath := range filePaths {
        if file, err := os.Open(filePath); err == nil {
            buf := bufio.NewReader(file)
            for {
                line, err := buf.ReadString('\n')
                if err != nil {
                    break
                }
                if repo == "valyala/fastjson" {
                    j, err := fastjson.Parse(line)
                    if j == nil || err != nil {
                        t.Fail()
                        t.Log(filePath)
                    }
                } else {
                    j := Loads(line)
                    if j == nil {
                        t.Fail()
                        t.Log(filePath)
                    }
                }
            }
        }
    }
}

func benchmark6(t *testing.T, repo string) {
    filePaths := []string{}
    filePaths = append(filePaths, "testdata/twitter.json")
    for _, filePath := range filePaths {
        if file, err := os.Open(filePath); err == nil {
            buf := bufio.NewReader(file)
            for {
                line, err := buf.ReadString('\n')
                if err != nil {
                    break
                }
                if repo == "valyala/fastjson" {
                    j, err := fastjson.Parse(line)
                    if j == nil || err != nil {
                        t.Fail()
                        t.Log(filePath)
                    }
                } else {
                    j := Loads(line)
                    if j == nil {
                        t.Fail()
                        t.Log(filePath)
                    }
                }
            }
        }
    }
}

func TestJsonBenchmark(t *testing.T) {
    for i := 0; i < 10; i++ {
        benchmark1(t, "zhengyuhong/json")
        benchmark2(t, "zhengyuhong/json")
        benchmark3(t, "zhengyuhong/json")
        benchmark4(t, "zhengyuhong/json")
        benchmark5(t, "zhengyuhong/json")
        benchmark6(t, "zhengyuhong/json")
    }
}

func TestFastJsonBenchmark(t *testing.T) {
    for i := 0; i < 10; i++ {
        benchmark1(t, "valyala/fastjson")
        benchmark2(t, "valyala/fastjson")
        benchmark3(t, "valyala/fastjson")
        benchmark4(t, "valyala/fastjson")
        benchmark5(t, "valyala/fastjson")
        benchmark6(t, "valyala/fastjson")
    }
}

func TestSwap(t *testing.T) {
    {
        j := NewJson(1)
        i := NewJson("1")
        j.Swap(i)
        if Dumps(j) != `"1"` {
            t.Fail()
            t.Log()
        }
        if Dumps(i) != `1` {
            t.Fail()
            t.Log()
        }
    }
}

func TestValidate(t *testing.T) {
    {
        s := ""
        if validate(s) {
            t.Fail()
            t.Log()
        }
    }
}

func TestType(t *testing.T) {
    {
        j := NewJson(1)
        if j.Type() != IntType {
            t.Fail()
            t.Log()
        }
    }
}

func TestCopy(t *testing.T) {
    {
        j := NewObject()
        j.Set("a", NewJson("A"))
        i := j.Copy()
        if Dumps(j) != Dumps(i) {
            t.Fail()
            t.Log()
        }
        j.Set("a", NewJson("a"))
        if Dumps(j) == Dumps(i) {
            t.Fail()
            t.Log()
        }
    }
}

func TestReverse(t *testing.T) {
    {
        arr := NewArray()
        arr.Append(NewJson(1))
        arr.Append(NewJson(2))
        arr.Reverse()
        if !arr.Get(0).Equal(NewJson(2)) {
            t.Fail()
            t.Log()
        }
    }
}
