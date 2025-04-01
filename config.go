package goconfig

import (
	"encoding"
	"encoding/json"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
)

const (
	// DefaultSep is the default separator used for environment variable names
	DefaultSep string = "_"
	// DefaultArrSep is the default separator used for array values
	DefaultArrSep string = ","
)

// Load loads configuration from environment variables into the provided struct.
// It uses default options and is a convenience wrapper around New().Load().
// The struct should be a pointer to a struct with fields tagged with "env" or "alias" tags.
func Load(s any, options ...Option) error {
	c := New(options...)
	return c.Load(s)
}

// New creates a new configuration loader with the provided options.
// If no options are provided, default values will be used.
func New(options ...Option) *Loader {
	c := &Loader{
		sep:                  DefaultSep,
		prefix:               "",
		arraySep:             DefaultArrSep,
		fieldNameTransformer: UperCaseTransformer,
	}

	for _, opt := range options {
		opt(c)
	}

	return c
}

// Loader handles the loading of configuration from environment variables into structs.
// It provides customization options for prefix, separators, and field name transformation.
type Loader struct {
	prefix               string
	sep                  string
	arraySep             string
	fieldNameTransformer func(name string) string
}

// Load loads environment variables into the provided struct.
// The struct should be a pointer to a struct with fields tagged with "env" or "alias" tags.
// Returns an error if the loading process fails.
func (c *Loader) Load(s any) error {
	_, err := c.recursiveLoadToStruct(s, nil)
	return err
}

// nolint:gocyclo
func (c *Loader) recursiveLoadToStruct(s any, prefix []string) (found bool, err error) {
	vPtr := reflect.ValueOf(s)

	if vPtr.Kind() != reflect.Ptr {
		return false, errors.Errorf("should be a pointer to %T", s)
	}

	defer func() {
		if p := recover(); p != nil {
			err = errors.Errorf(
				"cannot load to struct %T (prefix=%s). panic: %v",
				s, strings.Join(prefix, c.sep), p,
			)
		}
	}()

	return c.loopOverFields(
		c.getDirectType(reflect.TypeOf(s)),
		vPtr.Elem(),
		prefix,
	)
}

func (c *Loader) loopOverFields(
	t reflect.Type,
	v reflect.Value,
	prefix []string,
) (bool, error) {
	n := v.NumField()
	found := false

	for i := 0; i < n; i++ {
		foundField, err := c.loadToField(
			t.Field(i),
			v.Field(i),
			prefix,
		)
		if err != nil {
			return false, err
		}

		if foundField {
			found = true
		}
	}

	return found, nil
}

func (c *Loader) loadToField(
	tf reflect.StructField,
	vf reflect.Value,
	prefix []string,
) (found bool, err error) {
	if !vf.CanSet() {
		return false, nil
	}

	t := c.getDirectType(tf.Type)
	envKey, nPrefix := c.buildEnvKey(tf, prefix)
	envVal, exist := os.LookupEnv(envKey)

	defer func() {
		if p := recover(); p != nil {
			err = errors.Errorf(
				"cannot load to %s field (prefix=%s). panic: %v",
				t.String(),
				strings.Join(nPrefix, c.sep),
				p,
			)
		}
	}()

	if exist {
		set, err1 := c.setFieldVal(
			vf,
			envVal,
		)
		if err1 != nil {
			return false, errors.Wrapf(err1, "cannot set field %s value", envKey)
		}

		if set {
			return true, nil
		}
	}

	if c.isStruct(t.Kind()) {
		return c.setStructVal(vf, nPrefix)
	}

	return false, nil
}

func (c *Loader) getFieldName(tf reflect.StructField) (name string, exactly bool) {
	if tag, ok := tf.Tag.Lookup("env"); ok {
		return tag, true
	}

	// use alias name instead of field name
	if tag, ok := tf.Tag.Lookup("alias"); ok {
		return tag, false
	}

	name = tf.Name
	if c.fieldNameTransformer != nil {
		name = c.fieldNameTransformer(name)
	}

	return name, false
}

func (c *Loader) buildEnvKey(tf reflect.StructField, parentKeys []string) (string, []string) {
	joinKeys := func(names ...string) string {
		arr := []string{}

		if c.prefix != "" {
			arr = append(arr, c.prefix)
		}

		for _, name := range names {
			// skip empty name
			if name == "" {
				continue
			}

			arr = append(arr, name)
		}

		return strings.Join(arr, c.sep)
	}

	if tf.Anonymous {
		return joinKeys(parentKeys...), parentKeys
	}

	name, exactly := c.getFieldName(tf)
	parentKeys = append(parentKeys, name)

	if !exactly {
		return joinKeys(parentKeys...), parentKeys
	}

	return name, parentKeys
}

func (*Loader) isTextUnmarshaler(fval reflect.Value) (encoding.TextUnmarshaler, bool) {
	if fval.Kind() != reflect.Ptr && !fval.CanAddr() {
		return nil, false
	}

	if fval.Kind() != reflect.Ptr {
		fval = fval.Addr()
	}

	if fval.Type().NumMethod() == 0 {
		return nil, false
	}

	if !fval.CanInterface() {
		return nil, false
	}

	if u, ok := fval.Interface().(encoding.TextUnmarshaler); ok {
		return u, true
	}

	return nil, false
}

func (c *Loader) setFieldVal(fval reflect.Value, envVal string) (set bool, err error) {
	if !fval.CanSet() {
		return false, errors.Errorf("%s field is cannot be set", fval.Type().Name())
	}

	if fval.Kind() == reflect.Pointer && fval.IsNil() {
		fval.Set(reflect.New(fval.Type().Elem()))
	}

	fval = c.getDirectVal(fval)
	kind := fval.Kind()

	if v, ok := c.isTextUnmarshaler(fval); ok {
		return true, v.UnmarshalText([]byte(envVal))
	}

	switch {
	case c.isString(kind):
		fval.SetString(envVal)
		return true, nil
	case c.isBool(kind):
		return true, c.setBoolVal(fval, envVal)
	case c.isDuration(fval):
		return true, c.setDurationVal(fval, envVal)
	case c.isInt(kind):
		return true, c.setIntVal(fval, envVal)
	case c.isUint(kind):
		return true, c.setUintVal(fval, envVal)
	case c.isFloat(kind):
		return true, c.setFloatVal(fval, envVal)
	case c.isSliceField(kind):
		return true, c.setSliceValue(fval, envVal)
	case c.isMap(kind):
		return true, c.setMapVal(fval, envVal)
	case c.isStruct(kind):
		return false, nil
	default:
		return false, errors.Errorf("unsupported %s", fval.Type().Name())
	}
}

func (*Loader) isSliceField(kind reflect.Kind) bool {
	return kind == reflect.Slice
}

func (c *Loader) setSliceValue(vf reflect.Value, evnVal string) error {
	var err error

	parts := strings.Split(evnVal, c.arraySep)
	if len(parts) == 0 {
		return nil
	}

	slice := reflect.MakeSlice(vf.Type(), len(parts), len(parts))

	for i := range parts {
		v := slice.Index(i)

		if _, err = c.setFieldVal(v, parts[i]); err != nil {
			return errors.Wrapf(err, "cannot set slice value")
		}
	}

	vf.Set(slice)

	return nil
}

func (*Loader) setDurationVal(vf reflect.Value, envVal string) error {
	d, err := time.ParseDuration(envVal)
	if err != nil {
		return err
	}

	vf.Set(reflect.ValueOf(d))

	return nil
}

func (*Loader) isDuration(vf reflect.Value) bool {
	return vf.Type().AssignableTo(reflect.TypeOf(time.Duration(0)))
}

func (*Loader) isInt(kind reflect.Kind) bool {
	switch kind {
	case reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64:
		return true
	default:
		return false
	}
}

func (*Loader) setIntVal(vf reflect.Value, raw string) error {
	i, err := strconv.ParseInt(raw, 10, 64)
	if err != nil {
		return err
	}

	vf.SetInt(i)

	return nil
}

func (*Loader) isUint(kind reflect.Kind) bool {
	switch kind {
	case reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64,
		reflect.Uintptr:
		return true
	default:
		return false
	}
}

func (*Loader) setUintVal(vf reflect.Value, raw string) error {
	i, err := strconv.ParseUint(raw, 10, 64)
	if err != nil {
		return err
	}

	vf.SetUint(i)

	return nil
}

func (*Loader) isFloat(kind reflect.Kind) bool {
	switch kind {
	case reflect.Float32, reflect.Float64:
		return true
	default:
		return false
	}
}

func (*Loader) setFloatVal(vf reflect.Value, raw string) error {
	num, err := strconv.ParseFloat(raw, 64)
	if err != nil {
		return err
	}

	vf.SetFloat(num)

	return nil
}

func (*Loader) isBool(kind reflect.Kind) bool {
	return kind == reflect.Bool
}

func (*Loader) setBoolVal(vf reflect.Value, envVal string) error {
	v, err := strconv.ParseBool(envVal)
	if err != nil {
		return err
	}

	vf.SetBool(v)

	return nil
}

func (*Loader) isString(kind reflect.Kind) bool {
	return kind == reflect.String
}

func (*Loader) isStruct(kind reflect.Kind) bool {
	return kind == reflect.Struct
}

func (*Loader) isMap(kind reflect.Kind) bool {
	return kind == reflect.Map
}

func (*Loader) setMapVal(vf reflect.Value, jsonStr string) error {
	return json.Unmarshal([]byte(jsonStr), vf.Addr().Interface())
}

func (*Loader) getDirectType(t reflect.Type) reflect.Type {
	realType := t
	for realType.Kind() == reflect.Pointer {
		realType = realType.Elem()
	}

	return realType
}

func (*Loader) getDirectVal(v reflect.Value) reflect.Value {
	realVal := v

	for realVal.Kind() == reflect.Pointer {
		realVal = realVal.Elem()
	}

	return realVal
}

func (c *Loader) setStructVal(vf reflect.Value, prefix []string) (found bool, err error) {
	newVf := vf
	needSet := false

	if vf.Kind() == reflect.Pointer && vf.IsNil() {
		newVf = reflect.New(vf.Type().Elem())
		needSet = true
	} else {
		newVf = vf.Addr()
	}

	found, err = c.recursiveLoadToStruct(
		newVf.Interface(),
		prefix,
	)
	if err != nil {
		return false, err
	}

	if found && needSet {
		vf.Set(newVf)
	}

	return found, nil
}
