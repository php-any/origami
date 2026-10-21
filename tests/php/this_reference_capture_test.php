<?php
class ThisReferenceCapture {
    public function run() {
        $target =& $this;
        $other =& $this;
        $closure = static function () use (&$target) { $target = null; };
        $closure();
        if ($target !== null || $other !== $this || !($this instanceof ThisReferenceCapture)) {
            throw new Exception('Receiver reference rebinding failed');
        }
    }
}
(new ThisReferenceCapture())->run();
echo "this reference capture passed\n";
