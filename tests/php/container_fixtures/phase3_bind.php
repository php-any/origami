<?php

namespace tests\php\container_fixtures;

interface Phase3Bind_CacheInterface {
    public function get(string $key): ?string;
}

#[\Container\Bind(abstract: Phase3Bind_CacheInterface::class)]
class Phase3Bind_AnnotatedCache implements Phase3Bind_CacheInterface {
    public function get(string $key): ?string {
        return 'annotated:' . $key;
    }
}
