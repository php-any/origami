<?php
final class SpawnFrameProbe {
    public static function run(): void {
        $channel=new Channel(2);
        foreach ([20,40] as $value) {
            spawn(function () use ($channel,$value): void { $local=$value+2; $channel->send($local); });
        }
        $values=[$channel->receive(),$channel->receive()]; sort($values);
        if ($values !== [22,42]) { throw new RuntimeException('spawn closure frame lost captured/local variables'); }
        $channel->close();
    }
}
SpawnFrameProbe::run();
echo "spawn frame OK\n";
