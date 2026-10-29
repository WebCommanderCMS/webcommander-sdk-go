package http

import (
	"encoding/json"
	"fmt"
	"maps"
	"reflect"
	"slices"
	"time"

	wcerrors "github.com/WebCommanderCMS/webcommander-sdk-go/v4/errors"
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/redact"
)

// SortDirection is a query-parameter sort direction; it serialises to its wire value.
type SortDirection string

// The two sort directions.
const (
	SortAsc  SortDirection = "asc"
	SortDesc SortDirection = "desc"
)

// BuildQueryParams normalises every parameter to its query-string text: bools, ISO-8601 instants, repeated slices.
func BuildQueryParams(params map[string]any) (map[string]any, error) {
	out := make(map[string]any, len(params))
	for _, key := range slices.Sorted(maps.Keys(params)) {
		value := params[key]
		if value == nil {
			continue
		}
		normalised, err := normaliseQueryValue(key, value, 0, map[uintptr]bool{})
		if err != nil {
			return nil, err
		}
		out[key] = normalised
	}
	return out, nil
}

// normaliseQueryValue normalises one value; seen names a cyclic slice rather than tripping the depth limit.
func normaliseQueryValue(key string, value any, depth int, seen map[uintptr]bool) (any, error) {
	if depth >= redact.MaxRedactionDepth {
		return nil, wcerrors.NewUsageError(fmt.Sprintf(
			"query param '%s': nested more than %d levels deep, which is too deep to "+
				"encode in a query string", key, redact.MaxRedactionDepth))
	}
	switch v := value.(type) {
	case bool:
		if v {
			return "true", nil
		}
		return "false", nil
	case string:
		return v, nil
	case json.Number:
		return v, nil // kept as a number; stringified at the wire
	case SortDirection:
		return string(v), nil
	case time.Time:
		return isoTime(v), nil
	case []byte:
		return string(v), nil
	case int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64,
		float32, float64:
		return v, nil
	}

	reflectValue := reflect.ValueOf(value)
	switch reflectValue.Kind() {
	case reflect.Map:
		return nil, wcerrors.NewUsageError(fmt.Sprintf(
			"query param '%s': nested objects cannot be encoded in a query string; flatten it "+
				"or send it in the body", key))
	case reflect.Slice, reflect.Array:
		// Self-reference before the depth rule: a cyclic []any is its own diagnosis.
		identity := uintptr(0)
		if reflectValue.Kind() == reflect.Slice {
			identity = reflectValue.Pointer()
			if seen[identity] {
				return nil, wcerrors.NewUsageError(fmt.Sprintf(
					"query param '%s': the value refers to itself, so it has no finite "+
						"query-string encoding", key))
			}
		}
		if depth >= 1 {
			return nil, wcerrors.NewUsageError(fmt.Sprintf(
				"query param '%s': a sequence inside a sequence has no query-string encoding; "+
					"flatten it to one list of scalars", key))
		}
		if identity != 0 {
			seen[identity] = true
		}
		list := make([]any, 0, reflectValue.Len())
		for i := 0; i < reflectValue.Len(); i++ {
			item := reflectValue.Index(i).Interface()
			if item == nil {
				continue
			}
			normalised, err := normaliseQueryValue(key, item, depth+1, seen)
			if err != nil {
				return nil, err
			}
			list = append(list, normalised)
		}
		return list, nil
	}
	return fmt.Sprint(value), nil
}

// isoTime always carries the UTC offset, with a fractional part only when the instant has one.
func isoTime(t time.Time) string {
	fraction := ""
	if nanoseconds := t.Nanosecond(); nanoseconds != 0 {
		if nanoseconds%1000 == 0 {
			fraction = fmt.Sprintf(".%06d", nanoseconds/1000)
		} else {
			fraction = fmt.Sprintf(".%09d", nanoseconds)
		}
	}
	return t.Format("2006-01-02T15:04:05") + fraction + t.Format("-07:00")
}
