<?php

namespace Spring\DTO;

/**
 * 统一响应 DTO
 */
class ResponseDTO {
    private $code;
    private $message;
    private $data;
    private $timestamp;
    
    public function __construct($code = 200, $message = 'success', $data = null) {
        $this->code = $code;
        $this->message = $message;
        $this->data = $data;
        $this->timestamp = time();
    }
    
    public function getCode() {
        return $this->code;
    }
    
    public function setCode($code) {
        $this->code = $code;
    }
    
    public function getMessage() {
        return $this->message;
    }
    
    public function setMessage($message) {
        $this->message = $message;
    }
    
    public function getData() {
        return $this->data;
    }
    
    public function setData($data) {
        $this->data = $data;
    }
    
    public function getTimestamp() {
        return $this->timestamp;
    }
    
    /**
     * 成功响应
     */
    public static function success($data = null, $message = 'success') {
        return new self(200, $message, $data);
    }
    
    /**
     * 失败响应
     */
    public static function error($code = 500, $message = 'error', $data = null) {
        return new self($code, $message, $data);
    }
    
    /**
     * 转换为数组
     */
    public function toArray() {
        return [
            'code' => $this->code,
            'message' => $this->message,
            'data' => $this->data,
            'timestamp' => $this->timestamp
        ];
    }
    
    /**
     * JSON 序列化
     */
    public function toJson() {
        return json_encode($this->toArray());
    }
}
