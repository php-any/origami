package redis

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
)

type status string
type serverError string

// Bound peer-controlled lengths, total decoded size and nesting before allocation.
const maxReplyBytes = 64 << 20
const maxReplyDepth = 128

func readReply(r *bufio.Reader) (any, error) {
	budget := maxReplyBytes
	return readValue(r, &budget, 0)
}
func readValue(r *bufio.Reader, budget *int, depth int) (any, error) {
	if depth > maxReplyDepth || *budget < 64 {
		return nil, fmt.Errorf("Redis reply exceeds decoding limits")
	}
	*budget -= 64
	line, err := r.ReadSlice('\n')
	if err != nil {
		return nil, fmt.Errorf("Redis response line: %w", err)
	}
	if len(line) < 3 || line[len(line)-2] != '\r' {
		return nil, fmt.Errorf("malformed Redis response line")
	}
	kind := line[0]
	text := string(line[1 : len(line)-2])
	*budget -= len(line)
	switch kind {
	case '+':
		return status(text), nil
	case '-':
		return serverError(text), nil
	case ':':
		n, e := strconv.ParseInt(text, 10, 64)
		if e != nil {
			return nil, fmt.Errorf("invalid Redis integer")
		}
		return n, nil
	case '$', '*':
		n, e := strconv.ParseInt(text, 10, 64)
		if e != nil || n < -1 {
			return nil, fmt.Errorf("invalid Redis response length")
		}
		if n == -1 {
			return nil, nil
		}
		if kind == '$' {
			if n > int64(*budget)-2 {
				return nil, fmt.Errorf("Redis bulk reply exceeds decoding limit")
			}
			*budget -= int(n) + 2
			b := make([]byte, int(n)+2)
			if _, e = io.ReadFull(r, b); e != nil {
				return nil, e
			}
			if b[n] != '\r' || b[n+1] != '\n' {
				return nil, fmt.Errorf("malformed Redis bulk terminator")
			}
			return string(b[:n]), nil
		}
		if n > int64(*budget)/64 {
			return nil, fmt.Errorf("Redis array reply exceeds decoding limit")
		}
		a := make([]any, int(n))
		for i := range a {
			a[i], e = readValue(r, budget, depth+1)
			if e != nil {
				return nil, e
			}
		}
		return a, nil
	default:
		return nil, fmt.Errorf("unsupported Redis reply type %q", kind)
	}
}
func writeCommand(w *bufio.Writer, args []string) error {
	if _, err := fmt.Fprintf(w, "*%d\r\n", len(args)); err != nil {
		return err
	}
	for _, s := range args {
		if _, err := fmt.Fprintf(w, "$%d\r\n", len(s)); err != nil {
			return err
		}
		if _, err := w.WriteString(s); err != nil {
			return err
		}
		if _, err := w.WriteString("\r\n"); err != nil {
			return err
		}
	}
	return nil
}
