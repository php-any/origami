<?php

require __DIR__.'/vendor/autoload.php';

$app = require_once __DIR__.'/bootstrap/app.php';

class TestHandler {
    protected function context() {
        try {
            return array_filter([
                'userId' => Illuminate\Support\Facades\Auth::id(),
            ]);
        } catch (Throwable) {
            echo "context caught\n";
            return [];
        }
    }
    public function run() {
        return $this->context();
    }
}

$h = new TestHandler();
var_export($h->run());
