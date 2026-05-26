<?php

namespace App\Controllers;

use Net\Annotation\Controller;
use Net\Annotation\Route;
use Net\Annotation\GetMapping;
use Net\Annotation\PostMapping;
use Net\Annotation\PutMapping;
use Net\Annotation\DeleteMapping;
use Database\DB;
use App\Models\Project;
use App\Models\ProjectEnvironment;
use App\Models\ProjectTool;
use App\Models\ToolLink;

#[Controller]
#[Route(prefix: '/api/projects')]
class ProjectController
{
    #[GetMapping(path: '/')]
    public function lists($request, $response): void
    {
        $projects = DB::bind(Project::class)->orderBy('display_order ASC')->get();
        $result = [];
        foreach ($projects as $project) {
            // 获取项目环境
            $envs = DB::bind(ProjectEnvironment::class)->where('project_id = ?', $project->id)
                ->orderBy('display_order ASC')->get();
            $envList = [];
            foreach ($envs as $env) {
                $envList[] = [
                    'id' => $env->id,
                    'environmentName' => $env->environmentName,
                    'url' => $env->url,
                    'status' => $env->status,
                    'statusColor' => $env->statusColor,
                    'displayOrder' => $env->displayOrder,
                ];
            }

            // 获取项目关联的工具
            $projectTools = DB::bind(ProjectTool::class)->where('project_id = ?', $project->id)
                ->orderBy('display_order ASC')->get();

            $toolList = [];
            foreach ($projectTools as $projectTool) {
                $tools = DB::bind(ToolLink::class)->where('id = ?', $projectTool->toolId)->get();
                if ($tools !== null && count($tools) > 0) {
                    $toolList[] = $tools[0];
                }
            }

            $result[] = [
                'id' => $project->id,
                'name' => $project->name,
                'description' => $project->description,
                'icon' => $project->icon,
                'displayOrder' => $project->displayOrder,
                'environments' => $envList,
                'tools' => $toolList,
            ];
        }

        $response->header('Content-Type', 'application/json; charset=utf-8');
        $response->write(json_encode($result));
    }

    #[PostMapping(path: '/')]
    public function create($request, $response): void
    {
        $body = $request->body();
        $data = json_decode($body);

        $project = new Project();
        $project->name = $data->name;
        $project->description = $data->description ?? null;
        $project->icon = $data->icon ?? null;
        $project->displayOrder = $data->displayOrder ?? 0;

        $result = DB::bind(Project::class)->insert($project);
        $projectId = $result->insertId;

        // 创建项目环境
        if ($data->environments !== null) {
            foreach ($data->environments as $env) {
                $envObj = new ProjectEnvironment();
                $envObj->projectId = $projectId;
                $envObj->environmentName = $env->environmentName;
                $envObj->url = $env->url;
                $envObj->status = $env->status ?? '运行中';
                $envObj->statusColor = $env->statusColor ?? 'green';
                $envObj->displayOrder = $env->displayOrder ?? 0;
                DB::bind(ProjectEnvironment::class)->insert($envObj);
            }
        }

        // 创建项目工具关联
        if ($data->tools !== null) {
            for ($i = 0; $i < count($data->tools); $i++) {
                $toolId = $data->tools[$i];
                $projectTool = new ProjectTool();
                $projectTool->projectId = $projectId;
                $projectTool->toolId = (int) $toolId;
                $projectTool->displayOrder = $i + 1;
                DB::bind(ProjectTool::class)->insert($projectTool);
            }
        }

        $response->header('Content-Type', 'application/json; charset=utf-8');
        $response->write(json_encode(['success' => true, 'id' => $projectId]));
    }

    #[PutMapping(path: '/{id}')]
    public function update($request, $response): void
    {
        $id = $request->pathValue('id');
        if ($id === null) {
            $path = $request->path();
            $pathParts = explode('/', $path);
            $id = $pathParts[count($pathParts) - 1];
        }

        $body = $request->body();
        $data = json_decode($body);

        $project = new Project();
        $project->name = $data->name;
        $project->description = $data->description ?? null;
        $project->icon = $data->icon ?? null;
        $project->displayOrder = $data->displayOrder ?? 0;

        $result = DB::bind(Project::class)->where('id = ?', (int) $id)->update($project);

        // 同步更新环境：先删除现有环境，再创建新环境
        DB::bind(ProjectEnvironment::class)->where('project_id = ?', (int) $id)->delete();
        if ($data->environments !== null && count($data->environments) > 0) {
            foreach ($data->environments as $env) {
                if ($env->environmentName !== null && $env->url !== null) {
                    $envObj = new ProjectEnvironment();
                    $envObj->projectId = (int) $id;
                    $envObj->environmentName = $env->environmentName;
                    $envObj->url = $env->url;
                    $envObj->status = $env->status ?? '运行中';
                    $envObj->statusColor = $env->statusColor ?? 'green';
                    $envObj->displayOrder = $env->displayOrder ?? 0;
                    DB::bind(ProjectEnvironment::class)->insert($envObj);
                }
            }
        }

        // 同步更新工具关联：先删除现有工具关联，再创建新关联
        DB::bind(ProjectTool::class)->where('project_id = ?', (int) $id)->delete();
        if ($data->tools !== null && count($data->tools) > 0) {
            for ($i = 0; $i < count($data->tools); $i++) {
                $toolId = $data->tools[$i];
                if ($toolId !== null) {
                    $projectTool = new ProjectTool();
                    $projectTool->projectId = (int) $id;
                    $projectTool->toolId = (int) $toolId;
                    $projectTool->displayOrder = $i + 1;
                    DB::bind(ProjectTool::class)->insert($projectTool);
                }
            }
        }

        $response->header('Content-Type', 'application/json; charset=utf-8');
        $response->write(json_encode(['success' => true, 'rowsAffected' => $result]));
    }

    #[DeleteMapping(path: '/{id}')]
    public function delete($request, $response): void
    {
        $id = $request->pathValue('id');
        if ($id === null) {
            $path = $request->path();
            $pathParts = explode('/', $path);
            $id = $pathParts[count($pathParts) - 1];
        }

        DB::bind(ProjectTool::class)->where('project_id = ?', (int) $id)->delete();
        DB::bind(ProjectEnvironment::class)->where('project_id = ?', (int) $id)->delete();
        $result = DB::bind(Project::class)->where('id = ?', (int) $id)->delete();

        $response->header('Content-Type', 'application/json; charset=utf-8');
        $response->write(json_encode(['success' => true, 'rowsAffected' => $result]));
    }

    #[PostMapping(path: '/{projectId}/environments')]
    public function createEnvironment($request, $response): void
    {
        $projectId = $request->pathValue('projectId');
        if ($projectId === null) {
            $path = $request->path();
            $pathParts = explode('/', $path);
            $projectId = $pathParts[count($pathParts) - 2];
        }

        $body = $request->body();
        $data = json_decode($body);

        $env = new ProjectEnvironment();
        $env->projectId = (int) $projectId;
        $env->environmentName = $data->environmentName;
        $env->url = $data->url;
        $env->status = $data->status ?? '运行中';
        $env->statusColor = $data->statusColor ?? 'green';
        $env->displayOrder = $data->displayOrder ?? 0;

        $result = DB::bind(ProjectEnvironment::class)->insert($env);

        $response->header('Content-Type', 'application/json; charset=utf-8');
        $response->write(json_encode(['success' => true, 'id' => $result->insertId]));
    }

    #[PutMapping(path: '/{projectId}/environments/{envId}')]
    public function updateEnvironment($request, $response): void
    {
        $projectId = $request->pathValue('projectId');
        $envId = $request->pathValue('envId');

        if ($envId === null || $projectId === null) {
            $path = $request->path();
            $pathParts = explode('/', $path);
            $envId = $pathParts[count($pathParts) - 1];
            $projectId = $pathParts[count($pathParts) - 3];
        }

        $body = $request->body();
        $data = json_decode($body);

        $env = new ProjectEnvironment();
        $env->environmentName = $data->environmentName;
        $env->url = $data->url;
        $env->status = $data->status ?? '运行中';
        $env->statusColor = $data->statusColor ?? 'green';
        $env->displayOrder = $data->displayOrder ?? 0;

        $result = DB::bind(ProjectEnvironment::class)->where('id = ? AND project_id = ?', (int) $envId, (int) $projectId)->update($env);

        $response->header('Content-Type', 'application/json; charset=utf-8');
        $response->write(json_encode(['success' => true, 'rowsAffected' => $result]));
    }

    #[DeleteMapping(path: '/{projectId}/environments/{envId}')]
    public function deleteEnvironment($request, $response): void
    {
        $projectId = $request->pathValue('projectId');
        $envId = $request->pathValue('envId');

        if ($envId === null || $projectId === null) {
            $path = $request->path();
            $pathParts = explode('/', $path);
            $envId = $pathParts[count($pathParts) - 1];
            $projectId = $pathParts[count($pathParts) - 3];
        }

        $result = DB::bind(ProjectEnvironment::class)->where('id = ? AND project_id = ?', (int) $envId, (int) $projectId)->delete();

        $response->header('Content-Type', 'application/json; charset=utf-8');
        $response->write(json_encode(['success' => true, 'rowsAffected' => $result]));
    }
}
