<?php

namespace App\Notifications;

use App\Models\Product;
use Filament\Notifications\Notification as FilamentNotification;
use Illuminate\Bus\Queueable;
use Illuminate\Notifications\Notification;

class LowStockNotification extends Notification
{
    use Queueable;

    public function __construct(public Product $product) {}

    /**
     * @return list<string>
     */
    public function via(object $notifiable): array
    {
        return ['database'];
    }

    /**
     * @return array<string, mixed>
     */
    public function toDatabase(object $notifiable): array
    {
        return FilamentNotification::make()
            ->title('低库存警告')
            ->body("商品「{$this->product->name}」库存仅剩 {$this->product->stock} 件（SKU: {$this->product->sku}）")
            ->warning()
            ->getDatabaseMessage();
    }
}
