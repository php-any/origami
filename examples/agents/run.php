<?php

/**
 * CLI 入口
 *
 * 用法:
 *   agents run.php pipeline
 *   agents run.php debate
 *   agents run.php handoff
 */

require __DIR__ . "/autoload.php";

// 加载 #[CliApplication] 引导类，触发扫描并注册 cli/Command 下的所有命令
require __DIR__ . "/cli/AgentApp.php";
