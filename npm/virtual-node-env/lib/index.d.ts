import { SpawnOptionsWithoutStdio } from "child_process";

declare var exec: (argv: string, options?: SpawnOptionsWithoutStdio) => void;

export default exec;
export { exec };
