/**
 * Usage:
 *
 * GIT_REF=refs/tags/v1.0.0 node npm/prepare.js
 */
const fs = require("fs");
const path = require("path");

const ref = process.env.GIT_REF; // refs/tags/v1.0.0

const arr = ref.split("/");
const version = arr[arr.length - 1].replace(/^v/, "");

console.log(`prepare publish to npm for: ${version}`);

const packages = fs
  .readdirSync(__dirname)
  .filter((v) => v.startsWith("virtual-node-env-"))
  .concat(["virtual-node-env"]);

for (const pkgName of packages) {
  const pkgPath = path.join(__dirname, pkgName, "package.json");

  const pkg = require(pkgPath);

  pkg.version = version;

  if (pkg.optionalDependencies) {
    for (const subDeps in pkg.optionalDependencies) {
      if (subDeps.startsWith("@axetroy/virtual-node-env-")) {
        pkg.optionalDependencies[subDeps] = version;
      }
    }
  }

  fs.writeFileSync(pkgPath, JSON.stringify(pkg, null, 2));

  if (pkgName.startsWith("virtual-node-env-")) {
    const fileMap = {
      "virtual-node-env-darwin-arm64": "virtual_node_env_darwin_arm64",
      "virtual-node-env-darwin-amd64": "virtual_node_env_darwin_amd64_v1",
      "virtual-node-env-linux-amd64": "virtual_node_env_linux_amd64_v1",
      "virtual-node-env-linux-arm64": "virtual_node_env_linux_arm64",
      "virtual-node-env-windows-amd64": "virtual_node_env_windows_amd64_v1",
      "virtual-node-env-windows-arm64": "virtual_node_env_windows_arm64",
    };

    if (pkgName in fileMap === false)
      throw new Error(`Can not found prebuild file for package '${pkgName}'`);

    const distFolder = fileMap[pkgName];

    const executableFileName =
      "virtual-node-env" + (pkgName.indexOf("windows") > -1 ? ".exe" : "");

    const executableFilePath = path.join(
      __dirname,
      "..",
      "dist",
      distFolder,
      executableFileName
    );

    fs.statSync(executableFilePath);

    fs.readdirSync(path.join(__dirname, "..", "dist", distFolder));

    fs.statSync(path.join(__dirname, pkgName, executableFileName));

    fs.copyFileSync(
      executableFilePath,
      path.join(__dirname, pkgName, executableFileName)
    );
  } else {
    fs.copyFileSync(
      path.join(__dirname, "..", "README.md"),
      path.join(__dirname, pkgName, "README.md")
    );

    fs.copyFileSync(
      path.join(__dirname, "..", "LICENSE"),
      path.join(__dirname, pkgName, "LICENSE")
    );
  }
}
