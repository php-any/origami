<?php

namespace tests\php;

/**
 * protowire 标准库 PHP 端功能验收测试。
 *
 * 覆盖：常量定义、基本类型编码/解析、多字段、嵌套消息、
 * group、packed、repeated、格式错误、深度限制、编码辅助函数。
 */

// ---------------------------------------------------------------------------
// 1. 常量定义
// ---------------------------------------------------------------------------

if (!defined('PROTOWIRE_WIRE_VARINT')) {
    Log::fatal('PROTOWIRE_WIRE_VARINT 常量未定义');
}
if (PROTOWIRE_WIRE_VARINT !== 0) {
    Log::fatal('PROTOWIRE_WIRE_VARINT expected 0, got ' . PROTOWIRE_WIRE_VARINT);
}
if (PROTOWIRE_WIRE_FIXED64 !== 1) {
    Log::fatal('PROTOWIRE_WIRE_FIXED64 expected 1, got ' . PROTOWIRE_WIRE_FIXED64);
}
if (PROTOWIRE_WIRE_LENGTH_DELIMITED !== 2) {
    Log::fatal('PROTOWIRE_WIRE_LENGTH_DELIMITED expected 2, got ' . PROTOWIRE_WIRE_LENGTH_DELIMITED);
}
if (PROTOWIRE_WIRE_START_GROUP !== 3) {
    Log::fatal('PROTOWIRE_WIRE_START_GROUP expected 3, got ' . PROTOWIRE_WIRE_START_GROUP);
}
if (PROTOWIRE_WIRE_END_GROUP !== 4) {
    Log::fatal('PROTOWIRE_WIRE_END_GROUP expected 4, got ' . PROTOWIRE_WIRE_END_GROUP);
}
if (PROTOWIRE_WIRE_FIXED32 !== 5) {
    Log::fatal('PROTOWIRE_WIRE_FIXED32 expected 5, got ' . PROTOWIRE_WIRE_FIXED32);
}
Log::info('常量定义 测试通过');

// ---------------------------------------------------------------------------
// 2. encode 函数
// ---------------------------------------------------------------------------

// varint
$v = protowire_encode_varint(42);
if ($v === '' || strlen($v) < 1) {
    Log::fatal('protowire_encode_varint(42) 返回空字符串');
}
// varint 0
$v0 = protowire_encode_varint(0);
if ($v0 === '' || strlen($v0) < 1) {
    Log::fatal('protowire_encode_varint(0) 返回空字符串');
}
// varint 150
$v150 = protowire_encode_varint(150);
if ($v150 === '') {
    Log::fatal('protowire_encode_varint(150) 返回空字符串');
}

// tag
$t = protowire_encode_tag(1, PROTOWIRE_WIRE_VARINT);
if ($t === '') {
    Log::fatal('protowire_encode_tag(1, 0) 返回空字符串');
}

// bytes
$b = protowire_encode_bytes('hello');
if ($b === '') {
    Log::fatal('protowire_encode_bytes("hello") 返回空字符串');
}

// fixed32
$f32 = protowire_encode_fixed32(12345);
if ($f32 === '' || strlen($f32) !== 4) {
    Log::fatal('protowire_encode_fixed32(12345) expected 4 bytes, got ' . strlen($f32));
}

// fixed64
$f64 = protowire_encode_fixed64(99999);
if ($f64 === '' || strlen($f64) !== 8) {
    Log::fatal('protowire_encode_fixed64(99999) expected 8 bytes, got ' . strlen($f64));
}

Log::info('encode 函数 测试通过');

// ---------------------------------------------------------------------------
// 3. 基本解析：varint 往返
// ---------------------------------------------------------------------------

$data = protowire_encode_tag(1, PROTOWIRE_WIRE_VARINT) . protowire_encode_varint(42);
$fields = protowire_parse_raw_fields($data);

if (count($fields) !== 1) {
    Log::fatal('基本 varint: expected 1 field, got ' . count($fields));
}
$f = $fields[0];
if ($f['number'] !== 1) {
    Log::fatal('基本 varint: expected number 1, got ' . var_export($f['number'], true));
}
if ($f['wire_type'] !== PROTOWIRE_WIRE_VARINT) {
    Log::fatal('基本 varint: expected wire_type 0, got ' . var_export($f['wire_type'], true));
}
if ($f['value'] !== 42) {
    Log::fatal('基本 varint: expected value 42, got ' . var_export($f['value'], true));
}

Log::info('基本 varint 解析 测试通过');

// ---------------------------------------------------------------------------
// 4. 空数据解析
// ---------------------------------------------------------------------------

$empty = protowire_parse_raw_fields('');
if (count($empty) !== 0) {
    Log::fatal('空数据: expected 0 fields, got ' . count($empty));
}

Log::info('空数据解析 测试通过');

// ---------------------------------------------------------------------------
// 5. 多字段解析
// ---------------------------------------------------------------------------

$data = '';
$data .= protowire_encode_tag(1, PROTOWIRE_WIRE_VARINT) . protowire_encode_varint(100);
$data .= protowire_encode_tag(2, PROTOWIRE_WIRE_FIXED32) . protowire_encode_fixed32(999);
$data .= protowire_encode_tag(3, PROTOWIRE_WIRE_FIXED64) . protowire_encode_fixed64(888);

$fields = protowire_parse_raw_fields($data);
if (count($fields) !== 3) {
    Log::fatal('多字段: expected 3 fields, got ' . count($fields));
}
if ($fields[0]['value'] !== 100 || $fields[0]['number'] !== 1) {
    Log::fatal('多字段: field 1 mismatch');
}
if ($fields[1]['value'] !== 999 || $fields[1]['number'] !== 2) {
    Log::fatal('多字段: field 2 mismatch');
}
if ($fields[2]['value'] !== 888 || $fields[2]['number'] !== 3) {
    Log::fatal('多字段: field 3 mismatch');
}

Log::info('多字段解析 测试通过');

// ---------------------------------------------------------------------------
// 6. length-delimited 字节
// ---------------------------------------------------------------------------

$payload = 'hello protobuf';
$data = protowire_encode_tag(4, PROTOWIRE_WIRE_LENGTH_DELIMITED)
      . protowire_encode_bytes($payload);

$fields = protowire_parse_raw_fields($data);
$f = $fields[0];
if ($f['value'] !== $payload) {
    Log::fatal('length-delimited: expected "' . $payload . '", got ' . var_export($f['value'], true));
}

Log::info('length-delimited 字节 测试通过');

// ---------------------------------------------------------------------------
// 7. 嵌套消息
// ---------------------------------------------------------------------------

// 内层: field 1 (varint) = 7
$inner = protowire_encode_tag(1, PROTOWIRE_WIRE_VARINT) . protowire_encode_varint(7);
// 外层: field 5 (length-delimited) 包含 inner
$data = protowire_encode_tag(5, PROTOWIRE_WIRE_LENGTH_DELIMITED)
      . protowire_encode_bytes($inner);

// 不配置 message_fields → 应返回原始字节
$fieldsRaw = protowire_parse_raw_fields($data);
if (is_string($fieldsRaw[0]['value']) === false) {
    Log::fatal('嵌套消息(无配置): expected string bytes');
}
Log::info('嵌套消息(默认字节) 测试通过');

// 配置 message_fields → 应递归解析
$fields = protowire_parse_raw_fields($data, [
    'message_fields' => [5 => true],
]);

$outer = $fields[0];
if (is_array($outer['value']) === false) {
    Log::fatal('嵌套消息(递归): value expected array, got ' . gettype($outer['value']));
}
$innerFields = $outer['value'];
if (count($innerFields) !== 1) {
    Log::fatal('嵌套消息(递归): expected 1 inner field, got ' . count($innerFields));
}
if ($innerFields[0]['value'] !== 7) {
    Log::fatal('嵌套消息(递归): inner value expected 7, got ' . var_export($innerFields[0]['value'], true));
}

Log::info('嵌套消息 递归解析 测试通过');

// ---------------------------------------------------------------------------
// 8. 多层嵌套 + 深度限制
// ---------------------------------------------------------------------------

// 构建 6 层嵌套，max_depth=5 → 应失败
$maxDepth = 5;
$nestedDepth = 6;

// 最内层: field 1 (varint) = 99
$payload = protowire_encode_tag(1, PROTOWIRE_WIRE_VARINT) . protowire_encode_varint(99);

for ($i = 0; $i < $nestedDepth; $i++) {
    $payload = protowire_encode_tag(1, PROTOWIRE_WIRE_LENGTH_DELIMITED)
             . protowire_encode_bytes($payload);
}

// 深度 5 时, 6 层嵌套应触达限制
$depthErrorOccurred = false;
try {
    $r = protowire_parse_raw_fields($payload, [
        'message_fields' => [1 => true],
        'max_depth' => $maxDepth,
    ]);
} catch (\Exception $e) {
    $depthErrorOccurred = true;
}
if ($depthErrorOccurred === false) {
    Log::fatal('深度限制: expected error for recursion depth > 5');
}

Log::info('深度限制 测试通过');

// ---------------------------------------------------------------------------
// 9. Group 分组
// ---------------------------------------------------------------------------

// Group (field 10):
//   field 1 (varint) = 123
//   field 2 (fixed32) = 456
// EndGroup
$data = '';
$data .= protowire_encode_tag(10, PROTOWIRE_WIRE_START_GROUP);
$data .= protowire_encode_tag(1, PROTOWIRE_WIRE_VARINT) . protowire_encode_varint(123);
$data .= protowire_encode_tag(2, PROTOWIRE_WIRE_FIXED32) . protowire_encode_fixed32(456);
$data .= protowire_encode_tag(10, PROTOWIRE_WIRE_END_GROUP);

$fields = protowire_parse_raw_fields($data);
if (count($fields) !== 1) {
    Log::fatal('Group: expected 1 field, got ' . count($fields));
}

$groupFields = $fields[0]['value'];
if (is_array($groupFields) === false) {
    Log::fatal('Group: value expected array, got ' . gettype($groupFields));
}
if (count($groupFields) !== 2) {
    Log::fatal('Group: expected 2 inner fields, got ' . count($groupFields));
}
if ($groupFields[0]['value'] !== 123) {
    Log::fatal('Group: inner field 1 expected 123, got ' . var_export($groupFields[0]['value'], true));
}
if ($groupFields[1]['value'] !== 456) {
    Log::fatal('Group: inner field 2 expected 456, got ' . var_export($groupFields[1]['value'], true));
}

Log::info('Group 分组 测试通过');

// ---------------------------------------------------------------------------
// 10. 嵌套 Group
// ---------------------------------------------------------------------------

$data = '';
$data .= protowire_encode_tag(20, PROTOWIRE_WIRE_START_GROUP);
$data .= protowire_encode_tag(21, PROTOWIRE_WIRE_START_GROUP);
$data .= protowire_encode_tag(1, PROTOWIRE_WIRE_VARINT) . protowire_encode_varint(77);
$data .= protowire_encode_tag(21, PROTOWIRE_WIRE_END_GROUP);
$data .= protowire_encode_tag(20, PROTOWIRE_WIRE_END_GROUP);

$fields = protowire_parse_raw_fields($data);
$outerGroup = $fields[0]['value'];
$innerGroup = $outerGroup[0]['value'];
if ($innerGroup[0]['value'] !== 77) {
    Log::fatal('嵌套 Group: inner value expected 77, got ' . var_export($innerGroup[0]['value'], true));
}

Log::info('嵌套 Group 测试通过');

// ---------------------------------------------------------------------------
// 11. Group 不匹配的 End 标签 (应报错)
// ---------------------------------------------------------------------------

$data = '';
$data .= protowire_encode_tag(10, PROTOWIRE_WIRE_START_GROUP);
$data .= protowire_encode_tag(1, PROTOWIRE_WIRE_VARINT) . protowire_encode_varint(1);
$data .= protowire_encode_tag(99, PROTOWIRE_WIRE_END_GROUP); // 错误: field 99 匹配 field 10

$groupMismatchError = false;
try {
    protowire_parse_raw_fields($data);
} catch (\Throwable $e) {
    $groupMismatchError = true;
}
if ($groupMismatchError === false) {
    Log::fatal('Group 不匹配: expected error for mismatched end group');
}

Log::info('Group 不匹配 End 标签 测试通过');

// ---------------------------------------------------------------------------
// 12. 截断的 Group (无 End 标签) — 应报错
// ---------------------------------------------------------------------------

$data = '';
$data .= protowire_encode_tag(10, PROTOWIRE_WIRE_START_GROUP);
$data .= protowire_encode_tag(1, PROTOWIRE_WIRE_VARINT) . protowire_encode_varint(1);
// 缺少 EndGroup

$groupTruncatedError = false;
try {
    protowire_parse_raw_fields($data);
} catch (\Throwable $e) {
    $groupTruncatedError = true;
}
if ($groupTruncatedError === false) {
    Log::fatal('截断 Group: expected error for missing end group');
}

Log::info('截断 Group 测试通过');

// ---------------------------------------------------------------------------
// 13. Packed Varint
// ---------------------------------------------------------------------------

// 手工拼 packed: field 7, wire type = length-delimited
$packedPayload = '';
$packedPayload .= protowire_encode_varint(1);
$packedPayload .= protowire_encode_varint(2);
$packedPayload .= protowire_encode_varint(3);
$packedPayload .= protowire_encode_varint(4);
$packedPayload .= protowire_encode_varint(5);

$data = protowire_encode_tag(7, PROTOWIRE_WIRE_LENGTH_DELIMITED)
      . protowire_encode_bytes($packedPayload);

$fields = protowire_parse_raw_fields($data, [
    'packed_fields' => [7 => true],
    'packed_element_type' => [7 => PROTOWIRE_WIRE_VARINT],
]);

$vals = $fields[0]['value'];
if (is_array($vals) === false) {
    Log::fatal('Packed varint: expected array, got ' . gettype($vals));
}
if (count($vals) !== 5) {
    Log::fatal('Packed varint: expected 5 elements, got ' . count($vals));
}
if ($vals !== [1, 2, 3, 4, 5]) {
    Log::fatal('Packed varint: expected [1,2,3,4,5], got ' . var_export($vals, true));
}

Log::info('Packed Varint 测试通过');

// ---------------------------------------------------------------------------
// 14. Packed Fixed32
// ---------------------------------------------------------------------------

$packedPayload = '';
$packedPayload .= protowire_encode_fixed32(10);
$packedPayload .= protowire_encode_fixed32(20);
$packedPayload .= protowire_encode_fixed32(30);

$data = protowire_encode_tag(8, PROTOWIRE_WIRE_LENGTH_DELIMITED)
      . protowire_encode_bytes($packedPayload);

$fields = protowire_parse_raw_fields($data, [
    'packed_fields' => [8 => true],
    'packed_element_type' => [8 => PROTOWIRE_WIRE_FIXED32],
]);

$vals = $fields[0]['value'];
if ($vals !== [10, 20, 30]) {
    Log::fatal('Packed fixed32: expected [10,20,30], got ' . var_export($vals, true));
}

Log::info('Packed Fixed32 测试通过');

// ---------------------------------------------------------------------------
// 15. Packed Fixed64
// ---------------------------------------------------------------------------

$packedPayload = '';
$packedPayload .= protowire_encode_fixed64(100);
$packedPayload .= protowire_encode_fixed64(200);

$data = protowire_encode_tag(9, PROTOWIRE_WIRE_LENGTH_DELIMITED)
      . protowire_encode_bytes($packedPayload);

$fields = protowire_parse_raw_fields($data, [
    'packed_fields' => [9 => true],
    'packed_element_type' => [9 => PROTOWIRE_WIRE_FIXED64],
]);

$vals = $fields[0]['value'];
if ($vals !== [100, 200]) {
    Log::fatal('Packed fixed64: expected [100,200], got ' . var_export($vals, true));
}

Log::info('Packed Fixed64 测试通过');

// ---------------------------------------------------------------------------
// 16. 空 Packed (零长度 payload)
// ---------------------------------------------------------------------------

$data = protowire_encode_tag(5, PROTOWIRE_WIRE_LENGTH_DELIMITED)
      . protowire_encode_bytes('');

$fields = protowire_parse_raw_fields($data, [
    'packed_fields' => [5 => true],
    'packed_element_type' => [5 => PROTOWIRE_WIRE_VARINT],
]);

$vals = $fields[0]['value'];
if (is_array($vals) === false || count($vals) !== 0) {
    Log::fatal('空 Packed: expected empty array, got ' . var_export($vals, true));
}

Log::info('空 Packed 测试通过');

// ---------------------------------------------------------------------------
// 17. 非 packed repeated (同一字段号出现多次)
// ---------------------------------------------------------------------------

$data = '';
for ($i = 0; $i < 3; $i++) {
    $data .= protowire_encode_tag(1, PROTOWIRE_WIRE_VARINT) . protowire_encode_varint(10 + $i);
}

$fields = protowire_parse_raw_fields($data);
if (count($fields) !== 3) {
    Log::fatal('非packed repeated: expected 3 fields, got ' . count($fields));
}
if ($fields[0]['value'] !== 10 || $fields[1]['value'] !== 11 || $fields[2]['value'] !== 12) {
    Log::fatal('非packed repeated: value mismatch');
}

Log::info('非 packed repeated 测试通过');

// ---------------------------------------------------------------------------
// 18. 格式错误数据：截断的 tag (应报错)
// ---------------------------------------------------------------------------

$malformedTagError = false;
try {
    protowire_parse_raw_fields("\xff\xff\xff\xff\xff\xff\xff\xff\xff\x01");
} catch (\Throwable $e) {
    $malformedTagError = true;
}
if ($malformedTagError === false) {
    Log::fatal('格式错误 tag: expected error');
}

Log::info('格式错误 tag 测试通过');

// ---------------------------------------------------------------------------
// 19. 截断的 length-delimited (应报错)
// ---------------------------------------------------------------------------

$truncatedError = false;
try {
    $data = protowire_encode_tag(3, PROTOWIRE_WIRE_LENGTH_DELIMITED)
          . protowire_encode_varint(100); // 声明 100 字节，但后面没有数据
    protowire_parse_raw_fields($data);
} catch (\Throwable $e) {
    $truncatedError = true;
}
if ($truncatedError === false) {
    Log::fatal('截断 length-delimited: expected error');
}

Log::info('截断 length-delimited 测试通过');

// ---------------------------------------------------------------------------
// 20. Packed 缺少 PackedElementType 配置 (应报错)
// ---------------------------------------------------------------------------

$packedPayload = protowire_encode_varint(1);
$data = protowire_encode_tag(7, PROTOWIRE_WIRE_LENGTH_DELIMITED)
      . protowire_encode_bytes($packedPayload);

$missingConfigError = false;
try {
    protowire_parse_raw_fields($data, [
        'packed_fields' => [7 => true],
        // 缺少 packed_element_type
    ]);
} catch (\Throwable $e) {
    $missingConfigError = true;
}
if ($missingConfigError === false) {
    Log::fatal('缺少 PackedElementType: expected error');
}

Log::info('缺少 PackedElementType 配置 测试通过');

// ---------------------------------------------------------------------------
// 21. Group 深度限制
// ---------------------------------------------------------------------------

$data = '';
$data .= protowire_encode_tag(10, PROTOWIRE_WIRE_START_GROUP);
$data .= protowire_encode_tag(11, PROTOWIRE_WIRE_START_GROUP);
$data .= protowire_encode_tag(1, PROTOWIRE_WIRE_VARINT) . protowire_encode_varint(1);
$data .= protowire_encode_tag(11, PROTOWIRE_WIRE_END_GROUP);
// 没有 EndGroup(10)，但深度限制应该在此之前触发

$groupDepthError = false;
try {
    protowire_parse_raw_fields($data, ['max_depth' => 2]);
} catch (\Throwable $e) {
    $groupDepthError = true;
}
if ($groupDepthError === false) {
    Log::fatal('Group 深度限制: expected error');
}

Log::info('Group 深度限制 测试通过');

// ---------------------------------------------------------------------------
// 22. 混合场景：消息内嵌 group + packed
// ---------------------------------------------------------------------------

// 内层消息: group(10): varint=9; field 3 packed: [1,2,3]
$inner = '';
$inner .= protowire_encode_tag(10, PROTOWIRE_WIRE_START_GROUP);
$inner .= protowire_encode_tag(1, PROTOWIRE_WIRE_VARINT) . protowire_encode_varint(9);
$inner .= protowire_encode_tag(10, PROTOWIRE_WIRE_END_GROUP);

$packedPayload = '';
$packedPayload .= protowire_encode_varint(1);
$packedPayload .= protowire_encode_varint(2);
$packedPayload .= protowire_encode_varint(3);
$inner .= protowire_encode_tag(3, PROTOWIRE_WIRE_LENGTH_DELIMITED)
        . protowire_encode_bytes($packedPayload);

// 外层: field 1 (message) 包含 inner
$data = protowire_encode_tag(1, PROTOWIRE_WIRE_LENGTH_DELIMITED)
      . protowire_encode_bytes($inner);

$fields = protowire_parse_raw_fields($data, [
    'message_fields' => [1 => true],
    'packed_fields' => [3 => true],
    'packed_element_type' => [3 => PROTOWIRE_WIRE_VARINT],
]);

$innerFields = $fields[0]['value'];
if (count($innerFields) !== 2) {
    Log::fatal('混合场景: expected 2 inner fields, got ' . count($innerFields));
}

// 检查 group
$groupFields = $innerFields[0]['value'];
if ($groupFields[0]['value'] !== 9) {
    Log::fatal('混合场景: group inner expected 9, got ' . var_export($groupFields[0]['value'], true));
}

// 检查 packed
$packedVals = $innerFields[1]['value'];
if ($packedVals !== [1, 2, 3]) {
    Log::fatal('混合场景: packed expected [1,2,3], got ' . var_export($packedVals, true));
}

Log::info('混合场景 (message + group + packed) 测试通过');

// ---------------------------------------------------------------------------
// 23. 截断的 fixed32 (应报错)
// ---------------------------------------------------------------------------

$truncFixed32Error = false;
try {
    $data = protowire_encode_tag(4, PROTOWIRE_WIRE_FIXED32) . "\x01\x02\x03"; // 只 3 字节
    protowire_parse_raw_fields($data);
} catch (\Throwable $e) {
    $truncFixed32Error = true;
}
if ($truncFixed32Error === false) {
    Log::fatal('截断 fixed32: expected error');
}

Log::info('截断 fixed32 测试通过');

// ---------------------------------------------------------------------------
// 24. 截断的 fixed64 (应报错)
// ---------------------------------------------------------------------------

$truncFixed64Error = false;
try {
    $data = protowire_encode_tag(5, PROTOWIRE_WIRE_FIXED64) . "\x01\x02\x03\x04"; // 只 4 字节
    protowire_parse_raw_fields($data);
} catch (\Throwable $e) {
    $truncFixed64Error = true;
}
if ($truncFixed64Error === false) {
    Log::fatal('截断 fixed64: expected error');
}

Log::info('截断 fixed64 测试通过');

// ---------------------------------------------------------------------------
// 完成
// ---------------------------------------------------------------------------

Log::info('protowire 标准库 PHP 验收测试全部通过');
