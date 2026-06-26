<?php
namespace tests;

$path = __DIR__;

Log::info("path = ", $path);

for (_, $file in scandir($path)) {
    $subDir = OS::path($path, $file);
    if(is_dir($subDir)) {
        for (_, $file in scandir($subDir)) {
            $filePath = OS::path($subDir, $file);
            if(!is_dir($filePath)) {
                try {
                    if($filePath->indexOf(".php") == $filePath->length - 4) {
                        Log::info("执行 {$filePath}");
                        include($filePath);
                    } else {
                        continue;
                    }
                } catch (Exception $e) {
                    Log::fatal("执行文件发生错误, file={$filePath}; error={$e->getMessage()}")
                } catch (Error $e) {
                    Log::fatal("执行文件发生Error, file={$filePath}; error=" + $e->getMessage())
                }
            }
        }
    }
}

Log::info("🎉 接口功能测试完成");