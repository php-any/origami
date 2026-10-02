<?php

namespace App\Filament\Resources\Products;

use App\Filament\Resources\Products\Pages\CreateProduct;
use App\Filament\Resources\Products\Pages\EditProduct;
use App\Filament\Resources\Products\Pages\ListProducts;
use App\Models\Product;
use Filament\Actions\BulkAction;
use Filament\Actions\BulkActionGroup;
use Filament\Actions\DeleteAction;
use Filament\Actions\DeleteBulkAction;
use Filament\Actions\EditAction;
use Filament\Forms\Components\FileUpload;
use Filament\Forms\Components\Select;
use Filament\Forms\Components\Textarea;
use Filament\Forms\Components\TextInput;
use Filament\Resources\Resource;
use Filament\Schemas\Schema;
use Filament\Schemas\Components\Section;
use Filament\Schemas\Components\Group;
use Filament\Tables\Columns\ImageColumn;
use Filament\Tables\Columns\TextColumn;
use Filament\Tables\Filters\SelectFilter;
use Filament\Tables\Table;
use Illuminate\Database\Eloquent\Collection;

class ProductResource extends Resource
{
    protected static ?string $model = Product::class;

    protected static string|\BackedEnum|null $navigationIcon = 'heroicon-o-cube';

    protected static string|\UnitEnum|null $navigationGroup = '商品管理';

    protected static ?string $navigationLabel = '商品';

    protected static ?string $modelLabel = '商品';

    protected static ?string $pluralModelLabel = '商品';

    protected static ?int $navigationSort = 2;

    protected static ?string $recordTitleAttribute = 'name';

    public static function form(Schema $schema): Schema
    {
        return $schema
            ->columns(['default' => 1, 'lg' => 3])
            ->components([
                Group::make([
                    Section::make('商品信息')
                        ->description('商品名称、编码及分类信息')
                        ->icon('heroicon-o-cube')
                        ->columns(2)
                        ->schema([
                            TextInput::make('name')
                                ->label('商品名称')
                                ->required()
                                ->maxLength(255),
                            TextInput::make('sku')
                                ->label('SKU')
                                ->required()
                                ->unique(ignoreRecord: true)
                                ->maxLength(255),
                            Select::make('category_id')
                                ->label('分类')
                                ->relationship('category', 'name')
                                ->searchable()
                                ->preload()
                                ->nullable()
                                ->columnSpanFull(),
                            Textarea::make('description')
                                ->label('商品描述')
                                ->rows(6)
                                ->columnSpanFull(),
                        ]),
                    Section::make('商品图片')
                        ->description('用于商品列表和详情展示')
                        ->icon('heroicon-o-photo')
                        ->schema([
                            FileUpload::make('image')
                                ->label('图片')
                                ->image()
                                ->disk('public')
                                ->directory('products')
                                ->imageEditor()
                                ->nullable(),
                        ]),
                ])->columnSpan(['default' => 1, 'lg' => 2]),
                Group::make([
                    Section::make('价格与库存')
                        ->description('设置销售价格与可用库存')
                        ->icon('heroicon-o-banknotes')
                        ->schema([
                            TextInput::make('price')
                                ->label('价格')
                                ->numeric()
                                ->required()
                                ->prefix('¥')
                                ->minValue(0),
                            TextInput::make('stock')
                                ->label('库存')
                                ->numeric()
                                ->required()
                                ->default(0)
                                ->minValue(0),
                        ]),
                    Section::make('发布设置')
                        ->description('上架商品可在前台展示')
                        ->icon('heroicon-o-adjustments-horizontal')
                        ->schema([
                            Select::make('status')
                                ->label('状态')
                                ->options([
                                    'active' => '上架',
                                    'inactive' => '下架',
                                ])
                                ->default('active')
                                ->required(),
                        ]),
                ])->columnSpan(1),
            ]);
    }

    public static function table(Table $table): Table
    {
        return $table
            ->columns([
                ImageColumn::make('image')
                    ->label('图片')
                    ->disk('public')
                    ->square(),
                TextColumn::make('name')
                    ->label('名称')
                    ->searchable()
                    ->sortable(),
                TextColumn::make('sku')
                    ->label('SKU')
                    ->searchable()
                    ->copyable(),
                TextColumn::make('category.name')
                    ->label('分类')
                    ->placeholder('—')
                    ->sortable(),
                TextColumn::make('price')
                    ->label('价格')
                    ->money('CNY')
                    ->sortable(),
                TextColumn::make('stock')
                    ->label('库存')
                    ->sortable()
                    ->color(fn (Product $record): string => $record->stock <= 10 ? 'danger' : 'success'),
                TextColumn::make('status')
                    ->label('状态')
                    ->badge()
                    ->formatStateUsing(fn (string $state): string => $state === 'active' ? '上架' : '下架')
                    ->color(fn (string $state): string => $state === 'active' ? 'success' : 'gray'),
                TextColumn::make('updated_at')
                    ->label('更新时间')
                    ->dateTime('Y-m-d H:i')
                    ->sortable()
                    ->toggleable(isToggledHiddenByDefault: true),
            ])
            ->filters([
                SelectFilter::make('status')
                    ->label('状态')
                    ->options([
                        'active' => '上架',
                        'inactive' => '下架',
                    ]),
                SelectFilter::make('category_id')
                    ->label('分类')
                    ->relationship('category', 'name'),
            ])
            ->recordActions([
                EditAction::make(),
                DeleteAction::make(),
            ])
            ->toolbarActions([
                BulkActionGroup::make([
                    BulkAction::make('activate')
                        ->label('批量上架')
                        ->icon('heroicon-o-check-circle')
                        ->color('success')
                        ->requiresConfirmation()
                        ->action(fn (Collection $records) => $records->each->update(['status' => 'active']))
                        ->deselectRecordsAfterCompletion(),
                    BulkAction::make('deactivate')
                        ->label('批量下架')
                        ->icon('heroicon-o-x-circle')
                        ->color('warning')
                        ->requiresConfirmation()
                        ->action(fn (Collection $records) => $records->each->update(['status' => 'inactive']))
                        ->deselectRecordsAfterCompletion(),
                    DeleteBulkAction::make(),
                ]),
            ]);
    }

    public static function getPages(): array
    {
        return [
            'index' => ListProducts::route('/'),
            'create' => CreateProduct::route('/create'),
            'edit' => EditProduct::route('/{record}/edit'),
        ];
    }
}
