<?php

require dirname(__DIR__) . '/vendor/autoload.php';

use Illuminate\Notifications\Action;
use Illuminate\Notifications\AnonymousNotifiable;
use Illuminate\Notifications\Notification;

/**
 * 不定义 Notification 子类（Origami 上 extends Notification 会触发父类加载失败）。
 * 验证 Notification 基类、Action、AnonymousNotifiable 路由。
 */
$notification = new Notification();
$notification->locale('en');
if ($notification->locale !== 'en') {
    echo "FAIL: locale property\n";
    exit(1);
}

$action = new Action('View', 'https://example.com');
if ($action->text !== 'View' || $action->url !== 'https://example.com') {
    echo "FAIL: Action\n";
    exit(1);
}

$anon = new AnonymousNotifiable();
$anon->route('mail', 'user@example.com');
if ($anon->routeNotificationFor('mail') !== 'user@example.com') {
    echo "FAIL: AnonymousNotifiable route\n";
    exit(1);
}

echo "PASS\n";
