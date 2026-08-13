<?php

require_once __DIR__ . '/../../examples/laravel/vendor/autoload.php';

class ExtendsVendor_Child extends Illuminate\Notifications\Notification
{
    public function via(object $notifiable): array
    {
        return ['mail'];
    }
}

$n = new ExtendsVendor_Child();
if ($n->locale('en')->locale !== 'en') {
    echo "FAIL locale\n";
    exit(1);
}

echo "PASS\n";
