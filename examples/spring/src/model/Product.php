<?php

namespace Spring\Model;

class Product {
    private ?int $id;
    private string $name;
    private float $price;
    private string $category;
    private string $description;

    public function __construct(?int $id = null, string $name = '', float $price = 0.0, string $category = '', string $description = '') {
        $this->id = $id;
        $this->name = $name;
        $this->price = $price;
        $this->category = $category;
        $this->description = $description;
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

    public function getPrice(): float {
        return $this->price;
    }

    public function setPrice(float $price): void {
        $this->price = $price;
    }

    public function getCategory(): string {
        return $this->category;
    }

    public function setCategory(string $category): void {
        $this->category = $category;
    }

    public function getDescription(): string {
        return $this->description;
    }

    public function setDescription(string $description): void {
        $this->description = $description;
    }

    /**
     * 转换为数组
     */
    public function toArray(): array {
        return [
            'id' => $this->id,
            'name' => $this->name,
            'price' => $this->price,
            'category' => $this->category,
            'description' => $this->description
        ];
    }

    /**
     * JSON 序列化
     */
    public function toJson(): string {
        return json_encode($this->toArray());
    }

    /**
     * 获取格式化的价格
     */
    public function getFormattedPrice(): string {
        return "¥" . number_format($this->price, 2);
    }
}
