<?php

namespace tests\php;

use Container\Annotation\Bind;
use Container\Annotation\Inject;
use Container\Annotation\Named;
use Container\Annotation\Singleton;
use Database\Annotation\Column;
use Database\Annotation\Table;
use Net\Annotation\GetMapping;
use Net\Annotation\Middleware;
use Cli\Annotation\Command;

/**
 * ReflectionAttribute::getArguments / newInstance 与 std 注解及用户属性的兼容性。
 */

#[\Attribute]
class ReflectionAttr_CustomMark
{
    public function __construct(
        public string $label,
        public int $count = 0,
    ) {
    }
}

#[Table(name: 'reflection_attr_users')]
class ReflectionAttr_Entity
{
    #[Column(name: 'id')]
    public int $id = 0;
}

#[Middleware(className: 'App\\AuthMiddleware')]
class ReflectionAttr_Service
{
}

class ReflectionAttr_Controller
{
    #[GetMapping(path: '/reflection/hello')]
    public function hello(): string
    {
        return 'ok';
    }

    #[Command(name: 'reflect:ping', description: 'ping command')]
    public function ping(): void
    {
    }
}

#[ReflectionAttr_CustomMark('mix', count: 7)]
#[ReflectionAttr_CustomMark('first')]
class ReflectionAttr_MixedArgs
{
}

#[Attribute(Attribute::TARGET_PROPERTY | Attribute::TARGET_PARAMETER)]
class ReflectionAttr_InjectConfig
{
    public function __construct(public string $key)
    {
    }
}

#[Attribute(Attribute::TARGET_FUNCTION)]
class ReflectionAttr_FunctionOnly
{
    public function __construct()
    {
    }
}

#[Singleton(name: 'loom.module')]
class ReflectionAttr_SingletonHolder
{
}

interface ReflectionAttr_BindContract
{
}

#[Bind(abstract: 'tests\\php\\ReflectionAttr_BindContract')]
class ReflectionAttr_BindHolder
{
}

class ReflectionAttr_ContainerCtor
{
    public function __construct(
        #[Inject(service: 'db.main')]
        #[Named(name: 'primary')]
        mixed $conn = null,
    ) {
    }
}

class ReflectionAttr_ParamTargets
{
    public function __construct(
        #[ReflectionAttr_InjectConfig('app.cache')]
        mixed $cache = null,
    ) {
    }
}

function assertArgsNonEmpty($attr, string $key, $expected): void
{
    $args = $attr->getArguments();
    if (!is_array($args) || !array_key_exists($key, $args)) {
        \Log::fatal("getArguments 缺少键 {$key}: " . var_export($args, true));
    }
    if ($args[$key] !== $expected) {
        \Log::fatal("getArguments[{$key}] 期望 " . var_export($expected, true) . '，实际 ' . var_export($args[$key], true));
    }
}

// --- Table ---
$entityRef = new \ReflectionClass(ReflectionAttr_Entity::class);
$tableAttrs = $entityRef->getAttributes(Table::class);
if (count($tableAttrs) !== 1) {
    \Log::fatal('Table 注解数量应为 1');
}
assertArgsNonEmpty($tableAttrs[0], 'name', 'reflection_attr_users');
$tableInst = $tableAttrs[0]->newInstance();
if ($tableInst->name !== 'reflection_attr_users') {
    \Log::fatal('Table newInstance()->name 错误: ' . var_export($tableInst->name, true));
}

// --- Column 回归 ---
$colProp = $entityRef->getProperty('id');
$colAttrs = $colProp->getAttributes(Column::class);
if (count($colAttrs) !== 1) {
    \Log::fatal('Column 注解数量应为 1');
}
assertArgsNonEmpty($colAttrs[0], 'name', 'id');
$colInst = $colAttrs[0]->newInstance();
if ($colInst->name !== 'id') {
    \Log::fatal('Column newInstance()->name 错误');
}

// --- Middleware + Singleton ---
$svcRef = new \ReflectionClass(ReflectionAttr_Service::class);
$mwAttrs = $svcRef->getAttributes(Middleware::class);
if (count($mwAttrs) !== 1) {
    \Log::fatal('Middleware 注解缺失');
}
assertArgsNonEmpty($mwAttrs[0], 'className', 'App\\AuthMiddleware');
$mwInst = $mwAttrs[0]->newInstance();
if ($mwInst->className !== 'App\\AuthMiddleware') {
    \Log::fatal('Middleware newInstance()->className 为空或错误');
}
// --- GetMapping / Command ---
$ctrlRef = new \ReflectionClass(ReflectionAttr_Controller::class);
$hello = $ctrlRef->getMethod('hello');
$mapAttrs = $hello->getAttributes(GetMapping::class);
if (count($mapAttrs) !== 1) {
    \Log::fatal('GetMapping 注解缺失');
}
assertArgsNonEmpty($mapAttrs[0], 'path', '/reflection/hello');
$mapInst = $mapAttrs[0]->newInstance();
if ($mapInst->path !== '/reflection/hello') {
    \Log::fatal('GetMapping newInstance()->path 错误: ' . var_export($mapInst->path, true));
}

$ping = $ctrlRef->getMethod('ping');
$cmdAttrs = $ping->getAttributes(Command::class);
if (count($cmdAttrs) !== 1) {
    \Log::fatal('Command 注解缺失');
}
assertArgsNonEmpty($cmdAttrs[0], 'name', 'reflect:ping');
assertArgsNonEmpty($cmdAttrs[0], 'description', 'ping command');
$cmdInst = $cmdAttrs[0]->newInstance();
if ($cmdInst->name !== 'reflect:ping' || $cmdInst->description !== 'ping command') {
    \Log::fatal('Command newInstance 属性错误');
}

// --- Container：Singleton / Bind / Inject / Named newInstance 不得抛 target 类型错误 ---
$singletonRef = new \ReflectionClass(ReflectionAttr_SingletonHolder::class);
$singletonAttr = $singletonRef->getAttributes(Singleton::class)[0];
assertArgsNonEmpty($singletonAttr, 'name', 'loom.module');
$singletonInst = $singletonAttr->newInstance();
if (($singletonInst->name ?? '') !== 'loom.module') {
    \Log::fatal('Singleton newInstance()->name 错误');
}

$bindRef = new \ReflectionClass(ReflectionAttr_BindHolder::class);
$bindAttr = $bindRef->getAttributes(Bind::class)[0];
assertArgsNonEmpty($bindAttr, 'abstract', 'tests\\php\\ReflectionAttr_BindContract');
$bindInst = $bindAttr->newInstance();
if (($bindInst->abstract ?? '') !== 'tests\\php\\ReflectionAttr_BindContract') {
    \Log::fatal('Bind newInstance()->abstract 错误');
}

$ctorRef = new \ReflectionClass(ReflectionAttr_ContainerCtor::class)->getConstructor();
$ctorParam = $ctorRef->getParameters()[0];
$injectAttr = $ctorParam->getAttributes(Inject::class)[0];
$injectInst = $injectAttr->newInstance();
if (($injectInst->service ?? '') !== 'db.main') {
    \Log::fatal('Inject newInstance()->service 错误');
}
$namedAttr = $ctorParam->getAttributes(Named::class)[0];
$namedInst = $namedAttr->newInstance();
if (($namedInst->name ?? '') !== 'primary') {
    \Log::fatal('Named newInstance()->name 错误');
}

// --- 参数属性 target：PROPERTY|PARAMETER 合法；仅 FUNCTION 的注解在参数上应失败 ---
$paramRef = new \ReflectionClass(ReflectionAttr_ParamTargets::class)->getConstructor()->getParameters()[0];
$configAttr = $paramRef->getAttributes(ReflectionAttr_InjectConfig::class)[0];
$configArgs = $configAttr->getArguments();
$configKey = $configArgs['key'] ?? ($configArgs[0] ?? null);
if ($configKey !== 'app.cache') {
    \Log::fatal('InjectConfig getArguments key 错误: ' . var_export($configArgs, true));
}
$configInst = $configAttr->newInstance();
if ($configInst->key !== 'app.cache') {
    \Log::fatal('InjectConfig newInstance()->key 错误');
}

class ReflectionAttr_FunctionOnlyParam
{
    public function __construct(
        #[ReflectionAttr_FunctionOnly]
        mixed $x = null,
    ) {
    }
}
$badParam = (new \ReflectionClass(ReflectionAttr_FunctionOnlyParam::class))->getConstructor()->getParameters()[0];
$badAttr = $badParam->getAttributes(ReflectionAttr_FunctionOnly::class)[0];
try {
    $badAttr->newInstance();
    \Log::fatal('仅 TARGET_FUNCTION 的注解在参数上 newInstance 应抛 Error');
} catch (\Error $e) {
}

// --- 用户自定义属性：混用位置参数与命名参数 ---
$mixedRef = new \ReflectionClass(ReflectionAttr_MixedArgs::class);
$customAttrs = $mixedRef->getAttributes(ReflectionAttr_CustomMark::class);
if (count($customAttrs) !== 2) {
    \Log::fatal('自定义重复属性数量应为 2');
}
$argsNamed = $customAttrs[0]->getArguments();
$labelArg = $argsNamed['label'] ?? ($argsNamed[0] ?? null);
if ($labelArg !== 'mix' || !isset($argsNamed['count']) || $argsNamed['count'] !== 7) {
    \Log::fatal('混合命名参数 getArguments 错误: ' . var_export($argsNamed, true));
}
$instNamed = $customAttrs[0]->newInstance();
if ($instNamed->label !== 'mix' || $instNamed->count !== 7) {
    \Log::fatal('用户属性 newInstance 命名参数失败');
}
$argsPos = $customAttrs[1]->getArguments();
if (!isset($argsPos[0]) || $argsPos[0] !== 'first') {
    \Log::fatal('位置参数 getArguments 错误: ' . var_export($argsPos, true));
}

// getArguments 不得包含编译器注入的 target（std 注解构造函数的第二形参）
foreach ([$tableAttrs[0], $mapAttrs[0], $mwAttrs[0]] as $stdAttr) {
    $a = $stdAttr->getArguments();
    if (array_key_exists('target', $a)) {
        \Log::fatal('getArguments 不得包含 target: ' . var_export($a, true));
    }
}

\Log::info('reflection_attribute_args 测试通过');
