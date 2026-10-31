package mongodb

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"runtime"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/php-any/origami/data"
	oruntime "github.com/php-any/origami/runtime"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/writeconcern"
)

type managerState struct {
	mu     sync.Mutex
	opts   *options.ClientOptions
	client *mongo.Client
}

func (s *managerState) connect() (*mongo.Client, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.client == nil {
		c, err := mongo.Connect(s.opts)
		if err != nil {
			return nil, err
		}
		s.client = c
	}
	return s.client, nil
}

func (s *managerState) close() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.client == nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = s.client.Disconnect(ctx)
}

type queryState struct {
	filter bson.D
	opts   *options.FindOptionsBuilder
}
type commandState struct{ document bson.D }
type bulkState struct {
	models   []mongo.WriteModel
	ordered  bool
	bypass   *bool
	executed bool
}
type concernState struct{ concern *writeconcern.WriteConcern }

func driverError(ctx data.Context, err error, result ...*mongo.BulkWriteResult) (data.GetValue, data.Control) {
	name := "RuntimeException"
	code := 0
	var ce mongo.CommandError
	var be mongo.BulkWriteException
	switch {
	case errors.As(err, &be):
		name = "BulkWriteException"
		if len(be.WriteErrors) > 0 {
			code = be.WriteErrors[0].Code
		}
	case errors.As(err, &ce):
		name = "CommandException"
		code = int(ce.Code)
		if code == 18 {
			name = "AuthenticationException"
		}
	case mongo.IsTimeout(err):
		name = "ConnectionTimeoutException"
	case mongo.IsNetworkError(err):
		name = "ConnectionException"
	}
	c, _ := ctx.GetVM().GetClass(exceptionNS + name)
	v := data.NewClassValue(c, ctx)
	if len(result) > 0 {
		v.InstanceSource = result[0]
	}
	v.SetProperty("message", data.NewStringValue(err.Error()))
	v.SetProperty("code", data.NewIntValue(code))
	return nil, data.NewErrorThrowFromClassValue(nil, v)
}
func namespace(v data.Value) (string, string, error) {
	s, e := stringArg(v, "namespace")
	if e != nil {
		return "", "", e
	}
	db, col, ok := strings.Cut(s, ".")
	if !ok || db == "" || col == "" || strings.ContainsRune(s, 0) {
		return "", "", fmt.Errorf("namespace must contain a database and collection separated by '.'")
	}
	return db, col, nil
}

func driverClasses() []*nativeClass {
	manager := &nativeClass{name: driverNS + "Manager"}
	manager.methods = []*nativeMethod{
		method("__construct", []string{"uri?", "uriOptions?", "driverOptions?"}, func(ctx data.Context) (data.GetValue, data.Control) {
			uri := "mongodb://127.0.0.1/"
			var err error
			if !isNull(arg(ctx, 0)) {
				uri, err = stringArg(arg(ctx, 0), "uri")
				if err != nil {
					return invalid(err)
				}
			}
			// URI query options are validated and interpreted by the official driver.
			// Array options are deliberately restricted instead of silently discarded.
			uri, err = mergeURIOptions(uri, arg(ctx, 1))
			if err != nil {
				return invalid(err)
			}
			if _, err = optionMap(arg(ctx, 2)); err != nil {
				return invalid(err)
			}
			opts := options.Client().ApplyURI(uri)
			if err = opts.Validate(); err != nil {
				return invalid(err)
			}
			s := &managerState{opts: opts}
			runtime.SetFinalizer(s, (*managerState).close)
			receiver(ctx).InstanceSource = s
			return data.NewNullValue(), nil
		}),
		method("executeQuery", []string{"namespace", "query", "options?"}, func(ctx data.Context) (data.GetValue, data.Control) {
			s, ok := state[*managerState](ctx)
			if !ok {
				return logic("uninitialized Manager")
			}
			db, col, err := namespace(arg(ctx, 0))
			if err != nil {
				return invalid(err)
			}
			q, ok := objectState[*queryState](arg(ctx, 1))
			if !ok {
				return invalidf("query must be a MongoDB\\Driver\\Query")
			}
			if _, err = optionMap(arg(ctx, 2)); err != nil {
				return invalid(err)
			}
			defer runtime.KeepAlive(s)
			client, err := s.connect()
			if err != nil {
				return driverError(ctx, err)
			}
			cur, err := client.Database(db).Collection(col).Find(oruntime.RequestContext(), q.filter, q.opts)
			if err != nil {
				return driverError(ctx, err)
			}
			return newCursor(ctx, cur, s), nil
		}),
		method("executeCommand", []string{"db", "command", "options?"}, func(ctx data.Context) (data.GetValue, data.Control) {
			s, ok := state[*managerState](ctx)
			if !ok {
				return logic("uninitialized Manager")
			}
			db, err := stringArg(arg(ctx, 0), "db")
			if err != nil {
				return invalid(err)
			}
			q, ok := objectState[*commandState](arg(ctx, 1))
			if !ok {
				return invalidf("command must be a MongoDB\\Driver\\Command")
			}
			if _, err = optionMap(arg(ctx, 2)); err != nil {
				return invalid(err)
			}
			defer runtime.KeepAlive(s)
			client, err := s.connect()
			if err != nil {
				return driverError(ctx, err)
			}
			var cur *mongo.Cursor
			switch strings.ToLower(q.document[0].Key) {
			case "find", "aggregate", "listcollections", "listindexes":
				cur, err = client.Database(db).RunCommandCursor(oruntime.RequestContext(), q.document)
			default:
				var d bson.D
				err = client.Database(db).RunCommand(oruntime.RequestContext(), q.document).Decode(&d)
				if err == nil {
					cur, err = mongo.NewCursorFromDocuments([]any{d}, nil, nil)
				}
			}
			if err != nil {
				return driverError(ctx, err)
			}
			return newCursor(ctx, cur, s), nil
		}),
		method("executeBulkWrite", []string{"namespace", "bulk", "options?"}, func(ctx data.Context) (data.GetValue, data.Control) {
			s, ok := state[*managerState](ctx)
			if !ok {
				return logic("uninitialized Manager")
			}
			db, col, err := namespace(arg(ctx, 0))
			if err != nil {
				return invalid(err)
			}
			b, ok := objectState[*bulkState](arg(ctx, 1))
			if !ok {
				return invalidf("bulk must be a MongoDB\\Driver\\BulkWrite")
			}
			if b.executed {
				return invalidf("BulkWrite objects may only be executed once")
			}
			if len(b.models) == 0 {
				return invalidf("Cannot do an empty bulk write")
			}
			collectionOptions := options.Collection()
			v := arg(ctx, 2)
			if wc, ok := objectState[*concernState](v); ok {
				collectionOptions.SetWriteConcern(wc.concern)
			} else {
				m, e := optionMap(v, "writeConcern")
				if e != nil {
					return invalid(e)
				}
				if w, exists := m["writeConcern"]; exists {
					wc, ok := objectState[*concernState](w)
					if !ok {
						return invalidf("writeConcern must be a MongoDB\\Driver\\WriteConcern")
					}
					collectionOptions.SetWriteConcern(wc.concern)
				}
			}
			opts := options.BulkWrite().SetOrdered(b.ordered)
			if b.bypass != nil {
				opts.SetBypassDocumentValidation(*b.bypass)
			}
			defer runtime.KeepAlive(s)
			client, err := s.connect()
			if err != nil {
				return driverError(ctx, err)
			}
			b.executed = true
			result, err := client.Database(db).Collection(col, collectionOptions).BulkWrite(oruntime.RequestContext(), b.models, opts)
			if err != nil {
				// Go driver v2.5 returns before normalizing ordered partial
				// results on a server write error. PHP still exposes these counts.
				var bulkErr mongo.BulkWriteException
				if result != nil && b.ordered && errors.As(err, &bulkErr) {
					result.Acknowledged = true
					result.MatchedCount -= result.UpsertedCount
				}
				return driverError(ctx, err, result)
			}
			return object(ctx, driverNS+"WriteResult", result), nil
		}),
	}
	query := &nativeClass{name: driverNS + "Query", methods: []*nativeMethod{method("__construct", []string{"filter", "queryOptions?"}, func(ctx data.Context) (data.GetValue, data.Control) {
		d, err := encodeDocument(arg(ctx, 0))
		if err != nil {
			return invalid(err)
		}
		m, err := optionMap(arg(ctx, 1), "projection", "sort", "limit", "skip", "batchSize", "hint", "comment", "allowDiskUse", "noCursorTimeout")
		if err != nil {
			return invalid(err)
		}
		opts := options.Find()
		for k, v := range m {
			switch k {
			case "projection", "sort":
				var d bson.D
				d, err = encodeDocument(v)
				if k == "projection" {
					opts.SetProjection(d)
				} else {
					opts.SetSort(d)
				}
			case "limit", "skip", "batchSize":
				var n int64
				n, err = intArg(v, k)
				if k != "limit" && n < 0 {
					err = fmt.Errorf("%s must be non-negative", k)
				}
				switch k {
				case "limit":
					opts.SetLimit(n)
				case "skip":
					opts.SetSkip(n)
				case "batchSize":
					if n > 1<<31-1 {
						err = fmt.Errorf("batchSize is too large")
					}
					opts.SetBatchSize(int32(n))
				}
			case "hint":
				if s, ok := v.(*data.StringValue); ok {
					opts.SetHint(s.Value)
				} else {
					var d bson.D
					d, err = encodeDocument(v)
					opts.SetHint(d)
				}
			case "comment":
				var x any
				x, err = encodeValue(v, false, 0)
				opts.SetComment(x)
			case "allowDiskUse", "noCursorTimeout":
				var b bool
				b, err = boolOption(m, k, false)
				if k == "allowDiskUse" {
					opts.SetAllowDiskUse(b)
				} else {
					opts.SetNoCursorTimeout(b)
				}
			}
			if err != nil {
				return invalid(err)
			}
		}
		receiver(ctx).InstanceSource = &queryState{d, opts}
		return data.NewNullValue(), nil
	})}}
	command := &nativeClass{name: driverNS + "Command", methods: []*nativeMethod{method("__construct", []string{"document", "commandOptions?"}, func(ctx data.Context) (data.GetValue, data.Control) {
		d, err := encodeDocument(arg(ctx, 0))
		if err != nil {
			return invalid(err)
		}
		if len(d) == 0 {
			return invalidf("command document must not be empty")
		}
		if _, err = optionMap(arg(ctx, 1)); err != nil {
			return invalid(err)
		}
		receiver(ctx).InstanceSource = &commandState{d}
		return data.NewNullValue(), nil
	})}}
	bulk := bulkClass()
	result := &nativeClass{name: driverNS + "WriteResult"}
	result.methods = append(result.methods, method("__construct", nil, func(ctx data.Context) (data.GetValue, data.Control) {
		return logic("WriteResult cannot be constructed directly")
	}))
	for _, name := range []string{"getInsertedCount", "getMatchedCount", "getModifiedCount", "getDeletedCount", "getUpsertedCount", "getUpsertedIds", "isAcknowledged"} {
		result.methods = append(result.methods, method(name, nil, func(ctx data.Context) (data.GetValue, data.Control) {
			r, ok := state[*mongo.BulkWriteResult](ctx)
			if !ok {
				return logic("uninitialized WriteResult")
			}
			if name == "isAcknowledged" {
				return data.NewBoolValue(r.Acknowledged), nil
			}
			if !r.Acknowledged {
				return logic("Cannot read an unacknowledged write result")
			}
			var n int64
			switch name {
			case "getInsertedCount":
				n = r.InsertedCount
			case "getMatchedCount":
				n = r.MatchedCount
			case "getModifiedCount":
				n = r.ModifiedCount
			case "getDeletedCount":
				n = r.DeletedCount
			case "getUpsertedCount":
				n = r.UpsertedCount
			case "getUpsertedIds":
				a := data.NewArrayValue(nil).(*data.ArrayValue)
				keys := make([]int64, 0, len(r.UpsertedIDs))
				for i := range r.UpsertedIDs {
					keys = append(keys, i)
				}
				sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })
				for _, i := range keys {
					if id, ok := r.UpsertedIDs[i]; ok {
						v, e := decodeValue(ctx, id, typeMap{}, false)
						if e != nil {
							return invalid(e)
						}
						a.SetIntKey(int(i), v)
					}
				}
				return a, nil
			}
			return data.NewIntValue(int(n)), nil
		}))
	}
	concern := &nativeClass{name: driverNS + "WriteConcern", constants: map[string]data.Value{"MAJORITY": data.NewStringValue("majority")}}
	concern.methods = []*nativeMethod{method("__construct", []string{"w", "wtimeout?", "journal?"}, func(ctx data.Context) (data.GetValue, data.Control) {
		w := arg(ctx, 0)
		wc := &writeconcern.WriteConcern{}
		switch v := w.(type) {
		case *data.StringValue:
			wc.W = v.Value
		case *data.IntValue:
			if v.Value < -1 {
				return invalidf("w must be >= -1")
			}
			wc.W = v.Value
		default:
			return invalidf("w must be an integer or string")
		}
		if !isNull(arg(ctx, 1)) {
			n, e := intArg(arg(ctx, 1), "wtimeout")
			if e != nil {
				return invalid(e)
			}
			if n != 0 {
				return invalidf("nonzero wtimeout is not supported; use timeoutMS in the Manager URI")
			}
		}
		if !isNull(arg(ctx, 2)) {
			b, ok := arg(ctx, 2).(*data.BoolValue)
			if !ok {
				return invalidf("journal must be boolean")
			}
			wc.Journal = &b.Value
		}
		if !wc.IsValid() {
			return invalidf("invalid write concern (w must be non-negative and w=0 cannot use journal=true)")
		}
		receiver(ctx).InstanceSource = &concernState{wc}
		return data.NewNullValue(), nil
	})}
	for _, name := range []string{"getW", "getWtimeout", "getJournal", "isDefault"} {
		concern.methods = append(concern.methods, method(name, nil, func(ctx data.Context) (data.GetValue, data.Control) {
			s, ok := state[*concernState](ctx)
			if !ok {
				return logic("uninitialized WriteConcern")
			}
			switch name {
			case "getW":
				switch v := s.concern.W.(type) {
				case int:
					return data.NewIntValue(v), nil
				case string:
					return data.NewStringValue(v), nil
				}
			case "getWtimeout":
				return data.NewIntValue(0), nil
			case "getJournal":
				if s.concern.Journal != nil {
					return data.NewBoolValue(*s.concern.Journal), nil
				}
			case "isDefault":
				return data.NewBoolValue(s.concern.W == nil && s.concern.Journal == nil), nil
			}
			return data.NewNullValue(), nil
		}))
	}
	classes := []*nativeClass{manager, query, command, bulk, result, concern}
	return append(classes, cursorClasses()...)
}
func mergeURIOptions(uri string, v data.Value) (string, error) {
	m, err := optionMap(v, "appname", "authSource", "authMechanism", "replicaSet", "connectTimeoutMS", "serverSelectionTimeoutMS", "socketTimeoutMS", "timeoutMS", "maxPoolSize", "minPoolSize", "retryReads", "retryWrites", "directConnection", "tls", "w", "journal", "readPreference")
	if err != nil {
		return "", err
	}
	if len(m) == 0 {
		return uri, nil
	}
	u, err := url.Parse(uri)
	if err != nil {
		return "", err
	}
	q := u.Query()
	for k, v := range m {
		switch k {
		case "retryReads", "retryWrites", "directConnection", "tls", "journal":
			if _, ok := v.(*data.BoolValue); !ok {
				return "", fmt.Errorf("%s must be boolean", k)
			}
		case "connectTimeoutMS", "serverSelectionTimeoutMS", "socketTimeoutMS", "timeoutMS", "maxPoolSize", "minPoolSize":
			n, e := intArg(v, k)
			if e != nil {
				return "", e
			}
			if n < 0 {
				return "", fmt.Errorf("%s must be non-negative", k)
			}
		case "w":
			switch v.(type) {
			case *data.IntValue, *data.StringValue:
			default:
				return "", fmt.Errorf("w must be string or integer")
			}
		default:
			if _, e := stringArg(v, k); e != nil {
				return "", e
			}
		}
		for old := range q {
			if strings.EqualFold(old, k) {
				delete(q, old)
			}
		}
		s := v.AsString()
		if b, ok := v.(*data.BoolValue); ok {
			if b.Value {
				s = "true"
			} else {
				s = "false"
			}
		}
		q.Set(k, s)
	}
	u.RawQuery = q.Encode()
	return u.String(), nil
}

func bulkClass() *nativeClass {
	c := &nativeClass{name: driverNS + "BulkWrite", interfaces: []string{"Countable"}}
	c.methods = []*nativeMethod{method("__construct", []string{"options?"}, func(ctx data.Context) (data.GetValue, data.Control) {
		m, e := optionMap(arg(ctx, 0), "ordered", "bypassDocumentValidation")
		if e != nil {
			return invalid(e)
		}
		ordered, e := boolOption(m, "ordered", true)
		if e != nil {
			return invalid(e)
		}
		s := &bulkState{ordered: ordered}
		if _, ok := m["bypassDocumentValidation"]; ok {
			b, e := boolOption(m, "bypassDocumentValidation", false)
			if e != nil {
				return invalid(e)
			}
			s.bypass = &b
		}
		receiver(ctx).InstanceSource = s
		return data.NewNullValue(), nil
	}),
		method("count", nil, func(ctx data.Context) (data.GetValue, data.Control) {
			s, ok := state[*bulkState](ctx)
			if !ok {
				return logic("uninitialized BulkWrite")
			}
			return data.NewIntValue(len(s.models)), nil
		}),
		method("insert", []string{"document"}, func(ctx data.Context) (data.GetValue, data.Control) {
			s, ok := state[*bulkState](ctx)
			if !ok {
				return logic("uninitialized BulkWrite")
			}
			if s.executed {
				return logic("BulkWrite has already been executed")
			}
			d, e := encodeDocument(arg(ctx, 0))
			if e != nil {
				return invalid(e)
			}
			var id any
			exists := false
			for _, v := range d {
				if v.Key == "_id" {
					id = v.Value
					exists = true
					break
				}
			}
			if !exists {
				id = bson.NewObjectID()
				d = append(bson.D{{Key: "_id", Value: id}}, d...)
			}
			s.models = append(s.models, mongo.NewInsertOneModel().SetDocument(d))
			if oid, ok := id.(bson.ObjectID); ok {
				return object(ctx, bsonNS+"ObjectId", &bsonValue{oid}), nil
			}
			return data.NewNullValue(), nil
		}),
		method("update", []string{"filter", "newObj", "updateOptions?"}, func(ctx data.Context) (data.GetValue, data.Control) {
			s, ok := state[*bulkState](ctx)
			if !ok {
				return logic("uninitialized BulkWrite")
			}
			if s.executed {
				return logic("BulkWrite has already been executed")
			}
			if a, ok := arg(ctx, 1).(*data.ArrayValue); ok && a.Len() > 0 {
				list := true
				for i, slot := range a.Range() {
					key, ok := slot.PHPArrayKey(i).(*data.IntValue)
					if !ok || key.Value != i {
						list = false
						break
					}
				}
				if list {
					return invalidf("update pipelines are not supported yet")
				}
			}
			f, e := encodeDocument(arg(ctx, 0))
			if e != nil {
				return invalid(e)
			}
			u, e := encodeDocument(arg(ctx, 1))
			if e != nil {
				return invalid(e)
			}
			m, e := optionMap(arg(ctx, 2), "multi", "upsert")
			if e != nil {
				return invalid(e)
			}
			multi, e := boolOption(m, "multi", false)
			if e != nil {
				return invalid(e)
			}
			upsert, e := boolOption(m, "upsert", false)
			if e != nil {
				return invalid(e)
			}
			var model mongo.WriteModel
			if len(u) > 0 && strings.HasPrefix(u[0].Key, "$") {
				for _, v := range u {
					if !strings.HasPrefix(v.Key, "$") {
						return invalidf("update document mixes operators and fields")
					}
				}
				if multi {
					model = mongo.NewUpdateManyModel().SetFilter(f).SetUpdate(u).SetUpsert(upsert)
				} else {
					model = mongo.NewUpdateOneModel().SetFilter(f).SetUpdate(u).SetUpsert(upsert)
				}
			} else {
				if multi {
					return invalidf("replacement document cannot be used with multi=true")
				}
				for _, v := range u {
					if strings.HasPrefix(v.Key, "$") {
						return invalidf("replacement document contains an update operator")
					}
				}
				model = mongo.NewReplaceOneModel().SetFilter(f).SetReplacement(u).SetUpsert(upsert)
			}
			s.models = append(s.models, model)
			return data.NewNullValue(), nil
		}),
		method("delete", []string{"filter", "deleteOptions?"}, func(ctx data.Context) (data.GetValue, data.Control) {
			s, ok := state[*bulkState](ctx)
			if !ok {
				return logic("uninitialized BulkWrite")
			}
			if s.executed {
				return logic("BulkWrite has already been executed")
			}
			f, e := encodeDocument(arg(ctx, 0))
			if e != nil {
				return invalid(e)
			}
			m, e := optionMap(arg(ctx, 1), "limit")
			if e != nil {
				return invalid(e)
			}
			limit := int64(0)
			if v, ok := m["limit"]; ok {
				if b, ok := v.(*data.BoolValue); ok {
					if b.Value {
						limit = 1
					}
				} else {
					limit, e = intArg(v, "limit")
				}
				if e != nil {
					return invalid(e)
				}
				if limit != 0 && limit != 1 {
					return invalidf("delete limit must be 0 or 1")
				}
			}
			if limit == 1 {
				s.models = append(s.models, mongo.NewDeleteOneModel().SetFilter(f))
			} else {
				s.models = append(s.models, mongo.NewDeleteManyModel().SetFilter(f))
			}
			return data.NewNullValue(), nil
		}),
	}
	return c
}
