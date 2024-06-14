package fp

import (
	"encoding"
	"fmt"
	"reflect"
	"strconv"
)

var (
	_encodingInterface = reflect.TypeOf((*encoding.TextUnmarshaler)(nil)).Elem()
	_err               = reflect.TypeOf((*error)(nil)).Elem()
)

func fromFieldNameToKey(f reflect.StructField) (fieldName string, required bool, defaultValue string) {
	tag := f.Tag
	fieldName = tag.Get("fp")
	required, _ = strconv.ParseBool(tag.Get("fp-req"))
	defaultValue = tag.Get("fp-def")
	return
}

func Parse(dst any, values map[string][]string) error {
	dstValue := reflect.ValueOf(dst)
	dstType := reflect.TypeOf(dst)

	if dstValue.Kind() != reflect.Pointer {
		return fmt.Errorf("dst type %v is not a pointer", dstType)
	}
	if dstValue.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("dst type %v is not a struct pointer", dstType)
	}
	for i := range dstValue.Elem().NumField() {
		field := dstValue.Elem().Field(i)
		structField := dstType.Elem().Field(i)

		if !structField.IsExported() {
			continue
		}
		fieldName, isRequired, defaultValue := fromFieldNameToKey(structField)
		ptrField := reflect.New(field.Type())

		if field.Kind() == reflect.Struct &&
			!field.Type().Implements(_encodingInterface) &&
			!ptrField.Type().Implements(_encodingInterface) {
			newPtr := reflect.New(field.Type())
			err := Parse(newPtr.Interface(), values)
			if err != nil {
				return err
			}
			field.Set(newPtr.Elem())
			continue
		}

		_, ok := values[fieldName]
		if !ok && isRequired {
			return fmt.Errorf("cannot find field %v", fieldName)
		}
		if !ok && defaultValue == "" {
			continue
		}
		if !ok && defaultValue != "" {
			values[fieldName] = []string{defaultValue}
		}
		content := values[fieldName]
		if len(content) == 0 && isRequired {
			return fmt.Errorf("field %v is empty", fieldName)
		}
		if len(content) == 0 {
			continue
		}

		// check if pointer implements encoding
		if ptrField.Type().Implements(_encodingInterface) {
			err := callDecoding(ptrField, content[0])
			if err != nil {
				return err
			}
			field.Set(ptrField.Elem())
			continue
		}

		switch field.Kind() {
		case reflect.Slice:
			newSlice := reflect.MakeSlice(field.Type(), len(content), len(content))
			for index, c := range content {
				sliceValue := newSlice.Index(index)
				newPtr := reflect.New(sliceValue.Type())
				err := parseString(newPtr.Interface(), c)
				if err != nil {
					return err
				}
				sliceValue.Set(newPtr.Elem())
			}
			field.Set(newSlice)
		default:
			d := reflect.New(field.Type())
			err := parseString(d.Interface(), content[0])
			if err != nil {
				return err
			}
			field.Set(d.Elem())
		}

	}
	return nil
}

func callDecoding(dst reflect.Value, s string) error {
	method := _encodingInterface.Method(0)
	arg := reflect.ValueOf([]byte(s))
	res := dst.MethodByName(method.Name).Call([]reflect.Value{arg})
	if len(res) != 1 {
		return fmt.Errorf("encoding.UnmarshalText should return only 1 value got %d", len(res))
	}
	if res[0].Type() != _err {
		return fmt.Errorf("encoding.UnmarshalText did not return an error got %v", reflect.TypeOf(res[0].Interface()))
	}
	if res[0].IsNil() {
		return nil
	}
	return res[0].Interface().(error)

}

func parseString(dst any, s string) error {
	dstValue := reflect.ValueOf(dst)
	dstType := dstValue.Type()
	if dstValue.Kind() != reflect.Pointer {
		return fmt.Errorf("dst type %v is not a pointer", dstType)
	}
	if !dstValue.Elem().CanSet() {
		return fmt.Errorf("cannot set")
	}

	if dstValue.Type().Implements(_encodingInterface) {
		err := callDecoding(dstValue, s)
		return err
	}
	switch dstValue.Elem().Kind() {
	case reflect.String:
		dstValue.Elem().SetString(s)
	case reflect.Bool:
		res, err := strconv.ParseBool(s)
		if err != nil {
			return err
		}
		dstValue.Elem().SetBool(res)
	case reflect.Uint8:
		res, err := strconv.ParseUint(s, 10, 8)
		if err != nil {
			return err
		}
		dstValue.Elem().SetUint(res)
	case reflect.Uint16:
		res, err := strconv.ParseUint(s, 10, 16)
		if err != nil {
			return err
		}
		dstValue.Elem().SetUint(res)
	case reflect.Uint32:
		res, err := strconv.ParseUint(s, 10, 32)
		if err != nil {
			return err
		}
		dstValue.Elem().SetUint(res)
	case reflect.Uint64, reflect.Uint:
		res, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			return err
		}
		dstValue.Elem().SetUint(res)
	case reflect.Int, reflect.Int64:
		res, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return err
		}
		dstValue.Elem().SetInt(res)
	case reflect.Int32:
		res, err := strconv.ParseInt(s, 10, 32)
		if err != nil {
			return err
		}
		dstValue.Elem().SetInt(res)
	case reflect.Int16:
		res, err := strconv.ParseInt(s, 10, 16)
		if err != nil {
			return err
		}
		dstValue.Elem().SetInt(res)
	case reflect.Int8:
		res, err := strconv.ParseInt(s, 10, 8)
		if err != nil {
			return err
		}
		dstValue.Elem().SetInt(res)
	case reflect.Float32:
		res, err := strconv.ParseFloat(s, 32)
		if err != nil {
			return err
		}
		dstValue.Elem().SetFloat(res)
	case reflect.Float64:
		res, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return err
		}
		dstValue.Elem().SetFloat(res)
	default:
		return fmt.Errorf("parsing string for type %v is not implemented:", dstType.Elem().String())
	}
	return nil
}
