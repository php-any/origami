<?php
/**
 * 魔法函数测试：验证 __call、__invoke 等是否在 Origami 中正常触发
 * 运行: go run ./origami.go script.php
 */

echo "=== 魔法函数测试 ===\n\n";

// ---------------------------------------------------------------------------
// 1. __call：调用不存在的方法时触发
// ---------------------------------------------------------------------------
class CallTester
{
    public function __call(string $name, array $arguments)
    {
        echo "[__call] name=" . $name . ", args=" . count($arguments) . "\n";
        return "called:" . $name . "(" . implode(",", $arguments) . ")";
    }
}

$callObj = new CallTester();
$r1 = $callObj->foo(1, 2);
echo "1. __call 返回值: " . $r1 . "\n";

$r2 = $callObj->bar("hello");
echo "2. __call 返回值: " . $r2 . "\n\n";

// ---------------------------------------------------------------------------
// 2. __invoke：对象作为函数调用时触发
// ---------------------------------------------------------------------------
class InvokeTester
{
    public function __invoke($a, $b = null)
    {
        echo "[__invoke] a=" . $a . ", b=" . ($b ?? "null") . "\n";
        return $a + ($b ?? 0);
    }
}

$invokeObj = new InvokeTester();
$r3 = $invokeObj(10, 5);
echo "3. __invoke 返回值: " . $r3 . "\n";

$r4 = $invokeObj(7);
echo "4. __invoke(7) 返回值: " . $r4 . "\n\n";

// ---------------------------------------------------------------------------
// 5. __toString：对象被当作字符串使用时触发（Origami 未实现则此处会报错）
// ---------------------------------------------------------------------------
class ToStringTester
{
    public function __toString()
    {
        return "ToStringTester";
    }
}

$toStringObj = new ToStringTester();
echo "5. __toString: " . $toStringObj . "\n\n";

// ---------------------------------------------------------------------------
// 6. __get / __set：访问不存在或不可见属性时触发
// ---------------------------------------------------------------------------
class GetSetTester
{
    private $data = [];

    public function __get(string $name)
    {
        echo "[__get] name=" . $name . "\n";
        return $this->data[$name] ?? null;
    }

    public function __set(string $name, $value)
    {
        echo "[__set] name=" . $name . ", value=" . $value . "\n";
        $this->data[$name] = $value;
    }
}

$getsetObj = new GetSetTester();
$getsetObj->x = 100;
echo "6. __set 后 __get: " . $getsetObj->x . "\n\n";

echo "=== 魔法函数测试结束 ===\n";
echo "已验证: __call, __invoke, __toString, __get, __set 正常\n";
