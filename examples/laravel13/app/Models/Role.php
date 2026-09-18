<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Attributes\Fillable;
use Spatie\Activitylog\LogOptions;
use Spatie\Activitylog\Traits\LogsActivity;

#[Fillable(['name', 'display_name', 'description'])]
class Role extends \Illuminate\Database\Eloquent\Model
{
    use LogsActivity;

    protected $table = 'roles';

    public function getActivitylogOptions(): LogOptions
    {
        return LogOptions::defaults()
            ->logFillable()
            ->logOnlyDirty()
            ->dontSubmitEmptyLogs();
    }

    public function permissions()
    {
        return $this->belongsToMany(Permission::class, 'role_permission', 'role_id', 'permission_id');
    }

    public function admins()
    {
        return $this->belongsToMany(Admin::class, 'admin_role', 'role_id', 'admin_id');
    }

    public function isSuperAdmin(): bool
    {
        return $this->name === 'super-admin';
    }
}
