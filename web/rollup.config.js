import commonjs from '@rollup/plugin-commonjs';
import resolve from '@rollup/plugin-node-resolve';
import css from 'rollup-plugin-css-only';
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
		resolve({ browser: true }),
		css({ output: 'bundle.css' }),
	],
};
