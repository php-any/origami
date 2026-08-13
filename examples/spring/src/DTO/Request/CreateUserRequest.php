<?php

namespace Spring\DTO\Request;

use Validation\Annotation\Email;
use Validation\Annotation\Min;
use Validation\Annotation\NotBlank;
use Validation\Annotation\Size;

class CreateUserRequest
{
    #[NotBlank(message: 'name 不能为空')]
    #[Size(min: 1, max: 100)]
    public string $name;

    #[NotBlank(message: 'email 不能为空')]
    #[Email]
    public string $email;

    #[Min(value: 0)]
    public int $age = 0;
}
