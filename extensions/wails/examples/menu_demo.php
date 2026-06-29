<?php
/**
 * Wails v3 菜单 + 对话框 + 快捷键 示例
 *
 * 演示 Wails v3 的:
 *   - 应用菜单栏构建 (Menu + MenuItem)
 *   - 键盘快捷键 (Keys::cmdOrCtrl / Keys::combo)
 *   - 文件对话框 (Dialog::openFile / openDirectory / saveFile)
 *   - 消息对话框 (Dialog::message: info / warning / error / question)
 *   - 窗口操作 (center / reload / fullscreen / quit)
 *
 * 对应 Go API:
 *   app.Dialog.OpenFile().SetTitle(...).AddFilter(...).Show()
 *   app.Dialog.Question().SetTitle(...).SetMessage(...).SetButtons(...).Show()
 */

use Wails\Application;
use Wails\Options\App;
use Wails\Menu\Menu;
use Wails\Menu\Keys;
use Wails\MenuItemType;
use Wails\Runtime\Window;
use Wails\Runtime\Dialog;
use Wails\Runtime\Log;
use Wails\Dialog\FileFilter;
use Wails\Dialog\OpenDialogOptions;
use Wails\Dialog\SaveDialogOptions;
use Wails\Dialog\MessageDialogOptions;
use Wails\DialogType;

// ── 1. 构建菜单栏 ──
$menu = new Menu();

// ====== File 菜单 ======
$fileMenu = new Menu();

$fileMenu->addText("&Open...", Keys::cmdOrCtrl("o"), function () {
    $opts = new OpenDialogOptions([
        'Title'   => 'Open File',
        'Filters' => [
            new FileFilter("All Files", "*.*"),
            new FileFilter("Text Files", "*.txt;*.md;*.csv"),
        ],
    ]);
    $result = Dialog::openFile($opts);
    if ($result) {
        Log::info("Opened file: " . $result);
    }
});

$fileMenu->addText("Open &Folder...", Keys::cmdOrCtrl("d"), function () {
    $result = Dialog::openDirectory(new OpenDialogOptions([
        'Title' => 'Choose Folder',
    ]));
    if ($result) {
        Log::info("Opened folder: " . $result);
    }
});

$fileMenu->addSeparator();

$fileMenu->addText("&Save As...", Keys::cmdOrCtrl("s"), function () {
    $result = Dialog::saveFile(new SaveDialogOptions([
        'Title'           => 'Save File',
        'DefaultFilename' => 'untitled.txt',
        'Filters'         => [
            new FileFilter("Text Files", "*.txt"),
            new FileFilter("All Files", "*.*"),
        ],
    ]));
    if ($result) {
        Log::info("Saved to: " . $result);
    }
});

$fileMenu->addSeparator();

$fileMenu->addText("E&xit", Keys::cmdOrCtrl("q"), function () {
    Application::quit();
});

// ====== Edit 菜单 ======
$editMenu = new Menu();

$editMenu->addText("&Undo",  Keys::cmdOrCtrl("z"), function () {
    Log::info("Undo");
});

$editMenu->addText("&Redo",  Keys::combo("z", ["shift"]), function () {
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

// ====== View 菜单 (Checkbox + Radio) ======
$viewMenu = new Menu();

$viewMenu->addCheckbox("Show &Toolbar", true, "", function ($data) {
    $status = $data->Checked ? "ON" : "OFF";
    Log::info("Toolbar: {$status}");
});

$viewMenu->addCheckbox("Show &Status Bar", false, "", function ($data) {
    $status = $data->Checked ? "ON" : "OFF";
    Log::info("Status Bar: {$status}");
});

$viewMenu->addSeparator();

$viewMenu->addRadio("&Small Icons",  false, "", function () {
    Log::info("View mode: Small Icons");
});

$viewMenu->addRadio("&Large Icons",  true, "", function () {
    Log::info("View mode: Large Icons");
});

$viewMenu->addRadio("&Details",      false, "", function () {
    Log::info("View mode: Details");
});

// ====== Help 菜单 ======
$helpMenu = new Menu();

$helpMenu->addText("&About", "", function () {
    Dialog::message(new MessageDialogOptions([
        'Type'    => DialogType::INFO,
        'Title'   => 'About',
        'Message' => "Origami + Wails v3 Demo\nBuilt with the Origami PHP-like runtime",
        'Buttons' => ['OK'],
    ]));
});

$helpMenu->addText("&Documentation", "", function () {
    // Browser::openURL("https://wails.io/docs");
    Log::info("Opening documentation...");
});

// ====== 组装顶层菜单栏 ======
$menu->addSubMenu("&File", $fileMenu);
$menu->addSubMenu("&Edit", $editMenu);
$menu->addSubMenu("&View", $viewMenu);
$menu->addSubMenu("&Help", $helpMenu);

// ── 2. 配置应用 ──
$options = new App([
    'Title'  => '📋 Wails v3 Menu + Dialog Demo',
    'Width'  => 900,
    'Height' => 640,
    'Menu'   => $menu,
]);

// ── 3. 生命周期 ──
$options->onDomReady(function () {
    Window::center();
    Log::info("Menu and dialog demo is ready");
    Log::info("Try shortcuts: Ctrl+O = Open, Ctrl+Q = Quit");
});

// 关闭前确认对话框
$options->onBeforeClose(function () {
    $result = Dialog::message(new MessageDialogOptions([
        'Type'          => DialogType::QUESTION,
        'Title'         => 'Confirm Exit',
        'Message'       => 'Are you sure you want to quit?',
        'Buttons'       => ['Quit', 'Cancel'],
        'DefaultButton' => 'Quit',
    ]));
    return $result !== 'Quit';
});

// ── 4. 启动 ──
Application::run($options);
