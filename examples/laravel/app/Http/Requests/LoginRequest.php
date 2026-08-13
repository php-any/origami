<?php

namespace App\Http\Requests;

use Validation\Annotation\NotBlank;
use Validation\Annotation\Size;

/**
 * 登录请求 DTO（类似 Laravel FormRequest）
 */
class LoginRequest
{
    #[NotBlank(message: 'email 不能为空')]
    #[Size(min: 3, max: 150)]
    public string $email = '';

    #[NotBlank(message: 'password 不能为空')]
    #[Size(min: 6, max: 64)]
    public string $password = '';
}
