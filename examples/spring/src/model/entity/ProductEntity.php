<?php

namespace Spring\Model\Entity;

use Database\Annotation\Column;
use Database\Annotation\GeneratedValue;
use Database\Annotation\Id;
use Database\Annotation\Table;

/**
 * 商品实体，映射 products 表。
 */
#[Table("products")]
class ProductEntity {
    #[Id]
    #[GeneratedValue("AUTO")]
    #[Column("id", nullable: false)]
    public int $id;

    #[Column("name", nullable: false, length: 200)]
    public string $name;

    #[Column("price", nullable: false)]
    public float $price;

    #[Column("category", nullable: false, length: 100)]
    public string $category;

    #[Column("description", nullable: true)]
    public ?string $description;

    #[Column("created_at", nullable: true)]
    public ?string $created_at;
}
