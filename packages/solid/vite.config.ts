// SPDX-FileCopyrightText: 2025 caixw
//
// SPDX-License-Identifier: MIT

import { defineConfig } from 'vite';
import dts from 'vite-plugin-dts';
import solidPlugin from 'vite-plugin-solid';
import { viteStaticCopy } from 'vite-plugin-static-copy';

// https://vitejs.dev/config/
export default defineConfig({
    resolve: {
        extensions: ['.tsx', '.ts']
    },
    plugins: [
        solidPlugin(),
        dts({
            entryRoot: './src',
            insertTypesEntry: true,
            rollupTypes: true,
            exclude: [
                'node_modules/**',
                '**/lib/**',
                './src/**/*.spec.ts',
            ]
        }),
        viteStaticCopy({
            targets: [
                { src: '../../LICENSE', dest: '../' },
                { src: '../../.browserslistrc', dest: '../' },
            ]
        })
    ],

    build: {
        minify: true,
        sourcemap: false,
        outDir: './lib',
        target: 'ESNext',
        lib: {
            entry: {
                'index': './src/index.ts',
            },
            formats: ['es'],
            fileName: (_, name) => `${name}.js`
        }
    }
});
