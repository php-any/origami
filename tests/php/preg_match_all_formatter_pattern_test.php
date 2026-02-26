<?php

namespace tests\php;

/**
 * 直接使用 OutputFormatter::formatAndWrap 中的正则：
 *   preg_match_all("#<(($openTagRegex) | /($closeTagRegex)?)>#ix", $message, $matches, PREG_OFFSET_CAPTURE);
 * 对比 Origami 下 preg_match_all 的 $matches 结构与 PHP CLI 是否一致。
 */

$openTagRegex = '[a-z](?:[^\\\\<>]*+ | \\\\.)*';
$closeTagRegex = '[a-z][^<>]*+';

$pattern = "#<(($openTagRegex) | /($closeTagRegex)?)>#ix";

$message = 'Test Console Application <info>1.0.0</info>';

$ok = preg_match_all($pattern, $message, $matches, PREG_OFFSET_CAPTURE);

Log::info('preg_match_all ok=' . ($ok === false ? 'false' : (string) $ok));
Log::info('preg_match_all matches=' . json_encode($matches));

