# Eldritch Implementation

## Differences Between [Eldritch](https://docs.realm.pub/user-guide/eldritch) and Gnome
Not all functions are implemented, but this is something that can change. Make an issue or PR if you would like a specific function implmented. I am currently only implementing functions that I need/want as I need them

Other small differences are listed below:
- [ ] `process.kill` takes an optional kill signal to send
- [ ] `assets` is backed by an embed.FS or any other fs.FS compatible interface
- [ ] `sys.set_env` sets an environment variable
- [ ] `exit(int)` function has been added to allow any script to kill the interpreter completely
- [ ] `quit()` has been added to allow any script to stop execution of itself
- [ ] `fallback(cmd)` stop execution of all scripts and load the given command string in place. First trying the path from assets and falling back to the filesystem 
- [ ] Global variables are preserved across script executions allowing for data to be passed around