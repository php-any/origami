<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Attributes\Fillable;

#[Fillable(['name', 'sku', 'description', 'price', 'stock', 'status'])]
class Product extends \Illuminate\Database\Eloquent\Model
{
    protected $table = 'products';

    protected function casts(): array
    {
        return [
            'price' => 'decimal:2',
            'stock' => 'integer',
            'status' => 'string',
        ];
    }

    /**
     * The order items for the product.
     */
    public function orderItems()
    {
        return $this->hasMany(OrderItem::class);
    }

    /**
     * Scope for active products.
     */
    public function scopeActive($query)
    {
        return $query->where('status', 'active');
    }
}
