// Package dto is the typed view over a body — what every DTO in this SDK is built from.
package dto

import (
	"encoding/json"
	"maps"
	"reflect"
	"strings"
	"sync"
	"unicode"
)

// Enveloped is implemented by a type the API wraps in a single key.
type Enveloped interface {
	EnvelopeKey() string
}

// BesideEnvelope is implemented by a payload with fields that sit outside its envelope; TopLevel returns their wire names.
type BesideEnvelope interface {
	TopLevel() []string
}

// Fields is the wire names a type declares, in declaration order.
func Fields(value any) []string {
	subject := reflect.TypeOf(value)
	for subject != nil && subject.Kind() == reflect.Pointer {
		subject = subject.Elem()
	}
	if subject == nil || subject.Kind() != reflect.Struct {
		return nil
	}
	declared := fieldsOf(subject)
	names := make([]string, len(declared))
	for at, field := range declared {
		names[at] = field.name
	}
	return names
}

// Read builds T from a decoded body.
func Read[T any](body any) any {
	var target T
	subject := reflect.ValueOf(&target).Elem()
	if subject.Kind() != reflect.Struct || !readInto(subject, body) {
		return nil
	}
	return &target
}

// ToMap is the body this payload stands for, envelope included.
func ToMap(value any) map[string]any {
	subject := reflect.ValueOf(value)
	for subject.Kind() == reflect.Pointer {
		if subject.IsNil() {
			return nil
		}
		subject = subject.Elem()
	}
	if subject.Kind() != reflect.Struct {
		return nil
	}

	record := writeRecord(subject)
	enveloped, wrapped := subject.Interface().(Enveloped)
	if !wrapped {
		return record
	}
	body := map[string]any{enveloped.EnvelopeKey(): record}
	// Moved out rather than never written in, so one walk over the fields serves both.
	if beside, hasOutside := subject.Interface().(BesideEnvelope); hasOutside {
		for _, name := range beside.TopLevel() {
			if held, set := record[name]; set {
				delete(record, name)
				body[name] = held
			}
		}
	}
	return body
}

// readInto reports whether the body was a mapping at all.
func readInto(target reflect.Value, body any) bool {
	source, isObject := body.(map[string]any)
	if !isObject {
		return false
	}
	if enveloped, wrapped := target.Interface().(Enveloped); wrapped {
		if inner, isObject := source[enveloped.EnvelopeKey()].(map[string]any); isObject {
			merged := make(map[string]any, len(inner)+len(source))
			maps.Copy(merged, inner)
			// Merged second, so a key spelled both inside and outside reads as the outer one.
			for key, value := range source {
				if key != enveloped.EnvelopeKey() {
					merged[key] = value
				}
			}
			source = merged
		}
	}
	for _, field := range fieldsOf(target.Type()) {
		raw, present := source[field.name]
		if !present {
			continue
		}
		if value, isDeclaredType := readValue(raw, field.holds); isDeclaredType {
			target.Field(field.index).Set(value)
		}
	}
	return true
}

// readValue reports whether the value was the declared type; one that was not leaves the zero.
func readValue(raw any, holds reflect.Type) (reflect.Value, bool) {
	if raw == nil {
		return reflect.Value{}, false
	}
	switch holds.Kind() {
	case reflect.Pointer:
		if holds.Elem().Kind() == reflect.Struct {
			nested := reflect.New(holds.Elem())
			if !readInto(nested.Elem(), raw) {
				return reflect.Value{}, false
			}
			return nested, true
		}
		scalar, isDeclaredType := readScalar(raw, holds.Elem())
		if !isDeclaredType {
			return reflect.Value{}, false
		}
		held := reflect.New(holds.Elem())
		held.Elem().Set(scalar)
		return held, true

	case reflect.Slice:
		items, isList := raw.([]any)
		if !isList {
			return reflect.Value{}, false
		}
		// Length is preserved: a mistyped element reads as its zero, so the third address stays third.
		list := reflect.MakeSlice(holds, len(items), len(items))
		for at, item := range items {
			if value, isDeclaredType := readValue(item, holds.Elem()); isDeclaredType {
				list.Index(at).Set(value)
			}
		}
		return list, true

	case reflect.Interface:
		// Handed over as it arrived: a guessed type here would read nil for a body the store legitimately sends.
		return reflect.ValueOf(raw), true
	}
	return reflect.Value{}, false
}

func readScalar(raw any, holds reflect.Type) (reflect.Value, bool) {
	switch holds.Kind() {
	case reflect.String:
		if text, isText := raw.(string); isText {
			return reflect.ValueOf(text), true
		}
	case reflect.Bool:
		if flag, isFlag := raw.(bool); isFlag {
			return reflect.ValueOf(flag), true
		}
	case reflect.Int64:
		if whole, isWhole := wholeNumber(raw); isWhole {
			return reflect.ValueOf(whole), true
		}
	case reflect.Float64:
		if amount, isNumber := fractionalNumber(raw); isNumber {
			return reflect.ValueOf(amount), true
		}
	}
	return reflect.Value{}, false
}

// wholeNumber reads a decoded JSON number as an integer.
func wholeNumber(raw any) (int64, bool) {
	switch number := raw.(type) {
	case json.Number:
		whole, err := number.Int64()
		return whole, err == nil
	case float64:
		// One Go type covers every JSON number, so only a fraction says it was not whole.
		return int64(number), number == float64(int64(number))
	case int:
		return int64(number), true
	case int64:
		return number, true
	}
	return 0, false
}

// fractionalNumber widens an integer to a float, JSON having one number type; nothing is lost.
func fractionalNumber(raw any) (float64, bool) {
	switch number := raw.(type) {
	case json.Number:
		amount, err := number.Float64()
		return amount, err == nil
	case float64:
		return number, true
	case int:
		return float64(number), true
	case int64:
		return float64(number), true
	}
	return 0, false
}

func writeRecord(subject reflect.Value) map[string]any {
	record := map[string]any{}
	for _, field := range fieldsOf(subject.Type()) {
		held := subject.Field(field.index)
		if isUnset(held) {
			continue
		}
		record[field.name] = writeValue(held)
	}
	return record
}

// writeValue writes a nested DTO as its record, without the envelope the store does not expect there.
func writeValue(held reflect.Value) any {
	switch held.Kind() {
	case reflect.Interface, reflect.Pointer:
		if held.IsNil() {
			return nil
		}
		return writeValue(held.Elem())
	case reflect.Struct:
		return writeRecord(held)
	case reflect.Slice:
		items := make([]any, held.Len())
		for at := range items {
			items[at] = writeValue(held.Index(at))
		}
		return items
	}
	return held.Interface()
}

// isUnset reports a field left out of the body.
func isUnset(held reflect.Value) bool {
	switch held.Kind() {
	case reflect.Pointer, reflect.Slice, reflect.Map, reflect.Interface:
		return held.IsNil()
	case reflect.String:
		return held.Len() == 0
	}
	return false
}

type declaration struct {
	index int
	name  string
	holds reflect.Type
}

// Resolved once per type: a listing page builds one Customer per row.
var declarations sync.Map // reflect.Type -> []declaration

func fieldsOf(structType reflect.Type) []declaration {
	if cached, resolved := declarations.Load(structType); resolved {
		return cached.([]declaration)
	}
	declared := make([]declaration, 0, structType.NumField())
	for index := 0; index < structType.NumField(); index++ {
		field := structType.Field(index)
		if !field.IsExported() {
			continue
		}
		declared = append(declared, declaration{
			index: index, name: declaredName(field), holds: field.Type,
		})
	}
	declarations.Store(structType, declared)
	return declared
}

// WireTag spells a key no mechanical rule reaches, such as `wire:"office365"`; WireName is then skipped.
const WireTag = "wire"

// declaredName is the key a field travels under: its WireTag, else the spelling of its Go name.
func declaredName(field reflect.StructField) string {
	if named := field.Tag.Get(WireTag); named != "" {
		return named
	}
	return WireName(field.Name)
}

// WireName spells a Go field name as the API does: CustomerID → customer_id, AddressLine1 → address_line_1.
func WireName(name string) string {
	letters := []rune(name)
	var spelled strings.Builder
	spelled.Grow(len(letters) + 4)
	for at, letter := range letters {
		if at > 0 && startsWord(letters, at) {
			spelled.WriteByte('_')
		}
		spelled.WriteRune(unicode.ToLower(letter))
	}
	return spelled.String()
}

func startsWord(letters []rune, at int) bool {
	if unicode.IsDigit(letters[at]) {
		// A digit run starts a word (address_line_1) unless a lowercase letter follows its first digit (b2b_configuration).
		return !unicode.IsDigit(letters[at-1]) &&
			!(at+1 < len(letters) && unicode.IsLower(letters[at+1]))
	}
	if !unicode.IsUpper(letters[at]) {
		return false
	}
	if !unicode.IsUpper(letters[at-1]) {
		return true
	}
	return at+1 < len(letters) && unicode.IsLower(letters[at+1])
}
