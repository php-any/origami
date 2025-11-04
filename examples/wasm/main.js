(() => {
  const resultEl = () => document.getElementById("result");
  const codeEl = () => document.getElementById("code");
  const selectEl = () => document.getElementById("exampleSelect");

  // 示例集合：名称到文件路径
  const EXAMPLE_FILES = {
    "Hello + add": "./snippets/hello.zy",
    "Array map/filter": "./snippets/arrays.zy",
    "Function + return": "./snippets/functions.zy",
    "String operations": "./snippets/strings.zy",
    "Control structures": "./snippets/control.zy",
  };

  async function bootWasm() {
    if (!window.WebAssembly) {
      resultEl().textContent = "Error: 当前浏览器不支持 WebAssembly";
      return;
    }
    const go = new Go();
    const resp = await fetch("main.wasm");
    const bytes = await resp.arrayBuffer();
    const mod = await WebAssembly.instantiate(bytes, go.importObject);
    go.run(mod.instance);
  }

  function bindUI() {
    const noop = () => {};

    // 填充示例下拉
    const sel = selectEl();
    if (sel) {
      Object.keys(EXAMPLE_FILES).forEach((k) => {
        const opt = document.createElement("option");
        opt.value = k;
        opt.textContent = k;
        sel.appendChild(opt);
      });
      sel.addEventListener("change", () => {
        const v = sel.value;
        if (v && v !== "__current__") {
          fetch(EXAMPLE_FILES[v])
            .then((r) => r.text())
            .then((t) => {
              codeEl().value = t;
              resultEl().textContent = "";
              noop();
            })
            .catch((e) => {
              resultEl().textContent = "Error: 加载示例失败 - " + e;
            });
        }
      });
    }

    // 将控制台输出同步到右侧面板（一次性挂载）
    (function attachConsoleOnce() {
      if (window.__origamiConsolePatched) return;
      window.__origamiConsolePatched = true;
      const append = (msg) => {
        const el = resultEl();
        if (!el) return;
        el.textContent += (msg == null ? "" : String(msg)) + "\n";
        el.scrollTop = el.scrollHeight;
      };
      const origLog = console.log.bind(console);
      const origErr = console.error.bind(console);
      console.log = (...args) => {
        origLog(...args);
        try {
          append(args.join(" "));
        } catch (e) {}
      };
      console.error = (...args) => {
        origErr(...args);
        try {
          append("[err] " + args.join(" "));
        } catch (e) {}
      };
    })();

    const run = async () => {
      const code = codeEl().value;
      try {
        resultEl().textContent = "";
        // 为兼容 wasm_exec.js 行缓冲，追加换行触发刷新
        const codeToRun = code + "\n" + 'echo "\n";';
        const out = window.origamiRun
          ? window.origamiRun(codeToRun)
          : "Error: WASM 未就绪";
        if (out && String(out).trim() && String(out) !== "--") {
          const el = resultEl();
          el.textContent += String(out);
        }
        noop();
      } catch (e) {
        resultEl().textContent = "Error: " + (e && e.message ? e.message : e);
      }
    };

    document.getElementById("runBtn").addEventListener("click", run);
    document.getElementById("clearBtn").addEventListener("click", () => {
      resultEl().textContent = "";
      console.clear();
    });
    document.addEventListener("keydown", (e) => {
      if ((e.metaKey || e.ctrlKey) && e.key === "Enter") {
        e.preventDefault();
        run();
      }
    });

    // 取消高亮相关逻辑
    codeEl().addEventListener("input", () => {});
  }

  window.addEventListener("DOMContentLoaded", async () => {
    bindUI();
    await bootWasm();
  });
})();
