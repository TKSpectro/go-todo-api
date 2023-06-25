package utils

import "reflect"

// Some blessed person had this banger ready https://github.com/sas1024/gorm-loggable/issues/18#issuecomment-535024656
// pointerType is made to chase for value through all the way following,
// leading pointers until reach the deep final value which is not a pointer
func pointerType(t reflect.Type) reflect.Type {
	for {
		if t.Kind() != reflect.Ptr {
			return t
		}

		t = t.Elem()
	}
}
