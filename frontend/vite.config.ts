import { defineConfig, loadEnv } from 'vite'
import react from '@vitejs/plugin-react'

// https://vitejs.dev/config/
export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, process.cwd(), '')
  const port = parseInt(env.VITE_PORT || '3000', 10)

  return {
    plugins: [react()],
    server: {
      host: true,
      port: port
    },
    preview: {
      host: true,
      port: port
    }
  }
})
