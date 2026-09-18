<?php

#[Deprecated(message: 'use new_fn()', since: 'test:1.0.0')]
function old_fn(): string
{
    return 'ok';
}

echo old_fn();
