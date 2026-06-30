<?php
namespace App\Pages;

use Fyne\Container;
use Fyne\Widget\Label;
use Fyne\Widget\Card;
use Fyne\Widget\Form;
use Fyne\Widget\Entry;
use Fyne\Widget\Check;
use Fyne\Widget\Button;

/**
 * 个人页 — 头像区 + 表单设置 + 菜单
 */
class ProfilePage {
    static function build(): Container {
        // 头像信息
        $avatar = new Card('admin', 'ID: origami_2024', Container::newVBox([]));

        // 设置表单
        $form = new Form();
        $form->append('用户名',   new Entry());
        $form->append('邮箱',     new Entry());
        $form->append('深色模式', new Check(''));
        $form->append('推送通知', new Check(''));

        // 菜单列表
        $menus = [
            new Card('📷 相册',   '', Container::newVBox([])),
            new Card('📦 收藏',   '',  Container::newVBox([])),
            new Card('💳 钱包',   '余额 ¥128.00',  Container::newVBox([])),
            new Card('⚙ 设置',    '账号与安全',    Container::newVBox([])),
            new Card('ℹ 关于',    'v1.0.0',        Container::newVBox([])),
        ];

        return Container::newVBox(array_merge(
            [$avatar, Container::newVBox([]), $form, Container::newVBox([]), new Button('💾 保存', function() {})],
            [Container::newVBox([])],
            $menus
        ));
    }
}
