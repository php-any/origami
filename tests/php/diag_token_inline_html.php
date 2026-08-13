<?php
// 诊断：Origami token_get_all 的 T_INLINE_HTML 是否与 PHP 一致（Blade 依赖此 ID）
$s = "@auth\n";
$tokens = token_get_all($s);
$t0 = $tokens[0];
$id = is_array($t0) ? $t0[0] : -1;
echo "token[0]=$id T_INLINE_HTML=" . T_INLINE_HTML . " match=" . ($id == T_INLINE_HTML ? 'yes' : 'no') . "\n";
