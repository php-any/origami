<?php
namespace tests\php;

class ArrowVar_Handler {
    public function handleException($e) {
        Log::info('got: ' . get_class($e) . ' msg=' . $e->getMessage());
    }
    protected function forwardsTo($method) {
        return fn (...$arguments) => $this->{$method}(...$arguments);
    }
    public function run() {
        $cb = $this->forwardsTo('handleException');
        $cb(new \Exception('hello-arrow'));
    }
}

(new ArrowVar_Handler())->run();
Log::info('arrow_variadic 测试通过');
