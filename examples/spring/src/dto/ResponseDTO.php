<?php

namespace Spring\DTO;

/**
 * 统一响应 DTO
 */
class ResponseDTO {
    private int $code;
    private string $message;
    private mixed $data;
    private int $timestamp;

    public function __construct(int $code = 200, string $message = 'success', mixed $data = null) {
        $this->code = $code;
        $this->message = $message;
        $this->data = $data;
        $this->timestamp = time();
    }

    public function getCode(): int {
        return $this->code;
    }

    public function setCode(int $code): void {
        $this->code = $code;
    }

    public function getMessage(): string {
        return $this->message;
    }

    public function setMessage(string $message): void {
        $this->message = $message;
    }

    public function getData(): mixed {
        return $this->data;
    }

    public function setData(mixed $data): void {
        $this->data = $data;
    }

    public function getTimestamp(): int {
        return $this->timestamp;
    }

    /**
     * 成功响应
     */
    public static function success(mixed $data = null, string $message = 'success'): self {
        return new self(200, $message, $data);
    }

    /**
     * 失败响应
     */
    public static function error(int $code = 500, string $message = 'error', mixed $data = null): self {
        return new self($code, $message, $data);
    }

    /**
     * 转换为数组
     */
    public function toArray(): array {
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
    public function toJson(): string {
        return json_encode($this->toArray());
    }
}
