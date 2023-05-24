package util

import (
	"crypto/md5"
	"fmt"
	"github.com/shirou/gopsutil/cpu"
	"reflect"
	"strconv"
	"time"
)

func CPUPercent() float64 {
	// 获取 CPU 占用率
	percent, err := cpu.Percent(time.Second, false)
	if err != nil {
		panic(err)
	}
	return percent[0]
}

func Encryption(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}

func StructToMap(value interface{}) map[string]interface{} {
	v := reflect.Indirect(reflect.ValueOf(value))
	if v.Kind() != reflect.Struct {
		panic("value must be a struct")
	}
	vTyp := v.Type()
	m := make(map[string]interface{})
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		//匿名字段或非暴露字段
		if !vTyp.Field(i).IsExported() || vTyp.Field(i).Anonymous {
			continue
		}
		fKey := vTyp.Field(i).Name
		var fValue interface{}
		switch field.Kind() {
		case reflect.Struct:
			fValue = StructToMap(field.Interface())
		default:
			fValue = reflect.Indirect(field).Interface()
		}
		m[fKey] = fValue
	}
	return m
}

func MapToStruct(src map[string]string, dest interface{}) {
	destPValue := reflect.ValueOf(dest)
	if destPValue.Kind() != reflect.Pointer || destPValue.IsNil() {
		panic("dest must be a non null pointer")
	}
	destValue := reflect.Indirect(destPValue)
	if destValue.Kind() != reflect.Struct {
		panic("dest must be a struct type pointer")
	}
	for k, v := range src {
		field := destValue.FieldByName(k)
		if field.Kind() == reflect.Pointer {
			field.Set(reflect.New(field.Type().Elem()))
			field = field.Elem()
		}
		switch field.Kind() {
		case reflect.Int:
			i, err := strconv.Atoi(v)
			if err != nil {
				panic(err)
			}
			field.Set(reflect.ValueOf(i))
		case reflect.String:
			field.Set(reflect.ValueOf(v))
		default:
			panic(fmt.Sprintf("can't marshal %v", field.Kind()))
		}
	}
}
