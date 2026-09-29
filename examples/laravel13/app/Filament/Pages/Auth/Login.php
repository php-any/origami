<?php

namespace App\Filament\Pages\Auth;

use Filament\Auth\Pages\Login as BaseLogin;
use Illuminate\Contracts\Support\Htmlable;
use Illuminate\Support\HtmlString;

class Login extends BaseLogin
{
    public function getSubheading(): string | Htmlable | null
    {
        return new HtmlString(
            '默认账号：<code>admin@example.com</code> / <code>password</code><br>'.
            '普通管理员：<code>manager@example.com</code> / <code>password</code>'
        );
    }
}
