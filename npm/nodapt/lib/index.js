const spawn = require("child_process").spawn;
const os = require("os");
const path = require("path");

const platform = os.platform();
const arch = os.arch();

const ERR_NOT_SUPPORT = new Error("nodapt does not support your platform");

const platformMap = {
  win32: {
    arm64: "nodapt-windows-arm64",
    x64: "nodapt-windows-amd64",
  },
  darwin: {
    arm64: "nodapt-darwin-arm64",
    x64: "nodapt-darwin-amd64",
  },
  linux: {
    arm64: "nodapt-linux-arm64",
    x64: "nodapt-linux-amd64",
  },
};

const archMap = platformMap[platform];

if (!archMap) throw ERR_NOT_SUPPORT;

const prebuildPackageName = archMap[arch];

if (!prebuildPackageName) throw ERR_NOT_SUPPORT;

const binaryPackageDir = (() => {
  try {
    return path.dirname(
      require.resolve(`@axetroy/${prebuildPackageName}/package.json`)
    );
  } catch (err) {
    throw new Error(
      `Can't find the binary package "${prebuildPackageName}" in the node_modules, try to reinstall package.`
    );
  }
})();

const executableFileName = "nodapt" + (platform === "win32" ? ".exe" : "");

const executableFilePath = path.join(binaryPackageDir, executableFileName);

/**
 *
 * @param {Array<string>} argv
 * @param {import('child_process').SpawnOptionsWithoutStdio} [spawnOptions]
 * @returns
 */
function exec(argv, spawnOptions = {}) {
  const ps = spawn(executableFilePath, argv, {
    ...spawnOptions,
    stdout: "piped",
  });

  return ps;
}

module.exports = exec;
module.exports.exec = exec;
