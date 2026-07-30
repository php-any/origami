<?php

class ArmkBox
{
    public function headers(): array
    {
        return ['host' => 'localhost'];
    }

    public function payload(): array
    {
        return [];
    }

    public function response(): string
    {
        return 'HTML Response';
    }

    public function build(): array
    {
        return [
            'ip_address' => '127.0.0.1',
            'uri' => '/x',
            'method' => 'GET',
            'controller_action' => '',
            'middleware' => [],
            'headers' => $this->headers(),
            'payload' => $this->payload(),
            'response_status' => 200,
            'response' => $this->response(),
            'session' => [],
            'duration' => 1,
            'memory' => 1.0,
        ];
    }
}

$data = (new ArmkBox())->build();
if (($data['headers']['host'] ?? '') !== 'localhost') {
    Log::fatal('fail');
}
Log::info('array return many keys ok');
