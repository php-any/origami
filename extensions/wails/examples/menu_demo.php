<?php
/**
 * Wails 菜单 + 对话框 + 快捷键 示例
 *
 * 演示: 完整菜单栏构建、快捷键绑定、文件对话框、消息对话框
 *
 * 运行方式:
 *   在 Go main 中加载 wails 扩展，然后执行此脚本。
 */

use Wails\Application;
use Wails\Options\App;
use Wails\Menu\Menu;
use Wails\Menu\MenuItem;
use Wails\Menu\Keys;
use Wails\MenuItemType;
use Wails\Runtime\Window;
use Wails\Runtime\Dialog;
use Wails\Runtime\Log;
use Wails\Runtime\Browser;
use Wails\Dialog\FileFilter;
use Wails\Dialog\OpenDialogOptions;
use Wails\Dialog\MessageDialogOptions;
use Wails\DialogType;

// ── 1. 构建菜单栏 ──
$menu = new Menu();

// —— File 菜单 ——
$fileMenu = new Menu();
$fileMenu->addText("&Open File", Keys::cmdOrCtrl("o"), function () {
    $opts = new OpenDialogOptions([
        'Title'    => '选择文件',
        'Filters'  => [
            new FileFilter("所有文件", "*.*"),
            new FileFilter("文本文件", "*.txt;*.md"),
        ],
    ]);
    $result = Dialog::openFile($opts);
    if ($result) {
        Log::info("选择了文件: " . $result);
    }
});

$fileMenu->addText("Open &Directory", Keys::cmdOrCtrl("d"), function () {
    $result = Dialog::openDirectory(new OpenDialogOptions([
        'Title' => '选择目录',
    ]));
    if ($result) {
        Log::info("选择了目录: " . $result);
    }
});

$fileMenu->addSeparator();

$fileMenu->addText("&Save As...", Keys::cmdOrCtrl("s"), function () {
    $result = Dialog::saveFile(new \Wails\Dialog\SaveDialogOptions([
        'Title'           => '保存文件',
        'DefaultFilename' => 'untitled.txt',
        'Filters'         => [
            new FileFilter("文本文件", "*.txt"),
        ],
    ]));
    if ($result) {
        Log::info("保存到: " . $result);
    }
});

$fileMenu->addSeparator();

$fileMenu->addText("E&xit", Keys::cmdOrCtrl("q"), function () {
    Application::quit();
});

// —— Edit 菜单 ——
$editMenu = new Menu();
$editMenu->addText("&Undo",  Keys::cmdOrCtrl("z"), function () {
    Log::info("Undo");
});
$editMenu->addText("&Redo",  Keys::combo("y", ["shift"]), function () {
    Log::info("Redo");
});
$editMenu->addSeparator();
$editMenu->addText("Cu&t",   Keys::cmdOrCtrl("x"), function () {
    Log::info("Cut");
});
$editMenu->addText("&Copy",  Keys::cmdOrCtrl("c"), function () {
    Log::info("Copy");
});
$editMenu->addText("&Paste", Keys::cmdOrCtrl("v"), function () {
    Log::info("Paste");
});

// —— View 菜单 (Checkbox + Radio) ——
$viewMenu = new Menu();
$viewMenu->addCheckbox("Show &Toolbar", true, "", function ($data) {
    Log::info("Toolbar: " . ($data->Checked ? "ON" : "OFF"));
});
$viewMenu->addCheckbox("Show &Status Bar", false, "", function ($data) {
    Log::info("StatusBar: " . ($data->Checked ? "ON" : "OFF"));
});
$viewMenu->addSeparator();
$viewMenu->addRadio("&Small Icons",  false, "", function () {
    Log::info("ViewMode: Small");
});
$viewMenu->addRadio("&Large Icons",  true, "", function () {
    Log::info("ViewMode: Large");
});
$viewMenu->addRadio("&Details",      false, "", function () {
    Log::info("ViewMode: Details");
});

// —— Help 菜单 ——
$helpMenu = new Menu();
$helpMenu->addText("&About", "", function () {
    Dialog::message(new MessageDialogOptions([
        'Type'    => DialogType::INFO,
        'Title'   => '关于',
        'Message' => "Origami + Wails v2\nPowered by PHP-like runtime",
        'Buttons' => ['OK'],
    ]));
});
$helpMenu->addText("&Documentation", "", function () {
    Browser::openURL("https://wails.io/docs");
});

// —— 组装顶层菜单 ——
$menu->addSubMenu("&File", $fileMenu);
$menu->addSubMenu("&Edit", $editMenu);
$menu->addSubMenu("&View", $viewMenu);
$menu->addSubMenu("&Help", $helpMenu);

// ── 2. 右键上下文菜单 ──
$ctxMenu = new Menu();
$ctxMenu->addText("刷新", Keys::parse("F5"), function () {
    Window::reload();
});
$ctxMenu->addText("全屏", Keys::parse("F11"), function () {
    Window::fullscreen();
});
$ctxMenu->addSeparator();
$ctxMenu->addText("检查元素", Keys::combo("i", ["shift", "ctrl"]), function () {
    Log::info("Inspector (debug mode)");
});

// ── 3. 配置应用 ──
$options = new App([
    'Title'  => '📋 Wails Menu + Dialog Demo',
    'Width'  => 900,
    'Height' => 640,
    'Menu'   => $menu,
]);

$options->onDomReady(function () {
    Window::center();
    Log::info("菜单和对话框示例已就绪");
    Log::info("试试快捷键: Ctrl+O 打开文件, Ctrl+Q 退出");
});

$options->onBeforeClose(function () {
    $result = Dialog::message(new MessageDialogOptions([
        'Type'    => DialogType::QUESTION,
        'Title'   => '确认退出',
        'Message' => '确定要退出应用吗？',
        'Buttons' => ['确定', '取消'],
        'DefaultButton' => '确定',
        'CancelButton'  => '取消',
    ]));
    return $result !== '确定';
});

// ── 4. 启动 ──
Application::run($options);
