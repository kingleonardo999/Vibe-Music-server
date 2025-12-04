package util

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"time"
	"vibe-music-server/internal/pkg/cache"
)

func GetCache[T any](key string, data *T) bool {
	//cached, err := cache.Get(key)
	//if err == nil && cached != "" && json.Unmarshal([]byte(cached), data) == nil {
	//	return true
	//}
	return false
}

func SetCache[T any](key string, data T) {
	marshaled, err := json.Marshal(data)
	if err == nil {
		_ = cache.Set(key, string(marshaled))
	}
}

func DeleteCacheByPattern(pattern string) {
	ctx, cancel := context.WithTimeout(cache.Cache().Context(), 3*time.Second)
	defer cancel()
	iter := cache.Cache().Scan(ctx, 0, pattern, 0).Iterator()
	for iter.Next(cache.Cache().Context()) {
		_ = cache.Del(iter.Val())
	}
}

func GenKeyByPattern(pattern string, args ...any) string {
	for i, arg := range args {
		var s string
		if arg == nil {
			s = ""
		} else {
			rv := reflect.ValueOf(arg)
			if !rv.IsValid() {
				s = ""
			} else {
				// unwrap pointer chain safely
				for rv.Kind() == reflect.Ptr {
					if rv.IsNil() {
						rv = reflect.Value{}
						break
					}
					rv = rv.Elem()
				}
				if !rv.IsValid() {
					s = ""
				} else {
					switch rv.Kind() {
					case reflect.Slice, reflect.Map, reflect.Chan, reflect.Func, reflect.Interface:
						if rv.IsNil() {
							s = ""
						} else {
							s = fmt.Sprintf("%v", rv.Interface())
						}
					default:
						s = fmt.Sprintf("%v", rv.Interface())
					}
				}
			}
		}
		if i == 0 {
			pattern = fmt.Sprintf("%s:%s", pattern, s)
		} else {
			pattern = fmt.Sprintf("%s-%s", pattern, s)
		}
	}
	return pattern
}
