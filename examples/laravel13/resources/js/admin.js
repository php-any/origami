import '@tabler/core/dist/js/tabler.esm.js';

// Mirror Tabler's saved layout attributes into the isolated CSS scope roots.
function syncTemplateState() {
    document.querySelectorAll('.tabler-ui').forEach((root) => {
        ['data-bs-sidebar', 'data-bs-layout', 'data-bs-navbar'].forEach((name) => {
            const value = document.documentElement.getAttribute(name);
            if (value === null) root.removeAttribute(name);
            else root.setAttribute(name, value);
        });
    });
    const folded = (document.documentElement.getAttribute('data-bs-sidebar') ?? '').startsWith('folded');
    document.querySelectorAll('[data-bs-toggle="sidebar-folded"]').forEach((toggle) => {
        toggle.setAttribute('aria-pressed', String(folded));
    });
}

try {
    const sidebar = localStorage.getItem('tabler-sidebar');
    if (sidebar === 'folded' || sidebar === 'folded-hover') {
        document.documentElement.setAttribute('data-bs-sidebar', sidebar);
    }
} catch { /* Storage can be unavailable in private browsing. */ }

new MutationObserver(syncTemplateState).observe(document.documentElement, {
    attributes: true,
    attributeFilter: ['data-bs-sidebar', 'data-bs-layout', 'data-bs-navbar'],
});
document.addEventListener('DOMContentLoaded', syncTemplateState);
document.addEventListener('livewire:navigated', syncTemplateState);
syncTemplateState();

// Keep the current module visible in short desktop viewports without scrolling
// the main page. Mobile navigation retains Bootstrap's collapse behavior.
function revealCurrentNavigation() {
    if (!window.matchMedia('(min-width: 992px)').matches) return;
    const menu = document.getElementById('sidebar-menu');
    const current = menu?.querySelector('[aria-current="page"]');
    if (!current) return;
    const bounds = menu.getBoundingClientRect();
    const active = current.getBoundingClientRect();
    if (active.bottom > bounds.bottom) menu.scrollTop += active.bottom - bounds.bottom + 12;
    else if (active.top < bounds.top) menu.scrollTop -= bounds.top - active.top + 12;
}

document.addEventListener('DOMContentLoaded', revealCurrentNavigation);
document.addEventListener('livewire:navigated', revealCurrentNavigation);
