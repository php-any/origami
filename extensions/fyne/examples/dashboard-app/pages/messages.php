<?php
namespace App\Pages;

use Fyne\Container;
use Fyne\Widget\Label;
use Fyne\Widget\Card;

/**
 * 消息列表 — 卡片 + 红点徽章
 */
class MessagesPage {
    static function build(): Container {
        $items = [
            ['name' => '张伟',     'msg' => '好的，明天见！',          'time' => '10:32', 'badge' => '②'],
            ['name' => '设计组',   'msg' => '[图片]',                  'time' => '09:15', 'badge' => '⑤'],
            ['name' => '李工程师', 'msg' => '代码已提交，请 review',   'time' => '昨天',  'badge' => ''],
            ['name' => '王经理',   'msg' => '下午 3 点开会，别忘了',   'time' => '昨天',  'badge' => ''],
            ['name' => '技术支持', 'msg' => '问题已解决，请确认',       'time' => '周一',  'badge' => '①'],
            ['name' => '家人群',   'msg' => '妈妈: 周末回来吃饭吗？',   'time' => '周日',  'badge' => '99+'],
            ['name' => '项目组',   'msg' => '陈: 需求文档已更新',       'time' => '周六',  'badge' => ''],
            ['name' => '同学群',   'msg' => '下个月聚会，大家报名',     'time' => '周五',  'badge' => '③'],
        ];

        $cards = [];
        foreach ($items as $item) {
            $badge = $item['badge'] !== '' ? '  [' . $item['badge'] . ']' : '';
            $title = $item['name'] . '          ' . $item['time'] . $badge;
            $cards[] = new Card($title, $item['msg'], Container::newVBox([]));
        }

        return Container::newVBox($cards);
    }
}
