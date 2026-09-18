<?php

namespace App\Filament\Resources\Media;

use App\Filament\Resources\Media\Pages\CreateMedia;
use App\Filament\Resources\Media\Pages\EditMedia;
use App\Filament\Resources\Media\Pages\ListMedia;
use App\Models\Media;
use Filament\Actions\BulkActionGroup;
use Filament\Actions\DeleteAction;
use Filament\Actions\DeleteBulkAction;
use Filament\Actions\EditAction;
use Filament\Forms\Components\FileUpload;
use Filament\Forms\Components\TextInput;
use Filament\Resources\Resource;
use Filament\Schemas\Schema;
use Filament\Tables\Columns\TextColumn;
use Filament\Tables\Filters\SelectFilter;
use Filament\Tables\Table;

class MediaResource extends Resource
{
    protected static ?string $model = Media::class;

    protected static string|\BackedEnum|null $navigationIcon = 'heroicon-o-photo';

    protected static string|\UnitEnum|null $navigationGroup = '系统管理';

    protected static ?string $navigationLabel = '媒体库';

    protected static ?string $modelLabel = '媒体';

    protected static ?string $pluralModelLabel = '媒体';

    protected static ?int $navigationSort = 3;

    protected static ?string $recordTitleAttribute = 'name';

    public static function form(Schema $schema): Schema
    {
        return $schema
            ->components([
                TextInput::make('name')
                    ->label('名称')
                    ->maxLength(255)
                    ->helperText('留空则使用上传文件名'),
                FileUpload::make('upload')
                    ->label('文件')
                    ->disk('public')
                    ->directory('media')
                    ->required(fn (string $operation): bool => $operation === 'create')
                    ->visibleOn('create')
                    ->columnSpanFull(),
                TextInput::make('path')
                    ->label('路径')
                    ->disabled()
                    ->dehydrated(false)
                    ->visibleOn('edit'),
                TextInput::make('mime')
                    ->label('MIME')
                    ->disabled()
                    ->dehydrated(false)
                    ->visibleOn('edit'),
                TextInput::make('size')
                    ->label('大小（字节）')
                    ->disabled()
                    ->dehydrated(false)
                    ->visibleOn('edit'),
            ]);
    }

    public static function table(Table $table): Table
    {
        return $table
            ->columns([
                TextColumn::make('name')
                    ->label('名称')
                    ->searchable()
                    ->sortable(),
                TextColumn::make('mime')
                    ->label('MIME')
                    ->badge()
                    ->searchable(),
                TextColumn::make('size')
                    ->label('大小')
                    ->formatStateUsing(fn (?int $state): string => $state === null
                        ? '—'
                        : number_format($state / 1024, 1).' KB'),
                TextColumn::make('url')
                    ->label('URL')
                    ->url(fn (Media $record): string => $record->url)
                    ->openUrlInNewTab()
                    ->limit(40),
                TextColumn::make('uploader.name')
                    ->label('上传者')
                    ->placeholder('—'),
                TextColumn::make('created_at')
                    ->label('上传时间')
                    ->dateTime('Y-m-d H:i')
                    ->sortable(),
            ])
            ->filters([
                SelectFilter::make('mime')
                    ->label('MIME')
                    ->options(fn (): array => Media::query()
                        ->whereNotNull('mime')
                        ->distinct()
                        ->orderBy('mime')
                        ->pluck('mime', 'mime')
                        ->all()),
            ])
            ->defaultSort('created_at', 'desc')
            ->recordActions([
                EditAction::make(),
                DeleteAction::make(),
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
            'index' => ListMedia::route('/'),
            'create' => CreateMedia::route('/create'),
            'edit' => EditMedia::route('/{record}/edit'),
        ];
    }
}
