<?php

namespace App\Model\Entity;

use Database\Annotation\Column;
use Database\Annotation\GeneratedValue;
use Database\Annotation\Id;
use Database\Annotation\Table;

#[Table("users")]
class UserEntity {
    #[Id]
    #[GeneratedValue("AUTO")]
    #[Column("id", nullable: false)]
    public int $id;

    #[Column("name", nullable: false, length: 100)]
    public string $name;

    #[Column("email", nullable: false, length: 100)]
    public string $email;
}

#[Table("posts")]
class PostEntity {
    #[Id]
    #[GeneratedValue("AUTO")]
    #[Column("id", nullable: false)]
    public int $id;

    #[Column("user_id", nullable: false)]
    public int $user_id;

    #[Column("title", nullable: false, length: 200)]
    public string $title;

    #[Column("content", nullable: true)]
    public ?string $content;
}
