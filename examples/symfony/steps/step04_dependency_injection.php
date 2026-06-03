<?php

use Symfony\Component\DependencyInjection\ContainerBuilder;
use Symfony\Component\DependencyInjection\Definition;
use Symfony\Component\DependencyInjection\Reference;

// step04: symfony/dependency-injection

$container = new ContainerBuilder();

$container->setDefinition('logger', (new Definition(\stdClass::class))->setPublic(true));
$container->setDefinition('app.service', (new Definition(\stdClass::class))
    ->setPublic(true)
    ->addMethodCall('__construct', []));

$container->compile();
step_check('ContainerBuilder compile', $container->has('logger'));

$obj = $container->get('logger');
step_check('Container get service', $obj instanceof \stdClass);

$container2 = new ContainerBuilder();
$container2->register('inner', \stdClass::class)->setPublic(true);
$container2->register('outer', \stdClass::class)
    ->setPublic(true)
    ->addArgument(new Reference('inner'));
$container2->compile();
step_check('Container Reference', $container2->get('outer') instanceof \stdClass);

step_check('step04_dependency_injection', true, 'symfony/dependency-injection');
