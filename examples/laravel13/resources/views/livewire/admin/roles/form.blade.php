<div class="max-w-2xl">
    <div class="bg-white rounded-lg shadow p-6">
        <h2 class="text-lg font-semibold mb-6">{{ $roleId ? '编辑角色' : '新建角色' }}</h2>

        <form wire:submit="save">
            <div class="space-y-4">
                <div>
                    <label class="block text-sm font-medium text-gray-700 mb-1">名称</label>
                    <input type="text" wire:model="name" class="w-full px-3 py-2 border border-gray-300 rounded-md" @if($roleId === 1) readonly @endif>
                    @error('name') <span class="text-red-500 text-xs">{{ $message }}</span> @enderror
                </div>

                <div>
                    <label class="block text-sm font-medium text-gray-700 mb-1">显示名称</label>
                    <input type="text" wire:model="display_name" class="w-full px-3 py-2 border border-gray-300 rounded-md">
                    @error('display_name') <span class="text-red-500 text-xs">{{ $message }}</span> @enderror
                </div>

                <div>
                    <label class="block text-sm font-medium text-gray-700 mb-1">描述</label>
                    <textarea wire:model="description" rows="2" class="w-full px-3 py-2 border border-gray-300 rounded-md"></textarea>
                    @error('description') <span class="text-red-500 text-xs">{{ $message }}</span> @enderror
                </div>

                <div>
                    <label class="block text-sm font-medium text-gray-700 mb-3">权限</label>
                    @foreach($permissionGroups as $group => $perms)
                    <div class="mb-4">
                        <div class="text-sm font-semibold text-gray-600 mb-2">{{ $group }}</div>
                        <div class="grid grid-cols-2 gap-2">
                            @foreach($perms as $perm)
                            <label class="flex items-center">
                                <input type="checkbox" wire:model="permission_ids" value="{{ $perm->id }}" class="h-4 w-4 text-blue-600 border-gray-300 rounded">
                                <span class="ml-2 text-sm text-gray-700">{{ $perm->display_name }}</span>
                            </label>
                            @endforeach
                        </div>
                    </div>
                    @endforeach
                </div>

                <div class="flex space-x-3 pt-4">
                    <button type="submit" class="bg-blue-600 text-white px-4 py-2 rounded-md text-sm hover:bg-blue-700">保存</button>
                    <a href="{{ route('admin.roles.index') }}" class="bg-gray-300 text-gray-700 px-4 py-2 rounded-md text-sm hover:bg-gray-400">取消</a>
                </div>
            </div>
        </form>
    </div>
</div>
