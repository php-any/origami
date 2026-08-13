<?php
namespace MusicApp\Pages;

use Fyne\Container;
use Fyne\Widget\Label;
use Fyne\Widget\Button;
use Fyne\Widget\Entry;
use Fyne\Widget\Card;
use MusicApp\Services\Auth;

/**
 * 登录页面
 */
class LoginPage {
    static function show(\Fyne\Window $window): void {
        $window->setContent(self::build($window));
    }

    static function build(\Fyne\Window $window): Container {
        $titleLabel = new Label('🎵 Origami Music');
        $subtitleLabel = new Label('登录以继续');

        $usernameEntry = new Entry();
        $usernameEntry->setPlaceHolder('用户名');

        $passwordEntry = new Entry();
        $passwordEntry->setPlaceHolder('密码');

        $errorLabel = new Label('');
        $statusLabel = new Label('');

        $loginBtn = new Button('登  录', function()
            use ($usernameEntry, $passwordEntry, $errorLabel, $statusLabel, $window) {

            $u = $usernameEntry->getText();
            $p = $passwordEntry->getText();

            if ($u === '' || $p === '') {
                $errorLabel->setText('请输入用户名和密码');
                return;
            }

            $statusLabel->setText('登录中...');
            $auth = Auth::instance();

            if ($auth->tryLogin($u, $p)) {
                $errorLabel->setText('');
                $statusLabel->setText('登录成功！欢迎 ' . $auth->getUsername());
                HomePage::show($window);
            } else {
                $errorLabel->setText('用户名或密码错误');
                $statusLabel->setText('');
            }
        });

        $hintLabel = new Label('提示: admin/123456 或 demo/demo');

        $card = new Card('账号登录', '',
            Container::newVBox([
                $usernameEntry,
                $passwordEntry,
                $errorLabel,
                $statusLabel,
                $loginBtn,
            ])
        );

        return Container::newVBox([
            $titleLabel,
            $subtitleLabel,
            $card,
            $hintLabel,
        ]);
    }
}
