<?php

namespace App\Models;

use Database\Annotation\Column;
use Database\Annotation\GeneratedValue;
use Database\Annotation\Id;
use Database\Annotation\Table;

/**
 * 文章模型（类似 Laravel Eloquent Model）
 */
#[Table('posts')]
class Post
{
    #[Id]
    #[GeneratedValue('AUTO')]
    #[Column('id', nullable: false)]
    public int $id;

    #[Column('user_id', nullable: false)]
    public int $user_id;

    #[Column('title', nullable: false, length: 200)]
    public string $title;

    #[Column('body', nullable: false)]
    public string $body;

    #[Column('created_at', nullable: true)]
    public ?string $created_at;
}
