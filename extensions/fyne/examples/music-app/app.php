<?php
/**
 * Origami Music — 音乐播放器 App 入口
 *
 * 功能：登录 → 歌单浏览 → 歌曲播放
 *
 * 构建 APK:
 *   cd extensions/fyne
 *   fyne package -os android -appID com.origami.music -icon Icon.png ./cmd/music_app/
 *
 * 桌面调试:
 *   go build -o music_app.exe ./cmd/music_app/
 *   ./music_app.exe
 *
 * 测试账号: admin / 123456  或  demo / demo
 */

use Fyne\App;
use Fyne\Size;
use MusicApp\Pages\LoginPage;

$app = new App("com.origami.music");
$window = $app->newWindow("Origami Music");

// 从登录页开始
LoginPage::show($window);

$window->resize(new Size(400, 700));
$window->centerOnScreen();
$window->showAndRun();
