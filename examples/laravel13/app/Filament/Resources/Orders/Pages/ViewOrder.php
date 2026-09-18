<?php

namespace App\Filament\Resources\Orders\Pages;

use App\Filament\Resources\Orders\OrderResource;
use App\Filament\Resources\Orders\Pages\Concerns\HasOrderTransitionActions;
use Filament\Resources\Pages\ViewRecord;

class ViewOrder extends ViewRecord
{
    use HasOrderTransitionActions;

    protected static string $resource = OrderResource::class;

    protected function getHeaderActions(): array
    {
        return $this->getOrderTransitionActions();
    }
}
