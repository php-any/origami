<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Attributes\Fillable;

#[Fillable(['order_no', 'user_id', 'status', 'total_amount', 'remark', 'shipping_address', 'shipping_name', 'shipping_phone'])]
class Order extends \Illuminate\Database\Eloquent\Model
{
    protected $table = 'orders';

    const STATUS_PENDING = 'pending';
    const STATUS_PAID = 'paid';
    const STATUS_SHIPPED = 'shipped';
    const STATUS_COMPLETED = 'completed';
    const STATUS_CANCELLED = 'cancelled';

    protected function casts(): array
    {
        return [
            'total_amount' => 'decimal:2',
        ];
    }

    /**
     * The user that placed the order.
     */
    public function user()
    {
        return $this->belongsTo(User::class);
    }

    /**
     * The items in the order.
     */
    public function items()
    {
        return $this->hasMany(OrderItem::class);
    }

    /**
     * Get the order status display name.
     */
    public function getStatusDisplayAttribute()
    {
        $statusMap = [
            'pending' => '待付款',
            'paid' => '已付款',
            'shipped' => '已发货',
            'completed' => '已完成',
            'cancelled' => '已取消',
        ];
        return $statusMap[$this->status] ?? $this->status;
    }

    /**
     * Generate a unique order number.
     */
    public static function generateOrderNo(): string
    {
        return 'ORD' . date('YmdHis') . str_pad((string) mt_rand(1, 9999), 4, '0', STR_PAD_LEFT);
    }
}
