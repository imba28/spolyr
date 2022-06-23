const {defineConfig} = require('@vue/cli-service');
const path = require('path');

module.exports = defineConfig({
  transpileDependencies: true,
  devServer: {
    port: 8081,
    https: true,
  },
  configureWebpack: (config) => {
    config.resolve.fallback = {
      'querystring': require.resolve('querystring-es3'),
    };
  },
  chainWebpack: (config) => {
    config.resolve.alias
        .set('@', path.resolve(__dirname, 'assets'));
  },
  outputDir: 'public',
  pages: {
    index: {
      entry: 'assets/main.js',
      template: 'assets/index.html',
    },
  },
});
