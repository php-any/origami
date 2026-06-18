<?php

namespace Spring\Model\Entity;

use Database\Annotation\Table;

#[Table("orders")]
class OrderEntity {
    public int $id;
    public int $user_id;
    public int $product_id;
    public int $quantity;
    public float $total_price;
    public string $status;
    public ?string $created_at;

    public function toArray(): array {
        return [
            'id' => $this->id,
            'user_id' => $this->user_id,
            'product_id' => $this->product_id,
            'quantity' => $this->quantity,
            'total_price' => $this->total_price,
            'status' => $this->status,
            'created_at' => $this->created_at,
        ];
    }
}
