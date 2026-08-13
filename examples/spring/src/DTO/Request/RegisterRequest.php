<?php

namespace Spring\DTO\Request;

use Validation\Annotation\Email;
use Validation\Annotation\NotBlank;
use Validation\Annotation\Pattern;
use Validation\Annotation\Size;

class RegisterRequest
{
    #[NotBlank(message: 'username 不能为空')]
    #[Size(min: 2, max: 100, message: 'Username must be between 2 and 100 characters')]
    #[Pattern('/^[a-zA-Z0-9_-]+$/')]
    public string $username;

    #[NotBlank(message: 'password 不能为空')]
    #[Size(min: 6, max: 100, message: 'Password must be between 6 and 100 characters')]
    public string $password;

    #[NotBlank(message: 'email 不能为空')]
    #[Email]
    public string $email;
}
