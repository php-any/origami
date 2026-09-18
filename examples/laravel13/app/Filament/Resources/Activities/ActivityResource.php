<?php

namespace App\Filament\Resources\Activities;

use App\Filament\Resources\Activities\Pages\ListActivities;
use App\Filament\Resources\Activities\Pages\ViewActivity;
use Filament\Actions\ViewAction;
use Filament\Infolists\Components\TextEntry;
use Filament\Resources\Resource;
use Filament\Schemas\Components\Section;
use Filament\Schemas\Schema;
use Filament\Tables\Columns\TextColumn;
use Filament\Tables\Filters\SelectFilter;
use Filament\Tables\Table;
use Illuminate\Database\Eloquent\Model;
use Spatie\Activitylog\Models\Activity;

class ActivityResource extends Resource
{
    protected static ?string $model = Activity::class;

    protected static string|\BackedEnum|null $navigationIcon = 'heroicon-o-clipboard-document-list';

    protected static string|\UnitEnum|null $navigationGroup = '系统管理';

    protected static ?string $navigationLabel = '操作日志';

    protected static ?string $modelLabel = '操作日志';

    protected static ?string $pluralModelLabel = '操作日志';

    protected static ?int $navigationSort = 4;

    public static function form(Schema $schema): Schema
    {
        return $schema->components([]);
    }

    public static function infolist(Schema $schema): Schema
    {
        return $schema
            ->components([
                Section::make('日志详情')
                    ->columns(2)
                    ->schema([
                        TextEntry::make('description')
                            ->label('描述')
                            ->columnSpanFull(),
                        TextEntry::make('event')
                            ->label('事件')
                            ->badge()
                            ->placeholder('—'),
                        TextEntry::make('log_name')
                            ->label('日志名')
                            ->placeholder('—'),
                        TextEntry::make('subject_type')
                            ->label('对象类型')
                            ->placeholder('—'),
                        TextEntry::make('subject_id')
                            ->label('对象 ID')
                            ->placeholder('—'),
                        TextEntry::make('causer_type')
                            ->label('操作者类型')
                            ->placeholder('—'),
                        TextEntry::make('causer_id')
                            ->label('操作者 ID')
                            ->placeholder('—'),
                        TextEntry::make('properties')
                            ->label('变更内容')
                            ->formatStateUsing(fn ($state): string => is_string($state)
                                ? $state
                                : json_encode($state, JSON_UNESCAPED_UNICODE | JSON_PRETTY_PRINT) ?: '—')
                            ->columnSpanFull(),
                        TextEntry::make('created_at')
                            ->label('时间')
                            ->dateTime('Y-m-d H:i:s'),
                    ]),
            ]);
    }

    public static function table(Table $table): Table
    {
        return $table
            ->columns([
                TextColumn::make('created_at')
                    ->label('时间')
                    ->dateTime('Y-m-d H:i:s')
                    ->sortable(),
                TextColumn::make('description')
                    ->label('描述')
                    ->searchable()
                    ->limit(50),
                TextColumn::make('event')
                    ->label('事件')
                    ->badge()
                    ->placeholder('—'),
                TextColumn::make('subject_type')
                    ->label('对象')
                    ->formatStateUsing(fn (?string $state, Activity $record): string => $state
                        ? class_basename($state).'#'.$record->subject_id
                        : '—'),
                TextColumn::make('causer_type')
                    ->label('操作者')
                    ->formatStateUsing(fn (?string $state, Activity $record): string => $state
                        ? class_basename($state).'#'.$record->causer_id
                        : '系统'),
                TextColumn::make('log_name')
                    ->label('日志名')
                    ->toggleable(isToggledHiddenByDefault: true),
            ])
            ->filters([
                SelectFilter::make('event')
                    ->label('事件')
                    ->options([
                        'created' => '创建',
                        'updated' => '更新',
                        'deleted' => '删除',
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
            'index' => ListActivities::route('/'),
            'view' => ViewActivity::route('/{record}'),
        ];
    }

    public static function canCreate(): bool
    {
        return false;
    }

    public static function canEdit(Model $record): bool
    {
        return false;
    }

    public static function canDelete(Model $record): bool
    {
        return false;
    }
}
