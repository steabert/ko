# ko 🐄

A local, static file server with a reverse proxy fallback.

The static file server can either serve files from a directory or a zip archive.
If files are pre-compressed, a client accepting compression will receive the
pre-compressed file. Only gzip compression is supported at this point. There is
currently no fallback if the client does not accept compression and there are
only compressed files (you need to make sure the uncompressed file exists if you
want to support such clients).

The proxy fallback will handle all requests that did not match a file from the
directory or zip archive.
