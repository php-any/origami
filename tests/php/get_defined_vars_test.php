<?php

function get_defined_vars_probe($input)
{
    $nullable = null;
    extract(['dynamic' => 'loaded']);

    $vars = get_defined_vars();

    if ($vars['input'] !== 'value') {
        Log::fatal('get_defined_vars 未返回参数变量');
    }
    if (!array_key_exists('nullable', $vars) || $vars['nullable'] !== null) {
        Log::fatal('get_defined_vars 未保留值为 null 的变量');
    }
    if ($vars['dynamic'] !== 'loaded') {
        Log::fatal('get_defined_vars 未返回 extract 动态注入变量');
    }

    return array_diff_key($vars, ['input' => 1, 'nullable' => 1, 'dynamic' => 1, 'vars' => 1]);
}

$extra = get_defined_vars_probe('value');
if (count($extra) !== 0) {
    Log::fatal('get_defined_vars 返回了当前作用域之外的变量');
}

Log::info('get_defined_vars 测试通过');
