/**
 * Usage:
 *
 * NODE_AUTH_TOKEN=npm_xxxxx node npm/publish.js
 */

const fs = require("fs");
const path = require("path");
const spawn = require("child_process").spawn;

const packages = fs
  .readdirSync(__dirname)
  .filter((v) => v.startsWith("nodapt-"))
  .concat(["nodapt"]);

async function main() {
  for (const pkg of packages) {
    const cwd = path.join(__dirname, pkg);

    await new Promise((resolve, reject) => {
      const ps = spawn("npm", ["publish", "--access", "public"], {
        cwd: cwd,
        stdio: "inherit",
        env: process.env,
      });

      ps.on("close", (code) => {
        code === 0
          ? resolve()
          : reject(new Error(`Process exist with ${code}`));
      });
    });
  }
}

main().catch((err) => {
  console.error(err);
  process.exit(1);
});
