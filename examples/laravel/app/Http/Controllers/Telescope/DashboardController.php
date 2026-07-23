<?php

namespace App\Http\Controllers\Telescope;

use Laravel\Telescope\Telescope;
use Net\Http\Response;

/**
 * Telescope SPA shell（绕过 Blade：直接输出 layout HTML）。
 */
class DashboardController
{
    public function index(Response $response): void
    {
        bootstrap_telescope_http();

        $appName = (string) config('app.name', 'Origami');
        $cssFile = Telescope::$useDarkTheme ? 'app-dark.css' : 'app.css';
        $vars = [
            'path' => (string) config('telescope.path', 'telescope'),
            'timezone' => (string) config('app.timezone', 'UTC'),
            'recording' => !cache('telescope:pause-recording'),
        ];
        $json = json_encode($vars, JSON_UNESCAPED_UNICODE | JSON_UNESCAPED_SLASHES);

        $html = <<<HTML
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">
    <meta name="csrf-token" content="">
    <link rel="shortcut icon" href="/vendor/telescope/favicon.ico">
    <meta name="robots" content="noindex, nofollow">
    <title>Telescope - {$appName}</title>
    <link rel="preconnect" href="https://fonts.bunny.net">
    <link href="https://fonts.bunny.net/css?family=figtree:300,400,500,600" rel="stylesheet" />
    <link href="/vendor/telescope/{$cssFile}" rel="stylesheet" type="text/css">
</head>
<body>
<div id="telescope" v-cloak>
    <alert :message="alert.message"
           :type="alert.type"
           :auto-close="alert.autoClose"
           :confirmation-proceed="alert.confirmationProceed"
           :confirmation-cancel="alert.confirmationCancel"
           v-if="alert.type"></alert>

    <div class="container mb-5">
        <div class="d-flex align-items-stretch py-4 header">
            <router-link to="/" class="logo d-flex align-items-center">
                <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 80 80">
                    <path class="fill-primary" d="M0 40a39.87 39.87 0 0 1 11.72-28.28A40 40 0 1 1 0 40zm34 10a4 4 0 0 1-4-4v-2a2 2 0 1 0-4 0v2a4 4 0 0 1-4 4h-2a2 2 0 1 0 0 4h2a4 4 0 0 1 4 4v2a2 2 0 1 0 4 0v-2a4 4 0 0 1 4-4h2a2 2 0 1 0 0-4h-2zm24-24a6 6 0 0 1-6-6v-3a3 3 0 0 0-6 0v3a6 6 0 0 1-6 6h-3a3 3 0 0 0 0 6h3a6 6 0 0 1 6 6v3a3 3 0 0 0 6 0v-3a6 6 0 0 1 6-6h3a3 3 0 0 0 0-6h-3zm-4 36a4 4 0 1 0 0-8 4 4 0 0 0 0 8zM21 28a3 3 0 1 0 0-6 3 3 0 0 0 0 6z"></path>
                </svg>
                <h4 class="mb-0 ml-3"><strong>Laravel</strong> Telescope - {$appName}</h4>
            </router-link>

            <button class="btn btn-muted ml-auto mr-3 d-flex align-items-center py-2" v-on:click.prevent="toggleRecording" :title="recording ? 'Pause recording' : 'Resume recording'">
                <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" class="icon" fill="currentColor" v-if="recording">
                    <path d="M5.75 3a.75.75 0 00-.75.75v12.5c0 .414.336.75.75.75h1.5a.75.75 0 00.75-.75V3.75A.75.75 0 007.25 3h-1.5zM12.75 3a.75.75 0 00-.75.75v12.5c0 .414.336.75.75.75h1.5a.75.75 0 00.75-.75V3.75a.75.75 0 00-.75-.75h-1.5z" />
                </svg>
                <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" class="icon" fill="currentColor" v-else>
                    <path d="M6.3 2.841A1.5 1.5 0 004 4.11V15.89a1.5 1.5 0 002.3 1.269l9.344-5.89a1.5 1.5 0 000-2.538L6.3 2.84z" />
                </svg>
            </button>

            <button class="btn btn-muted mr-3 d-flex align-items-center py-2" v-on:click.prevent="clearEntries" title="Clear entries">
                <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" class="icon" fill="currentColor">
                    <path fill-rule="evenodd" d="M8.75 1A2.75 2.75 0 006 3.75v.443c-.795.077-1.584.176-2.365.298a.75.75 0 10.23 1.482l.149-.022.841 10.518A2.75 2.75 0 007.596 19h4.807a2.75 2.75 0 002.742-2.53l.841-10.52.149.023a.75.75 0 00.23-1.482A41.03 41.03 0 0014 4.193V3.75A2.75 2.75 0 0011.25 1h-2.5zM10 4c.84 0 1.673.025 2.5.075V3.75c0-.69-.56-1.25-1.25-1.25h-2.5c-.69 0-1.25.56-1.25 1.25v.325C8.327 4.025 9.16 4 10 4zM8.58 7.72a.75.75 0 00-1.5.06l.3 7.5a.75.75 0 101.5-.06l-.3-7.5zm4.34.06a.75.75 0 10-1.5-.06l-.3 7.5a.75.75 0 101.5.06l.3-7.5z" clip-rule="evenodd" />
                </svg>
            </button>

            <button class="btn btn-muted mr-3 d-flex align-items-center py-2" :class="{active: autoLoadsNewEntries}" v-on:click.prevent="autoLoadNewEntries" title="Auto load entries">
                <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" class="icon" fill="currentColor">
                    <path fill-rule="evenodd" d="M15.312 11.424a5.5 5.5 0 01-9.201 2.466l-.312-.311h2.433a.75.75 0 000-1.5H3.989a.75.75 0 00-.75.75v4.242a.75.75 0 001.5 0v-2.43l.31.31a7 7 0 0011.712-3.138.75.75 0 00-1.449-.39zm1.23-3.723a.75.75 0 00.219-.53V2.929a.75.75 0 00-1.5 0V5.36l-.31-.31A7 7 0 003.239 8.188a.75.75 0 101.448.389A5.5 5.5 0 0113.89 6.11l.311.31h-2.432a.75.75 0 000 1.5h4.243a.75.75 0 00.53-.219z" clip-rule="evenodd" />
                </svg>
            </button>

            <router-link to="/monitored-tags" class="btn btn-muted d-flex align-items-center py-2" title="Monitoring">
                <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" class="icon" fill="currentColor">
                    <path d="M10 12.5a2.5 2.5 0 100-5 2.5 2.5 0 000 5z" />
                    <path fill-rule="evenodd" d="M.664 10.59a1.651 1.651 0 010-1.186A10.004 10.004 0 0110 3c4.257 0 7.893 2.66 9.336 6.41.147.381.146.804 0 1.186A10.004 10.004 0 0110 17c-4.257 0-7.893-2.66-9.336-6.41zM14 10a4 4 0 11-8 0 4 4 0 018 0z" clip-rule="evenodd" />
                </svg>
            </router-link>
        </div>

        <div class="row mt-4">
            <div class="col-2 sidebar">
                <ul class="nav flex-column">
                    <li class="nav-item">
                        <router-link active-class="active" to="/requests" class="nav-link d-flex align-items-center"><span>Requests</span></router-link>
                    </li>
                    <li class="nav-item">
                        <router-link active-class="active" to="/commands" class="nav-link d-flex align-items-center"><span>Commands</span></router-link>
                    </li>
                    <li class="nav-item">
                        <router-link active-class="active" to="/schedule" class="nav-link d-flex align-items-center"><span>Schedule</span></router-link>
                    </li>
                    <li class="nav-item">
                        <router-link active-class="active" to="/jobs" class="nav-link d-flex align-items-center"><span>Jobs</span></router-link>
                    </li>
                    <li class="nav-item mt-3">
                        <router-link active-class="active" to="/batches" class="nav-link d-flex align-items-center"><span>Batches</span></router-link>
                    </li>
                    <li class="nav-item">
                        <router-link active-class="active" to="/cache" class="nav-link d-flex align-items-center"><span>Cache</span></router-link>
                    </li>
                    <li class="nav-item">
                        <router-link active-class="active" to="/dumps" class="nav-link d-flex align-items-center"><span>Dumps</span></router-link>
                    </li>
                    <li class="nav-item">
                        <router-link active-class="active" to="/events" class="nav-link d-flex align-items-center"><span>Events</span></router-link>
                    </li>
                    <li class="nav-item">
                        <router-link active-class="active" to="/exceptions" class="nav-link d-flex align-items-center"><span>Exceptions</span></router-link>
                    </li>
                    <li class="nav-item">
                        <router-link active-class="active" to="/gates" class="nav-link d-flex align-items-center"><span>Gates</span></router-link>
                    </li>
                    <li class="nav-item">
                        <router-link active-class="active" to="/client-requests" class="nav-link d-flex align-items-center"><span>HTTP Client</span></router-link>
                    </li>
                    <li class="nav-item">
                        <router-link active-class="active" to="/logs" class="nav-link d-flex align-items-center"><span>Logs</span></router-link>
                    </li>
                    <li class="nav-item">
                        <router-link active-class="active" to="/mail" class="nav-link d-flex align-items-center"><span>Mail</span></router-link>
                    </li>
                    <li class="nav-item">
                        <router-link active-class="active" to="/models" class="nav-link d-flex align-items-center"><span>Models</span></router-link>
                    </li>
                    <li class="nav-item">
                        <router-link active-class="active" to="/notifications" class="nav-link d-flex align-items-center"><span>Notifications</span></router-link>
                    </li>
                    <li class="nav-item">
                        <router-link active-class="active" to="/queries" class="nav-link d-flex align-items-center"><span>Queries</span></router-link>
                    </li>
                    <li class="nav-item">
                        <router-link active-class="active" to="/redis" class="nav-link d-flex align-items-center"><span>Redis</span></router-link>
                    </li>
                    <li class="nav-item">
                        <router-link active-class="active" to="/views" class="nav-link d-flex align-items-center"><span>Views</span></router-link>
                    </li>
                </ul>
            </div>
            <div class="col-10">
                <router-view></router-view>
            </div>
        </div>
    </div>
</div>
<script>
    window.Telescope = {$json};
</script>
<script src="/vendor/telescope/app.js"></script>
</body>
</html>
HTML;

        $response->html($html);
    }
}
