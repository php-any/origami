<?php

use Symfony\Component\EventDispatcher\EventDispatcher;
use Symfony\Component\EventDispatcher\GenericEvent;

// step03: symfony/event-dispatcher

$dispatcher = new EventDispatcher();
$called = false;
$dispatcher->addListener('app.test', function () use (&$called) {
    $called = true;
});
$dispatcher->dispatch(new GenericEvent(), 'app.test');
step_check('EventDispatcher listener', $called === true);

$order = [];
$dispatcher->addListener('app.order', function () use (&$order) {
    $order[] = 'second';
}, 0);
$dispatcher->addListener('app.order', function () use (&$order) {
    $order[] = 'first';
}, 10);
$dispatcher->dispatch(new GenericEvent(), 'app.order');
step_check('EventDispatcher priority', $order === ['first', 'second'], implode(',', $order));

step_check('step03_event_dispatcher', true, 'symfony/event-dispatcher');
