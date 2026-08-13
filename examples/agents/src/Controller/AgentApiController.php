<?php

namespace App\Controller;

use App\Service\AgentRunnerService;
use Net\Annotation\Controller;
use Net\Annotation\GetMapping;
use Net\Annotation\PostMapping;
use Net\Annotation\Route;
use Net\Http\Request;
use Net\Http\Response;

#[Controller]
#[Route(prefix: "/api")]
class AgentApiController
{
    public function __construct(
        private AgentRunnerService $runner,
    ) {}

    #[GetMapping(path: "/workflows")]
    public function workflows(Response $response): void
    {
        $response->success([
            "model" => $this->runner->getModel(),
            "workflows" => $this->runner->listWorkflows(),
        ]);
    }

    #[PostMapping(path: "/agents/run")]
    public function run(Request $request, Response $response): void
    {
        $body = $request->body();
        $workflow = is_array($body) ? ($body["workflow"] ?? "pipeline") : "pipeline";
        $input = is_array($body) ? trim($body["input"] ?? "") : "";

        if ($input === "") {
            $response->error("input 不能为空", 400);
            return;
        }

        try {
            $result = $this->runner->run($workflow, $input);
            $response->success($result);
        } catch (\Exception $e) {
            $response->error($e->getMessage(), 500);
        }
    }

    #[GetMapping(path: "/dev/status")]
    public function devStatus(Response $response): void
    {
        $version = getenv("ORIGAMI_DEV_VERSION") ?: "0";
        $response->success([
            "version" => $version,
            "hotReload" => true,
        ]);
    }
}
