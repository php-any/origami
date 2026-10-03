<?php
class RequestClosureReceiver {
    public int $count = 0;
    public function callback() { return function () { return ++$this->count; }; }
}
return (new RequestClosureReceiver())->callback();
