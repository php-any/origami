package openai

import (
	"context"
	"errors"
	api "github.com/openai/openai-go/v3"
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/parser"
	"github.com/php-any/origami/runtime"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestClientOptionsPreserveNestedPHPArrays(t *testing.T) {
	vm := runtime.NewVM(parser.NewParser())
	ctx := vm.CreateContext([]data.Variable{data.NewVariable("options", 0, nil)})
	options := data.NewArrayValueFromSlots(nil)
	options.SetStringKey("temperature", data.NewFloatValue(0.5))
	options.SetStringKey("stop", data.NewArrayValue([]data.Value{data.NewStringValue("END")}))
	schema := data.NewArrayValueFromSlots(nil)
	schema.SetStringKey("type", data.NewStringValue("object"))
	options.SetStringKey("schema", schema)
	ctx.SetVariableValue(data.NewVariable("options", 0, nil), options)
	result, ctl := clientOptions(ctx, 0)
	if ctl != nil || result["temperature"] != 0.5 || result["stop"].([]any)[0] != "END" || result["schema"].(map[string]any)["type"] != "object" {
		t.Fatalf("options: %#v / %v", result, ctl)
	}
}

func TestAllClientOperationsHonorCancellation(t *testing.T) {
	input := filepath.Join(t.TempDir(), "input.wav")
	if err := os.WriteFile(input, []byte("audio"), 0600); err != nil {
		t.Fatal(err)
	}
	for _, operation := range []string{"chat", "embeddings", "images", "speech", "transcription"} {
		t.Run(operation, func(t *testing.T) {
			started := make(chan struct{})
			stop := make(chan struct{})
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				_, _ = io.Copy(io.Discard, r.Body)
				if operation == "speech" {
					w.WriteHeader(200)
					w.Write([]byte("partial"))
					w.(http.Flusher).Flush()
				}
				close(started)
				select {
				case <-r.Context().Done():
				case <-stop:
				}
			}))
			defer server.Close()
			defer close(stop)
			client, err := newClient("test", server.URL+"/")
			if err != nil {
				t.Fatal(err)
			}
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			output := filepath.Join(t.TempDir(), "output.wav")
			done := make(chan error, 1)
			go func() {
				var err error
				switch operation {
				case "chat":
					_, err = client.chat(ctx, "model", []api.ChatCompletionMessageParamUnion{api.UserMessage("hello")}, nil)
				case "embeddings":
					_, err = client.embeddings(ctx, "model", api.EmbeddingNewParamsInputUnion{OfString: api.String("hello")}, nil)
				case "images":
					_, err = client.images(ctx, "model", "hello", nil)
				case "speech":
					err = client.speech(ctx, "model", "hello", "alloy", output, nil)
				case "transcription":
					_, err = client.transcription(ctx, "model", input, nil)
				}
				done <- err
			}()
			select {
			case <-started:
			case <-time.After(3 * time.Second):
				t.Fatal("request did not start")
			}
			cancel()
			select {
			case err := <-done:
				if !errors.Is(err, context.Canceled) {
					t.Fatalf("cancellation: %v", err)
				}
			case <-time.After(3 * time.Second):
				t.Fatal("operation ignored cancellation")
			}
			if _, err := os.Stat(output); !os.IsNotExist(err) {
				t.Fatal("canceled speech wrote a partial output")
			}
		})
	}
}
