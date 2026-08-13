<?php
/**
 * Dashboard — iOS 风格底部标签栏
 */
use Fyne\App;
use Fyne\Size;
use Fyne\Container;
use Fyne\Widget\Label;
use Fyne\Widget\BottomTabBar;
use App\Pages\MessagesPage;
use App\Pages\DiscoverPage;
use App\Pages\ProfilePage;

$app = new App("com.origami.dashboard");
$window = $app->newWindow("Dashboard");

function buildScreen(\Fyne\Window $win, string $tab): Container {
    // ── 内容区（先创建，因为 tabBar 回调需要引用 $content）──
    switch ($tab) {
        case 'messages': $content = MessagesPage::build(); $selectedIdx = 0; break;
        case 'discover': $content = DiscoverPage::build(); $selectedIdx = 1; break;
        case 'profile':  $content = ProfilePage::build();  $selectedIdx = 2; break;
        default:         $content = new Label('Unknown');   $selectedIdx = 0;
    }

    // ── iOS 风格底部标签栏 ──
    $tabBar = new BottomTabBar();
    $tabBar->append('消息', '✉', function() use ($win) {
        $win->setContent(buildScreen($win, 'messages'));
    });
    $tabBar->append('发现', '◎', function() use ($win) {
        $win->setContent(buildScreen($win, 'discover'));
    });
    $tabBar->append('我',   '☻', function() use ($win) {
        $win->setContent(buildScreen($win, 'profile'));
    });
    $tabBar->setSelected($selectedIdx);

    return Container::newBorder(null, $tabBar, null, null, $content);
}

$window->setContent(buildScreen($window, 'messages'));
$window->resize(new Size(420, 720));
$window->centerOnScreen();
$window->showAndRun();
