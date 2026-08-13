<?php

namespace Spring\Model\Entity;

use Database\Annotation\Column;
use Database\Annotation\GeneratedValue;
use Database\Annotation\Id;
use Database\Annotation\Table;

/**
 * 用户实体，映射 users 表。
 */
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

    #[Column("age", nullable: false)]
    public int $age;

    #[Column("created_at", nullable: true)]
    public ?string $created_at;
}
