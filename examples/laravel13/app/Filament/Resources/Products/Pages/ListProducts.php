<?php

namespace App\Filament\Resources\Products\Pages;

use App\Filament\Resources\Products\ProductResource;
use Filament\Actions\CreateAction;
use Filament\Resources\Pages\ListRecords;
use Filament\Schemas\Components\Tabs\Tab;

class ListProducts extends ListRecords
{
    protected static string $resource = ProductResource::class;

    public function getSubheading(): ?string
    {
        return '管理商品资料、上下架状态与库存，及时处理库存预警。';
    }

    public function getTabs(): array
    {
        return [
            'all' => Tab::make('全部商品'),
            'active' => Tab::make('在售商品')->modifyQueryUsing(fn ($query) => $query->where('status', 'active')),
            'inactive' => Tab::make('已下架')->modifyQueryUsing(fn ($query) => $query->where('status', 'inactive')),
            'low_stock' => Tab::make('库存预警')->modifyQueryUsing(fn ($query) => $query->where('status', 'active')->where('stock', '<=', 10)),
        ];
    }

    protected function getHeaderActions(): array
    {
        return [
            CreateAction::make()->label('新增商品')->icon('heroicon-o-plus'),
        ];
    }
}
