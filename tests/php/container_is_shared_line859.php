<?php
class MiniContainerLine859 {
    protected array $abstractAliases = [];

    public function getContextualConcrete(string $abstract): mixed {
        if (empty($this->abstractAliases[$abstract])) {
            return null;
        }
        return 'found';
    }
}

$c = new MiniContainerLine859();
echo $c->getContextualConcrete('App\\Http\\Kernel') === null ? "ok\n" : "fail\n";
