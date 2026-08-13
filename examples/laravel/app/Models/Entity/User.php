<?php

namespace App\Models\Entity;

use Database\Annotation\Column;
use Database\Annotation\GeneratedValue;
use Database\Annotation\Id;
use Database\Annotation\Table;

/**
 * #[Table] Entity：供 Database\migrate 同步 schema（与 Eloquent Model 双轨）。
 */
#[Table('users')]
class User
{
    #[Id]
    #[GeneratedValue('AUTO')]
    #[Column('id', nullable: false)]
    public int $id;

    #[Column('name', nullable: false, length: 100)]
    public string $name;

    #[Column('email', nullable: false, length: 150)]
    public string $email;

    #[Column('password', nullable: false, length: 255)]
    public string $password;

    #[Column('created_at', nullable: true)]
    public ?string $created_at;
}
