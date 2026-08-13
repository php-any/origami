<?php

namespace App\Command;

use Cli\Annotation\Command;

#[Command(name: "all", description: "运行全部测试")]
class AllCommand
{
    public function execute(): void
    {
        (new ChatCommand())->execute();
        (new JsonCommand())->execute();
        (new SchemaCommand())->execute();
        (new ErrorCommand())->execute();
        echo "\n🎉 全部测试完成\n";
    }
}
