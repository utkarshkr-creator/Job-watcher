# Verification Checklist

Use this checklist to verify that all fixes are working correctly.

## ✅ Pre-Deployment Checklist

### Local Testing
- [ ] Code compiles: `go build`
- [ ] Test sources work: `go run . --test-sources`
- [ ] Full run works: `go run .`
- [ ] At least 3-5 sources return jobs
- [ ] No compilation errors

### Configuration
- [ ] `config.yaml` has your keywords
- [ ] `config.yaml` has your locations
- [ ] At least 5 sources are enabled
- [ ] `.env` file exists with TG_TOKEN and TG_CHAT (for local testing)

### GitHub Setup
- [ ] Repository secrets are set:
  - [ ] `TG_TOKEN`
  - [ ] `TG_CHAT`
  - [ ] `GEMINI_API_KEY` (if using AI)
  - [ ] `RESUME_TEXT` (if using AI)
- [ ] Workflow permissions set to "Read and write"
  - Settings → Actions → General → Workflow permissions

## ✅ Post-Deployment Checklist

### After Pushing Changes
- [ ] Commit and push all changes
- [ ] GitHub Actions workflow appears in Actions tab
- [ ] No syntax errors in workflow files

### First Manual Run
- [ ] Go to Actions tab
- [ ] Click "Job Watcher" workflow
- [ ] Click "Run workflow"
- [ ] Wait for completion (should take 30-60 seconds)
- [ ] Check run status (should be green ✓)

### Check Logs
- [ ] Expand "Run Job Watcher" step
- [ ] See "🚀 Fetching jobs in parallel..."
- [ ] See multiple "✓ SourceName: X jobs" lines
- [ ] See "Total jobs fetched: X" (should be 1000+)
- [ ] See "New jobs matching keywords: X"
- [ ] No red error messages that stop execution

### Expected Output Example
```
🚀 Fetching jobs in parallel...
  ✓ Indeed: 96 jobs
  ✓ LinkedIn: 20 jobs
  ✓ Razorpay: 1 jobs
  ✓ HN Jobs: 21 jobs
  ✓ Reddit: 10 jobs
  ✓ Companies: 1006 jobs
  ✓ Shared Lists: 1515 jobs
  YC API returned status 403, trying HTML fallback
  ✓ YC Jobs: 15 jobs
  Instahyre API error: connection timeout
  ✓ Instahyre: 0 jobs

⏱️  Fetched in 45.2 seconds
Total jobs fetched: 2684
New jobs matching keywords: 15
```

### Telegram Notifications
- [ ] Received Telegram message (if new jobs found)
- [ ] Message format is correct
- [ ] Job links are clickable
- [ ] No duplicate notifications on second run

### Jobs.json Tracking
- [ ] `jobs.json` file created/updated
- [ ] Commit from `github-actions[bot]` appears
- [ ] File contains job records with timestamps

## ✅ Ongoing Monitoring (First 24 Hours)

### After 1 Hour
- [ ] Second automatic run completed
- [ ] No duplicate notifications
- [ ] `jobs.json` updated again

### After 6 Hours
- [ ] Multiple runs completed successfully
- [ ] Consistent job counts from sources
- [ ] No workflow failures

### After 24 Hours
- [ ] At least 20 runs completed
- [ ] Receiving relevant job notifications
- [ ] No rate limiting issues
- [ ] No repeated errors in logs

## ✅ Debug Workflow Test

### Run Debug Workflow
- [ ] Go to Actions tab
- [ ] Click "Debug Job Sources"
- [ ] Click "Run workflow"
- [ ] Wait for completion
- [ ] Check detailed output for each source
- [ ] Network connectivity tests pass

## ✅ Source-Specific Checks

### Indeed
- [ ] Returns 50-100+ jobs
- [ ] RSS feeds in config are valid
- [ ] Jobs are recent (within max_days_old)

### LinkedIn
- [ ] Returns 10-30+ jobs
- [ ] No 429 rate limit errors
- [ ] Jobs match keywords

### Companies
- [ ] Returns 500-1000+ jobs
- [ ] Takes 20-40 seconds (normal)
- [ ] Multiple companies listed

### Shared Lists
- [ ] Returns 1000-2000+ jobs
- [ ] Google Sheets/GitHub sources accessible
- [ ] Jobs are diverse

### HN Jobs
- [ ] Returns 10-30+ jobs
- [ ] Algolia API working
- [ ] Recent "Who's Hiring" thread found

### Reddit
- [ ] Returns 5-15+ jobs
- [ ] JSON API working
- [ ] Posts from r/cscareerquestions and r/forhire

### YC Jobs
- [ ] Returns 0-20 jobs (may vary)
- [ ] If 0, check logs for fallback attempts
- [ ] No workflow-breaking errors

### Instahyre
- [ ] Returns 0-10 jobs (may vary)
- [ ] If 0, check logs for API errors
- [ ] Fallback attempts logged

## ✅ Configuration Tuning

### If Too Many Jobs
- [ ] Narrow keywords in config.yaml
- [ ] Reduce max_experience_years
- [ ] Add more exclude_keywords
- [ ] Reduce max_days_old

### If Too Few Jobs
- [ ] Broaden keywords
- [ ] Increase max_experience_years
- [ ] Remove some exclude_keywords
- [ ] Increase max_days_old
- [ ] Enable more sources

### If Wrong Job Types
- [ ] Review and update keywords
- [ ] Add exclude_keywords for unwanted roles
- [ ] Adjust experience level filter
- [ ] Enable/disable specific sources

## ✅ AI Matching (If Enabled)

### AI Configuration
- [ ] `resume.txt` exists or `RESUME_TEXT` secret set
- [ ] AI provider configured (ollama or gemini)
- [ ] Model name is correct
- [ ] Threshold is reasonable (70-80)

### AI Execution
- [ ] See "🤖 Running AI Matcher..." in logs
- [ ] Jobs are scored
- [ ] Scores appear in job titles: "[AI: 85] Job Title"
- [ ] Only high-scoring jobs sent to Telegram

## ✅ Performance Checks

### Execution Time
- [ ] Total run time: 30-60 seconds (without AI)
- [ ] Total run time: 1-5 minutes (with AI)
- [ ] No timeouts
- [ ] Parallel fetching working

### Resource Usage
- [ ] Workflow completes within free tier limits
- [ ] No memory issues
- [ ] No rate limiting from job sites

## ✅ Error Handling

### Graceful Failures
- [ ] If one source fails, others continue
- [ ] Error messages are clear
- [ ] Workflow doesn't fail completely
- [ ] Partial results are still useful

### Common Errors Handled
- [ ] Network timeouts: Logged and skipped
- [ ] API changes: Fallback strategies used
- [ ] Rate limiting: Logged and continued
- [ ] JavaScript sites: Warning shown, skipped

## 🚨 Red Flags (Investigate Immediately)

- ❌ All sources return 0 jobs
- ❌ Workflow fails completely
- ❌ No Telegram notifications after 24 hours
- ❌ Same jobs notified repeatedly
- ❌ `jobs.json` not being committed
- ❌ Workflow times out
- ❌ Permission denied errors

## 📊 Success Metrics

After 24 hours, you should have:

- ✅ **20+ successful workflow runs**
- ✅ **1500-3000+ jobs fetched per run**
- ✅ **10-50+ new jobs notified** (depending on keywords)
- ✅ **5-7 sources working consistently**
- ✅ **No workflow failures**
- ✅ **Clear logs with no critical errors**
- ✅ **Telegram notifications arriving**
- ✅ **No duplicate notifications**

## 📝 Notes Section

Use this space to track issues or observations:

### Issues Found:
```
[Date] [Source] [Issue Description]
Example: 2024-04-19 YC Jobs - Returns 0 jobs, API may have changed
```

### Resolutions:
```
[Date] [Action Taken]
Example: 2024-04-19 Disabled YC Jobs in config.yaml
```

### Performance Notes:
```
[Date] [Observation]
Example: 2024-04-19 LinkedIn occasionally rate-limited, but recovers
```

---

## Final Verification

Once all checkboxes are complete:

- [ ] All critical sources working
- [ ] Receiving relevant job notifications
- [ ] No repeated errors
- [ ] Workflow runs reliably
- [ ] Documentation reviewed

**Status**: ⬜ Not Started | 🟡 In Progress | ✅ Complete | ❌ Issues Found

---

**Last Updated**: [Add date when you complete this checklist]

**Overall Status**: [Not Started / In Progress / Complete / Needs Attention]
