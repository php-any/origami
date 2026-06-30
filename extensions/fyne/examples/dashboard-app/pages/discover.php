<?php
namespace App\Pages;

use Fyne\Container;
use Fyne\Widget\Label;
use Fyne\Widget\Entry;
use Fyne\Widget\Card;

/**
 * 发现页 — 搜索栏 + 功能网格
 */
class DiscoverPage {
    static function build(): Container {
        $search = new Entry();
        $search->setPlaceHolder('搜索内容...');

        // 功能入口网格
        $features = ['📱 朋友圈', '🎬 视频号', '📦 小程序', '🛒 购物', '🎮 游戏', '📷 扫一扫'];
        $grid = [];
        foreach ($features as $f) {
            $grid[] = new Label($f);
        }

        return Container::newVBox([
            $search,
            Container::newVBox([]),
            Container::newGridWithColumns(3, $grid),
            Container::newVBox([]),
            new Card('📍 附近的人和动态', '', Container::newVBox([])),
            new Card('💰 支付与理财',     '', Container::newVBox([])),
            new Card('🎵 音乐与电台',     '', Container::newVBox([])),
        ]);
    }
}
