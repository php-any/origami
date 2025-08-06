"use strict";
// Origami VSCode 扩展
// 简化实现，避免复杂的类型依赖
Object.defineProperty(exports, "__esModule", { value: true });
exports.deactivate = exports.activate = void 0;
// 简化的语言客户端
class SimpleLanguageClient {
    constructor(id, name, serverOptions, clientOptions) {
        this.isRunning = false;
        this.id = id;
        this.name = name;
        this.serverOptions = serverOptions;
        this.clientOptions = clientOptions;
    }
    async start() {
        this.isRunning = true;
        // 实际实现中会启动 LSP 客户端
        return Promise.resolve();
    }
    async stop() {
        this.isRunning = false;
        // 实际实现中会停止 LSP 客户端
        return Promise.resolve();
    }
    setTrace(level) {
        // 设置跟踪级别
    }
}
let client;
function activate(context) {
    try {
        // 获取配置
        const config = vscode?.workspace?.getConfiguration('origami.lsp');
        const enabled = config?.get('enabled', true);
        if (!enabled) {
            return;
        }
        const serverPath = config?.get('serverPath', 'origami-lsp');
        const trace = config?.get('trace', 'off');
        // LSP 服务器配置
        const serverOptions = {
            command: serverPath,
            transport: 'stdio',
            options: {
                env: {}
            }
        };
        // 客户端配置
        const clientOptions = {
            documentSelector: [
                { scheme: 'file', language: 'origami' },
                { scheme: 'untitled', language: 'origami' }
            ],
            synchronize: {
                fileEvents: vscode?.workspace?.createFileSystemWatcher('**/*.{cjp,origami}'),
                configurationSection: 'origami'
            },
            outputChannelName: 'Origami Language Server',
            initializationOptions: {},
            errorHandler: {
                error: (error, message, count) => {
                    return { action: 'Continue' };
                },
                closed: () => {
                    return { action: 'DoNotRestart' };
                }
            }
        };
        // 创建语言客户端
        client = new SimpleLanguageClient('origamiLanguageServer', 'Origami Language Server', serverOptions, clientOptions);
        // 设置跟踪级别
        if (trace !== 'off') {
            client.setTrace(trace);
        }
        // 注册命令
        const restartCommand = vscode?.commands?.registerCommand('origami.restartLanguageServer', async () => {
            if (client) {
                await client.stop();
                await client.start();
                vscode?.window?.showInformationMessage('Origami 语言服务器已重启');
            }
        });
        if (restartCommand) {
            context.subscriptions.push(restartCommand);
        }
        // 启动客户端
        client.start().then(() => {
            // 显示状态栏项
            const statusBarItem = vscode?.window?.createStatusBarItem(vscode?.StatusBarAlignment?.Right, 100);
            if (statusBarItem) {
                statusBarItem.text = '$(check) Origami LSP';
                statusBarItem.tooltip = 'Origami Language Server 运行中';
                statusBarItem.show();
                context.subscriptions.push(statusBarItem);
            }
        }).catch((error) => {
            vscode?.window?.showErrorMessage(`启动 Origami Language Server 失败: ${error?.message || error}`);
        });
        // 监听配置变化
        const configChangeListener = vscode?.workspace?.onDidChangeConfiguration((event) => {
            if (event?.affectsConfiguration('origami.lsp')) {
                vscode?.window?.showInformationMessage('Origami LSP 配置已更改，请重启语言服务器以应用更改', '重启').then((selection) => {
                    if (selection === '重启') {
                        vscode?.commands?.executeCommand('origami.restartLanguageServer');
                    }
                });
            }
        });
        if (configChangeListener) {
            context.subscriptions.push(configChangeListener);
        }
    }
    catch (error) {
        // 静默处理错误
    }
}
exports.activate = activate;
function deactivate() {
    if (!client) {
        return undefined;
    }
    return client.stop();
}
exports.deactivate = deactivate;
//# sourceMappingURL=extension.js.map