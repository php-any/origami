<?php

$provider = new Illuminate\Auth\EloquentUserProvider(app('hash'), App\Models\Admin::class);
return $provider->createModel();
