<?php

namespace App\Models;

use Database\Annotation\Column;
use Database\Annotation\GeneratedValue;
use Database\Annotation\Id;
use Database\Annotation\Table;

/**
 * API 访问令牌（类似 Laravel personal_access_tokens）
 */
#[Table('api_tokens')]
class ApiToken
{
    #[Id]
    #[GeneratedValue('AUTO')]
    #[Column('id', nullable: false)]
    public int $id;

    #[Column('user_id', nullable: false)]
    public int $user_id;

    #[Column('token', nullable: false, length: 64)]
    public string $token;

    #[Column('created_at', nullable: true)]
    public ?string $created_at;
}
