<?php

class Target {
    public string $name;
}

function app() {
    return new Target();
}

app()->name;
