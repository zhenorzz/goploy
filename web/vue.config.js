module.exports = {
  devServer: {
    proxy: {
      '/api': {
        target: 'http://localhost:9091',
        changeOrigin: true, // 如果接口跨域，需要进行这个参数配置
        pathRewrite: {
          '^/api': '',
        },
      },
    },
  },
  chainWebpack: (config) => {
    config.module
        .rule('svg')
        .exclude.add(resolve('src/icons'))
        .end();

    config.module
        .rule('icons')
        .test(/\.svg$/)
        .include.add(resolve('src/icons'))
        .end()
        .use('svg-sprite-loader')
        .loader('svg-sprite-loader')
        .options({
          symbolId: 'icon-[name]',
        });
  },
};

const path = require('path');

/**
 * @param  {string} dir
 * @return {string}
 */
function resolve(dir) {
  return path.join(__dirname, './', dir);
}
