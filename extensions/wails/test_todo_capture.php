<?php
$todos  = [];
$nextId = 1;

$id = $nextId++;
$todos[$id] = ['id' => $id, 'text' => 'hello', 'done' => false];
$id = $nextId++;
$todos[$id] = ['id' => $id, 'text' => 'world', 'done' => true];

__wails_capture(array_values($todos));
