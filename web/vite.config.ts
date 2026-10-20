import { defineConfig, loadEnv, ConfigEnv, UserConfigExport } from 'vite'
import path from 'path'
import vue from '@vitejs/plugin-vue'
import vueI18n from '@intlify/vite-plugin-vue-i18n'
import { createSvgIconsPlugin } from 'vite-plugin-svg-icons'
import viteCompression from 'vite-plugin-compression'
import AutoImport from 'unplugin-auto-import/vite'
import Components from 'unplugin-vue-components/vite'
import { ElementPlusResolver } from 'unplugin-vue-components/resolvers'
// https://vitejs.dev/config/
export default ({ mode }: ConfigEnv): UserConfigExport => {
  process.env = {
    ...process.env,
    ...loadEnv(mode, process.cwd()),
    VITE_APP_VERSION: process.env.npm_package_version,
  }
  return defineConfig({
    resolve: {
      alias: [{ find: '@', replacement: path.resolve(__dirname, 'src') }],
    },
    plugins: [
      vue(),
      vueI18n({
        include: path.resolve(__dirname, './src/lang/**'),
      }),
      createSvgIconsPlugin({
        // Specify the icon folder to be cached
        iconDirs: [path.resolve(__dirname, './src/icons/svg')],
        // Specify symbolId format
        symbolId: 'icon-[dir]-[name]',
      }),
      viteCompression({ deleteOriginFile: true }),
      AutoImport({
        // eslintrc: {
        //   enabled: true,
        // },
        resolvers: [ElementPlusResolver()],
      }),
      Components({
        resolvers: [ElementPlusResolver()],
      }),
    ],
    build: {
      // Vite 8 (Rolldown) dropped the object form of output.manualChunks.
      // Keep element-plus and echarts in their own chunks via codeSplitting groups.
      rolldownOptions: {
        output: {
          codeSplitting: {
            groups: [
              {
                name: 'elementUI',
                test: /node_modules[\\/]element-plus[\\/]/,
                priority: 20,
              },
              {
                name: 'echarts',
                test: /node_modules[\\/]echarts[\\/]/,
                priority: 20,
              },
            ],
          },
        },
      },
    },
    server: {
      host: '0.0.0.0',
      port: 8000,
      proxy: {
        // change xxx-api/login => mock/login
        // detail: https://cli.vuejs.org/config/#devserver-proxy
        [process.env.VITE_APP_BASE_API]: {
          target: process.env.VITE_APP_PROXY_TARGET,
          changeOrigin: true,
          rewrite: (path) => {
            const reg = RegExp('^' + process.env.VITE_APP_BASE_API)
            return path.replace(reg, '')
          },
          ws: true,
        },
      },
    },
  })
}
