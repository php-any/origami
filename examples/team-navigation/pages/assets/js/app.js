// 首页脚本

// 初始化搜索引擎
document.addEventListener("DOMContentLoaded", function () {
  // 绑定搜索表单提交事件
  const searchForm = document.getElementById("searchForm");
  if (searchForm) {
    searchForm.addEventListener("submit", function (e) {
      e.preventDefault();
      performSearch();
    });
  }

  // 更新时间显示
  updateTime();
  setInterval(updateTime, 1000);
});

// 执行搜索
function performSearch() {
  const searchInput = document.getElementById("searchInput");
  const keyword = searchInput.value.trim();

  if (!keyword) {
    return;
  }

  // 使用全局变量 defaultSearchEngine（从服务器注入）
  if (typeof defaultSearchEngine === "undefined" || !defaultSearchEngine) {
    // 如果没有默认搜索引擎，使用百度搜索作为后备
    window.open("https://www.baidu.com/s?wd=" + encodeURIComponent(keyword), "_blank");
    return;
  }

  // 替换 URL 模板中的 {keyword} 占位符
  const searchUrl = defaultSearchEngine.urlTemplate.replace("{keyword}", encodeURIComponent(keyword));
  window.open(searchUrl, "_blank");
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
    const weekdays = ["星期日", "星期一", "星期二", "星期三", "星期四", "星期五", "星期六"];
    const year = now.getFullYear();
    const month = now.getMonth() + 1;
    const day = now.getDate();
    const weekday = weekdays[now.getDay()];
    dateEl.textContent = `${year}年${month}月${day}日 ${weekday}`;
  }
}

