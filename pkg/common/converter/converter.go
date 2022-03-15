package converter

import (
	"fmt"
	"reflect"
	"strconv"
)

type ConvertError struct {
	val  interface{}
	kind reflect.Kind
}

func (ce ConvertError) Error() string {
	return "转换失败"
}

type Convertable interface {
	Get() interface{}
}

type String struct {
	value interface{}
}

func (str *String) Get() interface{} {
	return str.value.(string)
}

type Float struct {
	value interface{}
}

func (ft *Float) Get() interface{} {
	switch ft.value.(type) {
	case float32:
		return float64(ft.value.(float32))
	case float64:
		return ft.value.(float64)
	default:
	}
	return nil
}

type Int struct {
	value interface{}
}

func (it *Int) Get() interface{} {
	switch it.value.(type) {

	}
	return nil
}

type Uint struct {
	value interface{}
}

func (uit *Uint) Get() interface{} {
	return nil
}

var converters = map[string]func(ct Convertable, kind reflect.Kind) interface{}{
	"string": convString,
	"float":  convFloat,
	"int":    convInt,
	"uint":   convUInt,
}

func convString(ct Convertable, kind reflect.Kind) interface{} {
	val := ct.Get().(string)
	switch kind {
	case reflect.String:
		return val
	case reflect.Float32, reflect.Float64:
		valFloat, error := strconv.ParseFloat(val, 2)
		if error == nil {
			return valFloat
		}
	case reflect.Int, reflect.Int8, reflect.Int16,
		reflect.Int32, reflect.Int64:
		valInt, error := strconv.ParseInt(val, 0, 64)
		if error == nil {
			return valInt
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		valUInt, error := strconv.ParseUint(val, 0, 64)
		if error == nil {
			return valUInt
		}
	case reflect.Bool:
		valBool, error := strconv.ParseBool(val)
		if error == nil {
			return valBool
		}
	}
	return nil
}
func convFloat(ct Convertable, kind reflect.Kind) interface{} {
	val := ct.Get().(float64)
	switch kind {
	case reflect.String:
		return fmt.Sprintf("%f", val)
	case reflect.Float32, reflect.Float64:
		return val
	case reflect.Int, reflect.Int8, reflect.Int16,
		reflect.Int32, reflect.Int64:
		return int64(val)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return uint64(val)
	case reflect.Bool:
		if val == 1.0 {
			return true
		} else if val == 0.0 {
			return false
		}
	}
	return nil
}
func convInt(ct Convertable, kind reflect.Kind) interface{} {
	val := ct.Get().(int64)
	switch kind {
	case reflect.String:
		return fmt.Sprintf("%d", val)
	case reflect.Int, reflect.Int8, reflect.Int16,
		reflect.Int32, reflect.Int64:
		return val
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return uint64(val)
	case reflect.Bool:
		if val == 1 {
			return true
		} else if val == 0 {
			return false
		}
	}
	return nil
}
func convUInt(ct Convertable, kind reflect.Kind) interface{} {
	val := ct.Get().(uint64)
	switch kind {
	case reflect.String:
		return fmt.Sprintf("%d", val)
	case reflect.Int, reflect.Int8, reflect.Int16,
		reflect.Int32, reflect.Int64:
		return int64(val)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return val
	case reflect.Bool:
		if val == 1 {
			return true
		} else if val == 0 {
			return false
		}
	}
	return nil
}
func Convert(val interface{}, kind reflect.Kind) (interface{}, error) {
	converterKey := ""
	var ct Convertable = nil
	switch val.(type) {
	case string:
		converterKey = "string"
		ct = &String{val}
		break
	case float32, float64:
		converterKey = "float"
		ct = &Float{val}
		break
	case int, int8, int16, int32, int64:
		converterKey = "int"
		ct = &Int{val}
		break
	case uint, uint8, uint16, uint32, uint64:
		converterKey = "uint"
		ct = &Uint{val}
		break
	default:
	}
	converter := converters[converterKey]
	if converter != nil {
		convertVal := converter(ct, kind)
		if convertVal != nil {
			return convertVal, nil
		}
	}
	return nil, ConvertError{val, kind}
}
