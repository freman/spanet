package server_test

// func TestReflection(t *testing.T) {
// 	s := &spanet.Spanet{}

// 	arg := interface{}(struct {
// 		Mode spanet.BlowerMode
// 	}{})

// 	rarg := reflect.TypeOf(arg)

// 	fn := interface{}(s.ControlBlower)
// 	rfn := reflect.TypeOf(fn)

// 	if rarg.NumField() != 1 {
// 		panic("invalid simple structure")
// 	}

// 	if rfn.NumIn() != 1 {
// 		panic("invalid simple func")
// 	}

// 	if rarg.Field(0).Type != rfn.In(0) {
// 		panic("argument mismatch")
// 	}

// // 	out := reflect.ValueOf(fn).Call([]reflect.Value{reflect.ValueOf(arg).Field(0)})

// // 	_, _ = inputs, outputs
// // }
