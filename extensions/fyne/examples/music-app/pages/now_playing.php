<?php
namespace MusicApp\Pages;

use Fyne\Container;
use Fyne\Widget\Label;
use Fyne\Widget\Button;
use Fyne\Widget\ProgressBar;
use Fyne\Widget\Card;
use MusicApp\Services\Player;
use MusicApp\Services\MusicData;

/**
 * 正在播放页面
 */
class NowPlayingPage {
    static function show(\Fyne\Window $window): void {
        $window->setContent(self::build($window));
    }

    static function build(\Fyne\Window $window): Container {
        $player = Player::instance();
        $song = $player->getCurrentSong();

        if ($song === null) {
            $backBtn = new Button('← 返回歌单', function() use ($window) {
                HomePage::show($window);
            });
            return Container::newVBox([
                new Label('🎵'),
                new Label('没有正在播放的歌曲'),
                new Label('从歌单中选择一首歌曲开始播放'),
                $backBtn,
            ]);
        }

        $titleLabel = new Label($song['title']);
        $artistLabel = new Label($song['artist']);
        $albumLabel = new Label('专辑: ' . $song['album']);
        $durationLabel = new Label('时长: ' . MusicData::formatDuration($song['duration']));

        $progress = new ProgressBar();
        $progress->setValue($player->getProgress());
        $player->progressBar = $progress;

        $timeLabel = new Label('00:00 / ' . MusicData::formatDuration($song['duration']));
        $player->timeLabel = $timeLabel;

        $prevBtn = new Button('⏮', function() use ($player) { $player->previous(); });
        $playBtn = new Button($player->getIsPlaying() ? '⏸' : '▶', function() use ($player) {
            $player->togglePlayPause();
        });
        $nextBtn = new Button('⏭', function() use ($player) { $player->next(); });
        $player->playBtn = $playBtn;

        $backBtn = new Button('← 返回首页', function() use ($window) {
            HomePage::show($window);
        });

        $infoCard = new Card('正在播放', '',
            Container::newVBox([
                new Label('🎵'),
                $titleLabel,
                $artistLabel,
                $albumLabel,
                $durationLabel,
            ])
        );

        return Container::newVBox([
            $backBtn,
            $infoCard,
            $progress,
            $timeLabel,
            Container::newHBox([$prevBtn, $playBtn, $nextBtn]),
        ]);
    }
}
