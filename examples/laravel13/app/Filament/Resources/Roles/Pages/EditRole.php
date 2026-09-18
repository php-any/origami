<?php

namespace App\Filament\Resources\Roles\Pages;

use App\Filament\Resources\Roles\RoleResource;
use App\Models\Role;
use Filament\Actions\DeleteAction;
use Filament\Resources\Pages\EditRecord;

class EditRole extends EditRecord
{
    protected static string $resource = RoleResource::class;

    protected function getHeaderActions(): array
    {
        return [
            DeleteAction::make()
                ->disabled(fn (): bool => $this->getRecord()->isSuperAdmin()),
        ];
    }

    protected function mutateFormDataBeforeSave(array $data): array
    {
        /** @var Role $record */
        $record = $this->getRecord();

        if ($record->isSuperAdmin()) {
            unset($data['name']);
        }

        return $data;
    }
}
