<?php
class PreviousHeaders {
    public function __construct(private ?string $referer) {}
    public function get(string $key): ?string {
        return $key === 'referer' ? $this->referer : null;
    }
}
class PreviousRequest {
    public PreviousHeaders $headers;
    public function __construct(?string $referer) { $this->headers = new PreviousHeaders($referer); }
    public function root(): string { return 'http://example.com'; }
}
$url = new Illuminate\Routing\UrlGenerator(null, new PreviousRequest('http://example.com/list'));
if ($url->previous() !== 'http://example.com/list') {
    throw new Exception('Previous URL referer missing');
}
$url->setRequest(new PreviousRequest(null));
if ($url->previous('/fallback') !== 'http://example.com/fallback'
    || $url->previous() !== 'http://example.com') {
    throw new Exception('Previous URL fallback failed');
}
echo "URL previous context OK\n";
