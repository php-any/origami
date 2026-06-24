(() => {
  const statusDot = document.getElementById('ws-status-dot');
  const statusText = document.getElementById('ws-status-text');
  const onlineCount = document.getElementById('online-count');
  const nicknameInput = document.getElementById('nickname');
  const connectBtn = document.getElementById('ws-connect');
  const disconnectBtn = document.getElementById('ws-disconnect');
  const chatLog = document.getElementById('chat-log');
  const chatInput = document.getElementById('chat-input');
  const chatSend = document.getElementById('chat-send');

  let ws = null;

  const wsUrl = () => {
    const protocol = location.protocol === 'https:' ? 'wss:' : 'ws:';
    return `${protocol}//${location.host}/ws/chat`;
  };

  const setConnected = (connected) => {
    statusDot.classList.toggle('online', connected);
    statusDot.classList.toggle('offline', !connected);
    statusText.textContent = connected ? '已连接' : '未连接';
    connectBtn.disabled = connected;
    disconnectBtn.disabled = !connected;
    nicknameInput.disabled = connected;
    chatInput.disabled = !connected;
    chatSend.disabled = !connected;
  };

  const appendMessage = (payload) => {
    const item = document.createElement('div');
    item.className = 'chat-item' + (payload.type === 'system' ? ' system' : '');

    const meta = document.createElement('div');
    meta.className = 'meta';
    const time = payload.time ? new Date(payload.time * 1000).toLocaleTimeString() : '';
    meta.textContent = payload.type === 'system'
      ? `[系统] ${time}`
      : `[${payload.from || '匿名'}] ${time}`;

    const body = document.createElement('div');
    body.textContent = payload.text || '';

    item.appendChild(meta);
    item.appendChild(body);
    chatLog.appendChild(item);
    chatLog.scrollTop = chatLog.scrollHeight;

    if (typeof payload.online === 'number') {
      onlineCount.textContent = `在线 ${payload.online} 人`;
    }
  };

  const connect = () => {
    if (ws) return;
    ws = new WebSocket(wsUrl());

    ws.onopen = () => {
      setConnected(true);
      appendMessage({ type: 'system', text: 'WebSocket 连接已建立', time: Math.floor(Date.now() / 1000) });
    };

    ws.onmessage = (event) => {
      try {
        appendMessage(JSON.parse(event.data));
      } catch (e) {
        appendMessage({ type: 'message', from: '服务端', text: event.data, time: Math.floor(Date.now() / 1000) });
      }
    };

    ws.onclose = () => {
      ws = null;
      setConnected(false);
      appendMessage({ type: 'system', text: '连接已断开', time: Math.floor(Date.now() / 1000) });
    };

    ws.onerror = () => {
      appendMessage({ type: 'system', text: '连接出错', time: Math.floor(Date.now() / 1000) });
    };
  };

  const disconnect = () => {
    if (ws) {
      ws.close();
      ws = null;
    }
  };

  const sendMessage = () => {
    const text = chatInput.value.trim();
    if (!text || !ws || ws.readyState !== WebSocket.OPEN) return;

    const payload = {
      type: 'message',
      from: nicknameInput.value.trim() || '访客',
      text,
      time: Math.floor(Date.now() / 1000),
    };
    ws.send(JSON.stringify(payload));
    chatInput.value = '';
  };

  connectBtn.addEventListener('click', connect);
  disconnectBtn.addEventListener('click', disconnect);
  chatSend.addEventListener('click', sendMessage);
  chatInput.addEventListener('keydown', (e) => {
    if (e.key === 'Enter') sendMessage();
  });

  setConnected(false);
})();
