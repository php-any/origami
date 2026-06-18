<?php

namespace Spring\Model\Entity;

use Database\Annotation\Table;

#[Table("users")]
class UserEntity {
    public int $id;
    public string $name;
    public string $email;
    public int $age;
    public ?string $created_at;

    public function toArray(): array {
        return [
            'id' => $this->id,
            'name' => $this->name,
            'email' => $this->email,
            'age' => $this->age,
            'created_at' => $this->created_at,
        ];
    }
}
