<?php

namespace Spring\DTO\Request;

use Validation\Annotation\Min;
use Validation\Annotation\Size;

class UpdateProductRequest
{
    #[Size(min: 1, max: 200)]
    public ?string $name = null;

    #[Min(value: 0, message: 'price 不能为负数')]
    public ?float $price = null;

    #[Size(max: 100)]
    public ?string $category = null;

    #[Size(max: 1000)]
    public ?string $description = null;

    public function toUpdateArray(): array
    {
        $data = [];
        if ($this->name !== null) {
            $data['name'] = $this->name;
        }
        if ($this->price !== null) {
            $data['price'] = $this->price;
        }
        if ($this->category !== null) {
            $data['category'] = $this->category;
        }
        if ($this->description !== null) {
            $data['description'] = $this->description;
        }
        return $data;
    }
}
