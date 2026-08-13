<?php

namespace tests\php;

/**
 * 验证 $GLOBALS / $_SESSION 可在同一请求内读写（bbs1org 用 $GLOBALS['__me_cache'] 缓存登录态）。
 * 跨 RequestVM 不串态由 std/php/fpm 的 Go 测试覆盖。
 */
$GLOBALS['__me_cache'] = ['id' => 42, 'name' => 'alice'];
$_SESSION['uid'] = 42;

if (!isset($GLOBALS['__me_cache']['id']) || (int)$GLOBALS['__me_cache']['id'] !== 42) {
	Log::fatal('$GLOBALS 读写失败');
}
if (!isset($_SESSION['uid']) || (int)$_SESSION['uid'] !== 42) {
	Log::fatal('$_SESSION 读写失败');
}

unset($GLOBALS['__me_cache']);
if (array_key_exists('__me_cache', $GLOBALS)) {
	Log::fatal('unset($GLOBALS[...]) 未生效');
}

Log::info('$GLOBALS/$_SESSION 请求内读写测试通过');
