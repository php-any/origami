<?php

namespace tests\basic;

echo "=== enum 方法测试 ===\n";

// 测试 enum 中的方法
enum Status: string {
    case OPEN = 'open';
    case CLOSED = 'closed';
    
    public function getLabel(): string {
        return match($this->value) {
            'open' => '开放',
            'closed' => '关闭',
            default => '未知',
        };
    }
    
    public static function fromLabel(string $label) {
        return match($label) {
            '开放' => Status::OPEN,
            '关闭' => Status::CLOSED,
            default => null,
        };
    }
}

// 测试实例方法
$status = Status::OPEN;
$label = $status->getLabel();
if($label === '开放') {
    Log::info("enum 实例方法测试通过");
} else {
    Log::fatal("enum 实例方法测试失败: " . $label);
}

// 测试静态方法
$status2 = Status::fromLabel('关闭');
if($status2 === Status::CLOSED) {
    Log::info("enum 静态方法测试通过");
} else {
    Log::fatal("enum 静态方法测试失败");
}

// 测试静态方法返回 null
$status3 = Status::fromLabel('不存在');
if($status3 === null) {
    Log::info("enum 静态方法返回 null 测试通过");
} else {
    Log::fatal("enum 静态方法返回 null 测试失败");
}

echo "=== enum 方法测试完成 ===\n";

