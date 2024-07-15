package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log/slog"
	"os"
	"reflect"
	"strconv"
	"strings"
)

type csvUtil struct {
	rawData          [][]string
	headers          []string
	tagMap           map[string]int
	fieldMap         map[string][2]string
	isSliceOfPointer bool
}

func LoadFile[T any](csvFile string, csvTag string) ([]T, error) {
	items := make([]T, 0)
	util := csvUtil{}
	err := util.load(csvFile, &items, csvTag)
	return items, err
}

func (c *csvUtil) load(csvFile string, slicePtr interface{}, csvTag string) error {
	f, err := os.Open(csvFile)
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			slog.Warn("%s close err: %v", csvFile, err)
		}
	}(f)
	if err != nil {
		return err
	}

	// 类型检查
	if !isSlicePtr(slicePtr) {
		return fmt.Errorf("slicePtr is not a slice pointer")
	}

	// 根据 csvTag 获取结构体的字段和 CSV 列名的映射
	if len(csvTag) == 0 {
		csvTag = "csv"
	}
	c.isSliceOfPointer = isPointerSlice(slicePtr)
	c.parseTypes(slicePtr, csvTag)

	// csv reader
	r := csv.NewReader(f)

	// read headers
	if c.headers, err = r.Read(); err != nil {
		return err
	}

	// read data
	c.rawData, err = r.ReadAll()
	if err != nil {
		return err
	}

	// 反射获取 slicePtr 的类型和值
	rt := reflect.TypeOf(slicePtr).Elem()
	rv := reflect.ValueOf(slicePtr).Elem()
	for _, record := range c.rawData {
		var newInstance reflect.Value
		if c.isSliceOfPointer {
			newInstance = reflect.New(rt.Elem().Elem())
		} else {
			newInstance = reflect.New(rt.Elem())
		}

		for headerIndex, headerValue := range c.headers {
			if index, ok := c.tagMap[headerValue]; ok {
				fieldValue := newInstance.Elem().Field(index)
				switch fieldValue.Kind() {
				case reflect.Int, reflect.Int64:
					var value int64
					if fieldValue.Kind() == reflect.Int64 {
						value, err = strconv.ParseInt(record[headerIndex], 10, 64)
					} else {
						var v int
						v, err = strconv.Atoi(record[headerIndex])
						value = int64(v)
					}
					if err != nil {
						return err
					}
					fieldValue.SetInt(value)
				case reflect.Float32:
					value, err := strconv.ParseFloat(record[headerIndex], 32)
					if err != nil {
						return err
					}
					fieldValue.SetFloat(value)
				case reflect.Float64:
					value, err := strconv.ParseFloat(record[headerIndex], 64)
					if err != nil {
						return err
					}
					fieldValue.SetFloat(value)
				case reflect.Slice:
					// handle types: []int,[]int64,[]float32,[]float64,[]string
					fType := c.fieldMap[headerValue][1]
					if len(record[headerIndex]) > 0 {
						fValues := strings.Split(record[headerIndex], ",")
						field := reflect.Indirect(newInstance).FieldByName(headerValue)
						for _, v := range fValues {
							fv, err := c.getFieldValue(fType, v)
							if err != nil {
								return err
							}
							field = reflect.Append(field, reflect.ValueOf(fv))
						}
						reflect.Indirect(newInstance).FieldByName(headerValue).Set(field)
					}
				case reflect.String:
					fieldValue.SetString(record[headerIndex])
				default:
					// unhandled types
					fmt.Println("unhandled type: ", fieldValue.Kind().String())
				}
			}
		}

		// set slice value
		if c.isSliceOfPointer {
			rv.Set(reflect.Append(rv, newInstance))
		} else {
			rv.Set(reflect.Append(rv, newInstance.Elem()))
		}
	}

	rv.Slice(0, rv.Len())
	return nil
}

func (c *csvUtil) getFieldValue(fieldType, fieldValue string) (any, error) {
	switch fieldType {
	case "[]int":
		return strconv.Atoi(fieldValue)
	case "[]int64":
		return strconv.ParseInt(fieldValue, 10, 64)
	case "[]float32":
		return strconv.ParseFloat(fieldValue, 32)
	case "[]float64":
		return strconv.ParseFloat(fieldValue, 64)
	case "[]string":
		return fieldValue, nil
	}

	return nil, fmt.Errorf("unhandled types: %s", fieldType)
}

func (c *csvUtil) parseTypes(slicePtr interface{}, csvTag string) {
	var rt reflect.Type // slice element type
	if c.isSliceOfPointer {
		rt = reflect.TypeOf(slicePtr).Elem().Elem().Elem()
	} else {
		rt = reflect.TypeOf(slicePtr).Elem().Elem()
	}

	c.tagMap = make(map[string]int, rt.NumField())
	c.fieldMap = make(map[string][2]string, rt.NumField())
	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		if tag := field.Tag.Get(csvTag); tag != "" {
			c.tagMap[tag] = i

			data := [2]string{}
			data[0] = field.Name
			data[1] = fmt.Sprintf("%v", field.Type)
			c.fieldMap[tag] = data
		}
	}
}

func isPointerSlice(slicePointer interface{}) bool {
	rt := reflect.TypeOf(slicePointer).Elem()
	return rt.Kind() == reflect.Slice && rt.Elem().Kind() == reflect.Ptr
}

// hasBOM 检查文件是否包含 UTF-8 BOM
func hasBOM(filePath string) (bool, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return false, err
	}
	defer func(file *os.File) {
		_ = file.Close() // ignore error
	}(f)

	// 创建一个带缓冲的读取器
	reader := bufio.NewReader(f)

	// 读取前三个字符
	b1, err := reader.Peek(1)
	if err != nil {
		return false, err
	}

	b2, err := reader.Peek(2)
	if err != nil {
		return false, err
	}

	b3, err := reader.Peek(3)
	if err != nil {
		return false, err
	}

	// 检查是否是 UTF-8 BOM
	return b1[0] == 0xEF && b2[1] == 0xBB && b3[2] == 0xBF, nil
}

func isSlicePtr(ptr interface{}) bool {
	return reflect.TypeOf(ptr).Elem().Kind() == reflect.Slice
}
