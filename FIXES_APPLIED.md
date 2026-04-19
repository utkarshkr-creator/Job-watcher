# Fixes Applied for GitHub Actions Parsing Issues

## Problem
Some job sources were returning 0 jobs when running in GitHub Actions, while they might work locally.

## Root Causes Identified

1. **JavaScript-Required Sites**: Sites like Wellfound, Naukri require JavaScript to render content
2. **API Endpoint Changes**: YC Jobs and Instahyre API endpoints may have changed
3. **Insufficient Error Handling**: Errors weren't being logged clearly
4. **No Fallback Mechanisms**: Single point of failure for each source

## Solutions Implemented

### 1. Enhanced YC Jobs Fetching (`extra_sources.go`)

**Changes:**
- Try multiple API endpoints (`/api/v1/jobs`, `/api/jobs`)
- Parse both array and object-wrapped JSON responses
- Improved HTML fallback with multiple regex patterns
- Better error messages showing which step failed
- Graceful degradation (returns empty array instead of breaking)

**Benefits:**
- More resilient to API changes
- Clear debugging information in logs
- Won't break other sources if YC fails

### 2. Enhanced Instahyre Fetching (`alternatives.go`)

**Changes:**
- Try primary search API endpoint first
- Fallback to alternative opportunities endpoint
- Added proper headers (X-Requested-With, Referer)
- Multiple search term attempts
- Detailed error logging

**Benefits:**
- Multiple fallback strategies
- Better chance of success
- Clear error messages for debugging

### 3. Improved Error Handling

**Changes:**
- All fetch functions now return empty arrays on error instead of propagating errors
- Added detailed logging with status codes
- Print warnings for JavaScript-required sites
- Continue processing other sources even if one fails

**Benefits:**
- Partial results are still useful
- Easy to identify which source is failing
- No cascading failures

### 4. Debug Workflow (`.github/workflows/debug-sources.yml`)

**New File:**
- Manual workflow to test each source individually
- Network connectivity tests
- Detailed output for troubleshooting

**Usage:**
```bash
# Go to Actions tab → Debug Job Sources → Run workflow
```

### 5. Test Script (`test_sources.go`)

**New File:**
- Test each source individually
- Show sample jobs from each source
- Summary of working vs failing sources
- Can be run locally or in CI

**Usage:**
```bash
go run . --test-sources
```

### 6. Comprehensive Documentation

**New Files:**
- `TROUBLESHOOTING.md` - Detailed debugging guide
- `FIXES_APPLIED.md` - This file

**Updated Files:**
- `README.md` - Added troubleshooting section

### 7. Workflow Improvements (`.github/workflows/job-watcher.yml`)

**Changes:**
- Added `DEBUG=true` environment variable
- Better comments explaining each step
- Maintained all existing functionality

## Testing Recommendations

### Local Testing
```bash
# Test all sources
go run . --test-sources

# Run full job watcher
go run .
```

### GitHub Actions Testing
1. **Manual Trigger**: Actions tab → Job Watcher → Run workflow
2. **Debug Mode**: Actions tab → Debug Job Sources → Run workflow
3. **Check Logs**: Expand "Run Job Watcher" step to see detailed output

## Expected Behavior After Fixes

### Working Sources (Should return jobs)
- ✅ Indeed (RSS feeds)
- ✅ LinkedIn (Guest API)
- ✅ Razorpay (Career page)
- ✅ HN Jobs (Algolia API)
- ✅ Reddit (JSON API)
- ✅ Company Career Pages (Direct scraping)
- ✅ Shared Lists (Google Sheets/GitHub)

### Potentially Limited Sources
- ⚠️ YC Jobs (May require JavaScript, has HTML fallback)
- ⚠️ Instahyre (API may require auth, has fallbacks)
- ⚠️ Triplebyte (Limited public access)

### Disabled by Default (Require JavaScript)
- ❌ Wellfound (Requires JavaScript)
- ❌ Naukri (Requires JavaScript)
- ❌ RemoteOK (Requires premium)

## Monitoring

### Check if fixes are working:

1. **View latest workflow run**:
   - Go to Actions tab
   - Click latest "Job Watcher" run
   - Check output for each source

2. **Look for these indicators**:
   - `✓ SourceName: X jobs` - Working
   - `Warning: SourceName: error` - Failed with reason
   - `Note: SourceName requires JavaScript` - Expected limitation

3. **Expected output example**:
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
   ```

## Next Steps

1. **Monitor for 24-48 hours** to see if sources are working consistently
2. **Check Telegram** for job notifications
3. **Review logs** if any source consistently returns 0 jobs
4. **Adjust config.yaml** to disable sources that don't work for you

## Rollback Plan

If issues persist, you can:

1. **Disable problematic sources** in `config.yaml`:
   ```yaml
   sources:
     ycjobs: false
     instahyre: false
   ```

2. **Revert changes**:
   ```bash
   git revert HEAD
   git push
   ```

3. **Use only reliable sources**:
   - Indeed (RSS)
   - LinkedIn (Guest API)
   - Company Career Pages
   - Shared Lists

## Additional Notes

- **Rate Limiting**: If you see 429 errors, the site is rate-limiting. Consider reducing frequency.
- **IP Blocking**: Some sites may block GitHub Actions IPs. This is expected for some sources.
- **API Changes**: Job sites change their APIs frequently. The fallback mechanisms help mitigate this.
- **JavaScript Sites**: Cannot be fully supported in GitHub Actions without headless browser (too slow/complex).

## Success Metrics

After these fixes, you should see:

- ✅ At least 5-7 sources returning jobs
- ✅ Total of 500+ jobs per run (depending on keywords)
- ✅ Clear error messages for failing sources
- ✅ No workflow failures (even if some sources fail)
- ✅ Telegram notifications working

## Support

If issues persist:

1. Run `go run . --test-sources` locally
2. Check `TROUBLESHOOTING.md`
3. Review GitHub Actions logs
4. Open an issue with logs and config (remove secrets)
