<?php

namespace Spring\DTO\Request;

use Validation\Annotation\Min;
use Validation\Annotation\NotBlank;
use Validation\Annotation\Size;

class CreateProductRequest
{
    #[NotBlank(message: 'name 不能为空')]
    #[Size(min: 1, max: 200)]
    public string $name;

    #[NotBlank(message: 'price 不能为空')]
    #[Min(value: 0, message: 'price 不能为负数')]
    public float $price;

    #[Size(max: 100)]
    public string $category = '未分类';

    #[Size(max: 1000)]
    public string $description = '';
}
