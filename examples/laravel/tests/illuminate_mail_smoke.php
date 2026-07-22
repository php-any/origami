<?php

require dirname(__DIR__) . '/vendor/autoload.php';

use Illuminate\Mail\Mailables\Address;
use Illuminate\Mail\Transport\ArrayTransport;

/**
 * 不调用 Transport::send（Symfony Envelope / Mime 构造在 Origami 上有缺口）。
 * 验证 Mailables\Address 与 ArrayTransport 本地收集器。
 */
$address = new Address('user@example.com', 'Origami');
if ($address->address !== 'user@example.com' || $address->name !== 'Origami') {
    echo "FAIL: Address\n";
    exit(1);
}

$transport = new ArrayTransport();
if ($transport->messages()->count() !== 0) {
    echo "FAIL: initial message count\n";
    exit(1);
}

$transport->flush();
if ($transport->messages()->count() !== 0) {
    echo "FAIL: flush\n";
    exit(1);
}

echo "PASS\n";
