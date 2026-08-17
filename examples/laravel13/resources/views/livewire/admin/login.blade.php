<div class="min-h-screen flex items-center justify-center bg-gray-100">
    <div class="bg-white rounded-lg shadow-lg p-8 w-full max-w-md">
        <h1 class="text-2xl font-bold text-center mb-8">管理后台登录</h1>

        @if(session('error'))
        <div class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded mb-4">
            {{ session('error') }}
        </div>
        @endif

        <form wire:submit="login">
            <div class="mb-4">
                <label class="block text-sm font-medium text-gray-700 mb-1">邮箱</label>
                <input type="email" wire:model="email" class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500" placeholder="admin@example.com">
                @error('email') <span class="text-red-500 text-xs">{{ $message }}</span> @enderror
            </div>

            <div class="mb-4">
                <label class="block text-sm font-medium text-gray-700 mb-1">密码</label>
                <input type="password" wire:model="password" class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500" placeholder="请输入密码">
                @error('password') <span class="text-red-500 text-xs">{{ $message }}</span> @enderror
            </div>

            <div class="mb-4 flex items-center">
                <input type="checkbox" wire:model="remember" class="h-4 w-4 text-blue-600 border-gray-300 rounded">
                <label class="ml-2 text-sm text-gray-700">记住我</label>
            </div>

            <button type="submit" class="w-full bg-blue-600 text-white py-2 px-4 rounded-md hover:bg-blue-700 transition duration-200">
                登录
            </button>
        </form>

        <p class="mt-4 text-sm text-gray-500 text-center">默认账号：admin@example.com / password</p>
    </div>
</div>
