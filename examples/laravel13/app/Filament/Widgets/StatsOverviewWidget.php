<?php

namespace App\Filament\Widgets;

use App\Models\Order;
use App\Models\Product;
use App\Models\User;
use Filament\Widgets\StatsOverviewWidget as BaseWidget;
use Filament\Widgets\StatsOverviewWidget\Stat;

class StatsOverviewWidget extends BaseWidget
{
    protected static ?int $sort = 1;

    protected function getStats(): array
    {
        $revenue = Order::query()
            ->where('status', '!=', Order::STATUS_CANCELLED)
            ->sum('total_amount');

        return [
            Stat::make('用户数', User::query()->count())
                ->description('前台用户')
                ->icon('heroicon-o-users'),
            Stat::make('订单数', Order::query()->count())
                ->description('全部订单')
                ->icon('heroicon-o-shopping-bag'),
            Stat::make('商品数', Product::query()->count())
                ->description('全部商品')
                ->icon('heroicon-o-cube'),
            Stat::make('营收', number_format((float) $revenue, 2))
                ->description('不含已取消')
                ->icon('heroicon-o-currency-yen'),
            Stat::make('待付款订单', Order::query()->where('status', Order::STATUS_PENDING)->count())
                ->description('待处理')
                ->icon('heroicon-o-clock')
                ->color('warning'),
        ];
    }
}
