<?php

namespace App\Filament\Resources\Orders\Pages;

use App\Filament\Resources\Orders\OrderResource;
use Filament\Schemas\Components\Tabs\Tab;
use Filament\Resources\Pages\ListRecords;

class ListOrders extends ListRecords
{
    protected static string $resource = OrderResource::class;

    public function getSubheading(): ?string
    {
        return '按状态处理订单，查看收货信息与订单明细。';
    }

    public function getTabs(): array
    {
        return [
            'all' => Tab::make('全部订单'),
            'pending' => Tab::make('待付款')->modifyQueryUsing(fn ($query) => $query->where('status', 'pending')),
            'paid' => Tab::make('待发货')->modifyQueryUsing(fn ($query) => $query->where('status', 'paid')),
            'shipped' => Tab::make('待收货')->modifyQueryUsing(fn ($query) => $query->where('status', 'shipped')),
            'completed' => Tab::make('已完成')->modifyQueryUsing(fn ($query) => $query->where('status', 'completed')),
            'cancelled' => Tab::make('已取消')->modifyQueryUsing(fn ($query) => $query->where('status', 'cancelled')),
        ];
    }
}
