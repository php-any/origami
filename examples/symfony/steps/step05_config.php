<?php

use Symfony\Component\Config\FileLocator;
use Symfony\Component\Config\Loader\LoaderResolver;
use Symfony\Component\DependencyInjection\Loader\PhpFileLoader;
use Symfony\Component\DependencyInjection\ContainerBuilder;

// step05: symfony/config + symfony/dependency-injection loader

$container = new ContainerBuilder();
$locator = new FileLocator([dirname(__DIR__) . '/config']);
$loader = new PhpFileLoader($container, $locator);
$loader->load('services.php');
$container->compile();

step_check('PhpFileLoader load', $container->has('app.greeting'));
step_check('service greeting value', $container->getParameter('app.greeting') === 'Hello Symfony');

step_check('step05_config', true, 'symfony/config');
