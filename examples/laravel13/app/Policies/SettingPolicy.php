<?php

namespace App\Policies;

use App\Models\Admin;
use App\Models\Setting;

class SettingPolicy extends AdminResourcePolicy
{
    public function viewAny(?Admin $admin): bool
    {
        return $this->allows($admin, 'settings.edit');
    }

    public function update(?Admin $admin, ?Setting $setting = null): bool
    {
        return $this->allows($admin, 'settings.edit');
    }
}
