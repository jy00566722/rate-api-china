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
  }),
  manifest: {
    name: '加密货币价格',
    version: '0.0.1',
    version_name: "0.0.1",
    description: '加密货币价格',
    icons: {
      16: 'icons/16.png',
      48: 'icons/48.png',
      128: 'icons/128.png',
    },
    permissions: ["storage","alarms",],
    incognito:"spanning",
    host_permissions:[
    
  ],
    
  },
});
