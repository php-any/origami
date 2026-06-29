// Origami Chat — 前端逻辑 (ES 模块)
// 通过 Wails v3 Events 与 PHP 后端双向通信
import { Events } from "/wails/runtime.js";

// ── DOM ──
const elChannels  = document.getElementById("channels");
const elUsers     = document.getElementById("users");
const elUserCount = document.getElementById("userCount");
const elMessages  = document.getElementById("messages");
const elInput     = document.getElementById("input");
const elNick      = document.getElementById("nick");
const elComposer  = document.getElementById("composer");
const elCurName   = document.getElementById("curChannelName");
const elCurTopic  = document.getElementById("curTopic");
const elToasts    = document.getElementById("toasts");

// ── 本地状态 ──
let nick     = "访客";
let channel  = "general";
let channels = [];
let users    = [];
let msgIds   = new Set(); // 去重

// ── 工具 ──
function normalizeMsg(raw) {
  if (!raw || typeof raw === "string") return null;
  // 兼容 Wails 事件包装
  const m = (raw.data && typeof raw.data === "object") ? raw.data : raw;
  if (!m || typeof m !== "object") return null;
  return {
    id:      m.id != null ? String(m.id) : `tmp-${Date.now()}-${Math.random().toString(36).slice(2, 7)}`,
    channel: m.channel ?? channel,
    author:  m.author ?? "?",
    text:    m.text ?? "",
    type:    m.type ?? "user",
    time:    m.time ?? "",
  };
}

function isNearBottom() {
  const gap = elMessages.scrollHeight - elMessages.scrollTop - elMessages.clientHeight;
  return gap < 80;
}
function now() {
  const d = new Date();
  return d.toLocaleTimeString("zh-CN", { hour: "2-digit", minute: "2-digit" });
}

function avatar(name, type) {
  if (type === "bot") return "🤖";
  if (type === "system") return "ℹ️";
  const emojis = ["😀","🦊","🐱","🐶","🐼","🦁","🐸","🐧","🦄","🐲"];
  let h = 0;
  for (const c of (name || "?")) h = (h * 31 + c.charCodeAt(0)) & 0xffff;
  return emojis[h % emojis.length];
}

function scrollBottom() {
  elMessages.scrollTop = elMessages.scrollHeight;
}

function showToast(author, text) {
  const el = document.createElement("div");
  el.className = "toast";
  el.innerHTML = `<div class="t-author">${esc(author)}</div><div class="t-text">${esc(text)}</div>`;
  elToasts.appendChild(el);
  setTimeout(() => el.remove(), 4000);
}

// ── 渲染 ──
function renderChannels() {
  elChannels.innerHTML = "";
  for (const ch of channels) {
    const li = document.createElement("li");
    li.className = ch.id === channel ? "active" : "";
    li.innerHTML = `<span class="hash">#</span>${esc(ch.name)}<span class="badge">${ch.count ?? 0}</span>`;
    li.addEventListener("click", () => {
      if (ch.id !== channel) Events.Emit("chat:switch", { channel: ch.id });
    });
    elChannels.appendChild(li);
  }
}

function renderUsers() {
  elUsers.innerHTML = "";
  elUserCount.textContent = users.length;
  for (const u of users) {
    const li = document.createElement("li");
    const isMe = u.name === nick;
    li.innerHTML = `
      <span class="dot${u.type === "bot" ? " bot" : ""}"></span>
      <span class="name">${esc(u.name)}</span>
      ${isMe ? '<span class="you">(你)</span>' : ""}`;
    elUsers.appendChild(li);
  }
}

function renderMessages(msgs, replace = false) {
  if (replace) {
    elMessages.innerHTML = "";
    msgIds.clear();
  }
  for (const raw of msgs) {
    const m = normalizeMsg(raw);
    if (m) appendMessage(m, false);
  }
  scrollBottom();
}

function appendMessage(m, scroll = true) {
  if (!m || !m.text && m.type !== "system") return;
  if (msgIds.has(m.id)) return;
  msgIds.add(m.id);

  const isSelf = m.author === nick && m.type === "user";
  const div = document.createElement("div");
  div.className = `msg ${m.type}${isSelf ? " self" : ""}`;

  if (m.type === "system") {
    div.innerHTML = `<div class="body"><div class="text">${esc(m.text)}</div></div>`;
  } else {
    div.innerHTML = `
      <div class="avatar">${avatar(m.author, m.type)}</div>
      <div class="body">
        <div class="meta">
          <span class="author">${esc(m.author)}</span>
          <span class="time">${esc(m.time || "")}</span>
        </div>
        <div class="text">${esc(m.text)}</div>
      </div>`;
  }
  elMessages.appendChild(div);
  if (scroll) scrollBottom();
}

function esc(s) {
  return String(s ?? "")
    .replace(/&/g, "&amp;").replace(/</g, "&lt;")
    .replace(/>/g, "&gt;").replace(/"/g, "&quot;");
}

// ── 发送 ──
function send(text) {
  text = text.trim();
  if (!text) return;
  Events.Emit("chat:send", { channel, nick, text, time: now() });
  elInput.value = "";
  elInput.focus();
}

elComposer.addEventListener("submit", (e) => {
  e.preventDefault();
  send(elInput.value);
});

elNick.addEventListener("change", () => {
  const v = elNick.value.trim();
  if (v && v !== nick) {
    Events.Emit("chat:nick", { nick: v });
  }
});
elNick.addEventListener("keydown", (e) => {
  if (e.key === "Enter") { elNick.blur(); elInput.focus(); }
});

// ── 后端事件 ──

// 完整状态（初始化 / 切频道）
Events.On("chat:state", (ev) => {
  const d = ev.data || {};
  if (d.nick)    { nick = d.nick; elNick.value = nick; }
  if (d.channel) { channel = d.channel; }
  if (d.channels) { channels = d.channels; renderChannels(); }
  if (d.users)    { users = d.users; renderUsers(); }
  if (d.messages) { renderMessages(d.messages, true); }

  const ch = channels.find((c) => c.id === channel);
  elCurName.textContent = ch?.name || channel;
  elCurTopic.textContent = ch?.topic ? `— ${ch.topic}` : "";
  elInput.placeholder = `在 #${ch?.name || channel} 发消息…  试试 /help`;
});

// 新消息（实时推送）
Events.On("chat:message", (ev) => {
  const m = normalizeMsg(ev?.data ?? ev);
  if (!m) return;

  const isCurrent = m.channel === channel;
  const fromOthers = m.author !== nick || m.type !== "user";

  if (isCurrent) {
    appendMessage(m, isNearBottom());
  } else if (fromOthers) {
    const ch = channels.find((c) => c.id === m.channel);
    showToast(m.author, `#${ch?.name ?? m.channel}: ${m.text}`);
  }

  const ch = channels.find((c) => c.id === m.channel);
  if (ch) { ch.count = (ch.count || 0) + 1; renderChannels(); }
});

// 在线用户更新
Events.On("chat:users", (ev) => {
  users = ev.data || [];
  renderUsers();
});

// 错误提示（应用内系统消息，不弹原生对话框）
Events.On("chat:error", (ev) => {
  const text = ev?.data?.text ?? ev?.text ?? "未知错误";
  appendMessage({
    id: "err-" + Date.now(),
    type: "system",
    text: "⚠️ " + text,
    channel,
    author: "",
    time: "",
  });
});

// ── 启动：向后端请求初始化 ──
window.addEventListener("DOMContentLoaded", () => {
  Events.Emit("chat:join", { nick: elNick.value.trim() || "访客" });
  elInput.focus();
});
