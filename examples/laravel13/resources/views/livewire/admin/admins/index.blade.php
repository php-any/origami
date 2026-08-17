<div>
    <div class="mb-4 flex items-center justify-between">
        <div class="flex space-x-2">
            <input type="text" wire:model.live.debounce.300ms="search" placeholder="搜索管理员..." class="px-3 py-2 border border-gray-300 rounded-md text-sm">
            <select wire:model.live="statusFilter" class="px-3 py-2 border border-gray-300 rounded-md text-sm">
                <option value="">全部状态</option>
                <option value="active">启用</option>
                <option value="inactive">禁用</option>
            </select>
        </div>
        <a href="{{ route('admin.admins.create') }}" class="bg-blue-600 text-white px-4 py-2 rounded-md text-sm hover:bg-blue-700">+ 新建管理员</a>
    </div>

    <div class="bg-white rounded-lg shadow overflow-hidden">
        <table class="min-w-full divide-y divide-gray-200">
            <thead class="bg-gray-50">
                <tr>
                    <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">ID</th>
                    <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">名称</th>
                    <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">邮箱</th>
                    <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">角色</th>
                    <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">状态</th>
                    <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">操作</th>
                </tr>
            </thead>
            <tbody class="bg-white divide-y divide-gray-200">
                @forelse($admins as $admin)
                <tr>
                    <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">{{ $admin->id }}</td>
                    <td class="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">{{ $admin->name }}</td>
                    <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{{ $admin->email }}</td>
                    <td class="px-6 py-4 whitespace-nowrap">
                        @foreach($admin->roles as $role)
                        <span class="px-2 py-1 text-xs rounded-full bg-indigo-100 text-indigo-800 mr-1">{{ $role->display_name }}</span>
                        @endforeach
                    </td>
                    <td class="px-6 py-4 whitespace-nowrap">
                        <button wire:click="toggleActive({{ $admin->id }})" @if($admin->id === auth('admin')->id()) disabled @endif
                            class="px-2 py-1 text-xs rounded-full {{ $admin->is_active ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800' }}">
                            {{ $admin->is_active ? '启用' : '禁用' }}
                        </button>
                    </td>
                    <td class="px-6 py-4 whitespace-nowrap text-sm">
                        <a href="{{ route('admin.admins.edit', $admin->id) }}" class="text-blue-600 hover:text-blue-800 mr-3">编辑</a>
                        @if($admin->id !== auth('admin')->id())
                        <button wire:click="delete({{ $admin->id }})" wire:confirm="确定删除该管理员？" class="text-red-600 hover:text-red-800">删除</button>
                        @endif
                    </td>
                </tr>
                @empty
                <tr>
                    <td colspan="6" class="px-6 py-4 text-center text-gray-500">暂无数据</td>
                </tr>
                @endforelse
            </tbody>
        </table>
        <div class="px-6 py-4">
            {{ $admins->links() }}
        </div>
    </div>
</div>
