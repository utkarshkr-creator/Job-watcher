# Quick Fix Guide

## 🚨 Common Issues & Fast Solutions

### Issue: Triplebyte or other source returns 0 jobs

**Quick Fix:**
Already fixed! Triplebyte is now disabled by default in `config.yaml`:
```yaml
sources:
  triplebyte: false  # Disabled - requires login
```

**Why**: Triplebyte was acquired by Karat and requires login. It's expected to fail.

**See also**: `MANAGING_SOURCES.md` for managing all sources

---

### Issue: YC Jobs returns 0 jobs

**Quick Fix:**
```yaml
# In config.yaml, temporarily disable it
sources:
  ycjobs: false
```

**Why**: YC may require JavaScript or have changed their API. The code now has fallbacks, but if it still fails, disable it.

---

### Issue: Instahyre returns 0 jobs

**Quick Fix:**
```yaml
# In config.yaml, temporarily disable it
sources:
  instahyre: false
```

**Why**: Instahyre API may require authentication. The code tries multiple endpoints, but may still fail.

---

### Issue: All sources return 0 jobs

**Quick Fix:**
1. Test locally first:
   ```bash
   go run . --test-sources
   ```

2. If local works but GitHub Actions doesn't:
   - Check GitHub Actions logs for specific errors
   - May be IP-based blocking
   - Try running workflow at different times

3. Focus on reliable sources:
   ```yaml
   sources:
     indeed: true
     linkedin: true
     companies: true
     sharedlists: true
     # Disable others
     ycjobs: false
     instahyre: false
     reddit: false
   ```

---

### Issue: Workflow fails completely

**Quick Fix:**
1. Check GitHub Secrets are set:
   - `TG_TOKEN`
   - `TG_CHAT`
   - `GEMINI_API_KEY` (if using AI)

2. Check workflow permissions:
   - Settings → Actions → General
   - Workflow permissions: "Read and write permissions"

3. Check Go version:
   ```yaml
   # In .github/workflows/job-watcher.yml
   go-version: '1.21'  # Should be 1.21 or higher
   ```

---

### Issue: No Telegram notifications

**Quick Fix:**
1. Test your bot token:
   ```bash
   curl https://api.telegram.org/bot<YOUR_TOKEN>/getMe
   ```

2. Test your chat ID:
   ```bash
   curl https://api.telegram.org/bot<YOUR_TOKEN>/sendMessage \
     -d "chat_id=<YOUR_CHAT_ID>" \
     -d "text=Test"
   ```

3. Verify secrets in GitHub:
   - Settings → Secrets and variables → Actions
   - Check `TG_TOKEN` and `TG_CHAT` are set

---

### Issue: Too many duplicate notifications

**Quick Fix:**
1. Check `jobs.json` is being committed:
   - Look for commits from `github-actions[bot]`
   - If missing, check workflow permissions

2. Increase retention:
   ```yaml
   # In config.yaml
   retention_days: 60  # Keep job history longer
   ```

---

### Issue: Getting senior-level jobs

**Quick Fix:**
```yaml
# In config.yaml, add more exclusions
exclude_keywords:
  - senior
  - sr.
  - sr 
  - lead
  - principal
  - staff
  - manager
  - director
  - architect
  - 5+ years
  - 4+ years
  - 3+ years
  - 5 years
  - 4 years
  - 3 years
  - experienced  # Add this
  - expert       # Add this
```

---

### Issue: Not enough jobs

**Quick Fix:**
1. Broaden keywords:
   ```yaml
   keywords:
     - software engineer
     - developer
     - programmer
     - engineer
     # Remove very specific terms
   ```

2. Increase experience range:
   ```yaml
   max_experience_years: 5  # Instead of 2
   ```

3. Increase date range:
   ```yaml
   max_days_old: 14  # Instead of 5
   ```

4. Enable more sources:
   ```yaml
   sources:
     indeed: true
     linkedin: true
     companies: true
     sharedlists: true
     ycjobs: true
     hnjobs: true
     reddit: true
   ```

---

### Issue: Workflow runs but no output in logs

**Quick Fix:**
1. Check if workflow is actually running:
   - Actions tab → Should see runs with timestamps

2. Check workflow file syntax:
   ```bash
   # Validate YAML
   yamllint .github/workflows/job-watcher.yml
   ```

3. Re-trigger manually:
   - Actions tab → Job Watcher → Run workflow

---

### Issue: "Permission denied" when committing jobs.json

**Quick Fix:**
1. Update workflow permissions:
   ```yaml
   # In .github/workflows/job-watcher.yml
   permissions:
     contents: write  # Must be 'write', not 'read'
   ```

2. Or in repository settings:
   - Settings → Actions → General
   - Workflow permissions: "Read and write permissions"

---

## 🧪 Testing Commands

```bash
# Test all sources locally
go run . --test-sources

# Run full job watcher locally
go run .

# Check Go version
go version  # Should be 1.21+

# Update dependencies
go mod download
go mod tidy
```

---

## 📊 Expected Results

After fixes, you should see:

- ✅ **Indeed**: 50-100 jobs
- ✅ **LinkedIn**: 10-30 jobs
- ✅ **Companies**: 500-1000 jobs
- ✅ **Shared Lists**: 1000-2000 jobs
- ✅ **HN Jobs**: 10-30 jobs
- ✅ **Reddit**: 5-15 jobs
- ⚠️ **YC Jobs**: 0-20 jobs (may vary)
- ⚠️ **Instahyre**: 0-10 jobs (may vary)

**Total**: 1500-3000+ jobs per run

---

## 🔍 Debug Workflow

Run the debug workflow to test sources:

1. Go to **Actions** tab
2. Click **Debug Job Sources**
3. Click **Run workflow**
4. Wait for completion
5. Check logs for detailed output

---

## 📞 Still Having Issues?

1. ✅ Read `TROUBLESHOOTING.md` for detailed guide
2. ✅ Check GitHub Actions logs
3. ✅ Test locally with `go run test_sources.go`
4. ✅ Review `FIXES_APPLIED.md` for what was changed
5. ✅ Open an issue with logs (remove secrets!)

---

## 🎯 Minimal Working Configuration

If all else fails, use this minimal config:

```yaml
# config.yaml - Minimal reliable setup
keywords:
  - software engineer
  - developer

locations:
  - remote
  - india

max_experience_years: 3

exclude_keywords:
  - senior
  - lead

sources:
  indeed: true
  linkedin: true
  companies: true
  sharedlists: true
  # Disable everything else
  remoteok: false
  razorpay: false
  wellfound: false
  naukri: false
  instahyre: false
  ycjobs: false
  hnjobs: false
  reddit: false
  triplebyte: false

ai:
  enabled: false  # Disable AI for simplicity

retention_days: 30
max_days_old: 7
```

This should give you 1000+ jobs from reliable sources.
