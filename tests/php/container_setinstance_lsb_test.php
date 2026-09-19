<?php

namespace tests\php;

/**
 * Laravel Container::$instance 声明在父类；子类 Application->setInstance()
 * 必须写入父类槽位，Container::getInstance() 才能读到。
 * 登录页 csrf_token() 走 app() → Container::getInstance()。
 */

class ContainerLsb_Box
{
    protected static $instance;

    public $tok = '';

    public static function getInstance()
    {
        return static::$instance ??= new static;
    }

    public static function setInstance($container = null)
    {
        return static::$instance = $container;
    }
}

class ContainerLsb_App extends ContainerLsb_Box
{
}

$boot = new ContainerLsb_App();
$boot->tok = 'BOOT';
ContainerLsb_Box::setInstance($boot);

$req = new ContainerLsb_App();
$req->tok = 'REQUEST';
$req->setInstance($req);

$got = ContainerLsb_Box::getInstance();
if (!($got instanceof ContainerLsb_App) || $got->tok !== 'REQUEST') {
    $label = is_object($got) ? ($got->tok ?? get_class($got)) : var_export($got, true);
    Log::fatal('Container::getInstance 未读到子类 setInstance: '.$label);
}

if (ContainerLsb_App::getInstance()->tok !== 'REQUEST') {
    Log::fatal('Application::getInstance 与 Container::getInstance 应共享 $instance');
}

Log::info('container_setinstance_lsb 测试通过');
