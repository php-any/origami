<?php

namespace App;

/**
 * 在这里开始，可以重新实现一个全新的 laravel 框架
 *  数据库需要常驻内存的建议放到 index.php 文件内，就能共享给 main 方法了
 *  main.php 内的调用类型存储信息会在请求结束后直接丢失
 */
function main($request, $response) {
    $path = $request->path();

    // 设置默认的 Content-Type
    $response->header("Content-Type", "text/html; charset=utf-8");

    // 静态资源处理：/assets/* 下的 css/js 等直接读取文件返回
    if (str_starts_with($path, "/assets/")) {
        $staticPath = "./pages" . $path;
        try {
            $content = file_get_contents($staticPath);
            $lowerPath = strtolower($path);
            $parts = explode(".", $lowerPath);
            $ext = "";
            if (count($parts) > 1) {
                $ext = $parts[count($parts) - 1];
            }

            if ($ext == "css") {
                $response->header("Content-Type", "text/css; charset=utf-8");
            } elseif ($ext == "js") {
                $response->header("Content-Type", "application/javascript; charset=utf-8");
            } elseif ($ext == "png") {
                $response->header("Content-Type", "image/png");
            } elseif ($ext == "jpg" || $ext == "jpeg") {
                $response->header("Content-Type", "image/jpeg");
            } elseif ($ext == "svg") {
                $response->header("Content-Type", "image/svg+xml");
            } elseif ($ext == "ico") {
                $response->header("Content-Type", "image/x-icon");
            }
            $response->write($content);
            return;
        } catch (\Exception $e) {
            Log::error("静态资源读取失败: " . $staticPath);
        }
    }

    // 处理根路径，默认加载 index
    if ($path == "/" || $path == "") {
        $path = "/index";
    }

    // 移除开头的斜杠
    $filePath = $path;
    if (strlen($filePath) > 0 && $filePath[0] == "/") {
        $filePath = substr($filePath, 1);
    }

    // 尝试加载 .html 文件
    $htmlPath = "./pages/" . $filePath . ".html";
    try {
        $html = file_get_contents($htmlPath);
        $response->header("Content-Type", "text/html; charset=utf-8");
        $response->write($html);
        return;
    } catch (\Exception $e) {
        Log::error("文件不存在或加载失败: " . $htmlPath);
    }

    // 如果文件不存在，返回 404
    $response->writeHeader(404);
    $response->write("404 Not Found: " . $path);
}
