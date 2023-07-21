package utils

import (
	"reflect"
)

// Some blessed person had this banger ready https://github.com/sas1024/gorm-loggable/issues/18#issuecomment-535024656
// pointerType is made to chase for value through all the way following,
// leading pointers until reach the deep final value which is not a pointer
func PointerType(t reflect.Type) reflect.Type {
	for {
		if t.Kind() != reflect.Ptr {
			return t
		}

		t = t.Elem()
	}
}

// CopyCommonFields copies the common fields from the struct
// pointed to srcp to the struct pointed to by destp.
// https://stackoverflow.com/questions/59556480/convert-a-type-struct-a-to-b
func CopyCommonFields(destp, srcp interface{}) {
	destv := reflect.ValueOf(destp).Elem()
	srcv := reflect.ValueOf(srcp).Elem()

	destt := destv.Type()
	for i := 0; i < destt.NumField(); i++ {
		sf := destt.Field(i)
		v := srcv.FieldByName(sf.Name)
		if !v.IsValid() || !v.Type().AssignableTo(sf.Type) {
			continue
		}
		destv.Field(i).Set(v)
	}
}

// Convert converts a given struct into another struct type (destStructType)
// copying the common fields from the struct pointed to srcp to the returned struct.
// Based on https://stackoverflow.com/questions/59556480/convert-a-type-struct-a-to-b and extended to be generic
func Convert[T interface{}](destStructType T, srcp interface{}) T {
	dest := reflect.New(reflect.TypeOf(destStructType)).Elem()
	srcv := reflect.ValueOf(srcp).Elem()

	destt := dest.Type()
	for i := 0; i < destt.NumField(); i++ {
		sf := destt.Field(i)
		v := srcv.FieldByName(sf.Name)
		if !v.IsValid() || !v.Type().AssignableTo(sf.Type) {
			continue
		}
		dest.Field(i).Set(v)
	}

	return dest.Interface().(T)
}
