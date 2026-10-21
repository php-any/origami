<?php

namespace HotReload\Service;

class GreetingService {
    public function greet(string $name): string {
        return 'service-v1: Hello, '.$name;
    }
}
