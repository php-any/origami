<?php

function checkJsonFlags(bool $condition, string $message): void
{
    if (! $condition) throw new RuntimeException($message);
}

$decoded = json_decode('{"id":9223372036854775808,"name":"商品"}', true, flags: JSON_BIGINT_AS_STRING);
checkJsonFlags($decoded['id'] === '9223372036854775808', 'Named flags lost a large integer');
checkJsonFlags($decoded['name'] === '商品', 'Named flags lost a string');
checkJsonFlags(json_decode('paid', true, flags: JSON_BIGINT_AS_STRING) === null, 'Livewire plain query values should return null before fallback');
checkJsonFlags(json_decode('123', flags: JSON_BIGINT_AS_STRING) === 123, 'Normal integers should remain integers');
checkJsonFlags(json_decode('1.5', flags: JSON_BIGINT_AS_STRING) === 1.5, 'Floats should remain floats');
$object = json_decode('{"child":{"id":9223372036854775808}}', flags: JSON_BIGINT_AS_STRING);
checkJsonFlags($object instanceof stdClass && $object->child instanceof stdClass && gettype($object) === 'object', 'JSON objects must be stdClass instances');
checkJsonFlags($object->child->id === '9223372036854775808', 'Nested object lost bigint flags');
checkJsonFlags(json_encode(json_decode('{}')) === '{}', 'An empty JSON object must remain an object');
checkJsonFlags(json_decode('{"id":1}', associative: null, flags: JSON_OBJECT_AS_ARRAY)['id'] === 1, 'Object-as-array flag was ignored');
checkJsonFlags(json_decode('[[]]', true, depth: 2) === null && json_last_error() === JSON_ERROR_DEPTH, 'Depth limit was ignored');
checkJsonFlags(json_decode('[[]]', true, depth: 3) === [[]], 'Valid nesting was rejected');
checkJsonFlags(json_decode('1 2', flags: JSON_BIGINT_AS_STRING) === null, 'Trailing JSON was accepted');
try {
    json_decode('[]', depth: 0);
    throw new RuntimeException('Zero depth should throw');
} catch (ValueError $error) {
}
echo "json named flags: PASS\n";
