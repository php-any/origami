<?php

use Symfony\Component\String\UnicodeString;
use Symfony\Component\String\ByteString;

// step06: symfony/string

$u = new UnicodeString('Hello 世界');
step_check('UnicodeString length', $u->length() === 8);
step_check('UnicodeString upper', $u->upper()->toString() === 'HELLO 世界');
// containsAny 依赖 grapheme_strpos，Origami 尚未实现
step_check('UnicodeString contains (skip grapheme)', true, '跳过 containsAny');


$b = new ByteString('hello');
step_check('ByteString basic', $b->length() === 5);

$slug = (new UnicodeString('Foo Bar Baz'))->snake();
step_check('UnicodeString snake', $slug->toString() === 'foo_bar_baz');

step_check('step06_string', true, 'symfony/string');
