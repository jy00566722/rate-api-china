import { defineConfig } from 'wxt';
import vue from '@vitejs/plugin-vue';

// See https://wxt.dev/api/config.html
export default defineConfig({
  imports: {
    addons: {
      vueTemplate: true,
    },
  },
  vite: () => ({
    plugins: [vue()],
    // resolvers: [
    //   // ElementPlusResolver(),

    //   // Auto import icon components
    //   // 自动导入图标组件
    //   IconsResolver({
    //     prefix: 'Icon',
    //   }),
    // ],
    build: {
      rollupOptions: {
        external: []
      }
    }
  }),
  manifest: {
    name: '汇率转换',
    version: '3.7.3.0',
    version_name: "3.7.3.0",
    description: '此插件为跨境卖家设计,用于显示跨境平台前端页面上商品价格对应的人民币价格。同时也提供手动输入数字转换货币的功能,货币种类共有163个',
    icons: {
      16: 'icons/16.png',
      48: 'icons/48.png',
      128: 'icons/128.png',
    },
    permissions: ["storage","alarms", "scripting"],
    incognito:"spanning",
    host_permissions:[
      "https://rate.lizudi.top/" 
  ],
    optional_host_permissions:[
      "https://*/*",
      "http://*/*"

  ],
    
  },
});
