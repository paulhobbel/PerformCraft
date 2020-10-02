package util

import (
	"reflect"
	"sync"
)

type StructInfo struct {
	tagName     string
	nameToIndex map[string]int
}

var structInfoCache sync.Map

func GetStructTagInfo(structType reflect.Type, tagName string) *StructInfo {
	if structInfo, ok := structInfoCache.Load(structType); ok {
		return structInfo.(*StructInfo)
	}

	structInfo := &StructInfo{
		tagName:     tagName,
		nameToIndex: make(map[string]int),
	}

	if structType.Kind() == reflect.Struct {
		for i := 0; i < structType.NumField(); i++ {
			field := structType.Field(i)
			fieldTag := field.Tag.Get(tagName)

			// Check if is private field
			if (field.PkgPath != "" && !field.Anonymous) || fieldTag == "-" {
				continue
			}

			structInfo.nameToIndex[fieldTag] = i

			// Fallback check field name
			if _, ok := structInfo.nameToIndex[field.Name]; !ok {
				structInfo.nameToIndex[field.Name] = i
			}
		}
	}

	cachedStructInfo, _ := structInfoCache.LoadOrStore(structType, structInfo)
	return cachedStructInfo.(*StructInfo)
}

func (s *StructInfo) FindIndexByName(name string) int {
	idx, ok := s.nameToIndex[name]
	if !ok {
		return -1
	}

	return idx
}
