import commonjs from '@rollup/plugin-commonjs';
import resolve from '@rollup/plugin-node-resolve';
import postcss from 'rollup-plugin-postcss';
import copy from 'rollup-plugin-copy';
import { terser } from 'rollup-plugin-terser';

export default {
    input: 'src/main.js',
    output: {
		file: 'assets/bundle.js',
		format: 'iife',
		plugins: [
			terser(),
		],
	},
    plugins: [
		commonjs(),
		resolve({
			browser: true,
		}),
		postcss({
			plugins: [
				require('tailwindcss'),
			],
			extract: true,
			minimize: true,
		}),
		copy({
			targets: [
				{ src: 'static/favicon.png', dest: 'assets' },
			],
		}),
	],
};
