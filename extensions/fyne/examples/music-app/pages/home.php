<?php
namespace MusicApp\Pages;

use Fyne\Container;
use Fyne\Widget\Label;
use Fyne\Widget\Button;
use Fyne\Widget\Entry;
use Fyne\Widget\Card;
use MusicApp\Services\Auth;
use MusicApp\Services\MusicData;

/**
 * 首页 — 歌单列表 + 搜索
 */
class HomePage {
    static function show(\Fyne\Window $window): void {
        $window->setContent(self::build($window));
    }

    static function build(\Fyne\Window $window): Container {
        $auth = Auth::instance();
        $musicData = MusicData::instance();

        // 顶部栏
        $userLabel = new Label('👤 ' . $auth->getUsername());
        $logoutBtn = new Button('退出', function() use ($window, $auth) {
            $auth->logout();
            LoginPage::show($window);
        });
        $topBar = Container::newHBox([$userLabel, $logoutBtn]);

        // 搜索框
        $searchEntry = new Entry();
        $searchEntry->setPlaceHolder('搜索歌曲或歌手...');
        $searchResultLabel = new Label('');

        $searchBtn = new Button('🔍 搜索', function()
            use ($searchEntry, $searchResultLabel, $musicData) {
            $kw = $searchEntry->getText();
            if ($kw === '') {
                $searchResultLabel->setText('请输入搜索关键词');
                return;
            }
            $results = $musicData->search($kw);
            $cnt = count($results);
            if ($cnt === 0) {
                $searchResultLabel->setText('未找到与 "' . $kw . '" 相关的结果');
            } else {
                $searchResultLabel->setText('找到 ' . $cnt . ' 首歌曲');
            }
        });

        $searchRow = Container::newHBox([$searchEntry, $searchBtn]);

        // 歌单列表
        $playlists = $musicData->getPlaylists();
        $cards = [];
        foreach ($playlists as $pl) {
            $info = $pl['cover'] . ' ' . $pl['name'] . ' (' . $pl['count'] . '首)';
            $openBtn = new Button('打开', function() use ($pl, $window) {
                PlaylistDetailPage::show($window, $pl);
            });
            $cards[] = new Card($info, $pl['desc'], Container::newVBox([$openBtn]));
        }

        return Container::newVBox(array_merge(
            [$topBar, new Label(''), $searchRow, $searchResultLabel, new Label('📋 推荐歌单')],
            $cards
        ));
    }
}
