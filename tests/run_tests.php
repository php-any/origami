<?php
namespace tests;

$path = __DIR__;

Log::info("path = ", $path);

$total = 0;
$failed = 0;

for (_, $file in scandir($path)) {
    $subDir = OS::path($path, $file);
    if(is_dir($subDir)) {
        for (_, $file in scandir($subDir)) {
            $filePath = OS::path($subDir, $file);
            if(!is_dir($filePath)) {
                if($filePath->indexOf(".php") != $filePath->length - 4) {
                    continue;
                }
                $total = $total + 1;
                $before = Log::fatalCount();
                $ok = true;
                try {
                    Log::info("执行 {$filePath}");
                    include($filePath);
                } catch (Exception $e) {
                    Log::error("执行文件发生错误, file={$filePath}; error={$e->getMessage()}");
                    $ok = false;
                } catch (Error $e) {
                    Log::error("执行文件发生Error, file={$filePath}; error=" . $e->getMessage());
                    $ok = false;
                }
                // Log::fatal 现在抛的是可捕获的异常，用例自己的 catch (\Throwable)
                // 可能把它吞掉，所以再用计数兜底，保证失败一定会被统计。
                if(Log::fatalCount() > $before) {
                    $ok = false;
                }
                if(!$ok) {
                    $failed = $failed + 1;
                    Log::error("用例失败: {$filePath}");
                }
            }
        }
    }
}

if($failed > 0) {
    Log::error("共 {$failed}/{$total} 个用例文件失败");
    exit(1);
}

Log::info("🎉 接口功能测试完成, 共 {$total} 个用例文件");
