{{-- Tabler's official vertical layout, backed by Filament's authorized navigation. --}}
@php
    $livewire ??= null;
    $renderHookScopes = $livewire?->getRenderHookScopes();
@endphp
<x-filament-panels::layout.base :livewire="$livewire">
    <div class="tabler-ui origami-admin" data-bs-theme="light">
    <div class="page min-vh-100">
        <a href="#fi-main-content" class="visually-hidden-focusable skip-link">跳转到主要内容</a>
        <aside class="navbar navbar-vertical navbar-expand-lg origami-sidebar d-print-none" aria-label="主导航">
            <div class="container-fluid">
                <button class="navbar-toggler" type="button" data-bs-toggle="collapse" data-bs-target="#sidebar-menu" aria-controls="sidebar-menu" aria-expanded="false" aria-label="展开导航">
                    <span class="navbar-toggler-icon"></span>
                </button>
                <div class="navbar-brand">
                    <a href="{{ filament()->getHomeUrl() }}" class="d-flex align-items-center gap-2 text-reset" aria-label="Origami Admin 首页">
                        <span class="origami-brand-mark"><x-filament::icon icon="heroicon-o-square-3-stack-3d" class="icon" /></span>
                        <span class="nav-link-title origami-brand-copy"><strong>Origami</strong><small>管理工作空间</small></span>
                    </a>
                </div>
                <nav class="collapse navbar-collapse" id="sidebar-menu" aria-label="业务模块">
                    <ul class="navbar-nav">
                        @foreach (filament()->getNavigation() as $group)
                            @if ($group->getLabel())
                                <li class="nav-section-title">{{ $group->getLabel() }}</li>
                            @endif
                            @foreach ($group->getItems() as $item)
                                <li @class(['nav-item', 'active' => $item->isActive()])>
                                    <a class="nav-link" href="{{ $item->getUrl() }}" @if ($item->isActive()) aria-current="page" @endif @if ($item->shouldOpenUrlInNewTab()) target="_blank" rel="noopener noreferrer" @endif title="{{ $item->getLabel() }}">
                                        <span class="nav-link-icon d-md-none d-lg-inline-block"><x-filament::icon :icon="$item->getIcon()" class="icon" /></span>
                                        <span class="nav-link-title">{{ $item->getLabel() }}</span>
                                        @if ($item->getBadge())<span class="badge bg-primary text-primary-fg ms-auto">{{ $item->getBadge() }}</span>@endif
                                    </a>
                                </li>
                            @endforeach
                        @endforeach
                    </ul>
                </nav>
                <div class="navbar-footer">
                    <div class="origami-workspace-label nav-link-title"><x-filament::icon icon="heroicon-o-square-3-stack-3d" class="icon" /><span>Origami Admin</span></div>
                </div>
            </div>
        </aside>

        <div class="page-wrapper">
            {{ \Filament\Support\Facades\FilamentView::renderHook(\Filament\View\PanelsRenderHook::TOPBAR_BEFORE, scopes: $renderHookScopes) }}
            @livewire(filament()->getTopbarLivewireComponent())
            {{ \Filament\Support\Facades\FilamentView::renderHook(\Filament\View\PanelsRenderHook::TOPBAR_AFTER, scopes: $renderHookScopes) }}
            {{ \Filament\Support\Facades\FilamentView::renderHook(\Filament\View\PanelsRenderHook::CONTENT_BEFORE, scopes: $renderHookScopes) }}
            <main class="page-body" id="fi-main-content" tabindex="-1">
                <div class="container-fluid">
                    <div class="filament-content">
                        {{ \Filament\Support\Facades\FilamentView::renderHook(\Filament\View\PanelsRenderHook::CONTENT_START, scopes: $renderHookScopes) }}
                        {{ $slot }}
                        {{ \Filament\Support\Facades\FilamentView::renderHook(\Filament\View\PanelsRenderHook::CONTENT_END, scopes: $renderHookScopes) }}
                    </div>
                </div>
            </main>
            {{ \Filament\Support\Facades\FilamentView::renderHook(\Filament\View\PanelsRenderHook::CONTENT_AFTER, scopes: $renderHookScopes) }}
            {{ \Filament\Support\Facades\FilamentView::renderHook(\Filament\View\PanelsRenderHook::FOOTER, scopes: $renderHookScopes) }}
        </div>
    </div>
    </div>
</x-filament-panels::layout.base>
