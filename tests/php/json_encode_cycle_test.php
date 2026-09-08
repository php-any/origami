<?php

namespace tests\php;

/**
 * json_encode：jsonSerialize() 返回 $this 时应报 JSON_ERROR_RECURSION，而不是死循环。
 */
class JsonEncodeCycle_Node implements \JsonSerializable
{
	public function jsonSerialize(): mixed
	{
		return $this;
	}
}

$encoded = json_encode(new JsonEncodeCycle_Node());
if ($encoded !== false) {
	Log::fatal('自引用 jsonSerialize 应返回 false, 实际: ' . var_export($encoded, true));
}
if (json_last_error() !== JSON_ERROR_RECURSION && json_last_error() !== JSON_ERROR_DEPTH) {
	Log::fatal('json_last_error 应为 RECURSION/DEPTH, 实际: ' . json_last_error());
}
Log::info('json_encode 循环引用测试通过');
