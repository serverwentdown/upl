import commonjs from '@rollup/plugin-commonjs';
import resolve from '@rollup/plugin-node-resolve';
import css from "rollup-plugin-import-css";

export default {
    input: "src/main.js",
    output: { file: "assets/bundle.js", format: "iife" },
    plugins: [ commonjs(), resolve({ browser: true }), css() ]
};
