<?php

use Illuminate\Support\Facades\Route;

Route::get('/', function () {
    return view('welcome');
});

Route::get('/origami-health', fn () => 'OK');

Route::redirect('/login', '/admin/login');
