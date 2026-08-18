<div>
    <div class="mb-4 flex items-center justify-between">
        <h2 class="text-lg font-semibold">角色列表</h2>
        <a href="{{ route('admin.roles.create') }}" class="bg-blue-600 text-white px-4 py-2 rounded-md text-sm hover:bg-blue-700">+ 新建角色</a>
    </div>

    <div class="bg-white rounded-lg shadow overflow-hidden">
        <table class="min-w-full divide-y divide-gray-200">
            <thead class="bg-gray-50">
                <tr>
                    <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">ID</th>
                    <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">名称</th>
                    <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">显示名称</th>
                    <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">管理员数</th>
                    <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">权限数</th>
                    <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">操作</th>
                </tr>
            </thead>
            <tbody class="bg-white divide-y divide-gray-200">
                @forelse($roles as $role)
                <tr>
                    <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">{{ $role->id }}</td>
                    <td class="px-6 py-4 whitespace-nowrap text-sm font-mono text-gray-900">{{ $role->name }}</td>
                    <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{{ $role->display_name }}</td>
                    <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{{ $role->admins_count }}</td>
                    <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{{ $role->permissions_count }}</td>
                    <td class="px-6 py-4 whitespace-nowrap text-sm">
                        <a href="{{ route('admin.roles.edit', $role->id) }}" class="text-blue-600 hover:text-blue-800 mr-3">编辑</a>
                        @if($role->name !== 'super-admin')
                        <button wire:click="delete({{ $role->id }})" wire:confirm="确定删除该角色？" class="text-red-600 hover:text-red-800">删除</button>
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
    </div>
</div>
