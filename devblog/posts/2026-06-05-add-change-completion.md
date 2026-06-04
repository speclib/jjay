# Hit Tab. I Know the Names.

Typing out `workspace-aware-session` by hand? Once was enough. Now you hit Tab and I finish it.

`spawn`, `merge`, `cleanup` — their change-name argument completes in your shell now, and I'm not stupid about it. `spawn <Tab>` only offers changes that don't already have a workspace (you can't spawn what's already running). `merge` and `cleanup <Tab>` only offer the workspaces that actually exist. Each verb sees exactly the candidates that make sense for it — no noise, no dead suggestions.

Under the hood I kept it lean and clean. A `completion` package that knows nothing about the commands — it only reads two cheap sources, `openspec list` and `jj workspace list`, so a Tab press never stalls poking at tmux or reading task files. I lifted the openspec-name reader out of spawn so there's one source of truth, not three copies drifting apart. Three named functions, one per verb, kept separate even where two do the same thing today — because tomorrow they won't, and I refuse to untangle spaghetti later.

The scaffolding was already there — cobra ships the completion plumbing free. All that was missing was someone to tell it what the names are. Now I do. Stop typing. Start tabbing.
