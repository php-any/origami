<?php

namespace Net\Http;


/**
 * server - Net\Http\Server 类
 * 
 * 此文件包含 server 模块的伪代码接口定义
 * 这些是自动生成的接口，仅用于参考，不包含具体实现
 */



/**
 * Net\Http\Server 类
 * Server 类
 */
class Server {

    /**
     * get 方法
     * 
     */
    public function get($path, $handle) {
        // 实现逻辑
    }

    /**
     * post 方法
     * 
     */
    public function post($path, $handle) {
        // 实现逻辑
    }

    /**
     * put 方法
     * 
     */
    public function put($path, $handle) {
        // 实现逻辑
    }

    /**
     * delete 方法
     * 
     */
    public function delete($path, $handle) {
        // 实现逻辑
    }

    /**
     * head 方法
     * 
     */
    public function head($path, $handle) {
        // 实现逻辑
    }

    /**
     * options 方法
     * 
     */
    public function options($path, $handle) {
        // 实现逻辑
    }

    /**
     * patch 方法
     * 
     */
    public function patch($path, $handle) {
        // 实现逻辑
    }

    /**
     * trace 方法
     * 
     */
    public function trace($path, $handle) {
        // 实现逻辑
    }

    /**
     * static 方法
     * 
     */
    public function static(string $prefix, string $dir) {
        // 实现逻辑
    }

    /**
     * any 方法
     *
     */
    public function any($handle) {
        // 实现逻辑
    }

    /**
     * flash 方法（启动时扫描注解路由并直接注册到 Server）
     *
     * 在启动阶段扫描指定目录下的所有 Controller 类及其注解路由，
     * 将每个路由直接注册到 Server（通过 get/post/put/delete 等方法），
     * 而不是使用 any 兜底路由。请求到达时直接匹配已注册路由，无需再次扫描。
     *
     * @param string $dir 应用目录路径（包含 controllers、middleware 等子目录）
     */
    public function flash(string $dir) {
        // 实现逻辑
    }

    /**
     * group 方法
     * 
     */
    public function group($prefix) : Net\Http\Server {
        // 实现逻辑
    }

    /**
     * middleware 方法
     * 
     */
    public function middleware($mid) {
        // 实现逻辑
    }

    /**
     * run 方法
     * 
     */
    public function run() {
        // 实现逻辑
    }

    /**
     * serveHTTP 方法
     * 
     */
    public function serveHTTP($param0, $param1) {
        // 实现逻辑
    }

}

