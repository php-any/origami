package mongodb

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/php-any/origami/data"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type bsonValue struct{ value any }

// BSON documents are ordered. A root PHP array is always a document; nested
// arrays are BSON arrays only when their integer keys are exactly 0..n-1.
func encodeDocument(v data.Value) (bson.D, error) {
	x, err := encodeValue(v, true, 0)
	if err != nil {
		return nil, err
	}
	d, ok := x.(bson.D)
	if !ok {
		return nil, fmt.Errorf("expected an array or object document")
	}
	return d, nil
}
func encodeValue(v data.Value, document bool, depth int) (any, error) {
	if depth > 100 {
		return nil, fmt.Errorf("BSON nesting exceeds 100 levels (possible recursion)")
	}
	switch v := v.(type) {
	case nil, *data.NullValue:
		return nil, nil
	case *data.BoolValue:
		return v.Value, nil
	case *data.StringValue:
		return v.Value, nil
	case *data.IntValue:
		return int64(v.Value), nil
	case *data.FloatValue:
		return v.Value, nil
	case *data.ArrayValue:
		list := !document
		for i, s := range v.Range() {
			k, ok := s.PHPArrayKey(i).(*data.IntValue)
			if !ok || k.Value != i {
				list = false
				break
			}
		}
		d := bson.D{}
		a := bson.A{}
		for i, s := range v.Range() {
			x, err := encodeValue(s.ReadValue(), false, depth+1)
			if err != nil {
				return nil, err
			}
			if list {
				a = append(a, x)
			} else {
				k := s.PHPArrayKey(i).AsString()
				if strings.ContainsRune(k, 0) {
					return nil, fmt.Errorf("BSON field names cannot contain NUL")
				}
				d = append(d, bson.E{Key: k, Value: x})
			}
		}
		if list {
			return a, nil
		}
		return d, nil
	case *data.ClassValue:
		if s, ok := v.InstanceSource.(*bsonValue); ok {
			return s.value, nil
		}
		if strings.HasPrefix(v.GetName(), driverNS) {
			return nil, fmt.Errorf("cannot encode %s as BSON", v.GetName())
		}
		if _, ok := v.Class.GetMethod("bsonSerialize"); ok {
			return nil, fmt.Errorf("BSON Serializable is not supported yet")
		}
		d := bson.D{}
		var err error
		v.RangeProperties(func(k string, value data.Value) bool {
			if strings.ContainsRune(k, 0) {
				return true
			}
			if p, ok := v.Class.GetProperty(k); ok && p.GetModifier() != data.ModifierPublic {
				return true
			}
			var x any
			x, err = encodeValue(value, false, depth+1)
			if err != nil {
				return false
			}
			d = append(d, bson.E{Key: k, Value: x})
			return true
		})
		return d, err
	default:
		return nil, fmt.Errorf("cannot encode %T as BSON", v)
	}
}

type typeMap struct{ rootArray, documentArray bool }

func parseTypeMap(v data.Value) (typeMap, error) {
	m, err := optionMap(v, "root", "document", "array")
	if err != nil {
		return typeMap{}, err
	}
	t := typeMap{}
	for k, v := range m {
		s, err := stringArg(v, k)
		if err != nil {
			return t, err
		}
		if k == "array" {
			if s != "array" {
				return t, fmt.Errorf("only array mapping is supported for BSON arrays")
			}
			continue
		}
		if s != "array" && s != "object" && s != "stdClass" {
			return t, fmt.Errorf("unsupported type map %s=%s", k, s)
		}
		if k == "root" {
			t.rootArray = s == "array"
		} else {
			t.documentArray = s == "array"
		}
	}
	return t, nil
}
func decodeValue(ctx data.Context, x any, t typeMap, root bool) (data.Value, error) {
	switch x := x.(type) {
	case nil:
		return data.NewNullValue(), nil
	case bool:
		return data.NewBoolValue(x), nil
	case string:
		return data.NewStringValue(x), nil
	case int32:
		return data.NewIntValue(int(x)), nil
	case int64:
		return data.NewIntValue(int(x)), nil
	case float64:
		return data.NewFloatValue(x), nil
	case bson.D:
		array := t.documentArray
		if root {
			array = t.rootArray
		}
		a := data.NewArrayValue(nil).(*data.ArrayValue)
		var o *data.ClassValue
		if !array {
			o = object(ctx, "stdClass", nil)
		}
		for _, e := range x {
			v, err := decodeValue(ctx, e.Value, t, false)
			if err != nil {
				return nil, err
			}
			if array {
				a.SetKey(data.NewStringValue(e.Key), v)
			} else {
				o.SetProperty(e.Key, v)
			}
		}
		if array {
			return a, nil
		}
		return o, nil
	case bson.A:
		a := make([]data.Value, 0, len(x))
		for _, e := range x {
			v, err := decodeValue(ctx, e, t, false)
			if err != nil {
				return nil, err
			}
			a = append(a, v)
		}
		return data.NewArrayValue(a), nil
	case bson.ObjectID:
		return object(ctx, bsonNS+"ObjectId", &bsonValue{x}), nil
	case bson.DateTime:
		return object(ctx, bsonNS+"UTCDateTime", &bsonValue{x}), nil
	case bson.Binary:
		return object(ctx, bsonNS+"Binary", &bsonValue{x}), nil
	case bson.Regex:
		return object(ctx, bsonNS+"Regex", &bsonValue{x}), nil
	case bson.Decimal128:
		return object(ctx, bsonNS+"Decimal128", &bsonValue{x}), nil
	case bson.Timestamp:
		return object(ctx, bsonNS+"Timestamp", &bsonValue{x}), nil
	case bson.MinKey:
		return object(ctx, bsonNS+"MinKey", &bsonValue{x}), nil
	case bson.MaxKey:
		return object(ctx, bsonNS+"MaxKey", &bsonValue{x}), nil
	default:
		return nil, fmt.Errorf("unsupported BSON type %T", x)
	}
}
func bsonClasses() []*nativeClass {
	var classes []*nativeClass
	for _, name := range []string{"ObjectId", "UTCDateTime", "Binary", "Regex", "Decimal128", "Timestamp", "MinKey", "MaxKey"} {
		c := &nativeClass{name: bsonNS + name, interfaces: []string{"JsonSerializable", bsonNS + "Type"}}
		params := []string{}
		switch name {
		case "ObjectId":
			params = []string{"id?"}
		case "UTCDateTime":
			params = []string{"milliseconds?"}
		case "Decimal128":
			params = []string{"value"}
		case "Binary":
			params = []string{"data", "type"}
		case "Regex":
			params = []string{"pattern", "flags?"}
		case "Timestamp":
			params = []string{"increment", "timestamp"}
		}
		c.methods = append(c.methods, method("__construct", params, func(ctx data.Context) (data.GetValue, data.Control) {
			var x any
			var err error
			v := arg(ctx, 0)
			switch name {
			case "ObjectId":
				if isNull(v) {
					x = bson.NewObjectID()
				} else {
					var s string
					s, err = stringArg(v, "id")
					if err == nil {
						x, err = bson.ObjectIDFromHex(s)
					}
				}
			case "UTCDateTime":
				if isNull(v) {
					x = bson.DateTime(time.Now().UnixMilli())
				} else {
					var n int64
					n, err = intArg(v, "milliseconds")
					if err == nil {
						x = bson.DateTime(n)
					}
				}
			case "Decimal128":
				var s string
				s, err = stringArg(v, "value")
				if err == nil {
					x, err = bson.ParseDecimal128(s)
				}
			case "Binary":
				var s string
				var n int64
				s, err = stringArg(v, "data")
				if err == nil {
					n, err = intArg(arg(ctx, 1), "type")
				}
				if err == nil && (n < 0 || n > 255) {
					err = fmt.Errorf("binary subtype must be between 0 and 255")
				}
				if err == nil && (n == 3 || n == 4) && len(s) != 16 {
					err = fmt.Errorf("UUID binary data must contain 16 bytes")
				}
				x = bson.Binary{Subtype: byte(n), Data: []byte(s)}
			case "Regex":
				var s, flags string
				s, err = stringArg(v, "pattern")
				if err == nil && !isNull(arg(ctx, 1)) {
					flags, err = stringArg(arg(ctx, 1), "flags")
				}
				if strings.ContainsRune(s, 0) || strings.ContainsRune(flags, 0) {
					err = fmt.Errorf("regex cannot contain NUL")
				}
				if err == nil {
					var sorted string
					for _, f := range "ilmsux" {
						if strings.ContainsRune(flags, f) {
							sorted += string(f)
						}
					}
					for _, f := range flags {
						if !strings.ContainsRune("ilmsux", f) {
							err = fmt.Errorf("invalid regex flag %c", f)
						}
					}
					x = bson.Regex{Pattern: s, Options: sorted}
				}
			case "Timestamp":
				var i, t int64
				i, err = intArg(v, "increment")
				if err == nil {
					t, err = intArg(arg(ctx, 1), "timestamp")
				}
				if i < 0 || i > 1<<32-1 || t < 0 || t > 1<<32-1 {
					err = fmt.Errorf("timestamp components must be uint32")
				}
				x = bson.Timestamp{T: uint32(t), I: uint32(i)}
			case "MinKey":
				x = bson.MinKey{}
			case "MaxKey":
				x = bson.MaxKey{}
			}
			if err != nil {
				return invalid(err)
			}
			receiver(ctx).InstanceSource = &bsonValue{x}
			return data.NewNullValue(), nil
		}))
		methods := []string{}
		switch name {
		case "ObjectId":
			methods = []string{"__toString", "getTimestamp"}
		case "UTCDateTime", "Decimal128":
			methods = []string{"__toString"}
		case "Binary":
			methods = []string{"__toString", "getData", "getType"}
		case "Regex":
			methods = []string{"__toString", "getPattern", "getFlags"}
		case "Timestamp":
			methods = []string{"__toString", "getIncrement", "getTimestamp"}
		}
		for _, mn := range methods {
			c.methods = append(c.methods, method(mn, nil, func(ctx data.Context) (data.GetValue, data.Control) {
				s, ok := state[*bsonValue](ctx)
				if !ok {
					return logic("uninitialized BSON value")
				}
				var out string
				switch x := s.value.(type) {
				case bson.ObjectID:
					if mn == "getTimestamp" {
						return data.NewIntValue(int(x.Timestamp().Unix())), nil
					}
					out = x.Hex()
				case bson.DateTime:
					out = strconv.FormatInt(int64(x), 10)
				case bson.Decimal128:
					out = x.String()
				case bson.Binary:
					if mn == "getType" {
						return data.NewIntValue(int(x.Subtype)), nil
					}
					out = string(x.Data)
				case bson.Regex:
					switch mn {
					case "getPattern":
						out = x.Pattern
					case "getFlags":
						out = x.Options
					default:
						out = "/" + x.Pattern + "/" + x.Options
					}
				case bson.Timestamp:
					switch mn {
					case "getIncrement":
						return data.NewIntValue(int(x.I)), nil
					case "getTimestamp":
						return data.NewIntValue(int(x.T)), nil
					default:
						out = fmt.Sprintf("[%d:%d]", x.I, x.T)
					}
				}
				return data.NewStringValue(out), nil
			}))
		}
		c.methods = append(c.methods, method("jsonSerialize", nil, bsonJSONSerialize))
		if name == "Binary" {
			c.constants = map[string]data.Value{}
			for k, v := range map[string]int{"TYPE_GENERIC": 0, "TYPE_FUNCTION": 1, "TYPE_OLD_BINARY": 2, "TYPE_OLD_UUID": 3, "TYPE_UUID": 4, "TYPE_MD5": 5, "TYPE_ENCRYPTED": 6, "TYPE_COLUMN": 7, "TYPE_SENSITIVE": 8, "TYPE_VECTOR": 9, "TYPE_USER_DEFINED": 128} {
				c.constants[k] = data.NewIntValue(v)
			}
		}
		classes = append(classes, c)
	}
	return classes
}

func bsonJSONSerialize(ctx data.Context) (data.GetValue, data.Control) {
	s, ok := state[*bsonValue](ctx)
	if !ok {
		return logic("uninitialized BSON value")
	}
	pair := func(key string, v data.Value) *data.ArrayValue {
		a := data.NewArrayValue(nil).(*data.ArrayValue)
		a.SetStringKey(key, v)
		return a
	}
	str := data.NewStringValue
	switch x := s.value.(type) {
	case bson.ObjectID:
		return pair("$oid", str(x.Hex())), nil
	case bson.DateTime:
		return pair("$date", pair("$numberLong", str(strconv.FormatInt(int64(x), 10)))), nil
	case bson.Decimal128:
		return pair("$numberDecimal", str(x.String())), nil
	case bson.Binary:
		a := pair("base64", str(base64.StdEncoding.EncodeToString(x.Data)))
		a.SetStringKey("subType", str(fmt.Sprintf("%02x", x.Subtype)))
		return pair("$binary", a), nil
	case bson.Regex:
		a := pair("pattern", str(x.Pattern))
		a.SetStringKey("options", str(x.Options))
		return pair("$regularExpression", a), nil
	case bson.Timestamp:
		a := pair("t", data.NewIntValue(int(x.T)))
		a.SetStringKey("i", data.NewIntValue(int(x.I)))
		return pair("$timestamp", a), nil
	case bson.MinKey:
		return pair("$minKey", data.NewIntValue(1)), nil
	case bson.MaxKey:
		return pair("$maxKey", data.NewIntValue(1)), nil
	}
	return invalidf("unsupported BSON value")
}

func bsonFunctions() []data.FuncStmt {
	var out []data.FuncStmt
	for _, name := range []string{"fromPHP", "toPHP", "fromJSON", "toCanonicalExtendedJSON", "toRelaxedExtendedJSON"} {
		params := []string{"value"}
		if name == "toPHP" {
			params = append(params, "typeMap?")
		}
		out = append(out, method(bsonNS+name, params, func(ctx data.Context) (data.GetValue, data.Control) {
			var raw []byte
			var err error
			if name == "fromPHP" {
				var d bson.D
				d, err = encodeDocument(arg(ctx, 0))
				if err == nil {
					raw, err = bson.Marshal(d)
				}
			} else {
				var s string
				s, err = stringArg(arg(ctx, 0), "value")
				raw = []byte(s)
			}
			if err != nil {
				return invalid(err)
			}
			if name == "fromPHP" {
				return data.NewStringValue(string(raw)), nil
			}
			var d bson.D
			if name == "fromJSON" {
				err = bson.UnmarshalExtJSON(raw, false, &d)
				if err == nil {
					raw, err = bson.Marshal(d)
				}
				if err != nil {
					return invalid(err)
				}
				return data.NewStringValue(string(raw)), nil
			}
			err = bson.Unmarshal(raw, &d)
			if err != nil {
				return invalid(err)
			}
			if name == "toPHP" {
				t, e := parseTypeMap(arg(ctx, 1))
				if e != nil {
					return invalid(e)
				}
				v, e := decodeValue(ctx, d, t, true)
				if e != nil {
					return invalid(e)
				}
				return v, nil
			}
			raw, err = bson.MarshalExtJSON(d, name != "toRelaxedExtendedJSON", false)
			if err != nil {
				return invalid(err)
			}
			return data.NewStringValue(string(raw)), nil
		}))
	}
	return out
}
