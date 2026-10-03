<?php
function option_check($ok, $message) { if (!$ok) throw new Exception($message); }
$hash = password_hash('options', PASSWORD_BCRYPT, ['cost' => 4]);
option_check(password_verify('options', $hash), 'bcrypt verify');
option_check(substr($hash, 0, 7) === '$2y$04$', 'bcrypt options ignored');
option_check(password_get_info($hash)['options']['cost'] === 4, 'bcrypt info');
option_check(password_get_info($hash)['algo'] === '2y', 'bcrypt algorithm identity');
option_check(password_get_info('invalid')['algo'] === null, 'unknown algorithm identity');
option_check(!password_needs_rehash($hash, PASSWORD_BCRYPT, ['cost' => 4]), 'unchanged cost');
option_check(password_needs_rehash($hash, PASSWORD_BCRYPT, ['cost' => 5]), 'changed cost');
option_check(get_debug_type(function () {}) === 'Closure', 'closure type');
echo "array options contract OK\n";
