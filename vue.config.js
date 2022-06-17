const {defineConfig} = require('@vue/cli-service');
const path = require('path');

module.exports = defineConfig({
  transpileDependencies: true,
  configureWebpack: (config) => {
    console.log(config.entry);
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
