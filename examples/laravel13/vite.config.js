import { defineConfig } from 'vite';
import laravel from 'laravel-vite-plugin';
import { bunny } from 'laravel-vite-plugin/fonts';
import tailwindcss from '@tailwindcss/vite';
import { readFileSync } from 'node:fs';
import { createRequire } from 'node:module';
import { fileURLToPath } from 'node:url';
import postcss from 'postcss';

const require = createRequire(import.meta.url);
const tablerEntry = fileURLToPath(new URL('./resources/css/admin/tabler.css', import.meta.url)).replaceAll('\\', '/');

// Keep the upstream stylesheet intact apart from its scope root. The boundary
// prevents Bootstrap resets from entering Filament's forms, tables and dialogs.
function tablerTemplate() {
    return {
        name: 'tabler-template',
        enforce: 'pre',
        load(id) {
            if (id.split('?')[0].replaceAll('\\', '/') !== tablerEntry) return;
            const source = require.resolve('@tabler/core/dist/css/tabler.css');
            this.addWatchFile(source);
            const root = postcss.parse(readFileSync(source, 'utf8'), { from: source });
            root.walkAtRules('charset', (rule) => rule.remove());
            root.walkComments((comment) => {
                if (comment.text.startsWith('# sourceMappingURL=')) comment.remove();
            });
            root.walkRules((rule) => {
                rule.selector = rule.selector.replace(/(^|[,\s])(?::root|html|body)(?=[\s,.#[:>+~]|$)/g, '$1:scope');
            });
            return `@scope (.tabler-ui) to (.filament-content) {\n${root.toString()}\n}`;
        },
    };
}

export default defineConfig({
    plugins: [
        tablerTemplate(),
        laravel({
            input: ['resources/css/app.css', 'resources/js/app.js', 'resources/css/filament/admin/theme.css', 'resources/css/admin/tabler.css', 'resources/js/admin.js'],
            refresh: true,
            fonts: [
                bunny('Instrument Sans', {
                    weights: [400, 500, 600],
                }),
            ],
        }),
        tailwindcss(),
    ],
    server: {
        watch: {
            ignored: ['**/storage/framework/views/**'],
        },
    },
});
