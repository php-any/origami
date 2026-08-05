<?php

namespace tests\php;

/**
 * 验证 strncasecmp 按长度、不区分大小写比较。
 */
$cases = [
	['SELECT *', 'select', 6, 0],
	['INSERT', 'select', 6, -1],
	['abc', 'ABCD', 3, 0],
	['abc', 'ab', 2, 0],
];
foreach ($cases as [$a, $b, $n, $want]) {
	$got = strncasecmp($a, $b, $n);
	$sign = $got <=> 0;
	$wantSign = $want <=> 0;
	if ($sign !== $wantSign) {
		Log::fatal("strncasecmp('$a','$b',$n) = $got, want sign $wantSign");
	}
}
Log::info('strncasecmp 测试通过');
