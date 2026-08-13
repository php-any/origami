<?php

namespace tests\php;

/**
 * 验证 html_entity_decode 基本解码。
 */
$got = html_entity_decode('&lt;a&gt;&amp;', ENT_QUOTES, 'UTF-8');
if ($got !== '<a>&') {
	Log::fatal("html_entity_decode 结果异常: $got");
}
Log::info('html_entity_decode 测试通过');
