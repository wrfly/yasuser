package utils

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

func ParseConfigEnv(conf interface{}, prefix []string) {
	duration := reflect.TypeOf(time.Second).Kind()

	confV, ok := conf.(reflect.Value)
	if !ok {
		confV = reflect.Indirect(reflect.ValueOf(conf))
	}

	confT := confV.Type()
	for i := 0; i < confV.NumField(); i++ {
		var (
			fieldStruct = confT.Field(i)
			field       = confV.Field(i)
			sName       = func() string {
				if y := fieldStruct.Tag.Get("yaml"); y != "" {
					return y
				}
				return fieldStruct.Name
			}()
			envName = strings.ToUpper(strings.Join(append(prefix, sName), "_"))
		)
		switch field.Kind() {
		case reflect.String:
			if v := os.Getenv(envName); v != "" {
				field.SetString(v)
			}
		case reflect.Int:
			if v := os.Getenv(envName); v != "" {
				vint, err := strconv.Atoi(v)
				if err != nil {
					errStr := fmt.Sprintf("convert %s error: %s", sName, err)
					logrus.Fatalf(AddLineNum(errStr))
				}
				field.SetInt(int64(vint))
			}
		case reflect.Bool:
			if v := os.Getenv(envName); v != "" {
				if strings.ToLower(v) == "true" {
					field.SetBool(true)
				} else {
					field.SetBool(false)
				}
			}
		case duration:
			if v := os.Getenv(envName); v != "" {
				d, err := time.ParseDuration(v)
				if err != nil {
					errStr := fmt.Sprintf("parse duration %s error: %s", sName, err)
					logrus.Fatalf(AddLineNum(errStr))
				}
				field.SetInt(int64(d))
			}
		case reflect.Struct:
			ParseConfigEnv(field, append(prefix, sName))
		}
	}
}

func DefaultConfig(conf interface{}) {
	duration := reflect.TypeOf(time.Second).Kind()

	confV, ok := conf.(reflect.Value)
	if !ok {
		confV = reflect.Indirect(reflect.ValueOf(conf))
	}

	confT := confV.Type()
	for i := 0; i < confV.NumField(); i++ {
		var (
			fieldStruct = confT.Field(i)
			field       = confV.Field(i)
			d           = fieldStruct.Tag.Get("default")
		)
		switch field.Kind() {
		case reflect.String:
			field.SetString(d)
		case reflect.Int:
			vint, err := strconv.Atoi(d)
			if err != nil {
				errStr := fmt.Sprintf("convert %s error: %s", fieldStruct.Name, err)
				logrus.Fatalf(AddLineNum(errStr))
			}
			field.SetInt(int64(vint))
		case reflect.Bool:
			if strings.ToLower(d) == "true" {
				field.SetBool(true)
			} else {
				field.SetBool(false)
			}
		case duration:
			d, err := time.ParseDuration(d)
			if err != nil {
				errStr := fmt.Sprintf("parse duration %s error: %s", fieldStruct.Name, err)
				logrus.Fatalf(AddLineNum(errStr))
			}
			field.SetInt(int64(d))
		case reflect.Struct:
			DefaultConfig(field)
		}
	}
}

func EnvConfigLists(conf interface{}, prefix []string) []string {
	envConfigLists := []string{}

	duration := reflect.TypeOf(time.Second).Kind()

	confV, ok := conf.(reflect.Value)
	if !ok {
		confV = reflect.Indirect(reflect.ValueOf(conf))
	}

	confT := confV.Type()
	for i := 0; i < confV.NumField(); i++ {
		var (
			fieldStruct = confT.Field(i)
			field       = confV.Field(i)
			d           = fieldStruct.Tag.Get("default")
			sName       = func() string {
				if y := fieldStruct.Tag.Get("yaml"); y != "" {
					return y
				}
				return fieldStruct.Name
			}()
			envName = strings.ToUpper(strings.Join(append(prefix, sName), "_"))
		)
		switch field.Kind() {
		case reflect.String:
			fallthrough
		case reflect.Int:
			fallthrough
		case duration:
			envConfigLists = append(envConfigLists, envName+"="+d)
		case reflect.Bool:
			e := envName + "=" + strings.ToLower(d)
			envConfigLists = append(envConfigLists, e)
		case reflect.Struct:
			envConfigLists = append(envConfigLists,
				EnvConfigLists(field, append(prefix, sName))...)
		}
	}

	return envConfigLists
}
