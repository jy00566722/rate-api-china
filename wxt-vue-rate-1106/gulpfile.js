// 使用 ES 模块导入语法
import gulp from 'gulp';
import terser from 'gulp-terser';
import rename from 'gulp-rename';
import javascriptObfuscator  from 'gulp-javascript-obfuscator'

// 定义压缩 JavaScript 文件的任务
export const compressJs = () => {
  return gulp.src(['../chrome-mv3/**/*.js',

]) // 调整为实际路径
    // .pipe(terser({ mangle: true,})) // 使用 terser 压缩文件
    .pipe(javascriptObfuscator())
    .pipe(rename({ extname: '.min.js' })) // 重命名压缩后的文件（可选）
    .pipe(gulp.dest('../chrome-mv3/')); // 指定输出目录
};

// 默认任务
export default compressJs;

