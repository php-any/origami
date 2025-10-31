// 管理后台脚本

let tools = [];
let projects = [];
let currentEditingTool = null;
let currentEditingProject = null;
let projectEnvironments = [];

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
async function loadTools() {
  try {
    const response = await fetch("/api/tools");
    tools = await response.json();
    renderToolsTable();
  } catch (error) {
    console.error("加载工具失败:", error);
    document.getElementById("toolsTableContainer").innerHTML =
      '<div class="empty-state">加载失败，请刷新重试</div>';
  }
}

// 加载项目
async function loadProjects() {
  try {
    const response = await fetch("/api/projects");
    projects = await response.json();
    renderProjectsTable();
  } catch (error) {
    console.error("加载项目失败:", error);
    document.getElementById("projectsTableContainer").innerHTML =
      '<div class="empty-state">加载失败，请刷新重试</div>';
  }
}

// 渲染工具表格
function renderToolsTable() {
  const container = document.getElementById("toolsTableContainer");
  if (tools.length === 0) {
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
          <th>顺序</th>
          <th>操作</th>
        </tr>
      </thead>
      <tbody>
        ${tools
          .map(
            (tool) => `
          <tr>
            <td>${tool.id}</td>
            <td>${tool.icon || "🔗"}</td>
            <td>${tool.name}</td>
            <td>${tool.category || "-"}</td>
            <td><a href="${
              tool.url
            }" target="_blank" style="color: var(--primary);">${
              tool.url
            }</a></td>
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
        `
          )
          .join("")}
      </tbody>
    </table>
  `;
}

// 渲染项目表格
function renderProjectsTable() {
  const container = document.getElementById("projectsTableContainer");
  if (projects.length === 0) {
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
        ${projects
          .map(
            (project) => `
          <tr>
            <td>${project.id}</td>
            <td>${project.name}</td>
            <td>${project.environments.length}</td>
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

  if (toolId) {
    title.textContent = "编辑工具链接";
    const tool = tools.find((t) => t.id === toolId);
    if (tool) {
      document.getElementById("toolId").value = tool.id;
      document.getElementById("toolName").value = tool.name;
      document.getElementById("toolUrl").value = tool.url;
      document.getElementById("toolIcon").value = tool.icon || "";
      document.getElementById("toolCategory").value = tool.category || "";
      document.getElementById("toolDescription").value = tool.description || "";
      document.getElementById("toolDisplayOrder").value =
        tool.displayOrder || 0;
    }
  } else {
    title.textContent = "添加工具链接";
    form.reset();
    document.getElementById("toolId").value = "";
  }

  modal.classList.add("active");
}

// 关闭工具编辑模态框
function closeToolModal() {
  document.getElementById("toolModal").classList.remove("active");
  currentEditingTool = null;
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
    await loadTools();
    alert("保存成功！");
  } catch (error) {
    console.error("保存失败:", error);
    alert("保存失败: " + error.message);
  }
}

// 编辑工具
function editTool(id) {
  openToolModal(id);
}

// 删除工具
async function deleteTool(id) {
  if (!confirm("确定要删除这个工具链接吗？")) return;

  try {
    const response = await fetch(`/api/tools/${id}`, {
      method: "DELETE",
    });

    if (!response.ok) throw new Error("删除失败");

    await loadTools();
    alert("删除成功！");
  } catch (error) {
    console.error("删除失败:", error);
    alert("删除失败: " + error.message);
  }
}

// 打开项目编辑模态框
function openProjectModal(projectId = null) {
  currentEditingProject = projectId;
  const modal = document.getElementById("projectModal");
  const title = document.getElementById("projectModalTitle");
  const form = document.getElementById("projectForm");

  if (projectId) {
    title.textContent = "编辑项目";
    const project = projects.find((p) => p.id === projectId);
    if (project) {
      document.getElementById("projectId").value = project.id;
      document.getElementById("projectName").value = project.name;
      document.getElementById("projectDisplayOrder").value =
        project.displayOrder || 0;
      projectEnvironments = JSON.parse(JSON.stringify(project.environments));
      renderEnvironments();
    }
  } else {
    title.textContent = "添加项目";
    form.reset();
    document.getElementById("projectId").value = "";
    projectEnvironments = [];
    renderEnvironments();
  }

  modal.classList.add("active");
}

// 关闭项目编辑模态框
function closeProjectModal() {
  document.getElementById("projectModal").classList.remove("active");
  currentEditingProject = null;
  projectEnvironments = [];
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
    <div style="background: rgba(15, 23, 42, 0.4); padding: 15px; border-radius: 8px; margin-bottom: 10px;">
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

// 保存项目
async function saveProject(event) {
  event.preventDefault();

  const formData = {
    name: document.getElementById("projectName").value,
    displayOrder:
      parseInt(document.getElementById("projectDisplayOrder").value) || 0,
    environments: projectEnvironments,
  };

  try {
    if (currentEditingProject) {
      // 更新 - 这里需要实现 PUT API
      alert("更新功能待实现，请使用 API 直接操作数据库");
    } else {
      // 创建
      const response = await fetch("/api/projects", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(formData),
      });

      if (!response.ok) throw new Error("创建失败");

      const result = await response.json();
      const projectId = result.id;

      // 创建环境
      for (const env of projectEnvironments) {
        await fetch(`/api/projects/${projectId}/environments`, {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify(env),
        });
      }
    }

    closeProjectModal();
    await loadProjects();
    alert("保存成功！");
  } catch (error) {
    console.error("保存失败:", error);
    alert("保存失败: " + error.message);
  }
}

// 编辑项目
function editProject(id) {
  openProjectModal(id);
}

// 删除项目
async function deleteProject(id) {
  if (!confirm("确定要删除这个项目吗？所有关联的环境也将被删除！")) return;

  try {
    const response = await fetch(`/api/projects/${id}`, {
      method: "DELETE",
    });

    if (!response.ok) throw new Error("删除失败");

    await loadProjects();
    alert("删除成功！");
  } catch (error) {
    console.error("删除失败:", error);
    alert("删除失败: " + error.message);
  }
}

// 初始化
document.addEventListener("DOMContentLoaded", function () {
  loadTools();
  loadProjects();

  // 点击模态框外部关闭
  document.querySelectorAll(".modal").forEach((modal) => {
    modal.addEventListener("click", function (e) {
      if (e.target === this) {
        this.classList.remove("active");
      }
    });
  });
});
