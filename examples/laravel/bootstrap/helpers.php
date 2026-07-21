<?php

/**
 * 容器辅助函数（CLI 使用 getInstance 容器）
 */
function app_make(string $abstract): mixed
{
    return \Container\Container::getInstance()->make($abstract);
}
