package mongodb

import (
	"context"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/parser"
	oruntime "github.com/php-any/origami/runtime"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func testContext(t *testing.T) data.Context {
	t.Helper()
	vm := oruntime.NewVM(parser.NewParser())
	vm.AddClass(&data.StdClass{})
	Load(vm)
	return vm.CreateContext(nil)
}
func call(t *testing.T, ctx data.Context, o *data.ClassValue, name string, args ...data.Value) (data.GetValue, data.Control) {
	t.Helper()
	m, ok := o.Class.GetMethod(name)
	if !ok {
		t.Fatalf("missing %s", name)
	}
	c := o.CreateContext(m.GetVariables())
	for i, v := range args {
		c.SetIndexZVal(i, data.NewZVal(v))
	}
	return m.Call(c)
}
func mustCall(t *testing.T, ctx data.Context, o *data.ClassValue, name string, args ...data.Value) data.Value {
	t.Helper()
	v, c := call(t, ctx, o, name, args...)
	if c != nil {
		t.Fatalf("%s: %v", name, c)
	}
	return v.(data.Value)
}
func construct(t *testing.T, ctx data.Context, name string, args ...data.Value) *data.ClassValue {
	t.Helper()
	v := object(ctx, driverNS+name, nil)
	mustCall(t, ctx, v, "__construct", args...)
	return v
}
func doc(pairs ...any) *data.ArrayValue {
	a := data.NewArrayValue(nil).(*data.ArrayValue)
	for i := 0; i < len(pairs); i += 2 {
		a.SetKey(data.NewStringValue(pairs[i].(string)), pairs[i+1].(data.Value))
	}
	return a
}

func TestBSONOrderTypesAndCycles(t *testing.T) {
	ctx := testContext(t)
	d := doc("z", data.NewIntValue(1), "a", data.NewArrayValue([]data.Value{data.NewIntValue(2)}))
	encoded, err := encodeDocument(d)
	if err != nil {
		t.Fatal(err)
	}
	if encoded[0].Key != "z" || encoded[1].Key != "a" {
		t.Fatal("field order lost")
	}
	if _, ok := encoded[1].Value.(bson.A); !ok {
		t.Fatal("list is not BSON array")
	}
	sparse := data.NewArrayValue(nil).(*data.ArrayValue)
	sparse.SetIntKey(4, data.NewStringValue("v"))
	v, err := encodeValue(sparse, false, 0)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(v, bson.D{{Key: "4", Value: "v"}}) {
		t.Fatalf("sparse: %#v", v)
	}
	cyclic := object(ctx, "stdClass", nil)
	cyclic.SetProperty("self", cyclic)
	if _, err = encodeDocument(cyclic); err == nil {
		t.Fatal("cycle accepted")
	}
}
func TestBulkSnapshotsAndFailedAppend(t *testing.T) {
	ctx := testContext(t)
	b := construct(t, ctx, "BulkWrite")
	d := doc("name", data.NewStringValue("before"))
	mustCall(t, ctx, b, "insert", d)
	d.SetStringKey("name", data.NewStringValue("after"))
	s := b.InstanceSource.(*bulkState)
	stored := s.models[0].(*mongo.InsertOneModel).Document.(bson.D)
	if stored[1].Value != "before" {
		t.Fatal("queued document was mutated")
	}
	if _, ok := d.LookupZValByStringKey("_id"); ok {
		t.Fatal("input was mutated")
	}
	_, ctl := call(t, ctx, b, "update", doc(), doc("field", data.NewIntValue(1)), doc("multi", data.NewBoolValue(true)))
	if ctl == nil || len(s.models) != 1 {
		t.Fatal("invalid update appended")
	}
}
func TestCursorBatchesAndRewind(t *testing.T) {
	ctx := testContext(t)
	cur, err := mongo.NewCursorFromDocuments([]any{bson.D{{Key: "n", Value: int64(1)}}, bson.D{{Key: "n", Value: int64(2)}}}, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	o := newCursor(ctx, cur, nil)
	mustCall(t, ctx, o, "setTypeMap", doc("root", data.NewStringValue("array")))
	it := mustCall(t, ctx, o, "getIterator").(*data.ClassValue)
	mustCall(t, ctx, it, "rewind")
	first := mustCall(t, ctx, it, "current").(*data.ArrayValue)
	slot, _ := first.LookupZValByStringKey("n")
	if slot.ReadValue().AsString() != "1" {
		t.Fatal("first row")
	}
	mustCall(t, ctx, it, "next")
	if _, ctl := call(t, ctx, it, "rewind"); ctl == nil {
		t.Fatal("rewind after advance succeeded")
	}
	mustCall(t, ctx, it, "next")
	if mustCall(t, ctx, it, "valid").(*data.BoolValue).Value {
		t.Fatal("cursor not exhausted")
	}
}

func TestCursorPHPForeach(t *testing.T) {
	ctx := testContext(t)
	ctx.GetVM().AddInterface(node.NewInterfaceStatement(nil, "Iterator", nil, nil))
	ctx.GetVM().AddInterface(node.NewInterfaceStatement(nil, "IteratorAggregate", nil, nil))
	cur, err := mongo.NewCursorFromDocuments([]any{bson.D{{Key: "n", Value: int64(3)}}, bson.D{{Key: "n", Value: int64(5)}}}, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	p := parser.NewParser()
	program, ctl := p.ParseString(`<?php $sum=0; foreach ($cursor as $row) { $sum += $row->n; }`, "mongodb_foreach.php")
	if ctl != nil {
		t.Fatalf("parse: %v", ctl)
	}
	ctx = ctx.GetVM().CreateContext(p.GetVariables())
	ctx.SetVariableByName("cursor", newCursor(ctx, cur, nil))
	_, ctl = program.GetValue(ctx)
	if ctl != nil {
		t.Fatalf("foreach: %v", ctl)
	}
	v, ok := ctx.GetVariableByName("sum")
	if !ok || v.AsString() != "8" {
		t.Fatalf("foreach sum: %v", v)
	}
}

// Opt-in test for a real MongoDB deployment. Only a uniquely named collection
// created by this test is removed. No existing application data is touched.
func TestMongoDBIntegration(t *testing.T) {
	uri := os.Getenv("ORIGAMI_MONGODB_URI")
	if uri == "" {
		t.Skip("set ORIGAMI_MONGODB_URI for real-server CRUD integration")
	}
	ctx := testContext(t)
	m := construct(t, ctx, "Manager", data.NewStringValue(uri))
	s := m.InstanceSource.(*managerState)
	if s.client != nil {
		t.Fatal("Manager eagerly connected")
	}
	if _, err := s.connect(); err != nil {
		t.Fatal(err)
	}
	defer s.close()
	name := "origami_driver_test_" + bson.NewObjectID().Hex()
	ns := data.NewStringValue("origami_test." + name)
	defer func() {
		cleanupCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_ = s.client.Database("origami_test").Collection(name).Drop(cleanupCtx)
	}()
	b := construct(t, ctx, "BulkWrite")
	id := mustCall(t, ctx, b, "insert", doc("name", data.NewStringValue("first"), "n", data.NewIntValue(1)))
	r := mustCall(t, ctx, m, "executeBulkWrite", ns, b).(*data.ClassValue)
	if mustCall(t, ctx, r, "getInsertedCount").AsString() != "1" {
		t.Fatal("insert count")
	}
	q := construct(t, ctx, "Query", doc("_id", id))
	cursor := mustCall(t, ctx, m, "executeQuery", ns, q).(*data.ClassValue)
	rows := mustCall(t, ctx, cursor, "toArray").(*data.ArrayValue)
	if rows.Len() != 1 {
		t.Fatal("query row count")
	}
	b = construct(t, ctx, "BulkWrite")
	mustCall(t, ctx, b, "update", doc("_id", id), doc("$set", doc("n", data.NewIntValue(2))))
	r = mustCall(t, ctx, m, "executeBulkWrite", ns, b).(*data.ClassValue)
	if mustCall(t, ctx, r, "getModifiedCount").AsString() != "1" {
		t.Fatal("update count")
	}
	b = construct(t, ctx, "BulkWrite")
	mustCall(t, ctx, b, "delete", doc("_id", id))
	r = mustCall(t, ctx, m, "executeBulkWrite", ns, b).(*data.ClassValue)
	if mustCall(t, ctx, r, "getDeletedCount").AsString() != "1" {
		t.Fatal("delete count")
	}
}

func TestURIOptions(t *testing.T) {
	ctx := testContext(t)
	m := construct(t, ctx, "Manager", data.NewStringValue("mongodb://127.0.0.1:1/"))
	if m.InstanceSource.(*managerState).client != nil {
		t.Fatal("Manager construction must not start connections")
	}
	uri, err := mergeURIOptions("mongodb://localhost/?retryWrites=true", doc("retryWrites", data.NewBoolValue(false), "serverSelectionTimeoutMS", data.NewIntValue(20)))
	if err != nil {
		t.Fatal(err)
	}
	opts := options.Client().ApplyURI(uri)
	if err = opts.Validate(); err != nil {
		t.Fatal(err)
	}
	if opts.RetryWrites == nil || *opts.RetryWrites {
		t.Fatal("URI option did not override")
	}
}
