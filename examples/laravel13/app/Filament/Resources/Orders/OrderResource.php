<?php

namespace App\Filament\Resources\Orders;

use App\Filament\Resources\Orders\Pages\ListOrders;
use App\Filament\Resources\Orders\Pages\ViewOrder;
use App\Models\Order;
use Filament\Actions\ViewAction;
use Filament\Infolists\Components\RepeatableEntry;
use Filament\Infolists\Components\TextEntry;
use Filament\Resources\Resource;
use Filament\Schemas\Components\Section;
use Filament\Schemas\Schema;
use Filament\Tables\Columns\TextColumn;
use Filament\Tables\Filters\SelectFilter;
use Filament\Tables\Table;

class OrderResource extends Resource
{
    protected static ?string $model = Order::class;

    protected static string|\BackedEnum|null $navigationIcon = 'heroicon-o-shopping-bag';

    protected static string|\UnitEnum|null $navigationGroup = '订单管理';

    protected static ?string $navigationLabel = '订单';

    protected static ?string $modelLabel = '订单';

    protected static ?string $pluralModelLabel = '订单';

    protected static ?int $navigationSort = 1;

    protected static ?string $recordTitleAttribute = 'order_no';

    public static function form(Schema $schema): Schema
    {
        return $schema->components([]);
    }

    public static function infolist(Schema $schema): Schema
    {
        return $schema
            ->components([
                Section::make('订单信息')
                    ->columns(2)
                    ->schema([
                        TextEntry::make('order_no')
                            ->label('订单号')
                            ->copyable(),
                        TextEntry::make('status_display')
                            ->label('状态')
                            ->badge(),
                        TextEntry::make('user.name')
                            ->label('用户')
                            ->placeholder('—'),
                        TextEntry::make('total_amount')
                            ->label('总金额')
                            ->money('CNY'),
                        TextEntry::make('created_at')
                            ->label('下单时间')
                            ->dateTime('Y-m-d H:i:s'),
                        TextEntry::make('updated_at')
                            ->label('更新时间')
                            ->dateTime('Y-m-d H:i:s'),
                        TextEntry::make('remark')
                            ->label('备注')
                            ->placeholder('—')
                            ->columnSpanFull(),
                    ]),
                Section::make('收货信息')
                    ->columns(2)
                    ->schema([
                        TextEntry::make('shipping_name')
                            ->label('收货人')
                            ->placeholder('—'),
                        TextEntry::make('shipping_phone')
                            ->label('电话')
                            ->placeholder('—'),
                        TextEntry::make('shipping_address')
                            ->label('地址')
                            ->placeholder('—')
                            ->columnSpanFull(),
                    ]),
                Section::make('订单明细')
                    ->schema([
                        RepeatableEntry::make('items')
                            ->label('')
                            ->schema([
                                TextEntry::make('product_name')
                                    ->label('商品'),
                                TextEntry::make('price')
                                    ->label('单价')
                                    ->money('CNY'),
                                TextEntry::make('quantity')
                                    ->label('数量'),
                                TextEntry::make('subtotal')
                                    ->label('小计')
                                    ->money('CNY'),
                            ])
                            ->columns(4),
                    ]),
            ]);
    }

    public static function table(Table $table): Table
    {
        return $table
            ->columns([
                TextColumn::make('order_no')
                    ->label('订单号')
                    ->searchable()
                    ->sortable()
                    ->copyable(),
                TextColumn::make('user.name')
                    ->label('用户')
                    ->searchable()
                    ->placeholder('—'),
                TextColumn::make('status_display')
                    ->label('状态')
                    ->badge()
                    ->color(fn (Order $record): string => match ($record->status) {
                        Order::STATUS_PENDING => 'warning',
                        Order::STATUS_PAID => 'info',
                        Order::STATUS_SHIPPED => 'primary',
                        Order::STATUS_COMPLETED => 'success',
                        Order::STATUS_CANCELLED => 'danger',
                        default => 'gray',
                    }),
                TextColumn::make('total_amount')
                    ->label('金额')
                    ->money('CNY')
                    ->sortable(),
                TextColumn::make('created_at')
                    ->label('下单时间')
                    ->dateTime('Y-m-d H:i')
                    ->sortable(),
            ])
            ->filters([
                SelectFilter::make('status')
                    ->label('状态')
                    ->options([
                        Order::STATUS_PENDING => '待付款',
                        Order::STATUS_PAID => '已付款',
                        Order::STATUS_SHIPPED => '已发货',
                        Order::STATUS_COMPLETED => '已完成',
                        Order::STATUS_CANCELLED => '已取消',
                    ]),
            ])
            ->defaultSort('created_at', 'desc')
            ->recordActions([
                ViewAction::make(),
            ]);
    }

    public static function getPages(): array
    {
        return [
            'index' => ListOrders::route('/'),
            'view' => ViewOrder::route('/{record}'),
        ];
    }

    public static function canCreate(): bool
    {
        return false;
    }
}
