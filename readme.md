# Preface 
Hola! I made this utility CLI to assist/expedite the use-cases of my day to day work. :) 

## Features: <br/> 

**[Copy root path to clipboard]** ✅ <br/>
- Command: **clip-file-contents**
- This command will copy the provided root path's contents recursively into the clipboard, except for .GIT, IDE generator (e.g _.git_, _.idea_), or irrelevant files
not inclusive to human generated content.
- It adds a header that identifies the filename associated with the contents.

**[Clip GPT Code Standards Preface]** ✅ <br/>
- Command: **clip-gpt-preface**
- This command copies a **code preface for Chat GPT** that improves code quality.
  ChatGPT code is usually "decent/passable" if configured to run on **version 4**, but the good engineering
  foundations usually still has something to be desired.


- It attempts to remediate the lack of: <br/>
  	_Defensive programming_, _testability_, _readability_, _modularity_ - and this
  	preface attempts to achieve that. It obviously will not be perfect,
  	but it gives tangible improvements (at least based on anecdotal experience).

**[Copy one folder to another]** ✅ <br/>
- This command copies one folder's contents to another, and at least has (not 100% enumerated here) the ff constraints:
  - **exclusions**: folder A may omit certain folders to copy into folder B
  - **wipe folder B**: folder B may be wiped clean, before folder A is copied into it; but it also has constraints on files not to wipe
- The main use-case for this feature is transferring one repository to another (folder A -> B), and preserving their respective trackers.

### Todo roadmap: <br/>
- Make files
- bash scripts to compile the binaries, and integrate into _.zshrc_, _.bashrc_
- Improved CLI docs (long, and short) in the Viper commands
- Automated testing CI on PRs, and merge requests (this is probably what I want to do next)
- Docker configuration (likely if we're going to evolve to have a backend)