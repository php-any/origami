<?php
$body = 'item=first&item=last&nested[value]=7';
$base = Symfony\Component\HttpFoundation\Request::create(
    'http://localhost/input', 'POST', ['item' => 'last'], [], [],
    ['CONTENT_TYPE' => 'application/x-www-form-urlencoded'], $body
);
$request = Illuminate\Http\Request::createFromBase($base);
if ($request->getContent() !== $body || $request->getContent() !== $body) {
    throw new Exception('Illuminate request must preserve the Symfony content cache');
}
if ($base->getContent() !== $body || $request->input('item') !== 'last') {
    throw new Exception('request conversion must preserve content and input');
}
echo "native request body cache: PASS\n";
