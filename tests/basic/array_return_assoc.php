<?php

class AraBox {
    public function build(): array
    {
        return [
            'ip_address' => '127.0.0.1',
            'uri' => '/x',
            'method' => 'GET',
        ];
    }
}

$box = new AraBox();
var_export($box->build());
echo "\nPASS\n";
