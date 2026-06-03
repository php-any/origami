<?php

use Symfony\Component\DependencyInjection\ContainerBuilder;

$container = new ContainerBuilder();
$container->setParameter('app.greeting', 'Hello Symfony');
$container->register('app.greeting', \stdClass::class)->setPublic(true);

return $container;
