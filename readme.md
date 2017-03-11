## Log

[![GoDoc](https://godoc.org/github.com/dist-ribut-us/log?status.svg)](https://godoc.org/github.com/dist-ribut-us/log)

### ToDo

Stack tracing. Log n up the stack or log the whole stack.

Buffered logs - particularly good for debugging. A log that writes to a buffer
and only on Commit is it dumped to the log.

Should there be an option to encrypt logs?

How to purge logs? We want to keep data around, but the file should be limited
in both size and time. Maybe when openning there's a chance to check it. Or
create a LogFile type. Give it a buffer and have it occasionally clean out the
file, the buffer can hold content until it's done.
