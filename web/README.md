# /web: sweet place for static assets

This directory contains files that should be served by server as-is if requested. For example accessing `https://refto.dev/something.html` will serve `./web/something.html` (only if server have not defined its own handler for `/something.html` and file `./web/something.html` exists)

Because contents of this directory will not be included in binary, you should deploy this directory to server's working dir.

Some cases when you want to use this dir:
   - serving frontend
   - serving anything that should be accessed as is (and not interfere with server's routes)
