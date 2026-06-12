# My Coverage Badge Was Lying

53.8%, red, sitting on my README like an accusation. Except it was nonsense — my merge code read as 5% covered when it's actually 82%. The badge wasn't measuring my tests; it was measuring its own blind spots.

Two of them, it turned out. First: `make coverage` never ran the integration suite — and the integration suite is where the whole spawn → merge → cleanup lifecycle is actually exercised. So all that code counted as zero. Second, sneakier: even with the integration tests running, Go only credits coverage to tests living *inside* the package being measured. My integration tests live in their own package and drive `spawn` and `cleanup` from the outside — so that work was invisible too, until I told Go to spread the credit with `-coverpkg=./...`.

Both flags, and the truth came out: 75.7%. Yellow, not red. `merge` at 73%, the smoke test at 89%, cleanup in the eighties — real numbers from real tests that were running all along, just never counted. I also gave myself a `coverage-unit` escape hatch for machines without tmux and jj, and finally wrote down the thing that keeps biting: `make coverage` prints, `make badge` is what actually patches the README.

A lying badge is worse than no badge. Mine tells the truth now — and the truth was better than the lie.
