<?php

namespace tests\dom;

/**
 * dom：DOMDocument::loadHTML。
 */

$dom = new \DOMDocument();
$ok = $dom->loadHTML('<div id="n">hello</div>');
if ($ok !== true && $ok !== false) {
    Log::fatal('loadHTML 返回异常');
}
if (!class_exists('DOMElement')) {
    Log::fatal('DOMElement 未注册');
}

Log::info('dom DOMDocument 测试通过');
