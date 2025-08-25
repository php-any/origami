<?php
namespace database\sql;


/**
 * conn - database\sql\Conn 类
 * 
 * 此文件包含 conn 模块的伪代码接口定义
 * 这些是自动生成的接口，仅用于参考，不包含具体实现
 */



/**
 * database\sql\Conn 类
 * Conn 类
 */
class Conn {

    /**
     * beginTx 方法
     * 
     */
    public static function beginTx($ctx, $opts) {
        // 实现逻辑
    }

    /**
     * close 方法
     * 
     */
    public static function close() {
        // 实现逻辑
    }

    /**
     * execContext 方法
     * 
     */
    public static function execContext($ctx, $query, $args) {
        // 实现逻辑
    }

    /**
     * pingContext 方法
     * 
     */
    public static function pingContext($ctx) {
        // 实现逻辑
    }

    /**
     * prepareContext 方法
     * 
     */
    public static function prepareContext($ctx, $query) {
        // 实现逻辑
    }

    /**
     * queryContext 方法
     * 
     */
    public static function queryContext($ctx, $query, $args) {
        // 实现逻辑
    }

    /**
     * queryRowContext 方法
     * 
     */
    public static function queryRowContext($ctx, $query, $args) {
        // 实现逻辑
    }

    /**
     * raw 方法
     * 
     */
    public static function raw($f) {
        // 实现逻辑
    }

}

