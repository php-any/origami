<?php

namespace App\Controllers;

use Net\Annotation\Controller;
use Net\Annotation\GetMapping;
use Database\DB;
use App\Models\Project;
use App\Models\ProjectEnvironment;
use App\Models\ProjectTool;
use App\Models\ToolLink;
use App\Models\SearchEngine;

#[Controller]
class HomeController
{
    #[GetMapping(path: '/')]
    public function index($request, $response): void
    {
        // 加载工具（用于首页工具区及收藏）
        $tools = DB::bind(ToolLink::class)->orderBy('display_order ASC')->get();
        $favoriteTools = [];
        foreach ($tools as $t) {
            if (($t->isFavorite ?? 0) == 1) {
                $favoriteTools[] = $t;
            }
        }

        // 加载项目并挂载环境与工具
        $projects = DB::bind(Project::class)->orderBy('display_order ASC')->get();
        foreach ($projects as $p) {
            $envs = DB::bind(ProjectEnvironment::class)->where('project_id = ?', $p->id)
                ->orderBy('display_order ASC')->get();
            $p->environments = $envs;

            $projectTools = DB::bind(ProjectTool::class)->where('project_id = ?', $p->id)
                ->orderBy('display_order ASC')->get();
            $pTools = [];
            foreach ($projectTools as $pt) {
                $ts = DB::bind(ToolLink::class)->where('id = ?', $pt->toolId)->get();
                if ($ts !== null && count($ts) > 0) {
                    $pTools[] = $ts[0];
                }
            }
            $p->tools = $pTools;
        }

        // 加载搜索引擎
        $searchEngines = DB::bind(SearchEngine::class)->orderBy('display_order ASC, id ASC')->get();
        $defaultSearchEngine = null;
        foreach ($searchEngines as $engine) {
            if (($engine->isDefault ?? 0) == 1) {
                $defaultSearchEngine = $engine;
                break;
            }
        }
        if ($defaultSearchEngine === null && $searchEngines !== null && count($searchEngines) > 0) {
            $defaultSearchEngine = $searchEngines[0];
        }

        $response->view('./src/views/index.html', [
            'projects' => $projects,
            'tools' => $tools,
            'favoriteTools' => $favoriteTools,
            'searchEngines' => $searchEngines,
            'defaultSearchEngine' => $defaultSearchEngine,
        ]);
    }

    #[GetMapping(path: '/admin')]
    public function admin($request, $response): void
    {
        // 加载工具
        $tools = DB::bind(ToolLink::class)->orderBy('display_order ASC')->get();

        // 加载项目并挂载环境与工具
        $projects = DB::bind(Project::class)->orderBy('display_order ASC')->get();
        foreach ($projects as $p) {
            $envs = DB::bind(ProjectEnvironment::class)->where('project_id = ?', $p->id)
                ->orderBy('display_order ASC')->get();
            $p->environments = $envs;

            $projectTools = DB::bind(ProjectTool::class)->where('project_id = ?', $p->id)
                ->orderBy('display_order ASC')->get();
            $pTools = [];
            foreach ($projectTools as $pt) {
                $ts = DB::bind(ToolLink::class)->where('id = ?', $pt->toolId)->get();
                if ($ts !== null && count($ts) > 0) {
                    $pTools[] = $ts[0];
                }
            }
            $p->tools = $pTools;
        }

        // 加载搜索引擎
        $searchEngines = DB::bind(SearchEngine::class)->orderBy('display_order ASC, id ASC')->get();

        $response->view('./src/views/admin.html', [
            'tools' => $tools,
            'projects' => $projects,
            'searchEngines' => $searchEngines,
        ]);
    }
}
