<?php

namespace Spring\Config;

/**
 * 应用配置类
 */
class AppConfig {
    
    // 应用名称
    const APP_NAME = 'Spring Demo';
    
    // 应用版本
    const APP_VERSION = '1.0.0';
    
    // 默认时区
    const DEFAULT_TIMEZONE = 'Asia/Shanghai';
    
    // 服务器配置
    const SERVER_HOST = '0.0.0.0';
    const SERVER_PORT = 8080;
    
    // API 配置
    const API_PREFIX = '/api';
    const API_VERSION = 'v1';
    
    // 认证配置
    const TOKEN_EXPIRY = 3600; // 1 小时
    
    // 分页配置
    const PAGE_SIZE = 20;
    const MAX_PAGE_SIZE = 100;
    
    // CORS 配置
    const ALLOWED_ORIGINS = ['*'];
    const ALLOWED_METHODS = ['GET', 'POST', 'PUT', 'DELETE', 'OPTIONS'];
    const ALLOWED_HEADERS = ['Content-Type', 'Authorization'];
    
    /**
     * 获取配置项
     */
    public static function get($key, $default = null) {
        $constants = [
            'app.name' => self::APP_NAME,
            'app.version' => self::APP_VERSION,
            'app.timezone' => self::DEFAULT_TIMEZONE,
            'server.host' => self::SERVER_HOST,
            'server.port' => self::SERVER_PORT,
            'api.prefix' => self::API_PREFIX,
            'api.version' => self::API_VERSION,
            'auth.token_expiry' => self::TOKEN_EXPIRY,
            'pagination.page_size' => self::PAGE_SIZE,
            'pagination.max_page_size' => self::MAX_PAGE_SIZE,
        ];
        
        return $constants[$key] ?? $default;
    }
}
