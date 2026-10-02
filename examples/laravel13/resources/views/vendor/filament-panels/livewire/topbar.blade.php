<header class="navbar navbar-expand-md d-print-none">
    <div class="container-fluid">
        <div class="navbar-nav flex-row d-none d-lg-flex">
            <button type="button" class="nav-link px-0 me-3" data-bs-toggle="sidebar-folded" aria-pressed="false" aria-label="切换侧栏宽度"><x-filament::icon icon="heroicon-o-bars-3" class="icon" /></button>
            <a href="{{ filament()->getHomeUrl() }}" class="nav-link">管理控制台</a>
        </div>
        <div class="navbar-nav flex-row order-md-last gap-3 align-items-center">
            <div class="filament-content">
                @if (filament()->isGlobalSearchEnabled())
                    @livewire(\Filament\Livewire\GlobalSearch::class)
                @endif
            </div>
            <div class="filament-content">
                @if (filament()->hasDatabaseNotifications())
                    @livewire(filament()->getDatabaseNotificationsLivewireComponent(), ['lazy' => filament()->hasLazyLoadedDatabaseNotifications()])
                @endif
            </div>
            <span class="small text-secondary d-none d-md-block">{{ filament()->getUserName(filament()->auth()->user()) }}</span>
            <div class="filament-content"><x-filament-panels::user-menu /></div>
        </div>
    </div>
    <div class="filament-content"><x-filament-actions::modals /></div>
</header>
