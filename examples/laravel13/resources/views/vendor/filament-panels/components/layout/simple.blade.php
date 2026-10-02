{{-- Tabler's sign-in card; the slot remains the official Livewire auth form. --}}
@props(['after' => null, 'heading' => null, 'subheading' => null])
@php
    $livewire ??= null;
@endphp
<x-filament-panels::layout.base :livewire="$livewire">
    <div class="tabler-ui" data-bs-theme="light">
    <div class="page page-center min-vh-100">
        <div class="container container-tight py-4">
            <div class="text-center mb-4">
                <span class="navbar-brand navbar-brand-autodark d-inline-flex gap-2">
                    <x-filament::icon icon="heroicon-o-square-3-stack-3d" class="icon text-primary" /> Origami Admin
                </span>
            </div>
            <div class="card card-md">
                <div class="card-body">
                    <main class="filament-content" id="fi-main-content" tabindex="-1">{{ $slot }}</main>
                </div>
            </div>
            <div class="text-center text-secondary mt-3">请使用管理员账号登录管理控制台</div>
        </div>
    </div>
    </div>
</x-filament-panels::layout.base>
