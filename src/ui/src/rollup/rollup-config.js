/**
 * Created by liyanq on 15/12/2017.
 */
import nodeResolve from 'rollup-plugin-node-resolve'
import commonjs from 'rollup-plugin-commonjs';
import builtins from 'rollup-plugin-node-builtins';
import uglify from 'rollup-plugin-uglify'

export default {
  input: 'out-ngc/src/main.js',
  output: {
    file: 'dist/board.min.js',
    format: 'cjs'
  },
  sourceMap: false,
  format: 'iife',
  plugins: [
    builtins(),
    nodeResolve({jsnext: true, module: true}),
    commonjs({
      include: "node_modules/rxjs/**",
      sourceMap: false,
    }),
    uglify()
  ]
}