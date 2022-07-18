[![Build Status](https://github.com/zhengyuhong/json/actions/workflows/go.yml/badge.svg)](https://github.com/zhengyuhong/json/actions/workflows/go.yml)
[![GoDoc](https://pkg.go.dev/badge/github.com/zhengyuhong/json)](https://pkg.go.dev/github.com/zhengyuhong/json)
[![codecov](https://codecov.io/gh/zhengyuhong/json/branch/main/graph/badge.svg)](https://codecov.io/gh/zhengyuhong/fastjson)
# Json
Simple JSON for Golang, like Python Json, easy to use.

## Features

  * Decodes JSON docuement with standard [encoding/json](https://golang.org/pkg/encoding/json/).
  * Parses arbitrary JSON without schema, reflection, struct magic and code generation


## Usage

Loads JSON docuement Like Python json lib:
```go
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
```
Dumps JSON Like Python json lib:
```go
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
```

Easy to access object or array:
```go
func TestSet(t *testing.T) {
    {
        j := NewObject()
        j.Set("a", NewJson("a"))
        j.Set("b", NewJson("b"))
        j.Set("c", NewJson("c"))
        if Dumps(j.Get("a")) != `"a"` {
            t.Fail()
            t.Log()
        }
    }
    {
        j := NewArray()
        j.Append(NewJson("a"))
        j.Append(NewJson("a"))
        j.Set(1, NewJson("b"))
    }
}
```

See also [examples](https://github.com/zhengyuhong/json/blob/main/json_test.go).
