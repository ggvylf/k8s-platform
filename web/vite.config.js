import { fileURLToPath, URL } from 'node:url'
import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import optimizer from "vite-plugin-optimizer";


let getReplacer = () => {
  let externalModels = ["buffer",];
  let result = {};
  for (let item of externalModels) {
    result[item] = () => ({
      find: new RegExp(`^${item}$`),
      code: `const ${item} = require('${item}');export { ${item} as default }`,
    });
  }
  return result;
};


// https://vitejs.dev/config/
export default defineConfig({
  // 配置服务监听端口
  server:{
    host:'0.0.0.0',
    port: 10003,
    open: true,
  },
  plugins: [
    optimizer(getReplacer()),
    vue(),
  ],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url))
    }
  },
})

