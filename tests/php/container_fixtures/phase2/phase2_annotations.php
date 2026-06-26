<?php

namespace tests\php\container_fixtures;

interface Phase2Anno_LoggerInterface {}

#[\Container\Bind(abstract: Phase2Anno_LoggerInterface::class)]
class Phase2Anno_FileLogger implements Phase2Anno_LoggerInterface {
    public string $channel = 'file';
}

#[\Container\Singleton]
class Phase2Anno_Mailer {
    public string $tag = 'mailer';
}

#[\Container\Component(name: 'phase2.transient')]
class Phase2Anno_TransientSvc {
    public string $tag = 'transient';
}

#[\Container\Singleton]
class Phase2Anno_Config {
    private static int $seq = 0;
    public int $id;
    public function __construct() {
        $this->id = ++self::$seq;
    }
}

class Phase2Anno_OrderService {
    private Phase2Anno_Mailer $mailer;
    private Phase2Anno_LoggerInterface $logger;

    public function __construct(
        Phase2Anno_Mailer $mailer,
        Phase2Anno_LoggerInterface $logger,
    ) {
        $this->mailer = $mailer;
        $this->logger = $logger;
    }

    public function mailerTag(): string {
        return $this->mailer->tag;
    }
    public function loggerChannel(): string {
        return $this->logger->channel;
    }
}
