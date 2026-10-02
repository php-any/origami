{{-- Tabler's official vertical layout, backed by Filament's authorized navigation. --}}
@php
    $livewire ??= null;
    $renderHookScopes = $livewire?->getRenderHookScopes();
@endphp
<x-filament-panels::layout.base :livewire="$livewire">
    <div class="tabler-ui" data-bs-theme="light">
    <div class="page min-vh-100">
        <a href="#fi-main-content" class="visually-hidden-focusable skip-link">跳转到主要内容</a>
        <aside class="navbar navbar-vertical navbar-expand-lg bg-dark d-print-none" data-bs-theme="dark" aria-label="主导航">
            <div class="container-fluid">
                <button class="navbar-toggler" type="button" data-bs-toggle="collapse" data-bs-target="#sidebar-menu" aria-controls="sidebar-menu" aria-expanded="false" aria-label="展开导航">
                    <span class="navbar-toggler-icon"></span>
                </button>
                <div class="navbar-brand">
                    <a href="{{ filament()->getHomeUrl() }}" class="d-flex align-items-center gap-2 text-reset" aria-label="Origami Admin 首页">
                        <x-filament::icon icon="heroicon-o-square-3-stack-3d" class="icon" />
                        <span class="nav-link-title">Origami Admin</span>
                    </a>
                    <button type="button" class="btn btn-action btn-sm d-none d-lg-inline-flex" data-bs-toggle="sidebar-folded" aria-pressed="false" aria-label="折叠侧栏">
                        <x-filament::icon icon="heroicon-o-bars-3-bottom-left" class="icon" />
                    </button>
                </div>
                <nav class="collapse navbar-collapse" id="sidebar-menu" aria-label="业务模块">
                    <ul class="navbar-nav pt-lg-3">
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
                    <div class="px-3 py-3 small text-secondary nav-link-title">Origami · 管理控制台</div>
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
            <footer class="footer footer-transparent d-print-none">
                <div class="container-fluid">
                    <div class="row text-center align-items-center flex-row-reverse">
                        <div class="col-lg-auto ms-lg-auto"><span class="text-secondary">{{ now()->format('Y 年 m 月 d 日') }}</span></div>
                        <div class="col-12 col-lg-auto mt-3 mt-lg-0"><span class="text-secondary">Origami Admin</span></div>
                    </div>
                </div>
            </footer>
            {{ \Filament\Support\Facades\FilamentView::renderHook(\Filament\View\PanelsRenderHook::FOOTER, scopes: $renderHookScopes) }}
        </div>
    </div>
    </div>
</x-filament-panels::layout.base>
