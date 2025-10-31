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
        loadingEl.textContent = "加载数据失败，请刷新页面重试";
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
      <div class="tool-card">
        <div class="tool-icon">${tool.icon || "🔗"}</div>
        <div class="tool-content">
          <h3 class="tool-name">${tool.name}</h3>
          <p class="tool-desc">${tool.description || ""}</p>
          <span class="tool-category">${tool.category || ""}</span>
        </div>
        <a href="${tool.url}" target="_blank" class="tool-link">访问 →</a>
      </div>
    `
      )
      .join("");
  }

  // 渲染项目卡片
  function renderProjects(projects) {
    const projectsGrid = document.getElementById("projectsGrid");
    if (!projectsGrid) return;

    projectsGrid.innerHTML = projects
      .map(
        (project) => `
      <div class="project-card">
        <h3 class="project-name">${project.name}</h3>
        <div class="env-list">
          ${project.environments
            .map(
              (env) => `
            <div class="env-item">
              <div class="env-info">
                <span class="env-name">${env.environmentName}</span>
                <span class="env-status status-${env.statusColor}">${env.status}</span>
              </div>
              <a href="${env.url}" target="_blank" class="env-link">
                ${env.url}
                <span class="link-icon">↗</span>
              </a>
            </div>
          `
            )
            .join("")}
        </div>
      </div>
    `
      )
      .join("");
  }

  // 设置最后更新时间
  const lastUpdateEl = document.getElementById("lastUpdate");
  if (lastUpdateEl) {
    const now = new Date();
    lastUpdateEl.textContent = now.toLocaleString("zh-CN");
  }

  // 搜索功能
  const searchInput = document.getElementById("searchInput");
  const searchResults = document.getElementById("searchResults");

  if (searchInput && searchResults) {
    searchInput.addEventListener("input", function (e) {
      const query = e.target.value.trim().toLowerCase();

      if (query.length === 0) {
        searchResults.classList.remove("active");
        searchResults.innerHTML = "";
        return;
      }

      searchResults.classList.add("active");

      const results = [];

      // 搜索工具
      allTools.forEach((tool) => {
        if (
          tool.name.toLowerCase().includes(query) ||
          (tool.description || "").toLowerCase().includes(query) ||
          (tool.category || "").toLowerCase().includes(query)
        ) {
          results.push({
            type: "工具",
            name: tool.name,
            desc: tool.description || "",
            link: tool.url,
          });
        }
      });

      // 搜索项目环境
      allProjects.forEach((project) => {
        project.environments.forEach((env) => {
          if (
            project.name.toLowerCase().includes(query) ||
            env.environmentName.toLowerCase().includes(query) ||
            env.url.toLowerCase().includes(query)
          ) {
            results.push({
              type: "项目环境",
              name: `${project.name} - ${env.environmentName}`,
              desc: env.url,
              link: env.url,
            });
          }
        });
      });

      // 显示搜索结果
      if (results.length > 0) {
        searchResults.innerHTML = results
          .map(
            (item) => `
          <div class="search-result-item" onclick="window.open('${item.link}', '_blank')">
            <div style="font-weight: 600; color: var(--text); margin-bottom: 4px;">
              ${item.name}
            </div>
            <div style="font-size: 0.85rem; color: var(--text-muted);">
              ${item.type} · ${item.desc}
            </div>
          </div>
        `
          )
          .join("");
      } else {
        searchResults.innerHTML = `
          <div class="search-result-item" style="text-align: center; color: var(--text-muted);">
            未找到相关结果
          </div>
        `;
      }
    });

    // 点击外部关闭搜索结果
    document.addEventListener("click", function (e) {
      if (
        !searchInput.contains(e.target) &&
        !searchResults.contains(e.target)
      ) {
        searchResults.classList.remove("active");
      }
    });
  }

  // 添加卡片点击效果
  const cards = document.querySelectorAll(".tool-card, .project-card");
  cards.forEach((card) => {
    card.addEventListener("click", function (e) {
      // 如果点击的不是链接，则添加点击反馈
      if (!e.target.closest("a")) {
        this.style.transform = "scale(0.98)";
        setTimeout(() => {
          this.style.transform = "";
        }, 150);
      }
    });
  });

  // 键盘快捷键支持
  document.addEventListener("keydown", function (e) {
    // Ctrl/Cmd + K 聚焦搜索框
    if ((e.ctrlKey || e.metaKey) && e.key === "k") {
      e.preventDefault();
      if (searchInput) {
        searchInput.focus();
      }
    }
  });

  // 平滑滚动
  document.querySelectorAll('a[href^="#"]').forEach((anchor) => {
    anchor.addEventListener("click", function (e) {
      e.preventDefault();
      const target = document.querySelector(this.getAttribute("href"));
      if (target) {
        target.scrollIntoView({
          behavior: "smooth",
          block: "start",
        });
      }
    });
  });

  // 页面加载完成后加载数据
  if (document.readyState === "loading") {
    document.addEventListener("DOMContentLoaded", loadData);
  } else {
    loadData();
  }
})();
