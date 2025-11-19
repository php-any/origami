/**
 * Toast 通知系统
 * 优雅的非阻塞式通知，替代 alert/confirm
 */

// Toast 容器
let toastContainer = null;

function initToastContainer() {
  if (!toastContainer) {
    toastContainer = document.createElement("div");
    toastContainer.id = "toast-container";
    toastContainer.className = "toast-container";
    document.body.appendChild(toastContainer);
  }
  return toastContainer;
}

/**
 * 显示 Toast 通知
 * @param {string} message - 消息内容
 * @param {string} type - 类型: 'success', 'error', 'warning', 'info'
 * @param {number} duration - 显示时长（毫秒），默认 3000
 */
function showToast(message, type = "info", duration = 3000) {
  const container = initToastContainer();

  const toast = document.createElement("div");
  toast.className = `toast toast-${type}`;

  // 图标映射
  const icons = {
    success: "✓",
    error: "✕",
    warning: "⚠",
    info: "ℹ",
  };

  toast.innerHTML = `
    <div class="toast-icon">${icons[type] || icons.info}</div>
    <div class="toast-message">${message}</div>
    <button class="toast-close" onclick="this.parentElement.remove()">×</button>
  `;

  container.appendChild(toast);

  // 触发动画
  setTimeout(() => toast.classList.add("show"), 10);

  // 自动移除
  setTimeout(() => {
    toast.classList.remove("show");
    setTimeout(() => toast.remove(), 300);
  }, duration);
}

/**
 * 显示成功通知
 */
function showSuccess(message, duration = 3000) {
  showToast(message, "success", duration);
}

/**
 * 显示错误通知
 */
function showError(message, duration = 4000) {
  showToast(message, "error", duration);
}

/**
 * 显示警告通知
 */
function showWarning(message, duration = 3500) {
  showToast(message, "warning", duration);
}

/**
 * 显示信息通知
 */
function showInfo(message, duration = 3000) {
  showToast(message, "info", duration);
}

/**
 * 确认对话框（使用自定义样式）
 * @param {string} message - 确认消息
 * @param {Function} onConfirm - 确认回调
 * @param {Function} onCancel - 取消回调（可选）
 */
function showConfirm(message, onConfirm, onCancel) {
  const overlay = document.createElement("div");
  overlay.className = "confirm-overlay";

  const dialog = document.createElement("div");
  dialog.className = "confirm-dialog";

  dialog.innerHTML = `
    <div class="confirm-icon">⚠</div>
    <div class="confirm-message">${message}</div>
    <div class="confirm-buttons">
      <button class="btn btn-secondary confirm-cancel">取消</button>
      <button class="btn btn-primary confirm-ok">确定</button>
    </div>
  `;

  overlay.appendChild(dialog);
  document.body.appendChild(overlay);

  // 动画
  setTimeout(() => overlay.classList.add("show"), 10);

  // 确认按钮
  dialog.querySelector(".confirm-ok").addEventListener("click", () => {
    overlay.classList.remove("show");
    setTimeout(() => {
      overlay.remove();
      if (onConfirm) onConfirm();
    }, 300);
  });

  // 取消按钮
  dialog.querySelector(".confirm-cancel").addEventListener("click", () => {
    overlay.classList.remove("show");
    setTimeout(() => {
      overlay.remove();
      if (onCancel) onCancel();
    }, 300);
  });

  // 点击遮罩层关闭
  overlay.addEventListener("click", (e) => {
    if (e.target === overlay) {
      overlay.classList.remove("show");
      setTimeout(() => {
        overlay.remove();
        if (onCancel) onCancel();
      }, 300);
    }
  });
}
