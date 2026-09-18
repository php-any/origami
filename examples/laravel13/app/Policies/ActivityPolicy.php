<?php

namespace App\Policies;

use App\Models\Admin;
use Spatie\Activitylog\Models\Activity;

class ActivityPolicy extends AdminResourcePolicy
{
    public function viewAny(?Admin $admin): bool
    {
        return $this->allows($admin, 'activity.view');
    }

    public function view(?Admin $admin, Activity $activity): bool
    {
        return $this->allows($admin, 'activity.view');
    }

    public function create(?Admin $admin): bool
    {
        return false;
    }

    public function update(?Admin $admin, Activity $activity): bool
    {
        return false;
    }

    public function delete(?Admin $admin, Activity $activity): bool
    {
        return false;
    }
}
