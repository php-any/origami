<?php

namespace App\Policies;

use App\Models\Admin;
use App\Models\Media;

class MediaPolicy extends AdminResourcePolicy
{
    public function viewAny(?Admin $admin): bool
    {
        return $this->allows($admin, 'media.view');
    }

    public function view(?Admin $admin, Media $media): bool
    {
        return $this->allows($admin, 'media.view');
    }

    public function create(?Admin $admin): bool
    {
        return $this->allows($admin, 'media.create');
    }

    public function update(?Admin $admin, Media $media): bool
    {
        return $this->allows($admin, 'media.edit');
    }

    public function delete(?Admin $admin, Media $media): bool
    {
        return $this->allows($admin, 'media.delete');
    }
}
