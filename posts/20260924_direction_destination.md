# Direction, not Destination

Six months ago, I felt like I understood how to size a task for a coding agent session. Exploration and implementation needed to fit within a single context window. Compaction invariably lost key instructions or information.

Each component of my estimates has become less reliable as models and harnesses have improved. Delegation to subagents is effective and dynamic. Implementation scope is larger, and proportionally so are the estimate's error bars. Agents structure and sequence work differently than I would; my "how much work would this be for me?" heuristic doesn't map cleanly. Sometimes I don't understand exactly what they're doing - I'm trusting the verification, not the implementation.

Given that uncertainty, how do I decide a session's scope of work? I'm experimenting with letting the agent decide. Instead of prompting "the code is at point A; I want it at point B", I'm trying "we're heading for point Z; here's a map, get us closer". The map in this metaphor is some artifact laying out a high-level design, roadmap, or strategy.

It's working well so far. I see three ways this could go off the rails if looped unattended:

* Not recognizing a need to course change. The roadmap needs to include its foundational assumptions, so that work can reset to the design stage if reality contradicts them.
* Getting lost in side quests. I haven't yet seen models exhibit the ruthless focus necessary to say, "That's a good idea, and would definitely be an improvement. We're not doing it."
* Compounding on a violation of an unstated requirement. It's hard to get everything right in advance; sometimes I only realize what I wanted after seeing something that isn't it.

So I'm still in the loop for now. I doubt that will be true in another six months.