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

/**
 * appScan 函数（启动时扫描：扫描目录下的注解路由，返回路由表，配合 Server::flash 使用）
 *
 * 在服务启动阶段调用，扫描指定目录下所有 #[Controller] 类的注解路由，
 * 返回路由描述数组。每个元素包含:
 *   - method: HTTP 方法（GET/POST/PUT/DELETE 等）
 *   - path: 完整路径（prefix + mapping path）
 *   - controller: 控制器类名
 *   - action: 方法名
 *   - middleware: 中间件/拦截器列表
 *
 * @param string $dir 应用目录路径
 * @return array 路由描述数组
 */
function appScan(string $dir): array {
    // 实现逻辑
}


