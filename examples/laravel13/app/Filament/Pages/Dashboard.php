<?php

namespace App\Filament\Pages;

use App\Filament\Widgets\LatestOrdersWidget;
use App\Filament\Widgets\StatsOverviewWidget;
use Filament\Pages\Dashboard as BaseDashboard;

class Dashboard extends BaseDashboard
{
    protected static string|\BackedEnum|null $navigationIcon = 'heroicon-o-home';

    protected static ?string $navigationLabel = '仪表盘';

    protected static ?string $title = '仪表盘';

    public function getWidgets(): array
    {
        return [
            StatsOverviewWidget::class,
            LatestOrdersWidget::class,
        ];
    }
}
