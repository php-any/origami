<?php

/**
 * reflect - Reflect 类
 * 
 * 此文件包含 reflect 模块的伪代码接口定义
 * 这些是自动生成的接口，仅用于参考，不包含具体实现
 */



/**
 * Reflect 类
 * Reflect 类
 */
class Reflect {

    /**
     * getClassInfo 方法
     * 
     */
    public function getClassInfo(string $className) : string {
        // 实现逻辑
    }

    /**
     * getMethodInfo 方法
     * 
     */
    public function getMethodInfo($className, $methodName) : string {
        // 实现逻辑
    }

    /**
     * getPropertyInfo 方法
     * 
     */
    public function getPropertyInfo($className, $propertyName) : string {
        // 实现逻辑
    }

    /**
     * listClasses 方法
     * 
     */
    public function listClasses() : array {
        // 实现逻辑
    }

    /**
     * listMethods 方法
     * 
     */
    public function listMethods($className) : array {
        // 实现逻辑
    }

    /**
     * listProperties 方法
     * 
     */
    public function listProperties($className) : array {
        // 实现逻辑
    }

    /**
     * getClassAnnotations 方法
     * 
     */
    public function getClassAnnotations($className) {
        // 实现逻辑
    }

    /**
     * getMethodAnnotations 方法
     * 
     */
    public function getMethodAnnotations($className, $methodName) {
        // 实现逻辑
    }

    /**
     * getPropertyAnnotations 方法
     * 
     */
    public function getPropertyAnnotations($className, $propertyName) {
        // 实现逻辑
    }

    /**
     * getAllAnnotations 方法
     * 
     */
    public function getAllAnnotations($className) {
        // 实现逻辑
    }

    /**
     * getAnnotationDetails 方法
     * 
     */
    public function getAnnotationDetails($className, $memberType, $memberName) {
        // 实现逻辑
    }

}

