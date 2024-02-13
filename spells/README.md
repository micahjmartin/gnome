# Spells

A list collection of useful [Eldritch](https://docs.realm.pub/user-guide/eldritch) functions. These enable advanced functionality within the current limitations of the tools

## List of spells
- [list_permissions](./perms.eldr) - List permissions of files. You cannot list permissions of a single file, you need to iterate the file tree to do so, this function handles that for you 
- [effective_perms](./perms.eldr) - Convert octal permissions to the effective permissions fo the current user (rwx). Eldritch will crash when opening or writing to a file with bad perms, therefor, it is vital to check if a user can access a file before reading or writing to it.
- [glob](./glob.eldr) a basic implementation of Glob