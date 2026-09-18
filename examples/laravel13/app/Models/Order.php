<?php

namespace App\Models;

use InvalidArgumentException;
use Illuminate\Database\Eloquent\Attributes\Fillable;
use Spatie\Activitylog\LogOptions;
use Spatie\Activitylog\Traits\LogsActivity;

#[Fillable(['order_no', 'user_id', 'status', 'total_amount', 'remark', 'shipping_address', 'shipping_name', 'shipping_phone'])]
class Order extends \Illuminate\Database\Eloquent\Model
{
    use LogsActivity;

    protected $table = 'orders';

    const STATUS_PENDING = 'pending';
    const STATUS_PAID = 'paid';
    const STATUS_SHIPPED = 'shipped';
    const STATUS_COMPLETED = 'completed';
    const STATUS_CANCELLED = 'cancelled';

    /** @var array<string, list<string>> */
    public const TRANSITIONS = [
        self::STATUS_PENDING => [self::STATUS_PAID, self::STATUS_CANCELLED],
        self::STATUS_PAID => [self::STATUS_SHIPPED, self::STATUS_CANCELLED],
        self::STATUS_SHIPPED => [self::STATUS_COMPLETED, self::STATUS_CANCELLED],
        self::STATUS_COMPLETED => [],
        self::STATUS_CANCELLED => [],
    ];

    protected function casts(): array
    {
        return [
            'total_amount' => 'decimal:2',
        ];
    }

    public function getActivitylogOptions(): LogOptions
    {
        return LogOptions::defaults()
            ->logFillable()
            ->logOnlyDirty()
            ->dontSubmitEmptyLogs();
    }

    public function user()
    {
        return $this->belongsTo(User::class);
    }

    public function items()
    {
        return $this->hasMany(OrderItem::class);
    }

    public function getStatusDisplayAttribute(): string
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

    public function allowedTransitions(): array
    {
        return self::TRANSITIONS[$this->status] ?? [];
    }

    public function canTransitionTo(string $status): bool
    {
        return in_array($status, $this->allowedTransitions(), true);
    }

    public function transitionTo(string $status): void
    {
        if (! $this->canTransitionTo($status)) {
            throw new InvalidArgumentException("Cannot transition order from {$this->status} to {$status}");
        }

        $this->status = $status;
        $this->save();
    }

    public static function generateOrderNo(): string
    {
        $prefix = Setting::getValue('order_prefix', 'ORD');

        return $prefix.date('YmdHis').str_pad((string) mt_rand(1, 9999), 4, '0', STR_PAD_LEFT);
    }
}
