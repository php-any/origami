// 管理后台脚本

let currentEditingTool = null;
let currentEditingProject = null;
let projectEnvironments = [];
let projectTools = [];

// 标签页切换
document.querySelectorAll(".tab").forEach((tab) => {
  tab.addEventListener("click", function () {
    document
      .querySelectorAll(".tab")
      .forEach((t) => t.classList.remove("active"));
    document
      .querySelectorAll(".tab-content")
      .forEach((c) => c.classList.remove("active"));

    this.classList.add("active");
    const tabName = this.dataset.tab;
    document.getElementById(tabName + "Tab").classList.add("active");
  });
});

// 加载工具链接
// 后端渲染，移除前端加载工具数据

// 加载项目
// 后端渲染，移除前端加载项目数据

// 渲染工具表格
function renderToolsTable() {
  // 确保 tools 是数组类型
  const toolsArray = Array.isArray(tools) ? tools : [];

  const container = document.getElementById("toolsTableContainer");
  if (toolsArray.length === 0) {
    container.innerHTML =
      '<div class="empty-state">暂无工具链接，点击上方按钮添加</div>';
    return;
  }

  container.innerHTML = `
    <table>
      <thead>
        <tr>
          <th>ID</th>
          <th>图标</th>
          <th>名称</th>
          <th>分类</th>
          <th>链接</th>
          <th>收藏</th>
          <th>顺序</th>
          <th>操作</th>
        </tr>
      </thead>
      <tbody>
        ${toolsArray
          .map((tool) => {
            const isImageIcon =
              tool.icon &&
              (tool.icon.startsWith("http://") ||
                tool.icon.startsWith("https://") ||
                tool.icon.startsWith("/") ||
                /\.(png|jpg|jpeg|gif|svg|webp|ico)$/i.test(tool.icon));
            const iconDisplay = isImageIcon
              ? `<img src="${tool.icon}" alt="${tool.name}" style="width: 20px; height: 20px; object-fit: contain; vertical-align: middle;">`
              : tool.icon || "🔗";
            return `
          <tr>
            <td>${tool.id}</td>
            <td>${iconDisplay}</td>
            <td>${tool.name}</td>
            <td>${tool.category || "-"}</td>
            <td><a href="${
              tool.url
            }" target="_blank" style="color: var(--primary);">${
              tool.url
            }</a></td>
            <td>${tool.isFavorite ? "⭐" : "-"}</td>
            <td>${tool.displayOrder || 0}</td>
            <td>
              <div class="action-buttons">
                <button class="btn btn-small" onclick="editTool(${
                  tool.id
                })">编辑</button>
                <button class="btn btn-small btn-danger" onclick="deleteTool(${
                  tool.id
                })">删除</button>
              </div>
            </td>
          </tr>
        `;
          })
          .join("")}
      </tbody>
    </table>
  `;
}

// 渲染项目表格
function renderProjectsTable() {
  // 确保 projects 是数组类型
  const projectsArray = Array.isArray(projects) ? projects : [];

  const container = document.getElementById("projectsTableContainer");
  if (projectsArray.length === 0) {
    container.innerHTML =
      '<div class="empty-state">暂无项目，点击上方按钮添加</div>';
    return;
  }

  container.innerHTML = `
    <table>
      <thead>
        <tr>
          <th>ID</th>
          <th>项目名称</th>
          <th>环境数量</th>
          <th>顺序</th>
          <th>操作</th>
        </tr>
      </thead>
      <tbody>
        ${projectsArray
          .map(
            (project) => `
          <tr>
            <td>${project.id}</td>
            <td>${project.name}</td>
            <td>${
              Array.isArray(project.environments)
                ? project.environments.length
                : 0
            }</td>
            <td>${project.displayOrder || 0}</td>
            <td>
              <div class="action-buttons">
                <button class="btn btn-small" onclick="editProject(${
                  project.id
                })">编辑</button>
                <button class="btn btn-small btn-danger" onclick="deleteProject(${
                  project.id
                })">删除</button>
              </div>
            </td>
          </tr>
        `
          )
          .join("")}
      </tbody>
    </table>
  `;
}

// 打开工具编辑模态框
function openToolModal(toolId = null) {
  currentEditingTool = toolId;
  const modal = document.getElementById("toolModal");
  const title = document.getElementById("toolModalTitle");
  const form = document.getElementById("toolForm");

  // 确保 tools 是数组类型
  const toolsArray = Array.isArray(tools) ? tools : [];

  if (toolId) {
    title.textContent = "编辑工具链接";
    const tool = toolsArray.find((t) => t.id === toolId);
    if (tool) {
      document.getElementById("toolId").value = tool.id;
      document.getElementById("toolName").value = tool.name;
      document.getElementById("toolUrl").value = tool.url;
      document.getElementById("toolIcon").value = tool.icon || "";
      document.getElementById("toolCategory").value = tool.category || "";
      document.getElementById("toolDescription").value = tool.description || "";
      document.getElementById("toolDisplayOrder").value =
        tool.displayOrder || 0;
      document.getElementById("toolIsFavorite").checked = tool.isFavorite == 1;
    }
  } else {
    title.textContent = "添加工具链接";
    form.reset();
    document.getElementById("toolId").value = "";
    document.getElementById("toolIsFavorite").checked = false;
  }

  modal.classList.add("active");
  // 阻止页面滚动，确保模态框相对于视口居中
  document.body.style.overflow = "hidden";
}

// 关闭工具编辑模态框
function closeToolModal() {
  document.getElementById("toolModal").classList.remove("active");
  currentEditingTool = null;
  // 恢复页面滚动
  document.body.style.overflow = "";
}

// 保存工具
async function saveTool(event) {
  event.preventDefault();

  const formData = {
    name: document.getElementById("toolName").value,
    url: document.getElementById("toolUrl").value,
    icon: document.getElementById("toolIcon").value,
    category: document.getElementById("toolCategory").value,
    description: document.getElementById("toolDescription").value,
    isFavorite: document.getElementById("toolIsFavorite").checked ? 1 : 0,
    displayOrder:
      parseInt(document.getElementById("toolDisplayOrder").value) || 0,
  };

  try {
    if (currentEditingTool) {
      // 更新
      const response = await fetch(`/api/tools/${currentEditingTool}`, {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(formData),
      });

      if (!response.ok) throw new Error("更新失败");
    } else {
      // 创建
      const response = await fetch("/api/tools", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(formData),
      });

      if (!response.ok) throw new Error("创建失败");
    }

    closeToolModal();
    showSuccess("保存成功！");
    setTimeout(() => window.location.reload(), 500);
  } catch (error) {
    console.error("保存失败:", error);
    showError("保存失败: " + error.message);
  }
}

// 编辑工具
function editTool(id) {
  openToolModal(id);
}

// 删除工具
async function deleteTool(id) {
  showConfirm("确定要删除这个工具链接吗？", async () => {
    try {
      const response = await fetch(`/api/tools/${id}`, {
        method: "DELETE",
      });

      if (!response.ok) throw new Error("删除失败");

      showSuccess("删除成功！");
      setTimeout(() => window.location.reload(), 500);
    } catch (error) {
      console.error("删除失败:", error);
      showError("删除失败: " + error.message);
    }
  });
}

// 打开项目编辑模态框
function openProjectModal(projectId = null) {
  currentEditingProject = projectId;
  const modal = document.getElementById("projectModal");
  const title = document.getElementById("projectModalTitle");
  const form = document.getElementById("projectForm");

  // 确保 projects 是数组类型
  const projectsArray = Array.isArray(projects) ? projects : [];

  if (projectId) {
    title.textContent = "编辑项目";
    // 确保ID类型一致（转换为数字进行比较）
    const project = projectsArray.find(
      (p) => p.id == projectId || String(p.id) === String(projectId)
    );
    if (project) {
      document.getElementById("projectId").value = project.id;
      document.getElementById("projectName").value = project.name;
      document.getElementById("projectDisplayOrder").value =
        project.displayOrder || 0;
      const descEl = document.getElementById("projectDescription");
      if (descEl) descEl.value = project.description || "";

      // 处理项目图标
      const icon = project.icon || "";
      if (icon && (icon.startsWith("http://") || icon.startsWith("https://"))) {
        // 是图片链接
        document.getElementById("projectIconUrl").value = icon;
        document.getElementById("projectIconEmoji").value = "";
      } else {
        // 是图标（Emoji）
        document.getElementById("projectIconEmoji").value = icon || "🚀";
        document.getElementById("projectIconUrl").value = "";
      }
      updateProjectIconPreview();

      projectEnvironments = JSON.parse(JSON.stringify(project.environments));
      projectTools = Array.isArray(project.tools)
        ? project.tools.map((t) => t.id)
        : [];
      renderEnvironments();
      renderProjectTools();
    }
  } else {
    title.textContent = "添加项目";
    form.reset();
    document.getElementById("projectId").value = "";
    const descEl = document.getElementById("projectDescription");
    if (descEl) descEl.value = "";
    document.getElementById("projectIconEmoji").value = "🚀";
    document.getElementById("projectIconUrl").value = "";
    updateProjectIconPreview();
    projectEnvironments = [];
    projectTools = [];
    renderEnvironments();
    renderProjectTools();
  }

  modal.classList.add("active");
  // 阻止页面滚动，确保模态框相对于视口居中
  document.body.style.overflow = "hidden";
}

// 关闭项目编辑模态框
function closeProjectModal() {
  document.getElementById("projectModal").classList.remove("active");
  currentEditingProject = null;
  projectEnvironments = [];
  projectTools = [];
  // 恢复页面滚动
  document.body.style.overflow = "";
}

// 渲染项目工具选择
function renderProjectTools() {
  const container = document.getElementById("projectToolsList");
  if (!container) return;

  // 确保 tools 和 projectTools 是数组类型
  const toolsArray = Array.isArray(tools) ? tools : [];
  const projectToolsArray = Array.isArray(projectTools) ? projectTools : [];

  container.innerHTML = toolsArray
    .map((tool) => {
      const isImageIcon =
        tool.icon &&
        (tool.icon.startsWith("http://") ||
          tool.icon.startsWith("https://") ||
          tool.icon.startsWith("/") ||
          /\.(png|jpg|jpeg|gif|svg|webp|ico)$/i.test(tool.icon));
      const iconDisplay = isImageIcon
        ? `<img src="${tool.icon}" alt="${tool.name}" style="width: 16px; height: 16px; object-fit: contain; vertical-align: middle;">`
        : tool.icon || "🔗";
      return `
    <label style="display: flex; align-items: center; gap: 8px; padding: 8px; border-radius: 4px; cursor: pointer; transition: background 0.2s;" 
           onmouseover="this.style.background='var(--bg-hover)'" 
           onmouseout="this.style.background='transparent'">
      <input type="checkbox" value="${tool.id}" 
             ${projectToolsArray.includes(tool.id) ? "checked" : ""} 
             onchange="toggleProjectTool(${tool.id}, this.checked)">
      ${iconDisplay}
      <span>${tool.name}</span>
      ${
        tool.category
          ? `<span style="color: var(--text-muted); font-size: 0.85rem;">(${tool.category})</span>`
          : ""
      }
    </label>
  `;
    })
    .join("");
}

// 切换项目工具
function toggleProjectTool(toolId, checked) {
  // 确保 projectTools 是数组类型
  const projectToolsArray = Array.isArray(projectTools) ? projectTools : [];

  if (checked) {
    if (!projectToolsArray.includes(toolId)) {
      projectToolsArray.push(toolId);
      projectTools = projectToolsArray; // 更新全局变量
    }
  } else {
    projectTools = projectToolsArray.filter((id) => id !== toolId);
  }
}

// 渲染环境列表
function renderEnvironments() {
  const container = document.getElementById("environmentsList");
  if (projectEnvironments.length === 0) {
    container.innerHTML =
      '<div class="empty-state" style="padding: 20px;">暂无环境配置</div>';
    return;
  }

  container.innerHTML = projectEnvironments
    .map(
      (env, index) => `
    <div>
      <div class="env-form-row">
        <div class="form-group">
          <label>环境名称</label>
          <input type="text" value="${
            env.environmentName || ""
          }" onchange="updateEnvironment(${index}, 'environmentName', this.value)" required>
        </div>
        <div class="form-group">
          <label>URL</label>
          <input type="url" value="${
            env.url || ""
          }" onchange="updateEnvironment(${index}, 'url', this.value)" required>
        </div>
      </div>
      <div class="env-form-row">
        <div class="form-group">
          <label>状态</label>
          <select onchange="updateEnvironment(${index}, 'status', this.value)">
            <option value="运行中" ${
              env.status === "运行中" ? "selected" : ""
            }>运行中</option>
            <option value="维护中" ${
              env.status === "维护中" ? "selected" : ""
            }>维护中</option>
            <option value="异常" ${
              env.status === "异常" ? "selected" : ""
            }>异常</option>
          </select>
        </div>
        <div class="form-group">
          <label>状态颜色</label>
          <select onchange="updateEnvironment(${index}, 'statusColor', this.value)">
            <option value="green" ${
              env.statusColor === "green" ? "selected" : ""
            }>绿色</option>
            <option value="yellow" ${
              env.statusColor === "yellow" ? "selected" : ""
            }>黄色</option>
            <option value="red" ${
              env.statusColor === "red" ? "selected" : ""
            }>红色</option>
          </select>
        </div>
        <div class="form-group">
          <label>顺序</label>
          <input type="number" value="${
            env.displayOrder || 0
          }" onchange="updateEnvironment(${index}, 'displayOrder', parseInt(this.value) || 0)">
        </div>
        <div class="form-group">
          <label>&nbsp;</label>
          <button type="button" class="btn btn-small btn-danger" onclick="removeEnvironment(${index})">删除</button>
        </div>
      </div>
      ${env.id ? `<input type="hidden" class="env-id" value="${env.id}">` : ""}
    </div>
  `
    )
    .join("");
}

// 更新环境
function updateEnvironment(index, field, value) {
  if (!projectEnvironments[index]) {
    projectEnvironments[index] = {};
  }
  projectEnvironments[index][field] = value;
  // 确保保留环境 ID（如果存在）
  if (projectEnvironments[index].id && field !== "id") {
    // id 已存在，保持不变
  }
}

// 添加环境
function addEnvironment() {
  projectEnvironments.push({
    environmentName: "",
    url: "",
    status: "运行中",
    statusColor: "green",
    displayOrder: 0,
  });
  renderEnvironments();
}

// 移除环境
function removeEnvironment(index) {
  projectEnvironments.splice(index, 1);
  renderEnvironments();
}

// 更新项目图标预览
function updateProjectIconPreview() {
  const iconUrl = document.getElementById("projectIconUrl").value.trim();
  const iconEmoji = document.getElementById("projectIconEmoji").value.trim();
  const previewContent = document.getElementById("projectIconPreviewContent");
  const previewWrapper = previewContent.parentElement;

  if (iconUrl) {
    // 显示图片
    if (previewContent.tagName === "IMG") {
      previewContent.src = iconUrl;
    } else {
      const img = document.createElement("img");
      img.id = "projectIconPreviewContent";
      img.src = iconUrl;
      img.style.width = "32px";
      img.style.height = "32px";
      img.style.objectFit = "contain";
      img.onerror = function () {
        // 图片加载失败，显示默认图标
        previewWrapper.innerHTML =
          '<div style="width: 32px; height: 32px; display: flex; align-items: center; justify-content: center; background: var(--primary-light); border-radius: var(--radius-sm); flex-shrink: 0;"><span id="projectIconPreviewContent" style="font-size: 18px;">🚀</span></div><span style="font-size: 0.85rem; color: var(--text-secondary);">预览</span>';
        updateProjectIconPreview();
      };
      previewWrapper.replaceChild(img, previewContent);
    }
  } else if (iconEmoji) {
    // 显示图标（Emoji）
    if (previewContent.tagName === "IMG") {
      const span = document.createElement("span");
      span.id = "projectIconPreviewContent";
      span.style.fontSize = "18px";
      span.textContent = iconEmoji;
      previewWrapper.replaceChild(span, previewContent);
    } else {
      previewContent.textContent = iconEmoji;
    }
  } else {
    // 默认图标
    if (previewContent.tagName === "IMG") {
      const span = document.createElement("span");
      span.id = "projectIconPreviewContent";
      span.style.fontSize = "18px";
      span.textContent = "🚀";
      previewWrapper.replaceChild(span, previewContent);
    } else {
      previewContent.textContent = "🚀";
    }
  }
}

// 保存项目
async function saveProject(event) {
  event.preventDefault();

  // 准备环境数据，确保所有字段都正确
  const environments = projectEnvironments
    .filter((env) => env.environmentName && env.url)
    .map((env) => ({
      environmentName: env.environmentName,
      url: env.url,
      status: env.status || "运行中",
      statusColor: env.statusColor || "green",
      displayOrder: env.displayOrder || 0,
    }));

  // 获取项目图标（优先使用图片链接，否则使用图标）
  const iconUrl = document.getElementById("projectIconUrl").value.trim();
  const iconEmoji = document.getElementById("projectIconEmoji").value.trim();
  const icon = iconUrl || iconEmoji || null;

  const formData = {
    name: document.getElementById("projectName").value,
    description: document.getElementById("projectDescription")?.value || "",
    icon: icon,
    displayOrder:
      parseInt(document.getElementById("projectDisplayOrder").value) || 0,
    environments: environments,
    tools: projectTools || [],
  };

  try {
    if (currentEditingProject) {
      // 更新项目
      const projectId = parseInt(currentEditingProject);
      if (isNaN(projectId)) {
        throw new Error("无效的项目ID");
      }

      const response = await fetch(`/api/projects/${projectId}`, {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(formData),
      });

      if (!response.ok) {
        const errorText = await response.text();
        throw new Error("更新失败: " + errorText);
      }
    } else {
      // 创建项目
      const response = await fetch("/api/projects", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(formData),
      });

      if (!response.ok) throw new Error("创建失败");
    }

    closeProjectModal();
    showSuccess("保存成功！");
    setTimeout(() => window.location.reload(), 500);
  } catch (error) {
    console.error("保存失败:", error);
    showError("保存失败: " + error.message);
  }
}

// 编辑项目
function editProject(id) {
  // 确保ID是数字类型
  const projectId = parseInt(id);
  if (isNaN(projectId)) {
    showError("无效的项目ID");
    return;
  }
  openProjectModal(projectId);
}

// 删除项目
async function deleteProject(id) {
  showConfirm("确定要删除这个项目吗？所有关联的环境也将被删除！", async () => {
    try {
      const response = await fetch(`/api/projects/${id}`, {
        method: "DELETE",
      });

      if (!response.ok) throw new Error("删除失败");

      showSuccess("删除成功！");
      setTimeout(() => window.location.reload(), 500);
    } catch (error) {
      console.error("删除失败:", error);
      showError("删除失败: " + error.message);
    }
  });
}

// 搜索引擎管理
let currentEditingSearchEngine = null;

function openSearchEngineModal(engineId = null) {
  currentEditingSearchEngine = engineId;
  const modal = document.getElementById("searchEngineModal");
  const title = document.getElementById("searchEngineModalTitle");
  const form = document.getElementById("searchEngineForm");

  // 确保 searchEngines 是数组类型
  const searchEnginesArray = Array.isArray(searchEngines) ? searchEngines : [];

  if (engineId) {
    title.textContent = "编辑搜索引擎";
    const engine = searchEnginesArray.find((e) => e.id === engineId);
    if (engine) {
      document.getElementById("searchEngineId").value = engine.id;
      document.getElementById("searchEngineName").value = engine.name;
      document.getElementById("searchEngineUrlTemplate").value =
        engine.urlTemplate;
      document.getElementById("searchEngineIcon").value = engine.icon || "";
      document.getElementById("searchEngineDisplayOrder").value =
        engine.displayOrder || 0;
      document.getElementById("searchEngineIsDefault").checked =
        engine.isDefault == 1;
    }
  } else {
    title.textContent = "添加搜索引擎";
    form.reset();
    document.getElementById("searchEngineId").value = "";
    document.getElementById("searchEngineIsDefault").checked = false;
  }

  modal.classList.add("active");
  // 阻止页面滚动，确保模态框相对于视口居中
  document.body.style.overflow = "hidden";
}

function closeSearchEngineModal() {
  document.getElementById("searchEngineModal").classList.remove("active");
  currentEditingSearchEngine = null;
  // 恢复页面滚动
  document.body.style.overflow = "";
}

async function saveSearchEngine(event) {
  event.preventDefault();

  const formData = {
    name: document.getElementById("searchEngineName").value,
    urlTemplate: document.getElementById("searchEngineUrlTemplate").value,
    icon: document.getElementById("searchEngineIcon").value,
    displayOrder:
      parseInt(document.getElementById("searchEngineDisplayOrder").value) || 0,
    isDefault: document.getElementById("searchEngineIsDefault").checked ? 1 : 0,
  };

  try {
    if (currentEditingSearchEngine) {
      const response = await fetch(
        `/api/search-engines/${currentEditingSearchEngine}`,
        {
          method: "PUT",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify(formData),
        }
      );

      if (!response.ok) throw new Error("更新失败");
    } else {
      const response = await fetch("/api/search-engines", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(formData),
      });

      if (!response.ok) throw new Error("创建失败");
    }

    closeSearchEngineModal();
    showSuccess("保存成功！");
    setTimeout(() => window.location.reload(), 500);
  } catch (error) {
    console.error("保存失败:", error);
    showError("保存失败: " + error.message);
  }
}

function editSearchEngine(id) {
  openSearchEngineModal(id);
}

async function deleteSearchEngine(id) {
  showConfirm("确定要删除这个搜索引擎吗？", async () => {
    try {
      const response = await fetch(`/api/search-engines/${id}`, {
        method: "DELETE",
      });

      if (!response.ok) throw new Error("删除失败");

      showSuccess("删除成功！");
      setTimeout(() => window.location.reload(), 500);
    } catch (error) {
      console.error("删除失败:", error);
      showError("删除失败: " + error.message);
    }
  });
}

// 初始化
document.addEventListener("DOMContentLoaded", function () {
  // 点击模态框外部关闭
  document.querySelectorAll(".modal").forEach((modal) => {
    modal.addEventListener("click", function (e) {
      if (e.target === this) {
        this.classList.remove("active");
        // 恢复页面滚动
        document.body.style.overflow = "";
      }
    });
  });
});
