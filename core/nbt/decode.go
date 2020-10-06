package nbt

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"errors"
	"fmt"
	"github.com/paulhobbel/performcraft/pkg/util"
	"io"
	"math"
	"reflect"
)

func Unmarshal(data []byte, v interface{}) error {
	return NewDecoder(bytes.NewReader(data)).Unmarshal(v)
}

func (u *Decoder) Unmarshal(v interface{}) error {
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Ptr {
		return errors.New("nbt: non-pointer passed to Unmarshal")
	}

	if err := u.initCompression(); err != nil {
		return fmt.Errorf("nbt: failed initiating decoder: %w", err)
	}

	tagType, tagName, err := u.readTag()
	if err != nil {
		return fmt.Errorf("nbt: failed reading root tag: %w", err)
	}

	err = u.unmarshalTag(tagType, tagName, val.Elem())
	if err != nil {
		return fmt.Errorf("nbt: failed to decode tag %q: %w", tagName, err)
	}
	return nil
}

func (u *Decoder) initCompression() error {
	compressionHead, err := u.reader.Peek(1)
	if err != nil {
		return fmt.Errorf("failed reading nbt head: %w", err)
	}

	switch compressionHead[0] {
	case 0x1f: // gzip
		gzipReader, err := gzip.NewReader(u.reader)
		if err != nil {
			return fmt.Errorf("failed initializing gzip reader: %w", err)
		}
		u.reader = bufio.NewReader(gzipReader)

	case 0x78: // zlib
		zlibReader, err := zlib.NewReader(u.reader)
		if err != nil {
			return fmt.Errorf("failed initializing zlib reader: %w", err)
		}

		u.reader = bufio.NewReader(zlibReader)
	}

	return nil
}

func (u *Decoder) unmarshalTag(tagType byte, tagName string, val reflect.Value) error {
	switch tagType {
	case TagEnd:
		return errors.New("end...")
	case TagByte:
		value, err := u.reader.ReadByte()
		if err != nil {
			return err
		}

		switch kind := val.Kind(); kind {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			val.SetInt(int64(value))
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			val.SetUint(uint64(value))
		case reflect.Interface:
			val.Set(reflect.ValueOf(value))
		default:
			return errors.New("cannot parse TagByte as " + kind.String())
		}

	case TagShort:
		value, err := u.readInt16()
		if err != nil {
			return err
		}

		switch kind := val.Kind(); kind {
		case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64:
			val.SetInt(int64(value))
		case reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			val.SetUint(uint64(value))
		case reflect.Interface:
			val.Set(reflect.ValueOf(value))
		default:
			return errors.New("cannot parse TagShort as " + kind.String())
		}

	case TagInt:
		value, err := u.readInt32()
		if err != nil {
			return err
		}

		switch kind := val.Kind(); kind {
		case reflect.Int, reflect.Int32, reflect.Int64:
			val.SetInt(int64(value))
		case reflect.Uint, reflect.Uint32, reflect.Uint64:
			val.SetUint(uint64(value))
		case reflect.Interface:
			val.Set(reflect.ValueOf(value))
		default:
			return errors.New("cannot parse TagInt as " + kind.String())
		}

	case TagFloat:
		valueInt, err := u.readInt32()
		if err != nil {
			return err
		}

		value := math.Float32frombits(uint32(valueInt))

		switch kind := val.Kind(); kind {
		case reflect.Float64:
			val.Set(reflect.ValueOf(float64(value)))
		case reflect.Float32, reflect.Interface:
			val.Set(reflect.ValueOf(value))
		default:
			return errors.New("cannot parse TagFloat as " + kind.String())
		}

	case TagLong:
		value, err := u.readInt64()
		if err != nil {
			return err
		}

		switch kind := val.Kind(); kind {
		case reflect.Int, reflect.Int64:
			val.SetInt(value)
		case reflect.Uint, reflect.Uint64:
			val.SetUint(uint64(value))
		case reflect.Interface:
			val.Set(reflect.ValueOf(value))
		default:
			return errors.New("cannot parse TagLong as " + kind.String())
		}

	case TagDouble:
		valueInt, err := u.readInt64()
		if err != nil {
			return err
		}

		value := math.Float64frombits(uint64(valueInt))

		switch kind := val.Kind(); kind {
		case reflect.Float64, reflect.Interface:
			val.Set(reflect.ValueOf(value))
		default:
			return errors.New("cannot parse TagDouble as " + kind.String())
		}

	case TagString:
		s, err := u.readString()
		if err != nil {
			return err
		}

		switch kind := val.Kind(); kind {
		case reflect.String:
			val.SetString(s)
		case reflect.Interface:
			val.Set(reflect.ValueOf(s))
		default:
			return errors.New("cannot parse TagString as " + kind.String())
		}

	case TagByteArray:
		length, err := u.readInt32()
		if err != nil {
			return err
		}

		if length < 0 {
			return errors.New("byte array length less than 0")
		}

		ba := make([]byte, length)
		if _, err := io.ReadFull(u.reader, ba); err != nil {
			return err
		}

		valType := val.Type()

		if valType == reflect.TypeOf(ba) {
			val.SetBytes(ba)
		} else if valType.Kind() == reflect.Interface {
			val.Set(reflect.ValueOf(ba))
		} else {
			return errors.New("cannot parse TagByteArray to " + valType.String() + " use []byte instead")
		}

	case TagIntArray:
		length, err := u.readInt32()
		if err != nil {
			return err
		}

		if length < 0 {
			return errors.New("int array length less than 0")
		}

		valType := val.Type()
		if valType.Kind() == reflect.Interface {
			valType = reflect.TypeOf([]int32{})
		} else if valType.Kind() != reflect.Slice {
			return errors.New("cannot parse TagIntArray to " + valType.String() + " it must be a int/int32 slice")
		} else if keyType := valType.Elem().Kind(); keyType != reflect.Int && keyType != reflect.Int32 {
			return errors.New("cannot parse TagIntArray to " + valType.String() + " it must be a int/int32 slice")
		}

		// Create new slice
		buf := reflect.MakeSlice(valType, int(length), int(length))

		// Fill slice
		for i := 0; i < int(length); i++ {
			arrVal, err := u.readInt32()
			if err != nil {
				return err
			}

			buf.Index(i).SetInt(int64(arrVal))
		}

		val.Set(buf)

	case TagLongArray:
		length, err := u.readInt32()
		if err != nil {
			return err
		}

		if length < 0 {
			return errors.New("long array length less than 0")
		}

		valType := val.Type()
		if valType.Kind() == reflect.Interface {
			valType = reflect.TypeOf([]int64{})
		} else if valType.Kind() != reflect.Slice {
			return errors.New("cannot parse TagIntArray to " + valType.String() + " it must be a int/int64 slice")
		} else if keyType := valType.Elem().Kind(); keyType != reflect.Int && keyType != reflect.Int64 {
			return errors.New("cannot parse TagIntArray to " + valType.String() + " it must be a int/int64 slice")
		}

		// Create new slice
		buf := reflect.MakeSlice(valType, int(length), int(length))

		// Fill slice
		for i := 0; i < int(length); i++ {
			arrVal, err := u.readInt64()
			if err != nil {
				return err
			}

			buf.Index(i).SetInt(arrVal)
		}

		val.Set(buf)

	case TagList:
		// Nameless type, so can't use readTag()
		listType, err := u.reader.ReadByte()
		if err != nil {
			return err
		}

		listLength, err := u.readInt32()
		if err != nil {
			return err
		}

		if listLength < 0 {
			return errors.New("list length less than 0")
		}

		var listVal reflect.Value
		valKind := val.Kind()

		switch valKind {
		case reflect.Interface:
			listVal = reflect.ValueOf(make([]interface{}, listLength))
		case reflect.Slice:
			listVal = reflect.MakeSlice(val.Type(), int(listLength), int(listLength))
		case reflect.Array:
			// Check if length is correct!
			if valLength := val.Len(); valLength < int(listLength) {
				return fmt.Errorf("TagList %s has length %d, however array %v only has length %d", tagName, listLength, val.Type(), valLength)
			}

			listVal = val
		default:
			return errors.New("cannot parse TagList to " + valKind.String())
		}

		for i := 0; i < int(listLength); i++ {
			if err := u.unmarshalTag(listType, "", listVal.Index(i)); err != nil {
				return err
			}
		}

		// We only create a value for other types than reflect.Array so
		if valKind != reflect.Array {
			val.Set(listVal)
		}

	case TagCompound:
		switch valKind := val.Kind(); valKind {
		case reflect.Map:
			if val.Type().Key().Kind() != reflect.String {
				return errors.New("cannot parse TagCompound as " + val.Type().String() + " expected key to be of type string")
			}

			if val.IsNil() {
				val.Set(reflect.MakeMap(val.Type()))
			}

			// Read tags inside of compound till TagEnd
			for {
				listTagType, listTagName, err := u.readTag()
				if err != nil {
					return err
				}

				if listTagType == TagEnd {
					break
				}

				tagVal := reflect.New(val.Type().Elem())
				if err = u.unmarshalTag(listTagType, listTagName, tagVal.Elem()); err != nil {
					return fmt.Errorf("failed to decode tag %q: %w", listTagName, err)
				}

				val.SetMapIndex(reflect.ValueOf(listTagName), tagVal.Elem())
			}

		case reflect.Interface:
			compoundVal := make(map[string]interface{})

			// Read tags inside of compound till TagEnd
			for {
				listTagType, listTagName, err := u.readTag()
				if err != nil {
					return err
				}

				if listTagType == TagEnd {
					break
				}

				var tagValue interface{}
				if err = u.unmarshalTag(listTagType, listTagName, reflect.ValueOf(&tagValue).Elem()); err != nil {
					return fmt.Errorf("failed to decode tag %q: %w", listTagName, err)
				}

				compoundVal[listTagName] = tagValue
			}

			val.Set(reflect.ValueOf(compoundVal))

		case reflect.Struct:
			structInfo := util.GetStructTagInfo(val.Type(), "nbt")

			// Read tags inside of compound till TagEnd
			for {
				listTagType, listTagName, err := u.readTag()
				if err != nil {
					return err
				}

				if listTagType == TagEnd {
					break
				}

				fieldIdx := structInfo.FindIndexByName(tagName)
				if fieldIdx != -1 {
					err = u.unmarshalTag(listTagType, listTagName, val.Field(fieldIdx))
					if err != nil {
						return fmt.Errorf("failed to decode tag %q: %w", listTagName, err)
					}
				} else {
					// TODO: Skip tags we don't care about
				}
			}
		}

	default:
		return fmt.Errorf("unknown Tag 0x%02x", tagType)
	}

	return nil
}

func (u *Decoder) readTag() (tagType byte, tagName string, err error) {
	tagType, err = u.reader.ReadByte()
	if err != nil {
		return
	}

	// Read the tag name
	if tagType != TagEnd {
		tagName, err = u.readString()
	}
	return
}

// readInt16 read a single signed, big endian 16 bit integer
func (u *Decoder) readInt16() (int16, error) {
	var data [2]byte
	_, err := io.ReadFull(u.reader, data[:])

	return int16(data[0])<<8 | int16(data[1]), err
}

// readInt32 read a single signed, big endian 32 bit integer
func (u *Decoder) readInt32() (int32, error) {
	var data [4]byte
	_, err := io.ReadFull(u.reader, data[:])

	return int32(data[0])<<24 | int32(data[1])<<16 |
		int32(data[2])<<8 | int32(data[3]), err
}

// readInt64 read a single signed, big endian 64 bit integer
func (u *Decoder) readInt64() (int64, error) {
	var data [8]byte
	_, err := io.ReadFull(u.reader, data[:])

	return int64(data[0])<<56 | int64(data[1])<<48 |
		int64(data[2])<<40 | int64(data[3])<<32 |
		int64(data[4])<<24 | int64(data[5])<<16 |
		int64(data[6])<<8 | int64(data[7]), err
}

// readString read a length prefixed UTF-8 string.
func (u *Decoder) readString() (string, error) {
	length, err := u.readInt16()
	if err != nil {
		return "", err
	}

	if length < 0 {
		return "", errors.New("string length less than 0")
	}

	if length == 0 {
		return "", nil
	}

	buf := make([]byte, length)
	_, err = io.ReadFull(u.reader, buf)
	return string(buf), err
}
