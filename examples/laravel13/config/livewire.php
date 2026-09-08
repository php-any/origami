<?php

return [
    /*
    |--------------------------------------------------------------------------
    | Class Namespace
    |--------------------------------------------------------------------------
    |
    | This value sets the root namespace for Livewire component classes in
    | your application. This value will change where component classes are
    | discovered and loaded from.
    |
    */

    'class_namespace' => 'App\\Livewire',

    /*
    |--------------------------------------------------------------------------
    | View Path
    |--------------------------------------------------------------------------
    |
    | This value is used to specify where Livewire component Blade templates are
    | stored when running `artisan make:livewire`. It is used when generating
    | new Livewire components.
    |
    */

    'view_path' => resource_path('views/livewire'),

    /*
    |--------------------------------------------------------------------------
    | Layout
    |--------------------------------------------------------------------------
    | The default layout view that will be used when rendering Livewire
    | components.
    */

    'layout' => 'layouts.admin',

    // Livewire 3 实际读取 component_layout；包默认是 layouts::app。
    // 与上面 layout 对齐，否则全页组件会去找不存在的 layouts::app。
    'component_layout' => 'layouts.admin',

    /*
    |--------------------------------------------------------------------------
    | Temporary File Upload Endpoint
    |--------------------------------------------------------------------------
    |
    | This value is the URI that Livewire will use as its temporary file upload
    | endpoint. The endpoint will handle the file upload request when files
    | are uploaded as part of a Livewire component's interaction.
    |
    */

    'temporary_file_upload' => [
        'disk' => null,
        'rules' => null,
        'directory' => null,
        'middleware' => null,
        'preview_mimes' => [
            'png', 'gif', 'bmp', 'svg', 'wav', 'mp4',
            'mov', 'avi', 'wmv', 'mp3', 'm4a',
            'jpg', 'jpeg', 'mpga', 'webp', 'wma',
        ],
        'max_size' => 5 * 1024,
        'image' => [
            'validate' => false,
            'width' => 1024,
            'height' => 1024,
        ],
        'timeout' => 60,
    ],

    /*
    |--------------------------------------------------------------------------
    | Render On Redirect
    |--------------------------------------------------------------------------
    |
    | This value determines if Livewire will render its components on redirects.
    | This can improve performance for requests with many components.
    |
    */

    'render_on_redirect' => false,

    /*
    |--------------------------------------------------------------------------
    | Eloquent Model Binding
    |--------------------------------------------------------------------------
    |
    | Previous versions of Livewire supported binding directly to eloquent model
    | properties using wire:model. However, since this behavior has been
    | removed, we've added a way to still use this if you want.
    |
    */

    'legacy_model_binding' => false,

    /*
    |--------------------------------------------------------------------------
    | Auto-inject Assets
    |--------------------------------------------------------------------------
    |
    | Determines whether Livewire's JavaScript and CSS are auto-injected
    | into the rendered page.
    |
    */

    'inject_assets' => true,

    /*
    |--------------------------------------------------------------------------
    | Navigate
    |--------------------------------------------------------------------------
    |
    | By default, Livewire renders with its "wire:navigate" feature.
    |
    */

    'navigate' => [
        'show_progress_bar' => true,
        'progress_bar_color' => '#2299dd',
    ],

    /*
    |--------------------------------------------------------------------------
    | SQL Morph Map
    |--------------------------------------------------------------------------
    |
    | The morph_map is an associative array which maps internal model names
    | to their actual class names.
    |
    */

    'morph_map' => [
        'user' => \App\Models\User::class,
        'admin' => \App\Models\Admin::class,
        'role' => \App\Models\Role::class,
        'permission' => \App\Models\Permission::class,
        'order' => \App\Models\Order::class,
        'product' => \App\Models\Product::class,
    ],

    /*
    |--------------------------------------------------------------------------
    | Pagination Theme
    |--------------------------------------------------------------------------
    |
    | This is the default pagination theme.
    |
    */

    'pagination_theme' => 'tailwind',

    /*
    |--------------------------------------------------------------------------
    | Back Button Cache
    |--------------------------------------------------------------------------
    |
    | The back button cache stores the rendered HTML of Livewire components
    | to be shown when the user navigates back via the browser's back button.
    |
    */

    'back_button_cache' => false,
];
