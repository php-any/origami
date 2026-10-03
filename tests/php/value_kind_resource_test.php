<?php
$values = [true, 1, 1.5, '', null, [], new stdClass(), function () {}];
$types = ['boolean', 'integer', 'double', 'string', 'NULL', 'array', 'object', 'object'];
foreach ($values as $index => $value) { if (gettype($value) !== $types[$index]) { throw new Exception('value kind mismatch'); } }
$resource = fopen(__FILE__, 'rb');
if (!is_resource($resource) || is_object($resource) || gettype($resource) !== 'resource') { throw new Exception('resource classified as object'); }
function object_only(object $value) {}
try { object_only($resource); throw new Exception('resource accepted as object'); } catch (TypeError $e) {}
fclose($resource);
if (is_resource($resource) || gettype($resource) !== 'resource (closed)' || get_resource_type($resource) !== 'Unknown') { throw new Exception('closed resource kind lost'); }
echo "value kind resource OK\n";
