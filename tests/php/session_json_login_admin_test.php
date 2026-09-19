<?php

namespace tests\php;

/**
 * Laravel 13 session serialization=json。
 * 登录后 payload 含 login_admin_<sha1> 键，json_decode(..., true) 必须能读出用户 id。
 */

$raw = '{"_token":"vqaZ2P9mDZ8sTOGAIHgsnXDuCPbtPdE3NlCUcx8E","_previous":{"url":"http://127.0.0.1:18086/admin/login","route":"filament.admin.auth.login"},"_flash":{"old":[],"new":[]},"login_admin_59ba36addc2b2f9401580f014c7f58ea4e30989d":1}';
$data = json_decode($raw, true);
$key = 'login_admin_59ba36addc2b2f9401580f014c7f58ea4e30989d';
if (!is_array($data)) {
    Log::fatal('session json_decode 应为数组, got '.gettype($data));
}
if (!array_key_exists($key, $data)) {
    Log::fatal('session json 丢失 login_admin 键: '.json_encode(array_keys($data)));
}
if ((int) $data[$key] !== 1) {
    Log::fatal('login_admin 用户 id 错误: '.var_export($data[$key], true));
}

Log::info('session_json_login_admin 测试通过');
