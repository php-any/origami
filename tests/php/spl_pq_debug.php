<?php
$pq = new SplPriorityQueue();
$pq->insert('low', 1);
$pq->insert('high', 10);
$pq->insert('mid', 5);
echo "extract1: " . $pq->extract() . "\n";
echo "extract2: " . $pq->extract() . "\n";
echo "extract3: " . $pq->extract() . "\n";
