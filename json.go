package json

import (
    "encoding/json"
    "fmt"
    "regexp"
    "sort"
    "strconv"
    "strings"
)

// Json Type
const (
    NullType   = iota // NullValue
    BoolType          // BoolValue
    IntType           // IntValue
    FloatType         // FloatValue
    StringType        // StringValue
    ArrayType         // ArrayValue
    ObjectType        // ObjectValue
)

const (
    epsilon = 10e-6
)

var None *Json // Just as read-only Parameter for NewJson, for example j := NewJson(None)

// Json Definition
type Json struct {
    valueType int8
    value     any
}

// Just as type of Parameter for NewJson
type ValueType interface {
    bool | int | int32 | int64 | float64 | string | []any | map[string]any | *Json
}

// Type
func (j *Json) Type() int8 {
    return j.valueType
}

// NewJson - create a Json with value (if provided) and return its pointer
func NewJson[T ValueType](value ...T) *Json {
    j := &Json{valueType: NullType, value: nil}
    for _, e := range value {
        var v any = e
        if t, ok := v.(*Json); ok && t == nil {
            break
        }
        if err := j.set(v); err != nil {
            return nil
        }
    }
    return j
}

// NewObject - create an empty ObjectValue and return its pointer
func NewObject() *Json {
    j := Loads(`{}`)
    return j
}

// NewArray - create an empty ArrayValue and return its pointer
func NewArray() *Json {
    j := Loads(`[]`)
    return j
}

// set - convert any to Json
func (j *Json) set(value any) error {
    switch value.(type) {
    case bool:
        j.valueType = BoolType
        j.value = value.(bool)
    case int:
        j.valueType = IntType
        value, _ := value.(int)
        j.value = int64(value)
    case int64:
        j.valueType = IntType
        j.value, _ = value.(int64)
    case float64:
        j.valueType = FloatType
        j.value, _ = value.(float64)
    case json.Number:
        s := value.(json.Number).String()
        if strings.Contains(s, ".") {
            j.valueType = FloatType
            j.value, _ = value.(json.Number).Float64()
        } else {
            j.valueType = IntType
            j.value, _ = value.(json.Number).Int64()
        }
    case string:
        j.valueType = StringType
        j.value = value.(string)
    case []any:
        j.valueType = ArrayType
        array := value.([]any)
        arrayValue := []*Json{}
        for _, value := range array {
            e := &Json{}
            if err := e.set(value); err != nil {
                return err
            }
            arrayValue = append(arrayValue, e)
        }
        j.value = arrayValue
    case map[string]any:
        j.valueType = ObjectType
        object := value.(map[string]any)
        objectValue := make(map[string]*Json, 0)
        for key, value := range object {
            v := &Json{}
            if err := v.set(value); err != nil {
                return err
            }
            objectValue[key] = v
        }
        j.value = objectValue
    case *Json:
        o, _ := value.(*Json)
        j.valueType = o.valueType
        j.value = o.value
    default:
        j.valueType = NullType
    }
    return nil
}

// Loads - Deserialize JSON document to Json struct
// PARAMS: s - string for Json document
func Loads(s string) *Json {
    if !validate(s) {
        return nil
    }
    j := &Json{}
    decoder := json.NewDecoder(strings.NewReader(s))
    decoder.UseNumber()
    var value any
    if err := decoder.Decode(&value); err != nil {
        return nil
    }
    if j.set(value) != nil {
        return nil
    }
    return j
}

// validate - check whether JSON document is valid for uncaught case in encoding/json
func validate(s string) bool {
    if len(s) == 0 {
        return false
    }
    if '0' <= s[0] && s[0] <= '9' {
        reg := regexp.MustCompile(`^(-?\d+)(\.\d+)?$`)
        if matched := reg.MatchString(s); !matched {
            return false
        }
    }
    return true
}

// Dumps - serialize Json struct to JSON document
func Dumps(j *Json) string {
    s := ""
    switch j.valueType {
    case IntType:
        intValue, _ := j.value.(int64)
        s = strconv.FormatInt(intValue, 10)
    case FloatType:
        floatValue, _ := j.value.(float64)
        s = strconv.FormatFloat(floatValue, 'f', -1, 64)
    case BoolType:
        boolValue, _ := j.value.(bool)
        if boolValue {
            s = "true"
        } else {
            s = "false"
        }
    case StringType:
        stringValue, _ := j.value.(string)
        s = `"` + strings.Replace(stringValue, `"`, `\"`, -1) + `"`
    case NullType:
        s = "null"
    case ArrayType:
        sl := make([]string, 0)
        arrayValue := j.value.([]*Json)
        for _, e := range arrayValue {
            sl = append(sl, Dumps(e))
        }
        s = fmt.Sprintf("[%s]", strings.Join(sl, ","))
    case ObjectType:
        sl := make([]string, 0)
        keys := []string{}
        objectValue, _ := j.value.(map[string]*Json)
        for key, _ := range objectValue {
            keys = append(keys, key)
        }
        sort.Strings(keys)
        for _, key := range keys {
            value, _ := objectValue[key]
            e := fmt.Sprintf(`"%s":%s`, key, Dumps(value))
            sl = append(sl, e)
        }
        s = fmt.Sprintf("{%s}", strings.Join(sl, ","))
    }
    return s
}

// Equal - return j == i
func (j *Json) Equal(i *Json) bool {
    if j.valueType != i.valueType {
        return false
    }
    if j.valueType == NullType {
        return true
    } else if j.valueType == BoolType {
        jv, _ := j.value.(bool)
        iv, _ := i.value.(bool)
        return jv == iv
    } else if j.valueType == IntType {
        jv, _ := j.value.(int64)
        iv, _ := i.value.(int64)
        return jv == iv
    } else if j.valueType == FloatType {
        jv, _ := j.value.(float64)
        iv, _ := i.value.(float64)
        delta := jv - iv
        if delta < epsilon && delta > -1*epsilon {
            return true
        }
    } else if j.valueType == StringType {
        jv, _ := j.value.(string)
        iv, _ := i.value.(string)
        return jv == iv
    } else if j.valueType == ArrayType {
        jv, _ := j.value.([]*Json)
        iv, _ := i.value.([]*Json)
        if len(jv) != len(iv) {
            return false
        }
        for index := 0; index < len(jv); index++ {
            if !jv[index].Equal(iv[index]) {
                return false
            }
        }
        return true
    } else if j.valueType == ObjectType {
        jv, _ := j.value.(map[string]*Json)
        iv, _ := i.value.(map[string]*Json)
        if len(jv) != len(iv) {
            return false
        }
        for key, value := range jv {
            if _, ok := iv[key]; !ok {
                return false
            }
            if !value.Equal(iv[key]) {
                return false
            }
        }
        return true
    }
    return false
}

// IsNull
func (j *Json) IsNull() bool {
    return j.valueType == NullType
}

// IsObject
func (j *Json) IsObject() bool {
    return j.valueType == ObjectType
}

// IsArray
func (j *Json) IsArray() bool {
    return j.valueType == ArrayType
}

// IsString
func (j *Json) IsString() bool {
    return j.valueType == StringType
}

// IsInt
func (j *Json) IsInt() bool {
    return j.valueType == IntType
}

// IsFloat
func (j *Json) IsFloat() bool {
    return j.valueType == FloatType
}

// IsBool
func (j *Json) IsBool() bool {
    return j.valueType == BoolType
}

// String - return the underlying JSON string for the Json.value. But call IsString() before String()
func (j *Json) String() string {
    stringValue, _ := j.value.(string)
    return stringValue
}

// Int - return the underlying JSON int for the Json.value. But call IsInt() before Int()
func (j *Json) Int() int64 {
    intValue, _ := j.value.(int64)
    return intValue
}

// Bool - return the underlying JSON bool for the Json.value. But call IsBool() before Bool()
func (j *Json) Bool() bool {
    boolValue, _ := j.value.(bool)
    return boolValue
}

// Float - return the underlying JSON float for the Json.value. But call IsFloat() before Float()
func (j *Json) Float() float64 {
    floatValue, _ := j.value.(float64)
    return floatValue
}

// Array - return the underlying JSON array for the Json.value. But call IsArray() before Array()
func (j *Json) Array() []*Json {
    arrayValue, _ := j.value.([]*Json)
    return arrayValue
}

// Object - return the underlying JSON object for the Json.value. But call IsObject() before Object()
func (j *Json) Object() map[string]*Json {
    objectValue, _ := j.value.(map[string]*Json)
    return objectValue
}

// Copy -  Deep copy operation on Json struct.
func (j *Json) Copy() *Json {
    return Loads(Dumps(j))
}

// Swap - swap j with i
func (j *Json) Swap(i *Json) {
    value, valueType := j.value, j.valueType
    j.value, j.valueType = i.value, i.valueType
    i.value, i.valueType = value, valueType
}

func toKey(keyOrIndex any) (string, bool) {
    if key, ok := keyOrIndex.(string); ok {
        return key, true
    }
    return "", false
}

func toValue(keyOrValue any) (*Json, bool) {
    if value, ok := keyOrValue.(*Json); ok {
        return value, true
    }
    return nil, false
}

func toIndex(keyOrIndex any) (int, bool) {
    if index, ok := keyOrIndex.(int); ok {
        return index, true
    } else if index, ok := keyOrIndex.(int8); ok {
        return int(index), true
    } else if index, ok := keyOrIndex.(int16); ok {
        return int(index), true
    } else if index, ok := keyOrIndex.(int32); ok {
        return int(index), true
    } else if index, ok := keyOrIndex.(int64); ok {
        return int(index), true
    }
    return 0, false
}

// Get - Get value by key(for ObjectValue) or index(for ArrayValue)
// PARAMS:
//   keyOrIndex - key for ObjectValue, index for ArrayValue
func (j *Json) Get(keyOrIndex any) *Json {
    if key, ok := toKey(keyOrIndex); ok && j.IsObject() {
        objectValue, _ := j.value.(map[string]*Json)
        if value, ok := objectValue[key]; ok {
            return value
        }
    } else if index, ok := toIndex(keyOrIndex); ok && j.IsArray() {
        arrayValue, _ := j.value.([]*Json)
        if index < len(arrayValue) {
            return arrayValue[index]
        }
    }
    return NewJson(None)
}

// Set - Set value by key(for ObjectValue) or index(for ArrayValue)
// PARAMS:
//   keyOrIndex
//     - key for ObjectValue if keyOrIndex is string
//     - index for ArrayValue if keyOrIndex is int
func (j *Json) Set(keyOrIndex any, value *Json) *Json {
    if key, ok := toKey(keyOrIndex); ok && j.IsObject() {
        objectValue, _ := j.value.(map[string]*Json)
        objectValue[key] = value
        j.value = objectValue
    } else if index, ok := toIndex(keyOrIndex); ok && j.IsArray() {
        arrayValue, _ := j.value.([]*Json)
        if index < len(arrayValue) {
            arrayValue[index] = value
            j.value = arrayValue
        }
    }
    return j
}

// SetDefault - ObjectValue.Set(key, value) if key not in ObjectValue
func (j *Json) SetDefault(key string, value *Json) {
    if j.IsObject() {
        objectValue, _ := j.value.(map[string]*Json)
        if _, ok := objectValue[key]; !ok {
            objectValue[key] = value
            j.value = objectValue
        }
    }
}

// Keys - a slice providing a view on ObjectValue's keys
func (j *Json) Keys() []string {
    keys := []string{}
    if j.IsObject() {
        objectValue, _ := j.value.(map[string]*Json)
        for key, _ := range objectValue {
            keys = append(keys, key)
        }
    }
    return keys
}

// Values - a slice providing a view on ObjectValue's values
func (j *Json) Values() []*Json {
    values := []*Json{}
    if j.IsObject() {
        objectValue, _ := j.value.(map[string]*Json)
        for _, value := range objectValue {
            values = append(values, value)
        }
    }
    return values
}

// Contains
//   - Contains(key) True if ObjectValue has key, else False
//   - Contains(value) True if ArrayValue has value, else False
func (j *Json) Contains(keyOrValue any) bool {
    if key, ok := toKey(keyOrValue); ok && j.IsObject() {
        objectValue, _ := j.value.(map[string]*Json)
        if _, ok := objectValue[key]; ok {
            return true
        }
    } else if value, ok := toValue(keyOrValue); ok {
        arrayValue, _ := j.value.([]*Json)
        for _, e := range arrayValue {
            if e.Equal(value) {
                return true
            }
        }
    }
    return false
}

// Clear - Remove all items from ObjectValue or ArrayValue
func (j *Json) Clear() {
    if j.IsObject() {
        j.value = make(map[string]*Json, 0)
    } else if j.IsArray() {
        j.value = []*Json{}
    }
}

// Update - update j ObjectValue with key value from i ObjectValue
func (j *Json) Update(i *Json) {
    if j.IsObject() && i.IsObject() {
        for _, key := range i.Keys() {
            j.Set(key, i.Get(key))
        }
    }
}

// Pop
//   - Pop(key) -> value, remove specified key on ObjectValue
//   - Pop(index) -> value, remove specified index in ArrayValue
//
func (j *Json) Pop(keyOrIndex any) {
    if key, ok := keyOrIndex.(string); ok && j.IsObject() {
        objectValue, _ := j.value.(map[string]*Json)
        delete(objectValue, key)
        j.value = objectValue
    } else if index, ok := toIndex(keyOrIndex); ok && j.IsArray() {
        arrayValue, _ := j.value.([]*Json)
        for i := index + 1; i < len(arrayValue); i++ {
            arrayValue[i-1] = arrayValue[i]
        }
        if index < len(arrayValue) {
            arrayValue = arrayValue[0 : len(arrayValue)-1]
        }
        j.value = arrayValue
    }
}

// Insert - insert value before index in ArrayValue
func (j *Json) Insert(index int, value *Json) {
    if j.IsArray() {
        arrayValue, _ := j.value.([]*Json)
        if index < len(arrayValue) {
            arrayValue = append(arrayValue, nil)
            for i := len(arrayValue) - 1; i > index; i-- {
                arrayValue[i] = arrayValue[i-1]
            }
            arrayValue[index] = value
        }
        j.value = arrayValue
    }
}

// Append - append Json to end of Json.value in ArrayValue
//
func (j *Json) Append(value *Json) {
    if j.IsArray() {
        arrayValue, _ := j.value.([]*Json)
        arrayValue = append(arrayValue, value)
        j.value = arrayValue
    }
}

// Extend - extend array by appending elements from the value in ArrayValue
//
func (j *Json) Extend(value *Json) {
    if value.IsArray() && j.IsArray() {
        arrayValue, _ := value.value.([]*Json)
        for _, e := range arrayValue {
            j.Append(e)
        }
    }
}

// Reverse - reverse in place in ArrayValue
func (j *Json) Reverse() {
    reverse := func(s []*Json) {
        for p, q := 0, len(s)-1; p < q; p, q = p+1, q-1 {
            s[p], s[q] = s[q], s[p]
        }
    }
    if j.IsArray() {
        arrayValue, _ := j.value.([]*Json)
        reverse(arrayValue)
        j.value = arrayValue
    }
}

// Remove - remove first occurrence of value if exists in ArrayValue
func (j *Json) Remove(value *Json) {
    if j.IsArray() {
        arrayValue, _ := j.value.([]*Json)
        for i, e := range arrayValue {
            if e.Equal(value) {
                for k := i + 1; k < len(arrayValue); k++ {
                    arrayValue[k-1] = arrayValue[k]
                }
                arrayValue = arrayValue[0 : len(arrayValue)-1]
                break
            }
        }
        j.value = arrayValue
    }
}

// Index - return first index of value in ArrayValue
func (j *Json) Index(value *Json) int {
    if j.IsArray() {
        arrayValue, _ := j.value.([]*Json)
        for i, e := range arrayValue {
            if e.Equal(value) {
                return i
            }
        }
    }
    return -1
}
