<?php

namespace App\Filament\Pages\Auth;

use Filament\Auth\Pages\Login as BaseLogin;
use Illuminate\Contracts\Support\Htmlable;
use Illuminate\Support\HtmlString;

class Login extends BaseLogin
{
    public function getHeading(): string | Htmlable
    {
        return '欢迎回来';
    }

    public function getSubheading(): string | Htmlable | null
    {
        return new HtmlString(
            '登录你的账号，开始今天的管理工作。'
        );
    }
}
