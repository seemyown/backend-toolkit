package ctxbinding

import (
	"fmt"
	"github.com/google/uuid"
	"reflect"
	"strconv"
	"sync"
	"time"
)

type parserFn func(any) (any, error)

var parsers = map[reflect.Type]parserFn{}
var mx sync.RWMutex

func init() {
	RegisterParser(int64(0), parseInt64)
	RegisterParser("", parseString)
	RegisterParser(true, parseBool)
	RegisterParser(float64(0), parseFloat64)
	RegisterParser(uuid.UUID{}, parseUUID)
	RegisterParser(time.Time{}, parseTime)
}

func RegisterParser(ex any, fn parserFn) {
	mx.Lock()
	defer mx.Unlock()
	typ := reflect.TypeOf(ex)
	parsers[typ] = fn
}

func getParser(typ reflect.Type) (parserFn, bool) {
	mx.RLock()
	defer mx.RUnlock()
	fn, ok := parsers[typ]
	return fn, ok
}

func parseInt64(val any) (any, error) {
	switch v := val.(type) {
	case int64:
		return v, nil
	case string:
		return strconv.ParseInt(v, 10, 64)
	default:
		return nil, fmt.Errorf("cannot parse %T to int64", val)
	}
}

func parseString(val any) (any, error) {
	switch v := val.(type) {
	case string:
		return v, nil
	case fmt.Stringer:
		return v.String(), nil
	default:
		return nil, fmt.Errorf("cannot parse %T to string", val)
	}
}

func parseBool(val any) (any, error) {
	switch v := val.(type) {
	case bool:
		return v, nil
	case string:
		return strconv.ParseBool(v)
	default:
		return nil, fmt.Errorf("cannot parse %T to bool", val)
	}
}

func parseFloat64(val any) (any, error) {
	switch v := val.(type) {
	case float64:
		return v, nil
	case string:
		return strconv.ParseFloat(v, 64)
	default:
		return nil, fmt.Errorf("cannot parse %T to float64", val)
	}
}

func parseUUID(val any) (any, error) {
	switch v := val.(type) {
	case string:
		return uuid.Parse(v)
	default:
		return nil, fmt.Errorf("cannot parse %T to uuid.UUID", val)
	}
}

func parseTime(val any) (any, error) {
	switch v := val.(type) {
	case string:
		return time.Parse(time.RFC3339, v)
	default:
		return nil, fmt.Errorf("cannot parse %T to time.Time", val)
	}
}
