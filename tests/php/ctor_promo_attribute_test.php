<?php

namespace tests\php;

/**
 * 构造器属性提升参数上的 #[Attr] 可出现在 public 之前（Laravel Auth Attempting 事件）。
 */
class CtorPromoAttr_Event
{
    public function __construct(
        public $guard,
        #[\SensitiveParameter] public $credentials,
        public $remember,
    ) {
    }
}

$e = new CtorPromoAttr_Event('admin', ['email' => 'a@b.c'], false);
if ($e->guard !== 'admin') {
    Log::fatal('提升属性 guard 失败: ' . var_export($e->guard, true));
}
if (!is_array($e->credentials) || ($e->credentials['email'] ?? '') !== 'a@b.c') {
    Log::fatal('提升属性 credentials 失败: ' . var_export($e->credentials, true));
}
if ($e->remember !== false) {
    Log::fatal('提升属性 remember 失败: ' . var_export($e->remember, true));
}

Log::info('构造器提升参数注解测试通过');
