<?php

require_once __DIR__.'/startup.php';

spl_autoload_register(function ($class) {
    $classes = [
        'HotReload\\Controller\\GreetingController' => '/src/Controller/GreetingController.php',
        'HotReload\\Service\\GreetingService' => '/src/Service/GreetingService.php',
        'HotReload\\View\\GreetingView' => '/src/View/GreetingView.php',
    ];
    if (isset($classes[$class])) {
        require_once __DIR__.$classes[$class];
    }
});
