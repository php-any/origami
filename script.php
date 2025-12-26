<?php

class Test {
        public function notify(Started $event): void
        {
            $subscribers = [
                new class($printer)
                {
                    public function notify(Configured $event): void
                    {
                        $this->printer()->setDecorated(
                            $event->configuration()->colors()
                        );
                    }
                },
                new class($printer)
                {
                    public function notify(Configured $event): void
                    {
                        $this->printer()->setDecorated(
                            $event->configuration()->colors()
                        );
                    }
                },
            ];
        }
}
$data = [
    1,2,
    3,4
];