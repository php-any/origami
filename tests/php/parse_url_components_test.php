<?php
$parts = parse_url('https://user:pass@127.0.0.1:8080/a%20b?q=one#tag');
if (!is_array($parts) || $parts['host'] !== '127.0.0.1' || $parts['port'] !== 8080 || $parts['path'] !== '/a%20b') { throw new Exception('URL components'); }
if (parse_url('https://[::1]:8443/a', 1) !== '[::1]' || parse_url('https://example.test:8443', 2) !== 8443) { throw new Exception('URL component types'); }
if (parse_url('https://example.test/?#')['query'] !== '' || parse_url('https://example.test/?#')['fragment'] !== '') { throw new Exception('empty URL components'); }
if (parse_url('https://example.test:65536') !== false) { throw new Exception('invalid port'); }
echo "URL components OK\n";
