<?php
namespace tests\php;
$html = "\n  <div wire:id=\"x\">hi</div>";
preg_match('/(?:\n\s*|^\s*)<([a-zA-Z0-9\-]+)/', $html, $matches, PREG_OFFSET_CAPTURE);
if (!count($matches)) {
  Log::fatal('no match');
}
if (!isset($matches[1][0]) || $matches[1][0] !== 'div') {
  Log::fatal('tag wrong: '.json_encode($matches));
}
if (!isset($matches[1][1]) || !is_int($matches[1][1])) {
  Log::fatal('offset wrong: '.json_encode($matches));
}
Log::info('PREG_OFFSET_CAPTURE ok tag='.$matches[1][0].' off='.$matches[1][1]);
