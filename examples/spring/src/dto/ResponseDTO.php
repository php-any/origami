<?php

namespace Spring\DTO;

/**
 * 统一响应 DTO（历史参考）
 *
 * 标准库已内置相同语义，请优先使用 Net\Http\Response：
 * - $response->success($data, $message = 'success', $status = 200)
 * - $response->error($message = 'error', $code = 500, $data = null)
 * - $server->onFormat(function ($code, $message, $data) { ... })  // 自定义信封
 *
 * 默认 JSON 信封：{ code, message, data, timestamp }
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

    public static function success(mixed $data = null, string $message = 'success'): self {
        return new self(200, $message, $data);
    }

    public static function error(int $code = 500, string $message = 'error', mixed $data = null): self {
        return new self($code, $message, $data);
    }

    public function toArray(): array {
        return [
            'code' => $this->code,
            'message' => $this->message,
            'data' => $this->data,
            'timestamp' => $this->timestamp
        ];
    }

    public function toJson(): string {
        return json_encode($this->toArray());
    }
}
