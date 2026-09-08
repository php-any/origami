<?php
/**
 * PHP 特征补充测试
 * 测试新增的 PHP 函数支持：
 * - 数学函数 (sqrt, exp, log, 三角函数, 双曲函数等)
 * - 字符串函数 (str_word_count, wordwrap, addslashes, nl2br等)
 * - 数组函数 (natsort, natcasesort, array_intersect_ukey)
 * - 日期时间函数 (date, mktime, getdate, checkdate)
 * - 输出函数 (print_r, printf, vprintf)
 * - 文件系统函数 (is_executable, touch, sys_get_temp_dir)
 * - 杂项函数 (usleep, base_convert, bindec)
 */

Log::info("===== PHP 数学函数测试 =====");

// sqrt
$v = sqrt(16);
assert($v == 4, "sqrt(16) 应为 4，实际为 {$v}");
Log::info("[PASS] sqrt(16) = {$v}");

// cbrt
$v = cbrt(27);
assert($v == 3, "cbrt(27) 应为 3，实际为 {$v}");
Log::info("[PASS] cbrt(27) = {$v}");

// exp
$v = exp(1);
assert(abs($v - 2.71828) < 0.001, "exp(1) 应约等于 2.71828");
Log::info("[PASS] exp(1) = {$v}");

// log (自然对数)
$v = log(exp(1));
assert(abs($v - 1) < 0.001, "log(e) 应约等于 1");
Log::info("[PASS] log(e) = {$v}");

// log10
$v = log10(1000);
assert($v == 3, "log10(1000) 应为 3");
Log::info("[PASS] log10(1000) = {$v}");

// log1p
$v = log1p(exp(1) - 1);
assert(abs($v - 1) < 0.001, "log1p(e-1) 应约等于 1");
Log::info("[PASS] log1p(e-1) = {$v}");

// expm1
$v = expm1(1);
assert(abs($v - (exp(1) - 1)) < 0.001, "expm1(1) 应约等于 e-1");
Log::info("[PASS] expm1(1) = {$v}");

// pi
$v = pi();
assert(abs($v - 3.14159) < 0.001, "pi() 应约等于 3.14159");
Log::info("[PASS] pi() = {$v}");

// sin
$v = sin(0);
assert($v == 0, "sin(0) 应为 0");
Log::info("[PASS] sin(0) = {$v}");

// cos
$v = cos(0);
assert($v == 1, "cos(0) 应为 1");
Log::info("[PASS] cos(0) = {$v}");

// tan
$v = tan(0);
assert($v == 0, "tan(0) 应为 0");
Log::info("[PASS] tan(0) = {$v}");

// acos
$v = acos(1);
assert($v == 0, "acos(1) 应为 0");
Log::info("[PASS] acos(1) = {$v}");

// asin
$v = asin(0);
assert($v == 0, "asin(0) 应为 0");
Log::info("[PASS] asin(0) = {$v}");

// atan
$v = atan(0);
assert($v == 0, "atan(0) 应为 0");
Log::info("[PASS] atan(0) = {$v}");

// atan2
$v = atan2(0, 1);
assert($v == 0, "atan2(0, 1) 应为 0");
Log::info("[PASS] atan2(0, 1) = {$v}");

// cosh
$v = cosh(0);
assert($v == 1, "cosh(0) 应为 1");
Log::info("[PASS] cosh(0) = {$v}");

// sinh
$v = sinh(0);
assert($v == 0, "sinh(0) 应为 0");
Log::info("[PASS] sinh(0) = {$v}");

// tanh
$v = tanh(0);
assert($v == 0, "tanh(0) 应为 0");
Log::info("[PASS] tanh(0) = {$v}");

// hypot
$v = hypot(3, 4);
assert($v == 5, "hypot(3, 4) 应为 5");
Log::info("[PASS] hypot(3, 4) = {$v}");

// fmod
$v = fmod(5, 2);
assert($v == 1, "fmod(5, 2) 应为 1");
Log::info("[PASS] fmod(5, 2) = {$v}");

// deg2rad
$v = deg2rad(180);
assert(abs($v - pi()) < 0.0001, "deg2rad(180) 应约等于 pi()");
Log::info("[PASS] deg2rad(180) = {$v}");

// rad2deg
$v = rad2deg(pi());
assert(abs($v - 180) < 0.0001, "rad2deg(pi) 应约等于 180");
Log::info("[PASS] rad2deg(pi) = {$v}");

// base_convert
$v = base_convert("FF", 16, 2);
assert($v == "11111111", "base_convert('FF', 16, 2) 应为 '11111111'");
Log::info("[PASS] base_convert('FF', 16, 2) = {$v}");

// bindec
$v = bindec("1010");
assert($v == 10, "bindec('1010') 应为 10");
Log::info("[PASS] bindec('1010') = {$v}");

Log::info("===== PHP 字符串函数测试 =====");

// str_word_count
$v = str_word_count("Hello World PHP");
assert($v == 3, "str_word_count('Hello World PHP') 应为 3");
Log::info("[PASS] str_word_count = {$v}");

// str_word_count format=1
$words = str_word_count("Hello World!", 1);
assert(is_array($words), "str_word_count format=1 应返回数组");
assert(count($words) == 2, "str_word_count('Hello World!', 1) 应返回2个单词");
Log::info("[PASS] str_word_count format=1 count = " . count($words));

// wordwrap
$wrapped = wordwrap("Hello World", 5, "\n");
assert(strpos($wrapped, "\n") !== false, "wordwrap 应包含换行符");
Log::info("[PASS] wordwrap 测试通过");

// str_rot13
$v = str_rot13("Hello");
assert($v == "Uryyb", "str_rot13('Hello') 应为 'Uryyb'");
Log::info("[PASS] str_rot13('Hello') = {$v}");

// addslashes
$v = addslashes("It's a test");
assert($v == "It\\'s a test", "addslashes 应转义单引号");
Log::info("[PASS] addslashes 测试通过");

// chop
$v = chop("Hello   ");
assert($v == "Hello", "chop('Hello   ') 应为 'Hello'");
Log::info("[PASS] chop 测试通过");

// quotemeta
$v = quotemeta("a.b*c");
assert($v == "a\\.b\\*c", "quotemeta('a.b*c') 应为 'a\\.b\\*c'");
Log::info("[PASS] quotemeta 测试通过");

// nl2br
$v = nl2br("line1\nline2");
assert(strpos($v, "<br />") !== false, "nl2br 应插入 <br />");
Log::info("[PASS] nl2br 测试通过");

// substr_compare
$v = substr_compare("abcdef", "cde", 2, 3);
assert($v == 0, "substr_compare('abcdef', 'cde', 2, 3) 应为 0");
Log::info("[PASS] substr_compare 测试通过");

// htmlentities
$v = htmlentities("<b>test</b>");
assert($v == "&lt;b&gt;test&lt;/b&gt;", "htmlentities 应转义 HTML");
Log::info("[PASS] htmlentities 测试通过");

// htmlspecialchars_decode
$v = htmlspecialchars_decode("&lt;b&gt;");
assert($v == "<b>", "htmlspecialchars_decode 应解码 HTML");
Log::info("[PASS] htmlspecialchars_decode 测试通过");

// soundex
$v = soundex("Euler");
assert($v == "E460", "soundex('Euler') 应为 'E460'");
Log::info("[PASS] soundex('Euler') = {$v}");

// similar_text
$v = similar_text("Hello", "Hello World");
assert($v == 5, "similar_text('Hello', 'Hello World') 应为 5");
Log::info("[PASS] similar_text = {$v}");

// str_shuffle
$s = "abcde";
$shuffled = str_shuffle($s);
assert(strlen($shuffled) == 5, "str_shuffle 长度应保持");
Log::info("[PASS] str_shuffle 测试通过");

Log::info("===== PHP 数组函数测试 =====");

// natsort
$arr = ["img12.png", "img10.png", "img2.png", "img1.png"];
natsort($arr);
$sorted = implode(", ", $arr);
assert($sorted == "img1.png, img2.png, img10.png, img12.png", "natsort 自然排序结果不正确: {$sorted}");
Log::info("[PASS] natsort: {$sorted}");

// natcasesort
$arr2 = ["B", "a", "C", "b"];
natcasesort($arr2);
$sorted2 = implode(", ", $arr2);
Log::info("[PASS] natcasesort: {$sorted2}");

// array_intersect_ukey
$array1 = ['blue' => 1, 'red' => 2, 'green' => 3, 'purple' => 4];
$array2 = ['green' => 5, 'blue' => 6, 'yellow' => 7, 'cyan' => 8];
$result = array_intersect_ukey($array1, $array2, function($key1, $key2) {
    if ($key1 == $key2) return 0;
    return ($key1 < $key2) ? -1 : 1;
});
$keys = array_keys($result);
assert(in_array("blue", $keys) && in_array("green", $keys), "array_intersect_ukey 应返回 blue 和 green");
Log::info("[PASS] array_intersect_ukey: " . implode(",", $keys));

Log::info("===== PHP 日期时间函数测试 =====");

// date
$dateStr = date("Y-m-d");
$year = intval(date("Y"));
assert($year >= 2000, "date('Y') 应返回当前年份，实际为 {$year}");
Log::info("[PASS] date('Y-m-d') = {$dateStr}");

// mktime
$ts = mktime(0, 0, 0, 1, 1, 2024);
assert($ts == 1704038400, "mktime(0,0,0,1,1,2024) 应为 1704038400");
Log::info("[PASS] mktime = {$ts}");

// checkdate
assert(checkdate(2, 29, 2024) === true, "checkdate(2,29,2024) 应为 true");
assert(checkdate(2, 29, 2023) === false, "checkdate(2,29,2023) 应为 false");
Log::info("[PASS] checkdate 测试通过");

// getdate
$g = getdate();
assert(is_array($g), "getdate() 应返回数组");
assert(array_key_exists("year", $g), "getdate() 应包含 year 键");
assert(array_key_exists("month", $g), "getdate() 应包含 month 键");
Log::info("[PASS] getdate 测试通过, year = " . $g["year"]);

Log::info("===== PHP 输出函数测试 =====");

// print_r
$printed = print_r(["a" => 1, "b" => 2], true);
assert(strpos($printed, "a") !== false && strpos($printed, "b") !== false, "print_r 应包含数组键");
Log::info("[PASS] print_r 测试通过");

// printf
$len = printf("%s is %d years old", "John", 30);
assert($len > 0, "printf 应返回输出长度");
Log::info("[PASS] printf 测试通过, 输出长度 = {$len}");

// vprintf
$len = vprintf("%s-%s", ["a", "b"]);
assert($len > 0, "vprintf 应返回输出长度");
Log::info("[PASS] vprintf 测试通过, 输出长度 = {$len}");

Log::info("===== PHP 文件系统函数测试 =====");

// sys_get_temp_dir
$tmp = sys_get_temp_dir();
assert(is_string($tmp) && strlen($tmp) > 0, "sys_get_temp_dir() 应返回非空字符串");
Log::info("[PASS] sys_get_temp_dir = {$tmp}");

// is_executable
$isExec = is_executable("/bin/ls");
Log::info("[PASS] is_executable('/bin/ls') = " . ($isExec ? "true" : "false"));

// touch
$tmpFile = $tmp . "/origami_test_touch_" . time() . ".txt";
$touched = touch($tmpFile);
assert($touched === true, "touch 应返回 true");
assert(is_file($tmpFile), "touch 后文件应存在");
@unlink($tmpFile);
Log::info("[PASS] touch 测试通过");

Log::info("===== PHP 杂项函数测试 =====");

// usleep
$start = microtime(true);
usleep(10000);
$elapsed = microtime(true) - $start;
assert($elapsed >= 0.009, "usleep(10000) 应至少等待 10ms");
Log::info("[PASS] usleep 测试通过");

Log::info("🎉 所有 PHP 特征补充测试通过");
