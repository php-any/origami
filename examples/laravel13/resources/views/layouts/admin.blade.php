<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>@yield('title', '管理后台') - Origami Admin</title>
    <script src="https://cdn.tailwindcss.com"></script>
    @livewireStyles
</head>
<body class="bg-gray-100">
    <div class="flex min-h-screen">
        <!-- Sidebar -->
        <aside class="w-64 bg-gray-900 text-white flex-shrink-0">
            <div class="p-4 text-xl font-bold border-b border-gray-700">
                🚀 Origami Admin
            </div>
            <nav class="p-4 space-y-1">
                <a href="{{ route('admin.dashboard') }}" class="block px-3 py-2 rounded hover:bg-gray-700 {{ request()->routeIs('admin.dashboard') ? 'bg-gray-700' : '' }}">
                    📊 仪表盘
                </a>

                @if(auth('admin')->user()?->hasPermission('admins.view'))
                <div class="mt-2 text-xs text-gray-400 uppercase px-3">系统管理</div>
                <a href="{{ route('admin.admins.index') }}" class="block px-3 py-2 rounded hover:bg-gray-700 {{ request()->routeIs('admin.admins.*') ? 'bg-gray-700' : '' }}">
                    👤 管理员管理
                </a>
                @endif

                @if(auth('admin')->user()?->hasPermission('users.view'))
                <a href="{{ route('admin.users.index') }}" class="block px-3 py-2 rounded hover:bg-gray-700 {{ request()->routeIs('admin.users.*') ? 'bg-gray-700' : '' }}">
                    👥 用户管理
                </a>
                @endif

                @if(auth('admin')->user()?->hasPermission('roles.view'))
                <a href="{{ route('admin.roles.index') }}" class="block px-3 py-2 rounded hover:bg-gray-700 {{ request()->routeIs('admin.roles.*') ? 'bg-gray-700' : '' }}">
                    🛡️ 角色管理
                </a>
                @endif

                @if(auth('admin')->user()?->hasPermission('permissions.view'))
                <a href="{{ route('admin.permissions.index') }}" class="block px-3 py-2 rounded hover:bg-gray-700 {{ request()->routeIs('admin.permissions.*') ? 'bg-gray-700' : '' }}">
                    🔑 权限管理
                </a>
                @endif

                <div class="mt-2 text-xs text-gray-400 uppercase px-3">业务管理</div>

                @if(auth('admin')->user()?->hasPermission('products.view'))
                <a href="{{ route('admin.products.index') }}" class="block px-3 py-2 rounded hover:bg-gray-700 {{ request()->routeIs('admin.products.*') ? 'bg-gray-700' : '' }}">
                    📦 产品管理
                </a>
                @endif

                @if(auth('admin')->user()?->hasPermission('orders.view'))
                <a href="{{ route('admin.orders.index') }}" class="block px-3 py-2 rounded hover:bg-gray-700 {{ request()->routeIs('admin.orders.*') ? 'bg-gray-700' : '' }}">
                    🛒 订单管理
                </a>
                @endif
            </nav>
        </aside>

        <!-- Main Content -->
        <div class="flex-1 flex flex-col">
            <!-- Header -->
            <header class="bg-white shadow-sm">
                <div class="flex items-center justify-between px-6 py-3">
                    <h1 class="text-lg font-semibold text-gray-800">@yield('title', '管理后台')</h1>
                    <div class="flex items-center space-x-4">
                        <span class="text-sm text-gray-600">{{ auth('admin')->user()?->name }}</span>
                        <a href="{{ route('admin.profile') }}" class="text-sm text-blue-600 hover:text-blue-800">个人设置</a>
                        <form method="POST" action="{{ route('admin.logout') }}">
                            @csrf
                            <button type="submit" class="text-sm text-red-600 hover:text-red-800">退出登录</button>
                        </form>
                    </div>
                </div>
            </header>

            <!-- Flash messages -->
            @if(session('message'))
            <div class="bg-green-100 border-l-4 border-green-500 text-green-700 p-4 m-6 mb-0">
                {{ session('message') }}
            </div>
            @endif

            @if(session('error'))
            <div class="bg-red-100 border-l-4 border-red-500 text-red-700 p-4 m-6 mb-0">
                {{ session('error') }}
            </div>
            @endif

            <!-- Content -->
            <main class="flex-1 p-6">
                {{ $slot ?? '' }}
                @yield('content')
            </main>
        </div>
    </div>
    @livewireScripts
</body>
</html>
