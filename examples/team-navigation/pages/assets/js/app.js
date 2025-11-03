// 团队导航页交互脚本
(function () {
  "use strict";

  let allTools = [];
  let allProjects = [];

  // 从 API 加载数据
  async function loadData() {
    try {
      // 加载工具链接
      const toolsResponse = await fetch("/api/tools");
      allTools = await toolsResponse.json();

      // 加载项目数据
      const projectsResponse = await fetch("/api/projects");
      allProjects = await projectsResponse.json();

      // 渲染工具卡片
      renderTools(allTools);

      // 渲染项目卡片
      renderProjects(allProjects);

      // 隐藏加载提示，显示内容
      const loadingEl = document.getElementById("loading");
      if (loadingEl) {
        loadingEl.style.display = "none";
      }

      const toolsSection = document.getElementById("toolsSection");
      const projectsSection = document.getElementById("projectsSection");
      if (toolsSection) toolsSection.style.display = "block";
      if (projectsSection) projectsSection.style.display = "block";
    } catch (error) {
      console.error("加载数据失败:", error);
      const loadingEl = document.getElementById("loading");
      if (loadingEl) {
        loadingEl.innerHTML = `
          <div style="color: #ea4335;">
            <p>加载数据失败，请刷新页面重试</p>
          </div>
        `;
      }
    }
  }

  // 渲染工具卡片
  function renderTools(tools) {
    const toolsGrid = document.getElementById("toolsGrid");
    if (!toolsGrid) return;

    toolsGrid.innerHTML = tools
      .map(
        (tool) => `
      <a href="${tool.url}" target="_blank" class="tool-card">
        <div class="tool-icon">${tool.icon || "🔗"}</div>
        <div class="tool-name">${tool.name}</div>
      </a>
    `
      )
      .join("");

    // 渲染收藏链接（只显示标记为收藏的工具，按显示顺序排序）
    const favoriteTools = tools
      .filter((tool) => tool.isFavorite == 1 || tool.isFavorite === 1)
      .sort((a, b) => (a.displayOrder || 0) - (b.displayOrder || 0));
    renderFavoriteLinks(favoriteTools);
  }

  // 渲染收藏链接
  function renderFavoriteLinks(favoriteTools) {
    const favoriteLinksEl = document.getElementById("favoriteLinks");
    if (!favoriteLinksEl) return;

    if (favoriteTools.length === 0) {
      favoriteLinksEl.style.display = "none";
      return;
    }

    favoriteLinksEl.style.display = "flex";
    favoriteLinksEl.innerHTML = favoriteTools
      .map(
        (tool) => `
      <a href="${tool.url}" target="_blank" class="favorite-link">
        <span class="favorite-link-icon">${tool.icon || "🔗"}</span>
        <span>${tool.name}</span>
      </a>
    `
      )
      .join("");
  }

  // 渲染项目卡片 - 简化设计
  function renderProjects(projects) {
    const projectsGrid = document.getElementById("projectsGrid");
    if (!projectsGrid) return;

    // 1) 数据清洗：去重（优先用 id 作为键，其次用 name），并过滤无效/空项目
    const seenKeys = new Set();
    const cleanedProjects = [];
    for (const raw of Array.isArray(projects) ? projects : []) {
      const project = raw || {};
      const name = (project.name || "").trim();
      // 跳过空名称的项目，避免渲染空白卡片
      if (!name) continue;
      const key =
        project.id != null ? `id:${project.id}` : `name:${name.toLowerCase()}`;
      if (seenKeys.has(key)) continue;
      seenKeys.add(key);
      cleanedProjects.push(project);
    }

    projectsGrid.innerHTML = cleanedProjects
      .map((project) => {
        // 获取第一个环境作为主要链接
        const firstEnv =
          project.environments && project.environments.length > 0
            ? project.environments[0]
            : null;
        const desc =
          project.description ||
          project.projectDescription ||
          project.desc ||
          "";

        return `
      <a href="${
        firstEnv ? firstEnv.url : "#"
      }" target="_blank" class="project-card">
        <h4 class="project-name">${project.name}</h4>
        ${desc ? `<p class="project-desc">${desc}</p>` : ""}
        
        ${
          project.tools && project.tools.length > 0
            ? `
        <div class="project-tools">
          ${project.tools
            .slice(0, 3)
            .map(
              (tool) => `
            <a href="${tool.url}" target="_blank" class="tool-tag" title="${
                tool.name
              }" onclick="event.stopPropagation();">
              <span class="tool-tag-icon">${tool.icon || "🔗"}</span>
              <span>${tool.name}</span>
            </a>
          `
            )
            .join("")}
        </div>
        `
            : ""
        }
        
        ${
          project.environments && project.environments.length > 0
            ? `
        <div class="project-environments">
          ${project.environments
            .slice(0, 3)
            .map(
              (env) => `
            <div class="env-item" onclick="event.stopPropagation();">
              <span class="env-name">${env.environmentName}</span>
              <a href="${
                env.url
              }" target="_blank" class="env-link" onclick="event.stopPropagation();">
                ${env.url.replace(/^https?:\/\//, "").split("/")[0]}
              </a>
            </div>
          `
            )
            .join("")}
        </div>
        `
            : ""
        }
      </a>
    `;
      })
      .join("");
  }

  // 更新时间显示
  function updateTime() {
    const now = new Date();
    const hour = String(now.getHours()).padStart(2, "0");
    const minute = String(now.getMinutes()).padStart(2, "0");

    const hourEl = document.getElementById("currentHour");
    const minuteEl = document.getElementById("currentMinute");
    const dateEl = document.getElementById("currentDate");

    if (hourEl) hourEl.textContent = hour;
    if (minuteEl) minuteEl.textContent = minute;

    if (dateEl) {
      const options = {
        year: "numeric",
        month: "long",
        day: "numeric",
        weekday: "long",
      };
      dateEl.textContent = now.toLocaleDateString("zh-CN", options);
    }

    // 更新页脚时间
    const lastUpdateEl = document.getElementById("lastUpdate");
    if (lastUpdateEl) {
      lastUpdateEl.textContent = now.toLocaleString("zh-CN");
    }
  }

  // 立即更新时间，然后每秒更新
  updateTime();
  setInterval(updateTime, 1000);

  // 百度搜索功能
  const baiduSearchForm = document.getElementById("baiduSearchForm");
  const baiduSearchInput = document.getElementById("baiduSearchInput");

  if (baiduSearchForm && baiduSearchInput) {
    // 表单提交验证 - 按 Enter 键搜索
    baiduSearchForm.addEventListener("submit", function (e) {
      const query = baiduSearchInput.value.trim();
      if (!query) {
        e.preventDefault();
        baiduSearchInput.focus();
        return false;
      }
    });

    // 自动聚焦搜索框（页面加载后）
    window.addEventListener("load", function () {
      setTimeout(function () {
        baiduSearchInput.focus();
      }, 100);
    });
  }

  // 键盘快捷键支持
  document.addEventListener("keydown", function (e) {
    // Ctrl/Cmd + K 聚焦搜索框
    if ((e.ctrlKey || e.metaKey) && e.key === "k") {
      e.preventDefault();
      if (baiduSearchInput) {
        baiduSearchInput.focus();
        baiduSearchInput.select();
      }
    }
    // Esc 键清除搜索框
    if (e.key === "Escape") {
      if (baiduSearchInput && document.activeElement === baiduSearchInput) {
        baiduSearchInput.blur();
        baiduSearchInput.value = "";
      }
    }
  });

  // 页面加载完成后加载数据
  if (document.readyState === "loading") {
    document.addEventListener("DOMContentLoaded", loadData);
  } else {
    loadData();
  }
})();
