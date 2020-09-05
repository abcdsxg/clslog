package clslog

import (
	"encoding/json"
	"reflect"

	"github.com/rs/zerolog"
)

type ZerologHook struct {
	cls       *ClsClient
	LogLevels []zerolog.Level
}

func NewZerologHook(cls *ClsClient, logLevels []zerolog.Level) *ZerologHook {
	return &ZerologHook{
		cls:       cls,
		LogLevels: logLevels,
	}
}

func (hook ZerologHook) Run(e *zerolog.Event, level zerolog.Level, msg string) {
	if !Contain(level, hook.LogLevels) {
		return
	}

	v := reflect.ValueOf(e)
	fieldBuf := reflect.Indirect(v).FieldByName("buf")
	buf := fieldBuf.Bytes()
	buf = append(buf, byte('}')) //buf does not contain } , needs to be completed

	fieldsMap := map[string]interface{}{}
	_ = json.Unmarshal(buf, &fieldsMap)

	fieldsMap["message"] = msg //buf does not contain message
	hook.cls.UploadMap(fieldsMap)
}

func Contain(obj interface{}, target interface{}) bool {
	targetValue := reflect.ValueOf(target)
	switch reflect.TypeOf(target).Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < targetValue.Len(); i++ {
			if targetValue.Index(i).Interface() == obj {
				return true
			}
		}
	case reflect.Map:
		if targetValue.MapIndex(reflect.ValueOf(obj)).IsValid() {
			return true
		}
	}

	return false
}
