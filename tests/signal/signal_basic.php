<?php

echo "SIGINT=" . SIGINT . "\n";
echo "SIGTERM=" . SIGTERM . "\n";

$ch = new Signal\Channel();
echo "channel created\n";

Signal\notify($ch, SIGINT, SIGTERM);
echo "notify ok\n";

Signal\stop($ch);
$ch->close();
echo "stop and close ok\n";

echo "done\n";
