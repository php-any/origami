<?php
namespace MusicApp\Services;

/**
 * 音乐数据服务 — 提供模拟的歌单和歌曲数据
 */
class MusicData {
    private $playlists = null;
    private $songs = null;

    function __construct() {
        $this->playlists = null;
        $this->songs = null;
    }

    /**
     * 获取实例
     */
    static function instance(): MusicData {
        static $inst = null;
        if ($inst === null) {
            $inst = new MusicData();
        }
        return $inst;
    }

    /**
     * 获取所有歌单
     */
    function getPlaylists(): array {
        if ($this->playlists !== null) {
            return $this->playlists;
        }
        $this->playlists = [
            ['id' => 1, 'name' => '华语流行精选', 'cover' => '🎵', 'count' => 42, 'desc' => '最受欢迎的华语歌曲'],
            ['id' => 2, 'name' => '欧美经典摇滚', 'cover' => '🎸', 'count' => 38, 'desc' => '经典摇滚永不褪色'],
            ['id' => 3, 'name' => '轻音乐 & 纯音乐', 'cover' => '🎹', 'count' => 56, 'desc' => '放松心情的纯音乐'],
            ['id' => 4, 'name' => '电子 & 舞曲',     'cover' => '🎧', 'count' => 64, 'desc' => '燃爆全场的电子音乐'],
            ['id' => 5, 'name' => '怀旧金曲',       'cover' => '📻', 'count' => 30, 'desc' => '那些年我们一起听过的歌'],
            ['id' => 6, 'name' => '日语动漫歌曲',   'cover' => '🎌', 'count' => 48, 'desc' => '经典动漫主题曲合集'],
        ];
        return $this->playlists;
    }

    /**
     * 获取指定歌单的歌曲列表
     */
    function getSongs(int $playlistId): array {
        $all = $this->getAllSongs();
        $result = [];
        foreach ($all as $song) {
            if ($song['playlist_id'] === $playlistId) {
                $result[] = $song;
            }
        }
        return $result;
    }

    /**
     * 获取所有歌曲
     */
    private function getAllSongs(): array {
        if ($this->songs !== null) {
            return $this->songs;
        }
        $this->songs = [
            // 华语流行精选 (playlist 1)
            ['id' => 1,  'playlist_id' => 1, 'title' => '七里香',              'artist' => '周杰伦',       'album' => '七里香',                     'duration' => 299],
            ['id' => 2,  'playlist_id' => 1, 'title' => '晴天',                'artist' => '周杰伦',       'album' => '叶惠美',                     'duration' => 269],
            ['id' => 3,  'playlist_id' => 1, 'title' => '起风了',              'artist' => '买辣椒也用券', 'album' => '起风了',                     'duration' => 325],
            ['id' => 4,  'playlist_id' => 1, 'title' => '光年之外',            'artist' => '邓紫棋',       'album' => '光年之外',                   'duration' => 235],
            ['id' => 5,  'playlist_id' => 1, 'title' => '平凡之路',            'artist' => '朴树',         'album' => '平凡之路',                   'duration' => 312],
            // 欧美经典摇滚 (playlist 2)
            ['id' => 6,  'playlist_id' => 2, 'title' => 'Bohemian Rhapsody',   'artist' => 'Queen',                  'album' => 'A Night at the Opera',       'duration' => 355],
            ['id' => 7,  'playlist_id' => 2, 'title' => 'Hotel California',    'artist' => 'Eagles',                 'album' => 'Hotel California',           'duration' => 391],
            ['id' => 8,  'playlist_id' => 2, 'title' => 'Stairway to Heaven',  'artist' => 'Led Zeppelin',           'album' => 'Led Zeppelin IV',            'duration' => 482],
            ['id' => 9,  'playlist_id' => 2, 'title' => 'Sweet Child O\' Mine','artist' => 'Guns N\' Roses',         'album' => 'Appetite for Destruction',   'duration' => 356],
            // 轻音乐 & 纯音乐 (playlist 3)
            ['id' => 10, 'playlist_id' => 3, 'title' => 'River Flows in You',  'artist' => 'Yiruma',                 'album' => 'First Love',                 'duration' => 187],
            ['id' => 11, 'playlist_id' => 3, 'title' => 'Kiss The Rain',       'artist' => 'Yiruma',                 'album' => 'From The Yellow Room',       'duration' => 264],
            ['id' => 12, 'playlist_id' => 3, 'title' => 'Canon in D',          'artist' => 'Johann Pachelbel',       'album' => 'Classical Best',             'duration' => 328],
            // 电子 & 舞曲 (playlist 4)
            ['id' => 13, 'playlist_id' => 4, 'title' => 'Faded',               'artist' => 'Alan Walker',            'album' => 'Different World',            'duration' => 212],
            ['id' => 14, 'playlist_id' => 4, 'title' => 'Wake Me Up',          'artist' => 'Avicii',                 'album' => 'True',                       'duration' => 249],
            ['id' => 15, 'playlist_id' => 4, 'title' => 'Don\'t Let Me Down',  'artist' => 'The Chainsmokers',       'album' => 'Collage',                    'duration' => 208],
            // 怀旧金曲 (playlist 5)
            ['id' => 16, 'playlist_id' => 5, 'title' => '海阔天空',            'artist' => 'Beyond',       'album' => '乐与怒',     'duration' => 326],
            ['id' => 17, 'playlist_id' => 5, 'title' => '同桌的你',            'artist' => '老狼',         'album' => '同桌的你',   'duration' => 285],
            ['id' => 18, 'playlist_id' => 5, 'title' => '朋友',                'artist' => '周华健',       'album' => '朋友',       'duration' => 278],
            // 日语动漫歌曲 (playlist 6)
            ['id' => 19, 'playlist_id' => 6, 'title' => '残酷な天使のテーゼ',  'artist' => '高橋洋子',     'album' => 'NEON GENESIS EVANGELION', 'duration' => 280],
            ['id' => 20, 'playlist_id' => 6, 'title' => '紅蓮華',              'artist' => 'LiSA',         'album' => '紅蓮華',       'duration' => 234],
            ['id' => 21, 'playlist_id' => 6, 'title' => 'Butter-Fly',          'artist' => '和田光司',     'album' => 'Butter-Fly',   'duration' => 255],
        ];
        return $this->songs;
    }

    /**
     * 搜索歌曲
     */
    function search(string $keyword): array {
        if ($keyword === '') {
            return [];
        }
        $results = [];
        foreach ($this->getAllSongs() as $song) {
            if (stripos($song['title'], $keyword) !== false ||
                stripos($song['artist'], $keyword) !== false) {
                $results[] = $song;
            }
        }
        return $results;
    }

    /**
     * 格式化时长为 mm:ss
     */
    static function formatDuration(int $seconds): string {
        $m = intdiv($seconds, 60);
        $s = $seconds % 60;
        return sprintf("%d:%02d", $m, $s);
    }
}
