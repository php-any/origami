<?php

/**
 * 最小问题代码（两层）：
 *
 * [1/2] MinimalDefinition：无 vendor，模仿 InputDefinition（setArguments([]) + addArguments + getArguments()）。
 *       PHP 与 Origami 均应得到 key=command。
 *
 * [2/2] 传参 symfony 时：用 Application + find('hello') + mergeApplicationDefinition()。
 *       完整失败路径：go run ./origami.go ./cli_test/app.php hello 张三
 *       （mergeApplicationDefinition 内 foreach app->getDefinition()->getArguments() 时 Origami 得 index=0 或 command=>null）
 *
 * 运行：php script.php  [symfony]  与  go run ./origami.go script.php  [symfony]
 */

// 完全模仿 InputDefinition：构造里 setDefinition，setArguments 里先清空再 addArguments
class MinimalDefinition
{
    public $arguments = [];

    public function __construct(array $definition = [])
    {
        $this->setDefinition($definition);
    }

    public function setDefinition(array $definition): void
    {
        $arguments = [];
        foreach ($definition as $item) {
            $arguments[] = $item;
        }
        $this->setArguments($arguments);
    }

    public function setArguments(array $arguments = []): void
    {
        $this->arguments = [];
        $this->addArguments($arguments);
    }

    public function addArguments(array $arguments = []): void
    {
        foreach ($arguments as $argument) {
            $this->addArgument($argument[0], $argument[1]);
        }
    }

    public function addArgument(string $name, $value): void
    {
        $this->arguments[$name] = $value;
    }

    public function getArguments(): array
    {
        return $this->arguments;
    }
}

// 与 Application::getDefaultInputDefinition() 一致：构造时传 [ [name, value], ... ]，内部 setArguments 再 addArguments，最后 getArguments() 应为 [ name => value ]
$def = new MinimalDefinition([['command', 'command-value']]);

$keys = [];
$values = [];
foreach ($def->getArguments() as $k => $v) {
    $keys[] = $k;
    $values[] = $v;
}

// 期望：唯一 key 是字符串 'command'，不是数字 0
if (count($keys) !== 1) {
    throw new \RuntimeException('期望 1 个元素，实际: ' . count($keys));
}
if ($keys[0] !== 'command') {
    throw new \RuntimeException(
        '期望 key=command，实际 key=' . (string) $keys[0]
    );
}

echo "[1/3] MinimalDefinition OK key=command\n";

// ---------------------------------------------------------------------------
// [2/3] 检查：方法返回数组时，调用方修改返回值，是否会“引用地”改到原数组
//      期待行为：返回的是值拷贝，调用方改 $b 不影响 $this->a
// ---------------------------------------------------------------------------
class ReturnArrayRefTest
{
    public array $a = [];

    public function __construct()
    {
        $this->a['command'] = 'orig';
    }

    public function getA(): array
    {
        return $this->a;
    }
}

$rt = new ReturnArrayRefTest();
$b = $rt->getA();
$b['command'] = 'changed';

if ($rt->a['command'] !== 'orig') {
    throw new \RuntimeException('返回数组被按引用修改：$rt->a[\"command\"]='.(string) $rt->a['command']);
}

echo "[2/3] ReturnArrayRefTest OK 返回数组不是引用\n";

// ---------------------------------------------------------------------------
// [3/3] 最小复现：mergeApplicationDefinition() 内会 foreach app->getDefinition()->getArguments()
//       Origami 下该 foreach 得到 index=0 和 index=command=>null，导致 addArgument(null) 报错
// 仅当传入参数 symfony 时执行：php script.php symfony  /  go run ./origami.go script.php symfony
// ---------------------------------------------------------------------------
$useSymfony = (count($_SERVER['argv'] ?? []) > 1 && $_SERVER['argv'][1] === 'symfony') || getenv('SCRIPT_SYMFONY');
if ($useSymfony && is_file(__DIR__ . '/cli_test/vendor/autoload.php')) {
    require __DIR__ . '/cli_test/vendor/autoload.php';
    $app = new \Go\Test\Application();
    $cmd = $app->find('hello');
    $cmd->mergeApplicationDefinition(); // 内部 setArguments(app->getDefinition()->getArguments())，foreach 时 Origami 得 key=0 及 command=>null
    echo "[2/2] mergeApplicationDefinition() OK\n";
}

echo "script.php 测试通过\n";
