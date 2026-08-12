<?php

namespace tests\php;

/**
 * 验证 exit(int) 退出码；用子进程跑 fixture，避免 exit 终止整个测试套件。
 * 注意：不能用 `go run` 断言非零码（go run 会把子进程非零码统一成 1）。
 */
$root = dirname(__DIR__, 2);
$bin = $root . '/tests/php/exit_status_fixtures/origami_exit_status_test_bin';
$go = getenv('GOTOOL') ?: 'go';

$descriptorspec = [1 => ['pipe', 'w'], 2 => ['pipe', 'w']];
$pipes = [];
$build = proc_open([$go, 'build', '-o', $bin, './zy.go'], $descriptorspec, $pipes, $root);
if ($build === false) {
	Log::fatal('无法 go build origami 用于 exit 测试');
}
$buildOut = stream_get_contents($pipes[1]) . stream_get_contents($pipes[2]);
fclose($pipes[1]);
fclose($pipes[2]);
$buildCode = proc_close($build);
if ($buildCode !== 0) {
	Log::fatal('go build 失败: ' . $buildOut);
}

$cases = [
	[__DIR__ . '/exit_status_fixtures/exit0.php', 0],
	[__DIR__ . '/exit_status_fixtures/exit7.php', 7],
];

foreach ($cases as [$script, $want]) {
	$pipes = [];
	$process = proc_open([$bin, $script], $descriptorspec, $pipes, $root);
	if ($process === false) {
		Log::fatal('无法启动 exit fixture 子进程: ' . $script);
	}
	$stdout = stream_get_contents($pipes[1]);
	$stderr = stream_get_contents($pipes[2]);
	fclose($pipes[1]);
	fclose($pipes[2]);
	$code = proc_close($process);
	if ($code !== $want) {
		Log::fatal("exit 退出码不符 want={$want} got={$code} stdout={$stdout} stderr={$stderr}");
	}
	if (is_string($stderr) && str_contains($stderr, '解析错误')) {
		Log::fatal('exit 不应被展示为解析错误: ' . $stderr);
	}
}

@unlink($bin);
Log::info('exit 退出码测试通过');
