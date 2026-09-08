<?php

namespace tests\php;

/**
 * json_encode 对象只输出公共属性，不把 protected 引用图编进 JSON。
 * Livewire effects.returns 里的 Redirector 依赖此语义。
 */
class JsonEncodePublicOnly_Box
{
	public string $visible = 'yes';
	protected string $hidden = 'secret';
	private string $private = 'no';
}

$json = json_encode(new JsonEncodePublicOnly_Box());
if ($json === false) {
	Log::fatal('json_encode 失败: ' . json_last_error_msg());
}
if (!str_contains($json, 'visible') || !str_contains($json, 'yes')) {
	Log::fatal('公共属性应出现在 JSON: ' . $json);
}
if (str_contains($json, 'hidden') || str_contains($json, 'secret') || str_contains($json, 'private')) {
	Log::fatal('非公共属性不应出现在 JSON: ' . $json);
}
Log::info('json_encode 仅公共属性测试通过');
