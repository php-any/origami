<?php

namespace HotReload\View;

class GreetingView {
    public function render(array $data): string {
        return require __DIR__.'/../../views/greeting.php';
    }
}
