<?php

namespace Spring\Model\Entity;

use Database\Annotation\Table;

#[Table("products")]
class ProductEntity {
    public int $id;
    public string $name;
    public float $price;
    public string $category;
    public string $description;
    public ?string $created_at;

    public function toArray(): array {
        return [
            'id' => $this->id,
            'name' => $this->name,
            'price' => $this->price,
            'category' => $this->category,
            'description' => $this->description,
            'created_at' => $this->created_at,
        ];
    }
}
