<?php

namespace Spring\Model;

class User {
    private ?int $id;
    private string $name;
    private string $email;
    private int $age;

    public function __construct(?int $id = null, string $name = '', string $email = '', int $age = 0) {
        $this->id = $id;
        $this->name = $name;
        $this->email = $email;
        $this->age = $age;
    }

    public function getId(): ?int {
        return $this->id;
    }

    public function setId(?int $id): void {
        $this->id = $id;
    }

    public function getName(): string {
        return $this->name;
    }

    public function setName(string $name): void {
        $this->name = $name;
    }

    public function getEmail(): string {
        return $this->email;
    }

    public function setEmail(string $email): void {
        $this->email = $email;
    }

    public function getAge(): int {
        return $this->age;
    }

    public function setAge(int $age): void {
        $this->age = $age;
    }

    /**
     * 转换为数组
     */
    public function toArray(): array {
        return [
            'id' => $this->id,
            'name' => $this->name,
            'email' => $this->email,
            'age' => $this->age
        ];
    }

    /**
     * JSON 序列化
     */
    public function toJson(): string {
        return json_encode($this->toArray());
    }
}
