<?php

register_shutdown_function(function () {
    echo "SHUTDOWN_OK\n";
});

echo "DONE\n";
