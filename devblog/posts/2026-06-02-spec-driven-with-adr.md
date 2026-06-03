# Decisions That Don't Vanish

The "use Go" decision was rotting in an archive folder. Not anymore.

Here's the problem: architectural decisions get buried with the change that made them. Future contributors dig through `archive/.../proposal.md` to find out *why* — if they find it at all. That's no way to run a flock. So I forked the schema into `spec-driven-with-adr` and gave it an `adr` artifact that writes to `openspec/adrs/` and lives *outside* the change lifecycle.

ADRs persist. They sit next to design, never get archived, stay discoverable. I had the agents write retroactive ones for the decisions we'd already made — ADR-001 for choosing Go, ADR-002 for the config context. Reasoning, preserved.

ADRs aren't mandatory and the artifact flow didn't change — proposal, specs, design, tasks, same as ever. This is additive. Architectural decisions just stopped disappearing.
