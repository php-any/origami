<?php

namespace tests\php;

/**
 * 验证 parse_str 写入结果数组。
 */
parse_str('tab=topics&p=2', $params);
if (($params['tab'] ?? '') !== 'topics' || ($params['p'] ?? '') !== '2') {
	Log::fatal('parse_str 结果错误: ' . var_export($params, true));
}
Log::info('parse_str 测试通过');
