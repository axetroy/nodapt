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
  .filter((v) => v.startsWith("nodapt-"))
  .concat(["nodapt"]);

for (const pkgName of packages) {
  const pkgPath = path.join(__dirname, pkgName, "package.json");

  const pkg = require(pkgPath);

  pkg.version = version;

  if (pkg.optionalDependencies) {
    for (const subDeps in pkg.optionalDependencies) {
      if (subDeps.startsWith("@axetroy/nodapt-")) {
        pkg.optionalDependencies[subDeps] = version;
      }
    }
  }

  fs.writeFileSync(pkgPath, JSON.stringify(pkg, null, 2));

  if (pkgName.startsWith("nodapt-")) {
    const fileMap = {
      "nodapt-darwin-arm64": "nodapt_darwin_arm64_v8.0",
      "nodapt-darwin-amd64": "nodapt_darwin_amd64_v1",
      "nodapt-linux-amd64": "nodapt_linux_amd64_v1",
      "nodapt-linux-arm64": "nodapt_linux_arm64_v8.0",
      "nodapt-windows-amd64": "nodapt_windows_amd64_v1",
      "nodapt-windows-arm64": "nodapt_windows_arm64_v8.0",
    };

    if (pkgName in fileMap === false)
      throw new Error(`Can not found prebuild file for package '${pkgName}'`);

    const destFolder = fileMap[pkgName];

    const executableFileName =
      "nodapt" + (pkgName.indexOf("windows") > -1 ? ".exe" : "");

    const executableFilePath = path.join(
      __dirname,
      "..",
      "dist",
      destFolder,
      executableFileName
    );

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
