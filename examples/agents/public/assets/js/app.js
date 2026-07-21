const workflowSelect = document.getElementById("workflow");
const inputEl = document.getElementById("input");
const runBtn = document.getElementById("run-btn");
const stepsEl = document.getElementById("steps");
const finalEl = document.getElementById("final");
const finalPre = finalEl.querySelector("pre");
const runStatus = document.getElementById("run-status");
const modelLabel = document.getElementById("model-label");
const hint = document.getElementById("hint");
const toast = document.getElementById("toast");

let devVersion = "0";

const defaultInputs = {
  pipeline: "折言（Origami）是一门用 Go 实现的 PHP 风格脚本语言",
  debate: "AI 辅助编程是否应该成为开发者的默认工作方式？",
  handoff: "为一个小型博客系统设计 REST API 端点列表",
};

async function loadWorkflows() {
  const res = await fetch("/api/workflows");
  const json = await res.json();
  const data = json.data || json;

  modelLabel.textContent = `模型: ${data.model || "-"}`;
  workflowSelect.innerHTML = "";

  (data.workflows || []).forEach((wf) => {
    const opt = document.createElement("option");
    opt.value = wf.id;
    opt.textContent = `${wf.name} — ${wf.description}`;
    workflowSelect.appendChild(opt);
  });

  updateHint();
}

function updateHint() {
  const id = workflowSelect.value;
  hint.textContent = defaultInputs[id] ? `示例: ${defaultInputs[id]}` : "";
  if (!inputEl.value) {
    inputEl.value = defaultInputs[id] || "";
  }
}

function renderSteps(steps) {
  stepsEl.innerHTML = "";
  (steps || []).forEach((step, idx) => {
    const card = document.createElement("div");
    card.className = "step";
    const title = step.agent || `Step ${idx + 1}`;
    const meta = [step.round, step.phase].filter(Boolean).join(" / ");
    card.innerHTML = `<h4>${title}${meta ? ` <span class="muted">(${meta})</span>` : ""}</h4><pre></pre>`;
    card.querySelector("pre").textContent = step.output || "";
    stepsEl.appendChild(card);
  });
}

function showToast(message) {
  toast.textContent = message;
  toast.classList.remove("hidden");
  setTimeout(() => toast.classList.add("hidden"), 3000);
}

async function runWorkflow() {
  const workflow = workflowSelect.value;
  const input = inputEl.value.trim();
  if (!input) {
    showToast("请输入内容");
    return;
  }

  runBtn.disabled = true;
  runStatus.textContent = "运行中...";
  stepsEl.innerHTML = "";
  finalEl.classList.add("hidden");

  try {
    const res = await fetch("/api/agents/run", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ workflow, input }),
    });
    const json = await res.json();
    if (!res.ok) {
      throw new Error(json.message || "请求失败");
    }

    const data = json.data || json;
    renderSteps(data.steps);
    finalPre.textContent = data.final || "";
    finalEl.classList.remove("hidden");
    runStatus.textContent = "完成";
  } catch (err) {
    runStatus.textContent = "失败";
    showToast(err.message || String(err));
  } finally {
    runBtn.disabled = false;
  }
}

async function pollDevStatus() {
  try {
    const res = await fetch("/api/dev/status");
    const json = await res.json();
    const data = json.data || json;
    if (data.version && data.version !== devVersion) {
      if (devVersion !== "0") {
        showToast("服务端已热重载，页面将刷新");
        setTimeout(() => location.reload(), 800);
      }
      devVersion = data.version;
    }
  } catch (_) {}
}

workflowSelect.addEventListener("change", updateHint);
runBtn.addEventListener("click", runWorkflow);

loadWorkflows();
setInterval(pollDevStatus, 2000);
