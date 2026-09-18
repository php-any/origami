<?php

namespace App\Filament\Resources\Orders\Pages\Concerns;

use App\Models\Order;
use Filament\Actions\Action;
use Filament\Notifications\Notification;

trait HasOrderTransitionActions
{
    /**
     * @return list<Action>
     */
    protected function getOrderTransitionActions(): array
    {
        return [
            Action::make('markPaid')
                ->label('标记已付款')
                ->icon('heroicon-o-banknotes')
                ->color('success')
                ->visible(fn (): bool => $this->getRecord()->canTransitionTo(Order::STATUS_PAID))
                ->requiresConfirmation()
                ->action(fn () => $this->transitionOrder(Order::STATUS_PAID, '订单已标记为已付款')),
            Action::make('ship')
                ->label('发货')
                ->icon('heroicon-o-truck')
                ->color('info')
                ->visible(fn (): bool => $this->getRecord()->canTransitionTo(Order::STATUS_SHIPPED))
                ->requiresConfirmation()
                ->action(fn () => $this->transitionOrder(Order::STATUS_SHIPPED, '订单已发货')),
            Action::make('complete')
                ->label('完成')
                ->icon('heroicon-o-check-badge')
                ->color('success')
                ->visible(fn (): bool => $this->getRecord()->canTransitionTo(Order::STATUS_COMPLETED))
                ->requiresConfirmation()
                ->action(fn () => $this->transitionOrder(Order::STATUS_COMPLETED, '订单已完成')),
            Action::make('cancel')
                ->label('取消订单')
                ->icon('heroicon-o-x-circle')
                ->color('danger')
                ->visible(fn (): bool => $this->getRecord()->canTransitionTo(Order::STATUS_CANCELLED))
                ->requiresConfirmation()
                ->action(fn () => $this->transitionOrder(Order::STATUS_CANCELLED, '订单已取消')),
        ];
    }

    protected function transitionOrder(string $status, string $message): void
    {
        /** @var Order $record */
        $record = $this->getRecord();

        if (! $record->canTransitionTo($status)) {
            Notification::make()
                ->danger()
                ->title('无法变更状态')
                ->body("当前状态不允许流转到 {$status}")
                ->send();

            return;
        }

        $record->transitionTo($status);
        $this->record->refresh();

        Notification::make()
            ->success()
            ->title($message)
            ->send();
    }
}
