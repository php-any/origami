<div>
    <div class="mb-4 flex items-center justify-between">
        <input type="text" wire:model.live.debounce.300ms="search" placeholder="搜索权限..." class="px-3 py-2 border border-gray-300 rounded-md text-sm w-64">
        <a href="{{ route('admin.permissions.create') }}" class="bg-blue-600 text-white px-4 py-2 rounded-md text-sm hover:bg-blue-700">+ 新建权限</a>
    </div>

    <div class="bg-white rounded-lg shadow overflow-hidden">
        <table class="min-w-full divide-y divide-gray-200">
            <thead class="bg-gray-50">
                <tr>
                    <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">ID</th>
                    <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">名称</th>
                    <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">显示名称</th>
                    <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">分组</th>
                    <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">描述</th>
                    <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">操作</th>
                </tr>
            </thead>
            <tbody class="bg-white divide-y divide-gray-200">
                @forelse($permissions as $perm)
                <tr>
                    <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">{{ $perm->id }}</td>
                    <td class="px-6 py-4 whitespace-nowrap text-sm font-mono text-gray-900">{{ $perm->name }}</td>
                    <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{{ $perm->display_name }}</td>
                    <td class="px-6 py-4 whitespace-nowrap text-sm">
                        <span class="px-2 py-1 text-xs rounded-full bg-blue-100 text-blue-800">{{ $perm->group }}</span>
                    </td>
                    <td class="px-6 py-4 text-sm text-gray-500">{{ $perm->description }}</td>
                    <td class="px-6 py-4 whitespace-nowrap text-sm">
                        <a href="{{ route('admin.permissions.edit', $perm->id) }}" class="text-blue-600 hover:text-blue-800 mr-3">编辑</a>
                        <button wire:click="delete({{ $perm->id }})" wire:confirm="确定删除该权限？" class="text-red-600 hover:text-red-800">删除</button>
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
