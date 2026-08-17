<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Attributes\Fillable;

#[Fillable(['name', 'display_name', 'description', 'group'])]
class Permission extends \Illuminate\Database\Eloquent\Model
{
    protected $table = 'permissions';

    /**
     * The roles that belong to the permission.
     */
    public function roles()
    {
        return $this->belongsToMany(Role::class, 'role_permission', 'permission_id', 'role_id');
    }
}
