<?php
namespace database\sql;


/**
 * tx - database\sql\Tx 类
 * 
 * 此文件包含 tx 模块的伪代码接口定义
 * 这些是自动生成的接口，仅用于参考，不包含具体实现
 */



/**
 * database\sql\Tx 类
 * Tx 类
 */
class Tx {

    /**
     * commit 方法
     * 
     */
    public static function commit() {
        // 实现逻辑
    }

    /**
     * exec 方法
     * 
     */
    public static function exec($query, $args) {
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
     * prepare 方法
     * 
     */
    public static function prepare($query) {
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
     * query 方法
     * 
     */
    public static function query($query, $args) {
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
     * queryRow 方法
     * 
     */
    public static function queryRow($query, $args) {
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
     * rollback 方法
     * 
     */
    public static function rollback() {
        // 实现逻辑
    }

    /**
     * stmt 方法
     * 
     */
    public static function stmt($stmt) {
        // 实现逻辑
    }

    /**
     * stmtContext 方法
     * 
     */
    public static function stmtContext($ctx, $stmt) {
        // 实现逻辑
    }

}

