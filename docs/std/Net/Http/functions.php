<?php

namespace Net\Http;


/**
 * functions - 标准库函数
 * 
 * 此文件包含 functions 模块的伪代码接口定义
 * 这些是自动生成的接口，仅用于参考，不包含具体实现
 */


/**
 * app 函数（开发模式：默认 hotReload=true，每请求 TempVM + 重新加载应用）
 */
function app($request, $response, string $filePath, string $fun, bool $hotReload) {
    // 实现逻辑
}

/**
 * appFlash 函数（生产模式：进程内只加载一次应用目录下的 main.php，按注解路由直接分发，不调用 main）
 */
function appFlash($request, $response, string $dir) {
    // 实现逻辑
}


