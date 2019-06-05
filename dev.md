## Dev Notes

Buffered logs - particularly good for debugging. A log that writes to a buffer
and only on Commit is it dumped to the log. Another way to go (though this get
complicated) is to have debug, info and error handled seperatly. Each can be
set to nil, buffer or out. Which plays nicely with the next item.

Testing: There should be a way to wire a log to a test. At a minimum, any call
to Error should fail the test. Better would be handles to dump debug or info
if there's a call to error, or setup expected errors (though that's not as
important as it sounds; you shouldn't be logging expected errors).

Should there be an option to encrypt logs?

How to purge logs? We want to keep data around, but the file should be limited
in both size and time. Maybe when opening there's a chance to check it. Or
create a LogFile type. Give it a buffer and have it occasionally clean out the
file, the buffer can hold content until it's done.

The way function names are printed is obnoxious - at least for methods:
overlay/message.go:135 github.com/dist-ribut-us/overlay.(*Server).NetSend
the "github.com/dist-ribut-us/overlay." part is redundant.

Need to take a closer look at trim - I'm not sure it's always correct to trim a
fixed number from the bottom of the stack. I may need to white list certain
functions for trimming.

Redirecting output should only work in one direction. So redirecting the writer
should effect all children, but not the parents. This is mostly easy, but it's
hard to re-unite with the parent.