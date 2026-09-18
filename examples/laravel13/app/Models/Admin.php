<?php

namespace App\Models;

use Filament\Models\Contracts\FilamentUser;
use Filament\Models\Contracts\HasName;
use Filament\Panel;
use Illuminate\Database\Eloquent\Attributes\Fillable;
use Illuminate\Database\Eloquent\Attributes\Hidden;
use Illuminate\Foundation\Auth\User as Authenticatable;
use Illuminate\Notifications\Notifiable;
use Spatie\Activitylog\LogOptions;
use Spatie\Activitylog\Traits\LogsActivity;

#[Fillable(['name', 'email', 'password', 'is_active'])]
#[Hidden(['password', 'remember_token'])]
class Admin extends Authenticatable implements FilamentUser, HasName
{
    use LogsActivity;
    use Notifiable;

    protected $table = 'admins';

    /** @var array<string, bool>|null */
    protected ?array $permissionCache = null;

    protected function casts(): array
    {
        return [
            'email_verified_at' => 'datetime',
            'password' => 'hashed',
            'is_active' => 'boolean',
        ];
    }

    public function canAccessPanel(Panel $panel): bool
    {
        return (bool) $this->is_active;
    }

    public function getFilamentName(): string
    {
        return $this->name;
    }

    public function getActivitylogOptions(): LogOptions
    {
        return LogOptions::defaults()
            ->logFillable()
            ->logOnlyDirty()
            ->dontSubmitEmptyLogs();
    }

    public function roles()
    {
        return $this->belongsToMany(Role::class, 'admin_role', 'admin_id', 'role_id');
    }

    public function hasRole(string $role): bool
    {
        return $this->roles()->where('name', $role)->exists();
    }

    public function hasPermission(string $permission): bool
    {
        if ($this->permissionCache === null) {
            $this->permissionCache = [];
            $this->loadMissing('roles.permissions');
            foreach ($this->roles as $role) {
                foreach ($role->permissions as $perm) {
                    $this->permissionCache[$perm->name] = true;
                }
            }
        }

        return isset($this->permissionCache[$permission]);
    }

    public function receivesBroadcastNotificationsOn(): string
    {
        return 'admins.'.$this->id;
    }
}
