<?php

namespace Spring\DTO\Request;

use Validation\Annotation\Pattern;
use Validation\Annotation\Size;

class LoginRequest
{
    #[Size(min: 2, max: 100, message: 'Username must be between 2 and 100 characters')]
    #[Pattern('/^[a-zA-Z0-9_-]+$/')]
    public string $username;
    #[Size(min: 6, max: 100, message: 'Password must be between 6 and 100 characters')]
    public string $password;
}
