# The Bitter Lesson for Harness Engineering

Everyone is building or customizing agent harnesses. I've noticed a pattern that I think is a mistake: defining rigidly prescriptive workflows.

It's usually a translation of how the work was done before AI. If it involved a team of humans, the workflow has distinct subagents playing corresponding roles. If it involved navigating a logical decision tree, the workflow has a concrete checklist or flowchart.

[The bitter lesson](http://www.incompleteideas.net/IncIdeas/BitterLesson.html) applies here, as it does to all things AI. Embedding our knowledge of how to do a task may produce better results now, but as models improve, the scaffold becomes a straightjacket.

Don't let instructions only grow, accreting rules and clarifications and edge cases. Anthropic recently [removed most of Claude Code's system prompt](https://claude.com/blog/the-new-rules-of-context-engineering-for-claude-5-generation-models). It's not so different from code: deletions are often more effective than additions.