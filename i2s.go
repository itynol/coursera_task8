package main

import (
	"errors"
	"reflect"
)

func sliceHandler(data interface{}, outType reflect.Type, outValueElem reflect.Value) error {
	if _, ok := data.([]interface{}); !ok {
		return errors.New("expected slice but get another struct")
	}
	typeElem := outType.Elem()
	sliceValue := reflect.New(typeElem).Elem()
	for _, value := range data.([]interface{}) {
		newValue := reflect.New(typeElem.Elem()).Elem()
		err := i2s(value, newValue.Addr().Interface())
		if err != nil {
			return err
		}
		sliceValue = reflect.Append(sliceValue, newValue)
	}
	outValueElem.Set(sliceValue)
	return nil
}

func i2s(data interface{}, out interface{}) error {
	outType := reflect.TypeOf(out)
	if outType.Kind() != reflect.Ptr && outType.Kind() != reflect.Slice {
		return errors.New("WTF")
	}
	outValue := reflect.ValueOf(out)
	outValueElem := outValue.Elem()
	if outType.Elem().Kind() == reflect.Slice {
		if err := sliceHandler(data, outType, outValueElem); err != nil {
			return err
		}
		return nil
	}
	customMap, ok := data.(map[string]interface{})
	if !ok {
		return errors.New("expected map, but get smt another")
	}
	for key, value := range customMap {
		elementField := outValueElem.FieldByName(key)
		switch value.(type) {
		case string:
			if elementField.Kind() != reflect.String {
				return errors.New("not string value but must be string")
			}
			elementField.SetString(value.(string))
		case float64:
			if elementField.Kind() != reflect.Int {
				return errors.New("not bool value but must be bool")
			}
			elementField.SetInt(int64(value.(float64)))
		case int:
			if elementField.Kind() != reflect.Int {
				return errors.New("not bool value but must be bool")
			}
			elementField.SetInt(value.(int64))
		case bool:
			if elementField.Kind() != reflect.Bool {
				return errors.New("not bool value but must be bool")
			}
			elementField.SetBool(value.(bool))
		default:
			err := i2s(value, elementField.Addr().Interface())
			if err != nil {
				return err
			}
		}
	}
	return nil
}
