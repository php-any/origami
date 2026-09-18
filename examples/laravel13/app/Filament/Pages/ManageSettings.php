<?php

namespace App\Filament\Pages;

use App\Models\Setting;
use Filament\Actions\Action;
use Filament\Forms\Components\TextInput;
use Filament\Forms\Components\Toggle;
use Filament\Notifications\Notification;
use Filament\Pages\Page;
use Filament\Schemas\Components\Actions;
use Filament\Schemas\Components\Form;
use Filament\Schemas\Schema;

/**
 * @property-read Schema $form
 */
class ManageSettings extends Page
{
    protected static string|\BackedEnum|null $navigationIcon = 'heroicon-o-cog-6-tooth';

    protected static string|\UnitEnum|null $navigationGroup = '系统管理';

    protected static ?string $navigationLabel = '系统设置';

    protected static ?string $title = '系统设置';

    protected static ?string $slug = 'settings';

    protected static ?int $navigationSort = 90;

    protected string $view = 'filament.pages.manage-settings';

    /**
     * @var array<string, mixed>|null
     */
    public ?array $data = [];

    public static function canAccess(): bool
    {
        $admin = auth('admin')->user();

        return $admin !== null && $admin->hasPermission('settings.edit');
    }

    public function mount(): void
    {
        $this->form->fill([
            'site_name' => Setting::getValue('site_name', 'Origami Shop'),
            'maintenance_mode' => (bool) Setting::getValue('maintenance_mode', false),
            'order_prefix' => Setting::getValue('order_prefix', 'ORD'),
        ]);
    }

    public function form(Schema $schema): Schema
    {
        return $schema
            ->components([
                Form::make([
                    TextInput::make('site_name')
                        ->label('站点名称')
                        ->required()
                        ->maxLength(255),
                    Toggle::make('maintenance_mode')
                        ->label('维护模式'),
                    TextInput::make('order_prefix')
                        ->label('订单号前缀')
                        ->required()
                        ->maxLength(20),
                ])
                    ->livewireSubmitHandler('save')
                    ->footer([
                        Actions::make([
                            Action::make('save')
                                ->label('保存设置')
                                ->submit('save')
                                ->keyBindings(['mod+s']),
                        ]),
                    ]),
            ])
            ->statePath('data');
    }

    public function save(): void
    {
        $data = $this->form->getState();

        Setting::setValue('site_name', $data['site_name'], 'general', '站点名称');
        Setting::setValue('maintenance_mode', (bool) $data['maintenance_mode'], 'general', '维护模式');
        Setting::setValue('order_prefix', $data['order_prefix'], 'orders', '订单号前缀');

        Notification::make()
            ->success()
            ->title('设置已保存')
            ->send();
    }
}
