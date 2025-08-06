// Origami VSCode 扩展
// 简化实现，避免复杂的类型依赖

// 基础类型定义
interface ExtensionContext {
  subscriptions: any[];
}

// 模拟 vscode API
declare const vscode: any;

// 简化的语言客户端
class SimpleLanguageClient {
  private id: string;
  private name: string;
  private serverOptions: any;
  private clientOptions: any;
  private isRunning: boolean = false;

  constructor(id: string, name: string, serverOptions: any, clientOptions: any) {
    this.id = id;
    this.name = name;
    this.serverOptions = serverOptions;
    this.clientOptions = clientOptions;
  }

  async start(): Promise<void> {
    this.isRunning = true;
    // 实际实现中会启动 LSP 客户端
    return Promise.resolve();
  }

  async stop(): Promise<void> {
    this.isRunning = false;
    // 实际实现中会停止 LSP 客户端
    return Promise.resolve();
  }

  setTrace(level: string): void {
    // 设置跟踪级别
  }
}

let client: SimpleLanguageClient | undefined;

export function activate(context: ExtensionContext) {
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
        error: (error: any, message: any, count: any) => {
          return { action: 'Continue' };
        },
        closed: () => {
          return { action: 'DoNotRestart' };
        }
      }
    };

    // 创建语言客户端
    client = new SimpleLanguageClient(
      'origamiLanguageServer',
      'Origami Language Server',
      serverOptions,
      clientOptions
    );

    // 设置跟踪级别
    if (trace !== 'off') {
      client.setTrace(trace);
    }

    // 注册命令
    const restartCommand = vscode?.commands?.registerCommand(
      'origami.restartLanguageServer',
      async () => {
        if (client) {
          await client.stop();
          await client.start();
          vscode?.window?.showInformationMessage('Origami 语言服务器已重启');
        }
      }
    );

    if (restartCommand) {
      context.subscriptions.push(restartCommand);
    }

    // 启动客户端
    client.start().then(() => {
      // 显示状态栏项
      const statusBarItem = vscode?.window?.createStatusBarItem(
        vscode?.StatusBarAlignment?.Right,
        100
      );
      if (statusBarItem) {
        statusBarItem.text = '$(check) Origami LSP';
        statusBarItem.tooltip = 'Origami Language Server 运行中';
        statusBarItem.show();
        context.subscriptions.push(statusBarItem);
      }
    }).catch((error: any) => {
      vscode?.window?.showErrorMessage(
        `启动 Origami Language Server 失败: ${error?.message || error}`
      );
    });

    // 监听配置变化
    const configChangeListener = vscode?.workspace?.onDidChangeConfiguration(
      (event: any) => {
        if (event?.affectsConfiguration('origami.lsp')) {
          vscode?.window?.showInformationMessage(
            'Origami LSP 配置已更改，请重启语言服务器以应用更改',
            '重启'
          ).then((selection: any) => {
            if (selection === '重启') {
              vscode?.commands?.executeCommand('origami.restartLanguageServer');
            }
          });
        }
      }
    );

    if (configChangeListener) {
      context.subscriptions.push(configChangeListener);
    }

  } catch (error) {
    // 静默处理错误
  }
}

export function deactivate(): Promise<void> | undefined {
  if (!client) {
    return undefined;
  }
  
  return client.stop();
}