<?php

namespace App\Controllers;

use Net\Annotation\Controller;
use Net\Annotation\Route;
use Net\Annotation\GetMapping;
use Net\Annotation\PostMapping;
use Net\Annotation\PutMapping;
use Net\Annotation\DeleteMapping;
use Database\DB;
use App\Models\ToolLink;
use Net\Http\Request;
use Net\Http\Response;

#[Controller]
#[Route(prefix: '/api/tools')]
class ToolController
{
    #[GetMapping(path: '/')]
    public function lists(Request $request, Response $response): void
    {
        $result = DB::bind(ToolLink::class)->orderBy('display_order ASC')->get();
        $response->header('Content-Type', 'application/json; charset=utf-8');
        $response->write(json_encode($result));
    }

    #[PostMapping(path: '/')]
    public function create(Request $request, Response $response): void
    {
        $body = $request->body();
        $data = json_decode($body);

        $tool = new ToolLink();
        $tool->name = $data->name;
        $tool->url = $data->url;
        $tool->icon = $data->icon ?? '';
        $tool->category = $data->category ?? '';
        $tool->description = $data->description ?? '';
        $tool->isFavorite = $data->isFavorite ?? 0;
        $tool->displayOrder = $data->displayOrder ?? 0;

        $result = DB::bind(ToolLink::class)->insert($tool);

        $response->header('Content-Type', 'application/json; charset=utf-8');
        $response->write(json_encode(['success' => true, 'id' => $result->insertId]));
    }

    #[PutMapping(path: '/{id}')]
    public function update(Request $request, Response $response): void
    {
        $id = $request->pathValue('id');
        if ($id === null) {
            $path = $request->path();
            $pathParts = explode('/', $path);
            $id = $pathParts[count($pathParts) - 1];
        }

        $body = $request->body();
        $data = json_decode($body);

        $tool = new ToolLink();
        $tool->name = $data->name;
        $tool->url = $data->url;
        $tool->icon = $data->icon ?? '';
        $tool->category = $data->category ?? '';
        $tool->description = $data->description ?? '';
        $tool->isFavorite = $data->isFavorite ?? 0;
        $tool->displayOrder = $data->displayOrder ?? 0;

        $result = DB::bind(ToolLink::class)->where('id = ?', (int) $id)->update($tool);

        $response->header('Content-Type', 'application/json; charset=utf-8');
        $response->write(json_encode(['success' => true, 'rowsAffected' => $result]));
    }

    #[DeleteMapping(path: '/{id}')]
    public function delete(Request $request, Response $response): void
    {
        $id = $request->pathValue('id');
        if ($id === null) {
            $path = $request->path();
            $pathParts = explode('/', $path);
            $id = $pathParts[count($pathParts) - 1];
        }

        $result = DB::bind(ToolLink::class)->where('id = ?', (int) $id)->delete();

        $response->header('Content-Type', 'application/json; charset=utf-8');
        $response->write(json_encode(['success' => true, 'rowsAffected' => $result]));
    }
}
