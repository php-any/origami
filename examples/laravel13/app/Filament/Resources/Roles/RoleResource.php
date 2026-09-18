<?php

namespace App\Filament\Resources\Roles;

use App\Filament\Resources\Roles\Pages\CreateRole;
use App\Filament\Resources\Roles\Pages\EditRole;
use App\Filament\Resources\Roles\Pages\ListRoles;
use App\Models\Role;
use Filament\Actions\BulkActionGroup;
use Filament\Actions\DeleteAction;
use Filament\Actions\DeleteBulkAction;
use Filament\Actions\EditAction;
use Filament\Forms\Components\CheckboxList;
use Filament\Forms\Components\Textarea;
use Filament\Forms\Components\TextInput;
use Filament\Resources\Resource;
use Filament\Schemas\Schema;
use Filament\Tables\Columns\TextColumn;
use Filament\Tables\Table;

class RoleResource extends Resource
{
    protected static ?string $model = Role::class;

    protected static string|\BackedEnum|null $navigationIcon = 'heroicon-o-user-group';

    protected static string|\UnitEnum|null $navigationGroup = '权限管理';

    protected static ?string $navigationLabel = '角色';

    protected static ?string $modelLabel = '角色';

    protected static ?string $pluralModelLabel = '角色';

    protected static ?int $navigationSort = 1;

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
                    ->disabled(fn (?Role $record): bool => $record?->isSuperAdmin() ?? false)
                    ->dehydrated(fn (?Role $record): bool => ! ($record?->isSuperAdmin() ?? false)),
                TextInput::make('display_name')
                    ->label('显示名称')
                    ->required()
                    ->maxLength(255),
                Textarea::make('description')
                    ->label('描述')
                    ->rows(3)
                    ->columnSpanFull(),
                CheckboxList::make('permissions')
                    ->label('权限')
                    ->relationship('permissions', 'display_name')
                    ->columns(2)
                    ->searchable()
                    ->bulkToggleable()
                    ->columnSpanFull(),
            ]);
    }

    public static function table(Table $table): Table
    {
        return $table
            ->columns([
                TextColumn::make('name')
                    ->label('标识')
                    ->searchable()
                    ->sortable(),
                TextColumn::make('display_name')
                    ->label('显示名称')
                    ->searchable()
                    ->sortable(),
                TextColumn::make('permissions_count')
                    ->label('权限数')
                    ->counts('permissions'),
                TextColumn::make('admins_count')
                    ->label('管理员数')
                    ->counts('admins'),
                TextColumn::make('updated_at')
                    ->label('更新时间')
                    ->dateTime('Y-m-d H:i')
                    ->sortable(),
            ])
            ->recordActions([
                EditAction::make(),
                DeleteAction::make()
                    ->disabled(fn (Role $record): bool => $record->isSuperAdmin()),
            ])
            ->toolbarActions([
                BulkActionGroup::make([
                    DeleteBulkAction::make(),
                ]),
            ]);
    }

    public static function getPages(): array
    {
        return [
            'index' => ListRoles::route('/'),
            'create' => CreateRole::route('/create'),
            'edit' => EditRole::route('/{record}/edit'),
        ];
    }
}
