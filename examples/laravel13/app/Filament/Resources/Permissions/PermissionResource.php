<?php

namespace App\Filament\Resources\Permissions;

use App\Filament\Resources\Permissions\Pages\CreatePermission;
use App\Filament\Resources\Permissions\Pages\EditPermission;
use App\Filament\Resources\Permissions\Pages\ListPermissions;
use App\Models\Permission;
use Filament\Actions\BulkActionGroup;
use Filament\Actions\DeleteAction;
use Filament\Actions\DeleteBulkAction;
use Filament\Actions\EditAction;
use Filament\Forms\Components\Textarea;
use Filament\Forms\Components\TextInput;
use Filament\Resources\Resource;
use Filament\Schemas\Schema;
use Filament\Tables\Columns\TextColumn;
use Filament\Tables\Filters\SelectFilter;
use Filament\Tables\Table;

class PermissionResource extends Resource
{
    protected static ?string $model = Permission::class;

    protected static string|\BackedEnum|null $navigationIcon = 'heroicon-o-key';

    protected static string|\UnitEnum|null $navigationGroup = '权限管理';

    protected static ?string $navigationLabel = '权限';

    protected static ?string $modelLabel = '权限';

    protected static ?string $pluralModelLabel = '权限';

    protected static ?int $navigationSort = 2;

    protected static ?string $recordTitleAttribute = 'display_name';

    public static function form(Schema $schema): Schema
    {
        return $schema
            ->components([
                TextInput::make('name')
                    ->label('标识')
                    ->required()
                    ->unique(ignoreRecord: true)
                    ->maxLength(255)
                    ->helperText('例如 products.view'),
                TextInput::make('display_name')
                    ->label('显示名称')
                    ->required()
                    ->maxLength(255),
                TextInput::make('group')
                    ->label('分组')
                    ->required()
                    ->maxLength(255),
                Textarea::make('description')
                    ->label('描述')
                    ->rows(3)
                    ->columnSpanFull(),
            ]);
    }

    public static function table(Table $table): Table
    {
        return $table
            ->columns([
                TextColumn::make('group')
                    ->label('分组')
                    ->badge()
                    ->sortable()
                    ->searchable(),
                TextColumn::make('name')
                    ->label('标识')
                    ->searchable()
                    ->sortable()
                    ->copyable(),
                TextColumn::make('display_name')
                    ->label('显示名称')
                    ->searchable()
                    ->sortable(),
                TextColumn::make('roles_count')
                    ->label('角色数')
                    ->counts('roles'),
                TextColumn::make('updated_at')
                    ->label('更新时间')
                    ->dateTime('Y-m-d H:i')
                    ->sortable()
                    ->toggleable(isToggledHiddenByDefault: true),
            ])
            ->filters([
                SelectFilter::make('group')
                    ->label('分组')
                    ->options(fn (): array => Permission::query()
                        ->whereNotNull('group')
                        ->distinct()
                        ->orderBy('group')
                        ->pluck('group', 'group')
                        ->all()),
            ])
            ->recordActions([
                EditAction::make(),
                DeleteAction::make(),
            ])
            ->toolbarActions([
                BulkActionGroup::make([
                    DeleteBulkAction::make(),
                ]),
            ])
            ->defaultSort('group');
    }

    public static function getPages(): array
    {
        return [
            'index' => ListPermissions::route('/'),
            'create' => CreatePermission::route('/create'),
            'edit' => EditPermission::route('/{record}/edit'),
        ];
    }
}
