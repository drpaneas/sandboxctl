# What is it?

This Go utility can take a JIRA project key and an API token and walk across the board to report on how connected outcomes are to stories, and stories back up to outcomes, assessing the tracebility of work in that project.

# Criteria checked
- Outcomes without features
- Features without outcomes 
- Features without epics
- Epics without features
- Epics without stories
- Stories without epics
- Summary:
	- Total untracked JIRA issues n out of x total JIRA issues.
	- Tracked work percentage: n%
	- Health Score: <letter grade>

# Can I use it?

 - This was built for the SANDBOX team to monitor their JIRA project,  but can be used by others with some modifications.  
 - After checking out the code, modify the `project` referenced in `config.go` from
   `SANDBOX` to your JIRA project key.
	 - You'll need it to authenticate with JIRA, you can set up an env variable (`JIRA_API_TOKEN`) in your
   shell with your token, or in a pinch, you can embed it into the code
   in `config.go` by setting it directly (don't commit that back and compromise your token).
   - If you use Mac OSX and want to create the variable, you'll need to edit an existing `~/.zshrc` file (or create it, if there isn't one) and add a new line to export the token, like `export JIRA_API_TOKEN="<yourtoken>"`.
   - JIRA Personal Access Tokens can be generated from your JIRA profile page.  
  - Run the script from its directory with `go run .`. 

# Notes
For best results, your local go version should be 1.21.1
