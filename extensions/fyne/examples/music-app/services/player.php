<?php
namespace MusicApp\Services;

use Fyne\Widget\Label;
use Fyne\Widget\Button;
use Fyne\Widget\ProgressBar;
use Fyne\Container;

/**
 * 音乐播放器服务 — 模拟播放控制
 */
class Player {
    private $isPlaying = false;
    private $currentSong = null;
    private $progress = 0.0;

    // UI 组件引用
    public $nowPlayingLabel = null;
    public $artistLabel = null;
    public $progressBar = null;
    public $timeLabel = null;
    public $playBtn = null;

    static function instance(): Player {
        static $inst = null;
        if ($inst === null) {
            $inst = new Player();
        }
        return $inst;
    }

    function getIsPlaying(): bool { return $this->isPlaying; }
    function getCurrentSong() { return $this->currentSong; }
    function getProgress(): float { return $this->progress; }

    function play(array $song): void {
        $this->currentSong = $song;
        $this->progress = 0.0;
        $this->isPlaying = true;
        $this->updateUI();
    }

    function togglePlayPause(): void {
        if ($this->currentSong === null) return;
        $this->isPlaying = !$this->isPlaying;
        $this->updateUI();
    }

    function stop(): void {
        $this->isPlaying = false;
        $this->progress = 0.0;
        $this->updateUI();
    }

    function next(): void {
        $this->progress = 0.0;
        $this->updateUI();
    }

    function previous(): void {
        $this->progress = 0.0;
        $this->updateUI();
    }

    function seek(float $value): void {
        $this->progress = $value;
        $this->updateUI();
    }

    private function updateUI(): void {
        if ($this->nowPlayingLabel !== null && $this->currentSong !== null) {
            $this->nowPlayingLabel->setText($this->currentSong['title']);
        }
        if ($this->artistLabel !== null && $this->currentSong !== null) {
            $this->artistLabel->setText($this->currentSong['artist'] . ' — ' . $this->currentSong['album']);
        }
        if ($this->progressBar !== null) {
            $this->progressBar->setValue($this->progress);
        }
        if ($this->timeLabel !== null && $this->currentSong !== null) {
            $dur = MusicData::formatDuration($this->currentSong['duration']);
            $cur = MusicData::formatDuration((int)($this->currentSong['duration'] * $this->progress));
            $this->timeLabel->setText($cur . ' / ' . $dur);
        }
        if ($this->playBtn !== null) {
            $this->playBtn->setText($this->isPlaying ? '⏸' : '▶');
        }
    }

    /**
     * 创建播放器控制栏 UI
     */
    static function createControlBar(): Container {
        $player = self::instance();

        $player->nowPlayingLabel = new Label('未在播放');
        $player->artistLabel = new Label('选择一首歌曲开始播放');
        $player->progressBar = new ProgressBar();
        $player->timeLabel = new Label('00:00 / 00:00');

        $prevBtn = new Button('⏮', function() use ($player) { $player->previous(); });
        $player->playBtn = new Button('▶', function() use ($player) { $player->togglePlayPause(); });
        $nextBtn = new Button('⏭', function() use ($player) { $player->next(); });

        return Container::newVBox([
            $player->nowPlayingLabel,
            $player->artistLabel,
            $player->progressBar,
            $player->timeLabel,
            Container::newHBox([$prevBtn, $player->playBtn, $nextBtn]),
        ]);
    }
}
