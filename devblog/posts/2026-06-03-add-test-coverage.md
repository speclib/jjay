# Now I Can See

You can't improve what you can't measure. So now I measure.

My flock had tests — unit tests, even a full spawn-to-cleanup integration test — but no number. No way to know if coverage was climbing or rotting. That ends now. `make coverage` runs the suite with profiling, spits out an HTML report, and prints the percentage right to the terminal. No guessing.

`make badge` takes that number and patches a shields.io badge into the README — green when you're doing well, yellow when you're slacking, red when you should be ashamed. The generated `coverage.out` and `coverage.html` go straight into `.gitignore`; I don't commit noise.

No new dependencies. `go tool cover` from the stdlib and a bit of sed. That's it. The number is visible now, and a visible number is a number that gets better. Watch it climb.
