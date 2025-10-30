(() => {
  const ready = () => {
    const t = document.createElement('div');
    t.textContent = 'JS 已加载 ✅';
    t.style.position = 'fixed';
    t.style.right = '16px';
    t.style.bottom = '16px';
    t.style.padding = '8px 10px';
    t.style.background = 'rgba(13,110,253,.9)';
    t.style.color = '#fff';
    t.style.borderRadius = '10px';
    t.style.fontSize = '12px';
    t.style.boxShadow = '0 6px 24px rgba(13,110,253,.35)';
    document.body.appendChild(t);
    setTimeout(() => t.remove(), 2400);
  };
  if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', ready);
  } else {
    ready();
  }
})();


