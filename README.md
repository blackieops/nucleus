# nucleus: self-hosted file storage backend

Nucleus helps you sync and store your files on your own infrastructure.

⚠️ **Warning:** this is a very early-stage project that is NOT feature complete,
safe to use, nor secure. Do not use this in production yet.

## Project Goals

Nucleus' primary goal is to be a fast, secure, and light-weight cloud file
storage and sync backend.

1. **Pluggable Storage** - Nucleus should allow configuring different backing
   storage implementations, for example with cloud object storage.
2. **Files only** - Nucleus does not want to replace your entire workgroups
   solution, your project management software, your video call software, your
   chat software... Nucleus is exclusively a file storage solution.
4. **[Nextcloud][nc]-compatible** - To easily fit into existing workflows, Nucleus
   should be reasonably compatible with Nextcloud client apps, exposing an API
   they expect so you can swap your Nextcloud server for Nucleus and pick up
   where you left off with minimal interruption.
5. **Declarative configuration** - Nucleus should not require mutating itself
   on-disk or writing any configuration.
6. **Scale-out, not up** - Nucleus should easily allow for scale-out
   infrastructure, support immutable filesystems, and fit right at home in a
   cloud-native environment.
7. **Gotta go fast** - Nucleus should be fast, efficient, and require as few
   server resources as possible.

[nc]: https://www.nextcloud.com
