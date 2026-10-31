package mongodb

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/php-any/origami/data"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// A wire fixture exercises the real Go driver's network path without requiring
// mongod in unit-test environments. It is not a substitute for the opt-in
// real-server integration test. Replies intentionally span cursor batches.
func wireFixture(t *testing.T) (string, <-chan bson.D) {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	commands := make(chan bson.D, 64)
	var wg sync.WaitGroup
	var mu sync.Mutex
	var conns []net.Conn
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			c, e := ln.Accept()
			if e != nil {
				return
			}
			mu.Lock()
			conns = append(conns, c)
			mu.Unlock()
			wg.Add(1)
			go func() {
				defer wg.Done()
				defer c.Close()
				for {
					_ = c.SetDeadline(time.Now().Add(15 * time.Second))
					header := make([]byte, 16)
					if _, e := io.ReadFull(c, header); e != nil {
						return
					}
					length := int(binary.LittleEndian.Uint32(header))
					if length < 16 || length > 16*1024*1024 {
						return
					}
					body := make([]byte, length-16)
					if _, e := io.ReadFull(c, body); e != nil {
						return
					}
					opcode := binary.LittleEndian.Uint32(header[12:])
					var raw []byte
					switch opcode {
					case 2004:
						if len(body) < 5 {
							return
						}
						end := 4
						for end < len(body) && body[end] != 0 {
							end++
						}
						if end+9 > len(body) {
							return
						}
						raw = body[end+9:]
					case 2013:
						if len(body) < 5 || body[4] != 0 {
							return
						}
						raw = body[5:]
						if len(raw) < 4 {
							return
						}
						n := int(binary.LittleEndian.Uint32(raw))
						if n > len(raw) || n < 5 {
							return
						}
						raw = raw[:n]
					default:
						return
					}
					var d bson.D
					if e := bson.Unmarshal(raw, &d); e != nil || len(d) == 0 {
						return
					}
					name := strings.ToLower(d[0].Key)
					reply := bson.D{{Key: "ok", Value: 1.0}}
					switch name {
					case "ismaster", "hello":
						reply = append(reply, bson.E{Key: "isWritablePrimary", Value: true}, bson.E{Key: "ismaster", Value: true}, bson.E{Key: "minWireVersion", Value: int32(0)}, bson.E{Key: "maxWireVersion", Value: int32(25)}, bson.E{Key: "maxBsonObjectSize", Value: int32(16777216)}, bson.E{Key: "maxMessageSizeBytes", Value: int32(48000000)}, bson.E{Key: "maxWriteBatchSize", Value: int32(100000)})
					case "insert":
						commands <- d
						if d[0].Value == "reject" {
							reply = append(reply, bson.E{Key: "n", Value: int32(0)}, bson.E{Key: "writeErrors", Value: bson.A{bson.D{{Key: "index", Value: int32(0)}, {Key: "code", Value: int32(11000)}, {Key: "errmsg", Value: "duplicate key fixture"}}}})
						} else {
							reply = append(reply, bson.E{Key: "n", Value: int32(1)})
						}
					case "update":
						commands <- d
						reply = append(reply, bson.E{Key: "n", Value: int32(1)}, bson.E{Key: "nModified", Value: int32(1)})
					case "delete":
						commands <- d
						reply = append(reply, bson.E{Key: "n", Value: int32(1)})
					case "find":
						commands <- d
						reply = append(reply, bson.E{Key: "cursor", Value: bson.D{{Key: "id", Value: int64(123)}, {Key: "ns", Value: "test.items"}, {Key: "firstBatch", Value: bson.A{bson.D{{Key: "n", Value: int64(1)}}}}}})
					case "getmore":
						commands <- d
						reply = append(reply, bson.E{Key: "cursor", Value: bson.D{{Key: "id", Value: int64(0)}, {Key: "ns", Value: "test.items"}, {Key: "nextBatch", Value: bson.A{bson.D{{Key: "n", Value: int64(2)}}}}}})
					case "failcommand":
						commands <- d
						reply = bson.D{{Key: "ok", Value: 0.0}, {Key: "code", Value: int32(13)}, {Key: "errmsg", Value: "fixture command rejected"}}
					default:
						commands <- d
					}
					payload, e := bson.Marshal(reply)
					if e != nil {
						return
					}
					var response []byte
					if opcode == 2004 {
						response = make([]byte, 36)
						binary.LittleEndian.PutUint32(response[12:], 1)
						binary.LittleEndian.PutUint32(response[32:], 1)
					} else {
						response = make([]byte, 21)
						binary.LittleEndian.PutUint32(response[12:], 2013)
					}
					response = append(response, payload...)
					binary.LittleEndian.PutUint32(response, uint32(len(response)))
					binary.LittleEndian.PutUint32(response[4:], 42)
					copy(response[8:12], header[4:8])
					if _, e = c.Write(response); e != nil {
						return
					}
				}
			}()
		}
	}()
	t.Cleanup(func() {
		_ = ln.Close()
		mu.Lock()
		for _, c := range conns {
			_ = c.Close()
		}
		mu.Unlock()
		wg.Wait()
	})
	return fmt.Sprintf("mongodb://%s/?directConnection=true&serverSelectionTimeoutMS=2000&timeoutMS=3000", ln.Addr()), commands
}

func TestDriverWireCRUDAndGetMore(t *testing.T) {
	uri, commands := wireFixture(t)
	ctx := testContext(t)
	manager := construct(t, ctx, "Manager", data.NewStringValue(uri))
	defer manager.InstanceSource.(*managerState).close()
	bulk := construct(t, ctx, "BulkWrite")
	mustCall(t, ctx, bulk, "insert", doc("name", data.NewStringValue("fixture")))
	mustCall(t, ctx, bulk, "update", doc(), doc("$set", doc("n", data.NewIntValue(2))))
	mustCall(t, ctx, bulk, "delete", doc(), doc("limit", data.NewIntValue(1)))
	result := mustCall(t, ctx, manager, "executeBulkWrite", data.NewStringValue("test.items"), bulk).(*data.ClassValue)
	for _, method := range []string{"getInsertedCount", "getModifiedCount", "getDeletedCount"} {
		if mustCall(t, ctx, result, method).AsString() != "1" {
			t.Fatalf("%s count", method)
		}
	}
	if _, ctl := call(t, ctx, manager, "executeBulkWrite", data.NewStringValue("test.items"), bulk); ctl == nil {
		t.Fatal("bulk reused")
	}
	query := construct(t, ctx, "Query", doc("active", data.NewBoolValue(true)), doc("batchSize", data.NewIntValue(1)))
	cursor := mustCall(t, ctx, manager, "executeQuery", data.NewStringValue("test.items"), query).(*data.ClassValue)
	rows := mustCall(t, ctx, cursor, "toArray").(*data.ArrayValue)
	if rows.Len() != 2 {
		t.Fatalf("expected getMore to produce 2 rows, got %d", rows.Len())
	}
	for i := 0; i < 2; i++ {
		v, _ := rows.At(i).ReadValue().(*data.ClassValue).GetProperty("n")
		if v.AsString() != fmt.Sprint(i+1) {
			t.Fatal("row order")
		}
	}
	command := construct(t, ctx, "Command", doc("ping", data.NewIntValue(1)))
	cursor = mustCall(t, ctx, manager, "executeCommand", data.NewStringValue("test"), command).(*data.ClassValue)
	if mustCall(t, ctx, cursor, "toArray").(*data.ArrayValue).Len() != 1 {
		t.Fatal("command cursor")
	}
	command = construct(t, ctx, "Command", doc("failCommand", data.NewIntValue(1)))
	if _, ctl := call(t, ctx, manager, "executeCommand", data.NewStringValue("test"), command); ctl == nil {
		t.Fatal("server error suppressed")
	} else {
		e := ctl.(*data.ThrowValue).Object
		if e.GetName() != exceptionNS+"CommandException" || mustCall(t, ctx, e, "getCode").AsString() != "13" {
			t.Fatal("command exception type/code")
		}
	}
	bulk = construct(t, ctx, "BulkWrite")
	mustCall(t, ctx, bulk, "insert", doc("_id", data.NewIntValue(1)))
	_, ctl := call(t, ctx, manager, "executeBulkWrite", data.NewStringValue("test.reject"), bulk)
	if ctl == nil {
		t.Fatal("duplicate key suppressed")
	}
	e := ctl.(*data.ThrowValue).Object
	if e.GetName() != exceptionNS+"BulkWriteException" || mustCall(t, ctx, e, "getCode").AsString() != "11000" {
		t.Fatal("bulk exception type/code")
	}
	partial := mustCall(t, ctx, e, "getWriteResult").(*data.ClassValue)
	if mustCall(t, ctx, partial, "getInsertedCount").AsString() != "0" {
		t.Fatal("partial write result")
	}
	seen := map[string]bool{}
	for len(commands) > 0 {
		d := <-commands
		seen[strings.ToLower(d[0].Key)] = true
	}
	for _, name := range []string{"insert", "update", "delete", "find", "getmore", "ping", "failcommand"} {
		if !seen[name] {
			t.Errorf("missing wire command %s", name)
		}
	}
}
