<?php
use Fyne\App;
use Fyne\Container;
use Fyne\Widget\Label;
use Fyne\Widget\Button;
use Fyne\Widget\Entry;

$app = new App("com.example.hello");
$window = $app->newWindow("Hello Fyne");

$label = new Label("Hello, Origami + Fyne!");
$entry = new Entry();
$entry->setPlaceHolder("输入你的名字...");

$button = new Button("问候", function() use ($label, $entry) {
    $name = $entry->getText();
    if ($name) {
        $label->setText("你好, " . $name . "!");
    } else {
        $label->setText("Hello, World!");
    }
});

$quitBtn = new Button("退出", function() use ($app) {
    $app->quit();
});

$content = Container::newVBox([$label, $entry, $button, $quitBtn]);
$window->setContent($content);
$window->resize(new \Fyne\Size(400, 300));
$window->centerOnScreen();
$window->showAndRun();
