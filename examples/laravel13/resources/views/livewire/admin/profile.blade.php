<div class="max-w-2xl space-y-6">
    @if(session('message'))
    <div class="bg-green-100 border-l-4 border-green-500 text-green-700 p-4">
        {{ session('message') }}
    </div>
    @endif

    @if(session('error'))
    <div class="bg-red-100 border-l-4 border-red-500 text-red-700 p-4">
        {{ session('error') }}
    </div>
    @endif

    <div class="bg-white rounded-lg shadow p-6">
        <h2 class="text-lg font-semibold mb-6">个人资料</h2>
        <form wire:submit="updateProfile">
            <div class="space-y-4">
                <div>
                    <label class="block text-sm font-medium text-gray-700 mb-1">名称</label>
                    <input type="text" wire:model="name" class="w-full px-3 py-2 border border-gray-300 rounded-md">
                    @error('name') <span class="text-red-500 text-xs">{{ $message }}</span> @enderror
                </div>
                <div>
                    <label class="block text-sm font-medium text-gray-700 mb-1">邮箱</label>
                    <input type="email" wire:model="email" class="w-full px-3 py-2 border border-gray-300 rounded-md">
                    @error('email') <span class="text-red-500 text-xs">{{ $message }}</span> @enderror
                </div>
                <button type="submit" class="bg-blue-600 text-white px-4 py-2 rounded-md text-sm hover:bg-blue-700">更新资料</button>
            </div>
        </form>
    </div>

    <div class="bg-white rounded-lg shadow p-6">
        <h2 class="text-lg font-semibold mb-6">修改密码</h2>
        <form wire:submit="changePassword">
            <div class="space-y-4">
                <div>
                    <label class="block text-sm font-medium text-gray-700 mb-1">当前密码</label>
                    <input type="password" wire:model="current_password" class="w-full px-3 py-2 border border-gray-300 rounded-md">
                    @error('current_password') <span class="text-red-500 text-xs">{{ $message }}</span> @enderror
                </div>
                <div>
                    <label class="block text-sm font-medium text-gray-700 mb-1">新密码</label>
                    <input type="password" wire:model="new_password" class="w-full px-3 py-2 border border-gray-300 rounded-md">
                    @error('new_password') <span class="text-red-500 text-xs">{{ $message }}</span> @enderror
                </div>
                <div>
                    <label class="block text-sm font-medium text-gray-700 mb-1">确认新密码</label>
                    <input type="password" wire:model="confirm_password" class="w-full px-3 py-2 border border-gray-300 rounded-md">
                    @error('confirm_password') <span class="text-red-500 text-xs">{{ $message }}</span> @enderror
                </div>
                <button type="submit" class="bg-green-600 text-white px-4 py-2 rounded-md text-sm hover:bg-green-700">修改密码</button>
            </div>
        </form>
    </div>
</div>
