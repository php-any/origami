package node

import (
	"github.com/php-any/origami/data"
	"io"
	"mime"
	"os"
	"path/filepath"
)

func multipartInput(ctx data.Context) (*data.ArrayValue, *data.ArrayValue) {
	post, files := data.NewArrayValueFromSlots(nil), data.NewArrayValueFromSlots(nil)
	request := getHTTPRequest(ctx)
	if request == nil || request.Method != "POST" {
		return post, files
	}
	// The request slots cache both views together, independently of which
	// superglobal Symfony reads first.
	postSlot, filesSlot := ctx.GetVM().EnsureGlobalZVal("_POST"), ctx.GetVM().EnsureGlobalZVal("_FILES")
	if postSlot.Defined && filesSlot.Defined {
		if value, ok := postSlot.ReadValue().(*data.ArrayValue); ok {
			post = value
		}
		if value, ok := filesSlot.ReadValue().(*data.ArrayValue); ok {
			files = value
		}
		return post, files
	}
	reader, err := request.MultipartReader()
	if err == nil {
		for {
			part, err := reader.NextPart()
			if err != nil {
				break
			}
			_, disposition, _ := mime.ParseMediaType(part.Header.Get("Content-Disposition"))
			field := disposition["name"]
			filename, upload := disposition["filename"]
			if !upload {
				body, err := io.ReadAll(part)
				if err == nil {
					data.SetFormField(post, field, data.NewStringValue(string(body)))
				}
				_ = part.Close()
				continue
			}
			code, temporary, size := 0, "", int64(0)
			if filename == "" {
				code = 4
			} else {
				file, err := os.CreateTemp("", "origami-upload-*")
				if err != nil {
					code = 6
				} else {
					temporary = file.Name()
					size, err = io.Copy(file, part)
					closeErr := file.Close()
					if err != nil || closeErr != nil {
						code = 7
						_ = os.Remove(temporary)
						temporary = ""
					}
					if code == 0 {
						host, ok := ctx.GetVM().(interface{ RegisterUploadedFile(string) bool })
						if !ok || !host.RegisterUploadedFile(temporary) {
							_ = os.Remove(temporary)
							temporary = ""
							code = 7
						}
					}
				}
			}
			_ = part.Close()
			path := data.FormFieldPath(field)
			if len(path) == 0 {
				continue
			}
			fileType := part.Header.Get("Content-Type")
			name := filepath.Base(filename)
			if code != 0 {
				name = filename
				fileType = ""
				size = 0
			}
			fields := []struct {
				name  string
				value data.Value
			}{
				{"name", data.NewStringValue(name)}, {"full_path", data.NewStringValue(filename)},
				{"type", data.NewStringValue(fileType)}, {"tmp_name", data.NewStringValue(temporary)},
				{"error", data.NewIntValue(code)}, {"size", data.NewIntValue(int(size))},
			}
			for _, field := range fields {
				keys := append([]string{path[0], field.name}, path[1:]...)
				data.SetFormFieldPath(files, keys, field.value)
			}
		}
	}
	if !postSlot.Defined {
		data.CowAssign(postSlot, post)
	}
	if !filesSlot.Defined {
		data.CowAssign(filesSlot, files)
	}
	return post, files
}
