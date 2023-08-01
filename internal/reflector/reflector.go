package reflector

import (
	"fmt"
	"path"
	"reflect"
	"strings"
)

var ErrNoName = fmt.Errorf("reflector: no name")

func TypeOf[V any](v V) (string, error) {
	t := reflect.TypeOf(v)
	if t == nil {
		t = reflect.TypeOf(new(V)).Elem()
	}
	return Type(t)
}

func Type(t reflect.Type) (string, error) {
	prefix, inner := innermost("", t)
	if inner.Name() == "" {
		return "", fmt.Errorf("%w: %s", ErrNoName, t)
	}
	return toString(prefix, inner), nil
}

func innermost(prefix string, t reflect.Type) (string, reflect.Type) {
	switch t.Kind() {
	case reflect.Ptr:
		return innermost(prefix+"*", t.Elem())
	case reflect.Slice:
		return innermost(prefix+"[]", t.Elem())
	case reflect.Map:
		key, innerKey := innermost("", t.Key())
		key = toString(key, innerKey)
		return innermost(prefix+"map["+key+"]", t.Elem())
	default:
		return prefix, t
	}
}

func toString(prefix string, t reflect.Type) string {
	typeName := t.String()
	typeParts := strings.SplitN(typeName, ".", 2)
	dir := ""
	if len(typeParts) == 2 {
		dir = typeParts[0]
		typeName = typeParts[1]
	}
	pkgPath := t.PkgPath()
	if pkgPath == "" {
		return prefix + typeName
	} else if strings.HasPrefix(pkgPath, "/") {
		return path.Dir(pkgPath) + "/" + dir + "." + prefix + typeName
	}
	return pkgPath + "." + prefix + typeName
}
