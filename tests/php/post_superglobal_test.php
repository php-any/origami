<?php

namespace tests\php;

/**
 * 验证 $_POST 在有 HTTP 请求体时应可读（通过模拟超全局写入路径的基本行为）。
 * 完整 HTTP ParseForm 路径由 examples/bbs1org 安装 POST 冒烟覆盖。
 */
$_POST['step'] = 'install';
if (($_POST['step'] ?? '') !== 'install') {
	Log::fatal('$_POST 赋值/读取失败');
}
Log::info('$_POST 基本读写测试通过');
