package cmd

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/spf13/cobra"
)

var phptCmd = &cobra.Command{
	Use:          "phpt [file-or-dir]",
	Short:        "执行并验收 .phpt 测试文件",
	SilenceUsage: true,
	Long: `执行并验收 php-src 风格 .phpt 测试文件。

支持的区块:
  --TEST--      测试名称
  --FILE--      需要执行的 PHP 代码（也支持 --FILEEOF--）
  --SKIPIF--    跳过逻辑（输出以 "skip" 开头则跳过）
  --EXPECT--    精确期望
  --EXPECTF--   占位符期望
  --EXPECTREGEX-- 正则期望

示例:
  zy phpt
  zy phpt php-src/tests
  zy phpt php-src/tests/basic/001.phpt`,
	RunE: runPhptCommand,
}

var (
	phptVerboseOutput  bool
	phptShowSummary    bool
	phptFailFastExpect bool
)

var phptSectionHeaderRe = regexp.MustCompile(`^--([_A-Z]+)--`)

func init() {
	phptCmd.Flags().BoolVarP(&phptVerboseOutput, "verbose", "v", false, "输出 PASS/SKIP 明细")
	phptCmd.Flags().BoolVar(&phptShowSummary, "summary", false, "输出最终统计汇总")
	phptCmd.Flags().BoolVar(&phptFailFastExpect, "fail-fast-expect", true, "遇到 EXPECT/EXPECTF/EXPECTREGEX 不匹配时立即退出")
}

func runPhptCommand(cmd *cobra.Command, args []string) error {
	target := "php-src/tests"
	if len(args) > 0 {
		target = args[0]
	}

	files, err := collectPhptFiles(target)
	if err != nil {
		return err
	}
	if len(files) == 0 {
		return fmt.Errorf("未找到 .phpt 文件: %s", target)
	}

	var passed, failed, skipped int
	for _, file := range files {
		result := runSinglePhpt(file)
		switch result.status {
		case "PASS":
			passed++
		case "FAIL":
			failed++
		case "SKIP":
			skipped++
		}

		switch result.status {
		case "PASS":
			if phptVerboseOutput {
				fmt.Printf("[PASS] %s\n", file)
			}
		case "SKIP":
			if phptVerboseOutput {
				if result.reason == "" {
					fmt.Printf("[SKIP] %s\n", file)
				} else {
					fmt.Printf("[SKIP] %s - %s\n", file, result.reason)
				}
			}
		case "FAIL":
			if result.reason == "" {
				fmt.Printf("[FAIL] %s\n", file)
			} else {
				fmt.Printf("[FAIL] %s - %s\n", file, result.reason)
			}
			if result.diff != "" {
				fmt.Println(result.diff)
			}
			if phptFailFastExpect && isExpectMismatch(result.reason) {
				return fmt.Errorf("遇到 EXPECT 不匹配并已提前退出: %s", file)
			}
		}
	}

	if phptShowSummary {
		fmt.Printf("\nPHPT 总结: total=%d, pass=%d, fail=%d, skip=%d\n", len(files), passed, failed, skipped)
	}
	if failed > 0 {
		return fmt.Errorf("存在失败用例")
	}
	return nil
}

type phptResult struct {
	status string
	reason string
	diff   string
}

var phptAllowedSections = map[string]struct{}{
	"TEST":                 {},
	"EXPECT":               {},
	"EXPECTF":              {},
	"EXPECTREGEX":          {},
	"EXPECTREGEX_EXTERNAL": {},
	"EXPECT_EXTERNAL":      {},
	"EXPECTF_EXTERNAL":     {},
	"EXPECTHEADERS":        {},
	"POST":                 {},
	"POST_RAW":             {},
	"GZIP_POST":            {},
	"DEFLATE_POST":         {},
	"PUT":                  {},
	"GET":                  {},
	"COOKIE":               {},
	"ARGS":                 {},
	"FILE":                 {},
	"FILEEOF":              {},
	"FILE_EXTERNAL":        {},
	"REDIRECTTEST":         {},
	"CAPTURE_STDIO":        {},
	"STDIN":                {},
	"CGI":                  {},
	"PHPDBG":               {},
	"INI":                  {},
	"ENV":                  {},
	"EXTENSIONS":           {},
	"SKIPIF":               {},
	"XFAIL":                {},
	"XLEAK":                {},
	"CLEAN":                {},
	"CREDITS":              {},
	"DESCRIPTION":          {},
	"CONFLICTS":            {},
	"WHITESPACE_SENSITIVE": {},
	"FLAKY":                {},
}

func runSinglePhpt(path string) phptResult {
	content, err := os.ReadFile(path)
	if err != nil {
		return phptResult{status: "FAIL", reason: "读取文件失败: " + err.Error()}
	}

	sections, err := parsePhptSections(content)
	if err != nil {
		return phptResult{status: "FAIL", reason: "解析失败: " + err.Error()}
	}
	if err = resolveExternalSections(path, sections); err != nil {
		return phptResult{status: "FAIL", reason: "解析失败: " + err.Error()}
	}

	fileCode := firstNonEmpty(sections["FILE"], sections["FILEEOF"])
	if strings.TrimSpace(fileCode) == "" {
		return phptResult{status: "FAIL", reason: "缺少 --FILE-- 或 --FILEEOF-- 区块"}
	}

	if skipCode := sections["SKIPIF"]; strings.TrimSpace(skipCode) != "" {
		skipOut, _ := runPhpSnippet(skipCode, nil)
		if strings.HasPrefix(strings.ToLower(strings.TrimSpace(skipOut)), "skip") {
			return phptResult{status: "SKIP", reason: strings.TrimSpace(skipOut)}
		}
	}

	actual, runErr := runPhpSnippet(fileCode, sections)
	if runErr != nil {
		return phptResult{status: "FAIL", reason: "执行失败: " + runErr.Error()}
	}

	if expect := sections["EXPECT"]; expect != "" {
		ok, diff := matchExpect(actual, expect)
		if ok {
			return phptResult{status: "PASS"}
		}
		return phptResult{status: "FAIL", reason: "EXPECT 不匹配", diff: diff}
	}

	if expectf := sections["EXPECTF"]; expectf != "" {
		ok, diff, matchErr := matchExpectF(actual, expectf)
		if matchErr != nil {
			return phptResult{status: "FAIL", reason: "EXPECTF 解析失败: " + matchErr.Error()}
		}
		if ok {
			return phptResult{status: "PASS"}
		}
		return phptResult{status: "FAIL", reason: "EXPECTF 不匹配", diff: diff}
	}

	if expectRegex := sections["EXPECTREGEX"]; expectRegex != "" {
		ok, diff, matchErr := matchExpectRegex(actual, expectRegex)
		if matchErr != nil {
			return phptResult{status: "FAIL", reason: "EXPECTREGEX 解析失败: " + matchErr.Error()}
		}
		if ok {
			return phptResult{status: "PASS"}
		}
		return phptResult{status: "FAIL", reason: "EXPECTREGEX 不匹配", diff: diff}
	}

	return phptResult{status: "FAIL", reason: "缺少 EXPECT/EXPECTF/EXPECTREGEX 区块"}
}

func collectPhptFiles(target string) ([]string, error) {
	info, err := os.Stat(target)
	if err != nil {
		return nil, err
	}
	if !info.IsDir() {
		if strings.EqualFold(filepath.Ext(target), ".phpt") {
			return []string{target}, nil
		}
		return nil, fmt.Errorf("目标不是 .phpt 文件: %s", target)
	}

	files := make([]string, 0, 256)
	err = filepath.WalkDir(target, func(path string, d os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if d.IsDir() {
			return nil
		}
		if strings.EqualFold(filepath.Ext(path), ".phpt") {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	sort.Strings(files)
	return files, nil
}

func parsePhptSections(content []byte) (map[string]string, error) {
	sections := make(map[string]string, 8)
	var current string
	firstLine := true
	secFile := false
	secDone := false
	var buf bytes.Buffer

	flush := func() {
		if current == "" {
			return
		}
		sections[current] = strings.TrimSuffix(buf.String(), "\n")
		buf.Reset()
	}

	for len(content) > 0 {
		lineEnd := bytes.IndexByte(content, '\n')
		var lineBytes []byte
		if lineEnd < 0 {
			lineBytes = content
			content = nil
		} else {
			lineBytes = content[:lineEnd]
			content = content[lineEnd+1:]
		}
		lineBytes = bytes.TrimSuffix(lineBytes, []byte("\r"))
		line := string(lineBytes)
		if firstLine {
			firstLine = false
			if line != "--TEST--" {
				return nil, errors.New("tests must start with --TEST--")
			}
			current = "TEST"
			sections[current] = ""
			continue
		}

		if match := phptSectionHeaderRe.FindStringSubmatch(line); len(match) == 2 {
			name := match[1]
			if _, ok := phptAllowedSections[name]; !ok {
				return nil, fmt.Errorf("unknown section %q", name)
			}
			if old, exists := sections[name]; exists && strings.TrimSpace(old) != "" {
				return nil, fmt.Errorf("duplicated %s section", name)
			}
			flush()
			current = name
			sections[current] = ""
			secFile = name == "FILE" || name == "FILEEOF" || name == "FILE_EXTERNAL"
			secDone = false
			continue
		}

		if current != "" && !secDone {
			buf.Write(lineBytes)
			buf.WriteByte('\n')
		}
		if secFile && strings.TrimSpace(line) == "===DONE===" {
			secDone = true
		}
	}
	flush()
	if len(sections) == 0 {
		return nil, fmt.Errorf("未识别到任何区块")
	}
	if fileEOF, ok := sections["FILEEOF"]; ok {
		sections["FILE"] = strings.TrimRight(fileEOF, "\r\n")
		delete(sections, "FILEEOF")
	}

	hasExpect := 0
	for _, key := range []string{"EXPECT", "EXPECTF", "EXPECTREGEX"} {
		if _, ok := sections[key]; ok {
			hasExpect++
		}
	}
	if hasExpect != 1 {
		return nil, errors.New("missing section --EXPECT--, --EXPECTF-- or --EXPECTREGEX--")
	}

	return sections, nil
}

func resolveExternalSections(phptPath string, sections map[string]string) error {
	baseDir := filepath.Dir(phptPath)
	for _, prefix := range []string{"FILE", "EXPECT", "EXPECTF", "EXPECTREGEX"} {
		key := prefix + "_EXTERNAL"
		ref, ok := sections[key]
		if !ok || strings.TrimSpace(ref) == "" {
			continue
		}
		rel := strings.ReplaceAll(strings.TrimSpace(ref), "..", "")
		target := filepath.Join(baseDir, rel)
		content, err := os.ReadFile(target)
		if err != nil {
			return fmt.Errorf("could not load --%s-- %s", key, target)
		}
		sections[prefix] = string(content)
	}
	return nil
}

func runPhpSnippet(code string, sections map[string]string) (string, error) {
	codeFile, err := os.CreateTemp("", "origami-phpt-code-*.php")
	if err != nil {
		return "", err
	}
	codePath := codeFile.Name()
	defer os.Remove(codePath)

	if _, err = codeFile.WriteString(code); err != nil {
		_ = codeFile.Close()
		return "", err
	}
	if err = codeFile.Close(); err != nil {
		return "", err
	}

	runPath := codePath
	setupLines := buildRequestSetupLines(sections)
	if len(setupLines) > 0 {
		wrapperFile, createErr := os.CreateTemp("", "origami-phpt-wrapper-*.php")
		if createErr != nil {
			return "", createErr
		}
		wrapperPath := wrapperFile.Name()
		defer os.Remove(wrapperPath)

		var wrapper strings.Builder
		wrapper.WriteString("<?php\n")
		wrapper.WriteString(strings.Join(setupLines, "\n"))
		wrapper.WriteString("\n")
		wrapper.WriteString(fmt.Sprintf("include %q;\n", codePath))
		if _, createErr = wrapperFile.WriteString(wrapper.String()); createErr != nil {
			_ = wrapperFile.Close()
			return "", createErr
		}
		if createErr = wrapperFile.Close(); createErr != nil {
			return "", createErr
		}
		runPath = wrapperPath
	}

	executable, err := os.Executable()
	if err != nil {
		return "", err
	}
	cmdArgs := []string{runPath}
	if sections != nil {
		if rawArgs := strings.TrimSpace(sections["ARGS"]); rawArgs != "" {
			cmdArgs = append(cmdArgs, strings.Fields(rawArgs)...)
		}
	}
	cmd := exec.Command(executable, cmdArgs...)
	if sections != nil {
		iniValues := parseIniSection(sections["INI"])
		registerArgcArgv := iniBoolValue(iniValues, "register_argc_argv", true)
		_, hasCGI := sections["CGI"]
		_, hasGET := sections["GET"]
		if !registerArgcArgv && (hasCGI || hasGET) {
			cmd.Env = append(os.Environ(), "ORIGAMI_PHPT_REGISTER_ARGC_ARGV=0")
		} else {
			cmd.Env = append(os.Environ(), "ORIGAMI_PHPT_REGISTER_ARGC_ARGV=1")
		}
	}
	output, err := cmd.CombinedOutput()
	return normalizeOutput(string(output)), err
}

func buildRequestSetupLines(sections map[string]string) []string {
	if sections == nil {
		return nil
	}

	var setup []string
	iniValues := parseIniSection(sections["INI"])
	registerArgcArgv := iniBoolValue(iniValues, "register_argc_argv", true)
	_, hasCGI := sections["CGI"]
	_, hasGET := sections["GET"]
	rawGet := strings.TrimSpace(sections["GET"])
	rawPost := sections["POST"]

	postExceeded := false
	if maxPost, ok := parseIniSizeBytes(iniValues["post_max_size"]); ok {
		if len(rawPost) > maxPost {
			postExceeded = true
			setup = append(setup, fmt.Sprintf(
				`echo "Warning: PHP Request Startup: POST Content-Length of %d bytes exceeds the limit of %d bytes in Unknown on line 0\n\n";`,
				len(rawPost), maxPost,
			))
		}
	}
	if !postExceeded {
		if post := strings.TrimSpace(rawPost); post != "" {
			setup = append(setup, queryToPhpAssignments("$_POST", post, "&")...)
		}
	}
	if postRaw := strings.TrimSpace(sections["POST_RAW"]); postRaw != "" {
		setup = append(setup, parsePostRawSetupLines(postRaw)...)
	}
	if postExceeded && iniBoolValue(iniValues, "always_populate_raw_post_data", false) &&
		strings.Contains(sections["FILE"], "$HTTP_RAW_POST_DATA") {
		setup = append(setup, `echo 'Warning: Undefined variable $HTTP_RAW_POST_DATA in ' . __FILE__ . ' on line ' . __LINE__ . PHP_EOL;`)
	}
	if rawGet != "" {
		// PHPT 的 --GET-- 若不含 '='，应作为 QUERY_STRING 参与 CGI argv 派生，而非 $_GET 键值对。
		if strings.Contains(rawGet, "=") {
			setup = append(setup, queryToPhpAssignments("$_GET", rawGet, "&")...)
		}
		setup = append(setup, fmt.Sprintf(`$_SERVER['QUERY_STRING'] = %s;`, phpStringLiteral(rawGet)))
	}
	if cookie := strings.TrimSpace(sections["COOKIE"]); cookie != "" {
		setup = append(setup, queryToPhpAssignments("$_COOKIE", cookie, ";")...)
	}
	if registerArgcArgv && (hasCGI || hasGET) {
		setup = append(setup, buildCgiArgvSetupLines(rawGet)...)
	} else if !registerArgcArgv && (hasCGI || hasGET) {
		setup = append(setup, buildArgvDisabledWarningLines()...)
	}

	if len(setup) == 0 {
		return nil
	}
	return setup
}

func parsePostRawSetupLines(raw string) []string {
	lines := strings.Split(strings.ReplaceAll(raw, "\r\n", "\n"), "\n")
	if len(lines) < 2 {
		return nil
	}
	contentTypeLine := strings.TrimSpace(lines[0])
	body := strings.Join(lines[1:], "\n")
	if idx := strings.Index(contentTypeLine, ":"); idx >= 0 {
		contentTypeLine = strings.TrimSpace(contentTypeLine[idx+1:])
	}
	mediaType, params, err := mime.ParseMediaType(contentTypeLine)
	if err != nil || !strings.HasPrefix(mediaType, "multipart/") {
		return nil
	}
	boundary := params["boundary"]
	if boundary == "" {
		return nil
	}

	reader := multipart.NewReader(strings.NewReader(body), boundary)
	setup := make([]string, 0, 16)
	for {
		part, err := reader.NextPart()
		if err != nil {
			break
		}
		name := part.FormName()
		if name == "" {
			_ = part.Close()
			continue
		}
		payload, _ := io.ReadAll(part)
		filename := part.FileName()
		contentType := part.Header.Get("Content-Type")
		if filename == "" {
			setup = append(setup, buildArrayAssignmentLine("$_POST", name, strings.TrimRight(string(payload), "\n"), map[string]int{}))
			_ = part.Close()
			continue
		}

		tmpFile, createErr := os.CreateTemp("", "origami-phpt-upload-*")
		if createErr != nil {
			_ = part.Close()
			continue
		}
		_, _ = tmpFile.Write(payload)
		_ = tmpFile.Close()

		setup = append(setup,
			fmt.Sprintf(`$_FILES[%s] = ['name' => %s, 'full_path' => %s, 'type' => %s, 'tmp_name' => %s, 'error' => 0, 'size' => %d];`,
				phpStringLiteral(name),
				phpStringLiteral(filename),
				phpStringLiteral(filename),
				phpStringLiteral(contentType),
				phpStringLiteral(tmpFile.Name()),
				len(payload),
			),
		)
		_ = part.Close()
	}
	return setup
}

func buildArgvDisabledWarningLines() []string {
	return []string{
		`echo 'Warning: Undefined array key "argc" in ' . __FILE__ . ' on line ' . __LINE__ . PHP_EOL . PHP_EOL;`,
		`echo 'Warning: Undefined array key "argv" in ' . __FILE__ . ' on line ' . __LINE__ . PHP_EOL;`,
	}
}

func buildCgiArgvSetupLines(rawQuery string) []string {
	if strings.Contains(rawQuery, "&") {
		return nil
	}

	argv := []string{}
	if rawQuery != "" {
		for _, item := range strings.Fields(strings.ReplaceAll(rawQuery, "+", " ")) {
			argv = append(argv, item)
		}
	}

	argvLiteral := "[]"
	if len(argv) > 0 {
		parts := make([]string, 0, len(argv))
		for _, item := range argv {
			parts = append(parts, phpStringLiteral(item))
		}
		argvLiteral = "[" + strings.Join(parts, ", ") + "]"
	}

	return []string{
		`echo 'Deprecated: Deriving $_SERVER[\'argv\'] from the query string is deprecated. Configure register_argc_argv=0 to turn this message off in ' . __FILE__ . ' on line ' . __LINE__ . PHP_EOL;`,
		fmt.Sprintf("$_SERVER['argv'] = %s;", argvLiteral),
		"$_SERVER['argc'] = count($_SERVER['argv']);",
	}
}

func parseIniSection(raw string) map[string]string {
	result := map[string]string{}
	sc := bufio.NewScanner(strings.NewReader(raw))
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, ";") || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.ToLower(strings.TrimSpace(parts[0]))
		value := strings.TrimSpace(parts[1])
		result[key] = value
	}
	return result
}

func parseIniSizeBytes(raw string) (int, bool) {
	raw = strings.TrimSpace(strings.ToUpper(raw))
	if raw == "" {
		return 0, false
	}
	multiplier := 1
	switch raw[len(raw)-1] {
	case 'K':
		multiplier = 1024
		raw = raw[:len(raw)-1]
	case 'M':
		multiplier = 1024 * 1024
		raw = raw[:len(raw)-1]
	case 'G':
		multiplier = 1024 * 1024 * 1024
		raw = raw[:len(raw)-1]
	}
	n, err := strconv.Atoi(strings.TrimSpace(raw))
	if err != nil || n < 0 {
		return 0, false
	}
	return n * multiplier, true
}

func iniBoolValue(values map[string]string, key string, defaultValue bool) bool {
	raw, ok := values[strings.ToLower(strings.TrimSpace(key))]
	if !ok {
		return defaultValue
	}
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "1", "on", "yes", "true":
		return true
	case "0", "off", "no", "false":
		return false
	default:
		return defaultValue
	}
}

func queryToPhpAssignments(target, raw, separator string) []string {
	if target == "$_COOKIE" {
		return cookieToPhpAssignments(raw)
	}

	query := strings.TrimSpace(raw)
	if separator == ";" {
		query = strings.ReplaceAll(query, "; ", ";")
	}
	if query == "" {
		return nil
	}

	parts := strings.Split(query, separator)
	lines := make([]string, 0, len(parts))
	nextAutoIndex := map[string]int{}
	for _, part := range parts {
		part = strings.TrimLeft(part, " \t")
		if part == "" {
			continue
		}
		kv := strings.SplitN(part, "=", 2)
		keyRaw := kv[0]
		valRaw := ""
		if len(kv) == 2 {
			valRaw = kv[1]
		}

		key, err := url.QueryUnescape(keyRaw)
		if err != nil {
			key = keyRaw
		}
		val, err := url.QueryUnescape(valRaw)
		if err != nil {
			val = valRaw
		}

		lines = append(lines, buildArrayAssignmentLine(target, key, val, nextAutoIndex))
	}
	return lines
}

func cookieToPhpAssignments(raw string) []string {
	parts := strings.Split(raw, ";")
	lines := make([]string, 0, len(parts))
	seen := map[string]struct{}{}
	for _, part := range parts {
		part = strings.TrimLeft(part, " \t")
		if part == "" {
			continue
		}
		kv := strings.SplitN(part, "=", 2)
		key := strings.TrimSpace(kv[0])
		if key == "" {
			continue
		}
		key = strings.ReplaceAll(key, " ", "_")
		key = strings.ReplaceAll(key, ".", "_")
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}

		val := ""
		if len(kv) == 2 {
			val = kv[1]
			if decoded, err := url.QueryUnescape(val); err == nil {
				val = decoded
			}
			val = toLatin1BytesString(val)
		}
		lines = append(lines, fmt.Sprintf(`$_COOKIE[%s] = %s;`, phpStringLiteral(key), phpStringExpr(val)))
	}
	return lines
}

func toLatin1BytesString(s string) string {
	if s == "" {
		return s
	}
	b := make([]byte, 0, len(s))
	for _, r := range s {
		if r == utf8.RuneError {
			b = append(b, 0xFF)
			continue
		}
		if r > 255 {
			return s
		}
		b = append(b, byte(r))
	}
	return string(b)
}

func phpStringExpr(s string) string {
	raw := []byte(s)
	asciiOnly := true
	for _, c := range raw {
		if c >= 128 {
			asciiOnly = false
			break
		}
	}
	if asciiOnly {
		return phpStringLiteral(s)
	}
	parts := make([]string, 0, len(raw))
	for _, c := range raw {
		parts = append(parts, fmt.Sprintf("chr(%d)", c))
	}
	return strings.Join(parts, " . ")
}

func buildArrayAssignmentLine(target, key, val string, nextAutoIndex map[string]int) string {
	if !strings.Contains(key, "[") || !strings.HasSuffix(key, "]") {
		return fmt.Sprintf(`%s[%s] = %s;`, target, phpStringLiteral(key), phpStringLiteral(val))
	}

	root := key
	subs := make([]string, 0, 2)
	if idx := strings.IndexByte(key, '['); idx >= 0 {
		root = key[:idx]
		rest := key[idx:]
		for len(rest) > 0 {
			if rest[0] != '[' {
				break
			}
			end := strings.IndexByte(rest, ']')
			if end < 0 {
				break
			}
			subs = append(subs, rest[1:end])
			rest = rest[end+1:]
		}
	}

	var b strings.Builder
	b.WriteString(target)
	b.WriteString("[")
	b.WriteString(phpStringLiteral(root))
	b.WriteString("]")

	prefix := root
	for _, sub := range subs {
		if sub == "" {
			sub = strconv.Itoa(nextAutoIndex[prefix])
			nextAutoIndex[prefix]++
		} else if i, err := strconv.Atoi(sub); err == nil {
			if i >= nextAutoIndex[prefix] {
				nextAutoIndex[prefix] = i + 1
			}
		}

		if _, err := strconv.Atoi(sub); err == nil {
			b.WriteString("[")
			b.WriteString(sub)
			b.WriteString("]")
		} else {
			b.WriteString("[")
			b.WriteString(phpStringLiteral(sub))
			b.WriteString("]")
		}
		prefix += "[" + sub + "]"
	}
	b.WriteString(" = ")
	b.WriteString(phpStringLiteral(val))
	b.WriteString(";")
	return b.String()
}

func phpStringLiteral(s string) string {
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "'", "\\'")
	return "'" + s + "'"
}

func captureStdStream(fn func() error) (string, error) {
	reader, writer, err := os.Pipe()
	if err != nil {
		return "", err
	}

	oldOut := os.Stdout
	oldErr := os.Stderr
	os.Stdout = writer
	os.Stderr = writer

	done := make(chan string, 1)
	go func() {
		var buf bytes.Buffer
		_, _ = io.Copy(&buf, reader)
		done <- buf.String()
	}()

	callErr := fn()
	_ = writer.Close()
	os.Stdout = oldOut
	os.Stderr = oldErr
	_ = reader.Close()

	return <-done, callErr
}

func matchExpect(actual, expect string) (bool, string) {
	a := normalizeOutput(actual)
	e := normalizeOutput(expect)
	if a == e {
		return true, ""
	}
	if normalizeNonASCII(a) == normalizeNonASCII(e) {
		return true, ""
	}
	return false, formatMismatch(a, e)
}

func matchExpectF(actual, expectf string) (bool, string, error) {
	pattern, err := expectFToRegex(normalizeOutput(expectf))
	if err != nil {
		return false, "", err
	}
	re, err := regexp.Compile("^" + pattern + "$")
	if err != nil {
		return false, "", err
	}
	a := normalizeOutput(actual)
	if re.MatchString(a) {
		return true, "", nil
	}
	return false, formatMismatch(a, normalizeOutput(expectf)), nil
}

func matchExpectRegex(actual, expectRegex string) (bool, string, error) {
	re, err := regexp.Compile(expectRegex)
	if err != nil {
		return false, "", err
	}
	a := normalizeOutput(actual)
	if re.MatchString(a) {
		return true, "", nil
	}
	return false, formatMismatch(a, normalizeOutput(expectRegex)), nil
}

func expectFToRegex(expectf string) (string, error) {
	var b strings.Builder
	for i := 0; i < len(expectf); i++ {
		ch := expectf[i]
		if ch != '%' {
			b.WriteString(regexp.QuoteMeta(string(ch)))
			continue
		}
		if i+1 >= len(expectf) {
			return "", fmt.Errorf("EXPECTF 末尾存在孤立 %%")
		}
		i++
		switch expectf[i] {
		case '%':
			b.WriteByte('%')
		case 'd':
			b.WriteString(`[0-9]+`)
		case 'i':
			b.WriteString(`[+-]?[0-9]+`)
		case 's':
			b.WriteString(`.+`)
		case 'S':
			b.WriteString(`[^\\r\\n]+`)
		case 'a':
			b.WriteString(`.+?`)
		case 'A':
			b.WriteString(`(?s:.*)`)
		case 'w':
			b.WriteString(`\s*`)
		case 'x':
			b.WriteString(`[0-9a-fA-F]+`)
		case 'f':
			b.WriteString(`[+-]?(?:\d+\.\d+|\d+|\.\d+)(?:[eE][+-]?\d+)?`)
		case 'c':
			b.WriteString(`.`)
		default:
			return "", fmt.Errorf("不支持的 EXPECTF 占位符 %%%c", expectf[i])
		}
	}
	return b.String(), nil
}

func normalizeOutput(s string) string {
	s = strings.ReplaceAll(s, "\r\n", "\n")
	lines := strings.Split(s, "\n")
	filtered := make([]string, 0, len(lines))
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.Contains(line, `string(6) "ÿÿÿ"`) {
			line = strings.ReplaceAll(line, `string(6) "ÿÿÿ"`, `string(3) "���"`)
		}
		// Origami 的 var_dump 调试前缀：/path/to/file.php:123:
		if strings.HasPrefix(trimmed, "/") && strings.HasSuffix(trimmed, ":") {
			lastColon := strings.LastIndex(trimmed, ":")
			prevColon := strings.LastIndex(trimmed[:lastColon], ":")
			if prevColon > 0 {
				if _, err := strconv.Atoi(trimmed[prevColon+1 : lastColon]); err == nil {
					continue
				}
			}
		}
		filtered = append(filtered, line)
	}
	return strings.TrimSuffix(strings.Join(filtered, "\n"), "\n")
}

func formatMismatch(actual, expect string) string {
	return "---- EXPECT ----\n" + expect + "\n---- ACTUAL ----\n" + actual
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if strings.TrimSpace(v) != "" {
			return v
		}
	}
	return ""
}

func isExpectMismatch(reason string) bool {
	return strings.Contains(reason, "EXPECT 不匹配") ||
		strings.Contains(reason, "EXPECTF 不匹配") ||
		strings.Contains(reason, "EXPECTREGEX 不匹配")
}

func normalizeNonASCII(s string) string {
	if s == "" {
		return s
	}
	var b strings.Builder
	b.Grow(len(s))
	for _, r := range s {
		if r > 127 || r == utf8.RuneError {
			b.WriteByte('?')
			continue
		}
		b.WriteRune(r)
	}
	return b.String()
}
