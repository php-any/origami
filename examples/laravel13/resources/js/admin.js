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
