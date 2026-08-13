<?php
namespace MusicApp\Pages;

use Fyne\Container;
use Fyne\Widget\Label;
use Fyne\Widget\Button;
use Fyne\Widget\Card;
use MusicApp\Services\MusicData;
use MusicApp\Services\Player;

/**
 * 歌单详情页 — 歌曲列表 + 播放控制
 */
class PlaylistDetailPage {
    static function show(\Fyne\Window $window, array $playlist): void {
        $window->setContent(self::build($window, $playlist));
    }

    static function build(\Fyne\Window $window, array $playlist): Container {
        $musicData = MusicData::instance();

        $headerLabel = new Label($playlist['cover'] . ' ' . $playlist['name']);
        $descLabel = new Label($playlist['desc'] . ' — ' . $playlist['count'] . ' 首');

        $backBtn = new Button('← 返回', function() use ($window) {
            HomePage::show($window);
        });

        $playAllBtn = new Button('▶ 播放全部', function() use ($playlist, $musicData) {
            $songs = $musicData->getSongs($playlist['id']);
            if (count($songs) > 0) {
                Player::instance()->play($songs[0]);
            }
        });

        $topBar = Container::newHBox([$backBtn, $playAllBtn]);

        // 歌曲列表
        $songs = $musicData->getSongs($playlist['id']);
        $songCards = [];
        foreach ($songs as $i => $song) {
            $idx = $i + 1;
            $title = '#' . $idx . ' ' . $song['title'];
            $subtitle = $song['artist'] . ' · ' . $song['album'] . ' · '
                      . MusicData::formatDuration($song['duration']);

            $playBtn = new Button('▶', function() use ($song) {
                Player::instance()->play($song);
            });

            $songCards[] = new Card($title, $subtitle, Container::newHBox([$playBtn]));
        }

        // 播放控制栏
        $controlBar = Player::createControlBar();

        $items = array_merge(
            [$headerLabel, $descLabel, $topBar, new Label('')],
            $songCards,
            [new Label(''), $controlBar]
        );

        return Container::newVBox($items);
    }
}
