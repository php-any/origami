<?php

namespace App\Filament\Widgets;

use App\Filament\Resources\Orders\OrderResource;
use App\Models\Order;
use Filament\Actions\Action;
use Filament\Tables\Columns\TextColumn;
use Filament\Tables\Table;
use Filament\Widgets\TableWidget as BaseWidget;

class LatestOrdersWidget extends BaseWidget
{
    protected static ?int $sort = 2;

    protected int|string|array $columnSpan = 'full';

    protected static ?string $heading = '最新订单';

    public static function canView(): bool
    {
        return OrderResource::canViewAny();
    }

    public function table(Table $table): Table
    {
        return $table
            ->description('最近创建的 5 笔订单')
            ->headerActions([Action::make('allOrders')->label('查看全部')->url(OrderResource::getUrl())->link()])
            ->emptyStateHeading('暂无订单')
            ->emptyStateDescription('客户创建订单后，订单信息会显示在这里。')
            ->emptyStateIcon('heroicon-o-shopping-bag')
            ->query(
                Order::query()
                    ->with('user')
                    ->latest()
                    ->limit(5)
            )
            ->columns([
                TextColumn::make('order_no')
                    ->label('订单号'),
                TextColumn::make('user.name')
                    ->label('用户')
                    ->placeholder('—'),
                TextColumn::make('status_display')
                    ->label('状态')
                    ->badge()
                    ->color(fn (Order $record): string => match ($record->status) {
                        'pending' => 'warning',
                        'paid' => 'info',
                        'shipped' => 'primary',
                        'completed' => 'success',
                        default => 'gray',
                    }),
                TextColumn::make('total_amount')
                    ->label('金额')
                    ->money('CNY'),
                TextColumn::make('created_at')
                    ->label('创建时间')
                    ->dateTime('Y-m-d H:i'),
            ])
            ->recordActions([
                Action::make('view')
                    ->label('查看')
                    ->url(fn (Order $record): string => OrderResource::getUrl('view', ['record' => $record])),
            ])
            ->paginated(false);
    }
}
