<?php

namespace Spring\Model\Entity;

use Database\Annotation\Column;
use Database\Annotation\GeneratedValue;
use Database\Annotation\Id;
use Database\Annotation\Table;

/**
 * 订单实体，映射 orders 表。
 */
#[Table("orders")]
class OrderEntity {
    #[Id]
    #[GeneratedValue("AUTO")]
    #[Column("id", nullable: false)]
    public int $id;

    #[Column("user_id", nullable: false)]
    public int $user_id;

    #[Column("product_id", nullable: false)]
    public int $product_id;

    #[Column("quantity", nullable: false)]
    public int $quantity;

    #[Column("total_price", nullable: false)]
    public float $total_price;

    #[Column("status", nullable: false, length: 50)]
    public string $status;

    #[Column("created_at", nullable: true)]
    public ?string $created_at;
}
