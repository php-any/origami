<div class="max-w-2xl">
    <div class="bg-white rounded-lg shadow p-6">
        <h2 class="text-lg font-semibold mb-6">{{ $userId ? '编辑用户' : '新建用户' }}</h2>

        <form wire:submit="save">
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

                <div>
                    <label class="block text-sm font-medium text-gray-700 mb-1">{{ $userId ? '密码（留空则不修改）' : '密码' }}</label>
                    <input type="password" wire:model="password" class="w-full px-3 py-2 border border-gray-300 rounded-md">
                    @error('password') <span class="text-red-500 text-xs">{{ $message }}</span> @enderror
                </div>

                <div class="flex space-x-3 pt-4">
                    <button type="submit" class="bg-blue-600 text-white px-4 py-2 rounded-md text-sm hover:bg-blue-700">保存</button>
                    <a href="{{ route('admin.users.index') }}" class="bg-gray-300 text-gray-700 px-4 py-2 rounded-md text-sm hover:bg-gray-400">取消</a>
                </div>
            </div>
        </form>
    </div>
</div>
