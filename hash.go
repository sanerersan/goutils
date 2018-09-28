package goutils

import (
	"crypto/md5"
	"fmt"
	"reflect"
)

//try to hash every type you can check wether two types are equal by this

func IsBaseType(t reflect.Type) bool {
	switch t.Kind() {
		case reflect.Bool,reflect.Int,reflect.Int8,reflect.Int16,
		reflect.Int32,reflect.Int64,reflect.Uint,reflect.Uint8,
		reflect.Uint16,reflect.Uint32,reflect.Uint64,reflect.Uintptr,
		reflect.Float32,reflect.Float64,reflect.Complex64,reflect.Complex128,
		reflect.Chan,reflect.String,reflect.UnsafePointer:
		return true;
	}

	return false;
}

func isTypeAlreadyRecord(str string,recordType map[string]int) bool {
	if _,ok := recordType[str]; ok {
		return true
	}
	recordType[str] = 1
	return false
}

//watch out loop back
func getDescString(t reflect.Type,recordType map[string]int) string {
	if IsBaseType(t) {
		return t.String()
	}

	switch t.Kind() {
	case reflect.Array:
		if isTypeAlreadyRecord(t.String(),recordType) {
			return t.String()
		}
		return fmt.Sprint(t.String(),t.Len(),getDescString(t.Elem(),recordType))
	case reflect.Func:
		str := fmt.Sprint(t.String(),t.NumIn(),t.NumOut())
		if isTypeAlreadyRecord(str,recordType) {
			return str
		}
		for i := 0;i < t.NumIn();i++ {
			str += getDescString(t.In(i),recordType)
		}
		for i := 0;i < t.NumOut();i++ {
			str += getDescString(t.Out(i),recordType)
		}

		return str
	case reflect.Map:
		if isTypeAlreadyRecord(t.String(),recordType) {
			return t.String()
		}
		return fmt.Sprint(t.String(),getDescString(t.Key(),recordType),getDescString(t.Elem(),recordType))
	case reflect.Ptr:
		if isTypeAlreadyRecord(t.String(),recordType) {
			return t.String()
		}
		return fmt.Sprint(t.String(),getDescString(t.Elem(),recordType))
	case reflect.Slice:
		if isTypeAlreadyRecord(t.String(),recordType) {
			return t.String()
		}
		return fmt.Sprint(t.String(),getDescString(t.Elem(),recordType))
	case reflect.Interface:
		str := fmt.Sprint(t.String(),t.NumMethod())
		if isTypeAlreadyRecord(str,recordType) {
			return str
		}
		for i := 0;i < t.NumMethod();i++ {
			str += getDescString(t.Method(i).Type,recordType)
		}
		return str
	case reflect.Struct:
		//cannot covert between structs which have different name and same fields
		str := fmt.Sprint(t.String(),t.NumField(),t.NumMethod())
		if isTypeAlreadyRecord(str,recordType) {
			return str
		}
		for i := 0;i < t.NumField();i++ {
			str += getDescString(t.Field(i).Type,recordType)
		}
		for i := 0;i < t.NumMethod();i++ {
			str += getDescString(t.Method(i).Type,recordType)
		}
		return str
	}

	return "unknown"
}

func HashType(st interface{}) (shortName string,hash string) {
	shortName = ""
	hash = ""
	if nil == st {
		return
	}
	recordType := make(map[string]int)
	t := reflect.TypeOf(st)
	desc := getDescString(t,recordType)
	sum := md5.Sum([]byte(desc))
	hash = ""
	for i := range sum {
		hash += fmt.Sprintf("%02x",sum[i])
	}
	shortName = t.String()
	return
}

func IsTypeEqual(st1 interface{},st2 interface{}) bool {
	if (nil == st1) && (nil == st2) {
		return true
	}
	if (nil == st1) || (nil == st2) {
		return false
	}

	sn1,h1 := HashType(st1)
	sn2,h2 := HashType(st2)

	if (sn1 == sn2) && (h1 == h2) {
		return true
	}

	return false
}