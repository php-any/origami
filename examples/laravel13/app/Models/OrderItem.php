<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Attributes\Fillable;

#[Fillable(['order_id', 'product_id', 'product_name', 'price', 'quantity', 'subtotal'])]
class OrderItem extends \Illuminate\Database\Eloquent\Model
{
    protected $table = 'order_items';

    protected function casts(): array
    {
        return [
            'price' => 'decimal:2',
            'quantity' => 'integer',
            'subtotal' => 'decimal:2',
        ];
    }

    /**
     * The order that the item belongs to.
     */
    public function order()
    {
        return $this->belongsTo(Order::class);
    }

    /**
     * The product for this item.
     */
    public function product()
    {
        return $this->belongsTo(Product::class);
    }
}
