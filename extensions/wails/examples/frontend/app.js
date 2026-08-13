// 前端目录示例的入口脚本 (ES 模块)。
// 通过 Wails v3 运行时与 Origami PHP 后端进行事件通信：
//   - 前端 Emit("todo:add" | "todo:toggle" | "todo:delete" | "todo:clear" | "todo:list")
//   - 后端推送 "todo:changed"，携带完整任务列表
import { Events } from "/wails/runtime.js";

const input = document.getElementById("input");
const addBtn = document.getElementById("add");
const clearBtn = document.getElementById("clear");
const list = document.getElementById("list");
const count = document.getElementById("count");

function add() {
  const text = input.value.trim();
  if (!text) return;
  Events.Emit("todo:add", text);
  input.value = "";
  input.focus();
}

addBtn.addEventListener("click", add);
input.addEventListener("keydown", (e) => { if (e.key === "Enter") add(); });
clearBtn.addEventListener("click", () => Events.Emit("todo:clear", null));

// 后端推送最新任务列表 → 重新渲染
Events.On("todo:changed", (ev) => render(ev.data || []));

function render(todos) {
  list.innerHTML = "";
  if (todos.length === 0) {
    const li = document.createElement("li");
    li.className = "empty";
    li.textContent = "暂无任务，添加一个吧 ✨";
    list.appendChild(li);
  } else {
    for (const t of todos) {
      const li = document.createElement("li");
      if (t.done) li.classList.add("done");

      const cb = document.createElement("input");
      cb.type = "checkbox";
      cb.checked = !!t.done;
      cb.addEventListener("change", () => Events.Emit("todo:toggle", t.id));

      const span = document.createElement("span");
      span.className = "text";
      span.textContent = t.text;

      const del = document.createElement("button");
      del.className = "del";
      del.textContent = "×";
      del.title = "删除";
      del.addEventListener("click", () => Events.Emit("todo:delete", t.id));

      li.append(cb, span, del);
      list.appendChild(li);
    }
  }
  const remaining = todos.filter((t) => !t.done).length;
  count.textContent = `${todos.length} 项任务 · ${remaining} 项未完成`;
}

// 页面就绪后向后端请求一次当前列表
window.addEventListener("DOMContentLoaded", () => Events.Emit("todo:list", null));
