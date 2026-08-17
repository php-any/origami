<div class="max-w-2xl">
    <div class="bg-white rounded-lg shadow p-6">
        <h2 class="text-lg font-semibold mb-6">{{ $adminId ? '编辑管理员' : '新建管理员' }}</h2>

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
                    <label class="block text-sm font-medium text-gray-700 mb-1">{{ $adminId ? '密码（留空则不修改）' : '密码' }}</label>
                    <input type="password" wire:model="password" class="w-full px-3 py-2 border border-gray-300 rounded-md">
                    @error('password') <span class="text-red-500 text-xs">{{ $message }}</span> @enderror
                </div>

                <div>
                    <label class="block text-sm font-medium text-gray-700 mb-1">角色</label>
                    <div class="space-y-2">
                        @foreach($roles as $role)
                        <label class="flex items-center">
                            <input type="checkbox" wire:model="role_ids" value="{{ $role->id }}" class="h-4 w-4 text-blue-600 border-gray-300 rounded">
                            <span class="ml-2 text-sm text-gray-700">{{ $role->display_name }}</span>
                        </label>
                        @endforeach
                    </div>
                </div>

                <div>
                    <label class="flex items-center">
                        <input type="checkbox" wire:model="is_active" class="h-4 w-4 text-blue-600 border-gray-300 rounded">
                        <span class="ml-2 text-sm text-gray-700">启用</span>
                    </label>
                </div>

                <div class="flex space-x-3 pt-4">
                    <button type="submit" class="bg-blue-600 text-white px-4 py-2 rounded-md text-sm hover:bg-blue-700">保存</button>
                    <a href="{{ route('admin.admins.index') }}" class="bg-gray-300 text-gray-700 px-4 py-2 rounded-md text-sm hover:bg-gray-400">取消</a>
                </div>
            </div>
        </form>
    </div>
</div>
