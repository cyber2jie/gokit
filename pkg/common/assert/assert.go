package assert

import (
	"gokit/pkg/common/converter"
	"reflect"
)

func IsTrue(val bool, msg string) {
	if !val {
		panic(msg)
	}
}
func Equals(leftVal interface{}, rightVal interface{}, strict bool) {
	leftType := reflect.TypeOf(leftVal)
	rightType := reflect.TypeOf(rightVal)
	equals := false
	if leftType == rightType && leftType.Comparable() {
		equals = leftVal == rightVal
	} else if !leftType.Comparable() {
		//nop
	}
	if leftType != rightType && strict {
		panic("类型不一致")
	} else {
		containStr := leftType.Kind() == reflect.String || rightType.Kind() == reflect.String
		allIsStr := leftType.Kind() == reflect.String && rightType.Kind() == reflect.String
		if containStr && !allIsStr {
			leftValStr, error := converter.Convert(leftVal, reflect.Float64)
			if error == nil {
				rightValStr, error := converter.Convert(rightVal, reflect.Float64)
				if error == nil {
					equals = leftValStr == rightValStr
				}
			}
		} else if leftType.ConvertibleTo(rightType) {
			equals = leftVal == rightVal
		}
	}
	if !equals {
		panic("数据不相等")
	}
}
