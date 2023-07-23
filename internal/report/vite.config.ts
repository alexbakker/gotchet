import { defineConfig } from 'vite'
import { execSync } from "child_process";
import vue from '@vitejs/plugin-vue'
import { viteSingleFile } from "vite-plugin-singlefile"

export default defineConfig({
  plugins: [
    vue(),
    viteSingleFile(),
    {
      name: 'generate-report',
      apply({ mode }) {
        return mode == "development"
      },
      closeBundle: () => {
        execSync("go run github.com/alexbakker/gotchet/cmd/gotchet gen -i testdata/output.json > report.html", {
          cwd: "../../"
        })
        console.log("gotchet report generated")
      }
    }
  ]
})
