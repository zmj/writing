# Feature Factories

A common assertion in discussions of AI in software development is that AI-generated codebases are messy, and trend increasingly so over time. This is user error: if you only ever ask the AI to write features, you're operating a feature factory.

The term "feature factory" was coined to describe organizations that exclusively prioritize and reward launching new features, neglecting maintainability. It's meant to evoke a sense of dirty, mechanical, disempowering work: a developer working the factory line doesn't have the agency to apply their craftsmanship to the product or the process.

The result is a codebase with little coherent structure or internal consistency. Each new feature added finds it easier to reimplement existing functionality rather than understand and extend it. This compounds over time. The lack of a consistent pattern to follow is itself a pattern that informs the next piece of work.

If you agreed with the first sentence, that description should be familiar to you.

Developers sometimes cope with that kind of working environment by smuggling in refactors along with feature work. To a small degree, that's healthy application of the campsite rule (leave the code a little better than you found it). When the refactor dominates the scope of the net new code, that's maladaptive behavior. It undermines the effort/reward estimate that the organization prioritizes by, and doesn't fix the dysfunction that maintenance is illegible and unappreciated.

Organizations don't want that behavior from developers, and developers don't want that behavior from coding agents. Even if agents did redefine task scope on their own authority, they wouldn't be good at it. The context that informs a human developer's refactor is rarely found next to the new code's site. It's somewhere else in time ("didn't another team do something like this last month?") or space ("I'm getting deja vu writing this, I must have seen it before"). Coding agents don't have those memories, and searching without direction would be inefficient use of limited context windows.

The good news is that the industry already understands how to avoid this failure mode: explicitly reserve capacity for maintenance. The same solution applies to a developer with a coding agent. There are many ways one could put that into practice - and they're all better than doing nothing.