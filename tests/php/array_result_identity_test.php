<?php
function result_check($ok, $message) {
    if (!$ok) { throw new Exception($message); }
}
function result_array(array $value): array { return $value; }
function check_result($value, $expected, $name) {
    result_check(is_array($value) && !is_object($value) && gettype($value) === 'array', $name . ' identity');
    result_check(result_array($value) === $expected, $name . ' keys/values');
    $copy = $value;
    $copy['changed'] = 1;
    result_check(!array_key_exists('changed', $value), $name . ' copy isolation');
}
$left = ['nested' => [1]];
$nothing = null;
$name = 'left';
check_result(compact('left', 'nothing'), ['left' => $left, 'nothing' => null], 'compact');
check_result(compact($name, ['nothing', ['left']]), ['left' => $left, 'nothing' => null], 'dynamic compact');
$vars = get_defined_vars();
result_check(is_array($vars) && !is_object($vars) && $vars['nothing'] === null, 'defined vars identity');
$vars['left']['nested'][0] = 2;
result_check($left['nested'][0] === 1, 'defined vars nested copy');
check_result(array_change_key_case(['UP' => 1, 8 => 2, '' => 3, '01' => 4]), ['up' => 1, 8 => 2, '' => 3, '01' => 4], 'key case');
check_result(array_count_values(['first', 8, 'first', '08', 8]), ['first' => 2, 8 => 2, '08' => 1], 'count values');
check_result(array_column([['id' => 8, 'v' => 'a'], ['v' => 'b'], ['id' => '01', 'v' => null], ['id' => 9]], 'v', 'id'), [8 => 'a', 9 => 'b', '01' => null], 'column');
check_result(array_unique([8 => 'a', 'label' => 'b', '' => 'a', 20 => 'c']), [8 => 'a', 'label' => 'b', 20 => 'c'], 'unique');
check_result(array_slice([8 => 'a', 'label' => 'b', '' => 'c'], 0, null, true), [8 => 'a', 'label' => 'b', '' => 'c'], 'slice');
check_result(unpack('C*', "\x01\x02"), [1 => 1, 2 => 2], 'unpack');
$source = ['word' => ['' => null, 8 => 3], '01' => 4, 12 => 5];
check_result(unserialize(serialize($source)), $source, 'serialization');
result_check(serialize($source) === 'a:3:{s:4:"word";a:2:{s:0:"";N;i:8;i:3;}s:2:"01";i:4;i:12;i:5;}', 'serialized key representation');
$bytes = ['quote' => "a\"b\x00c", 'utf8' => '中文'];
check_result(unserialize(serialize($bytes)), $bytes, 'binary serialization');
check_result(array_chunk(['label' => 1, 8 => 2, '' => 3], 2, true), [['label' => 1, 8 => 2], ['' => 3]], 'chunk');
check_result(array_diff(['keep' => 1, '' => 2, 8 => 3], [3]), ['keep' => 1, '' => 2], 'diff');
$alias = 7;
$referenced = ['key' => &$alias, 8 => 2];
$unique = array_unique($referenced);
$chunks = array_chunk($referenced, 2, true);
$unique['key'] = 9;
result_check($alias === 9 && $chunks[0]['key'] === 9, 'result reference aliases');
foreach (['array_unique', 'array_change_key_case', 'array_count_values'] as $function) {
    $rejected = false;
    try { $function(new stdClass()); } catch (TypeError $error) { $rejected = true; }
    result_check($rejected, $function . ' accepted object as array');
}
echo "array result identity OK\n";
