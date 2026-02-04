# Personalizing Your Job Watcher

This guide explains how to tweak the project to find jobs that are relevant to YOU. All the magic happens in `config.yaml` and `resume.txt`.

## 1. Defining Your Skills & Roles
**File:** `config.yaml` -> `keywords`

This is the most important part. The watcher looks for these words in job titles and descriptions.

-   **Languages**: Add the languages you know (e.g., "python", "java", "golang").
-   **Frameworks**: specific tools you use (e.g., "react", "django", "kubernetes").
-   **Roles**: Titles you are looking for (e.g., "SDE", "DevOps Engineer").

**Example:**
```yaml
keywords:
  - python
  - aws
  - devops engineer
  - backend developer
```

## 2. Choosing Locations
**File:** `config.yaml` -> `locations`

Filter where you want to work.
-   **Specific Cities**: "bangalore", "london", "nyc".
-   **Remote**: "remote", "wfh".

**Example:**
```yaml
locations:
  - remote
  - berlin
  - munich
```

## 3. Experience Level
**File:** `config.yaml` -> `max_experience_years` & `exclude_keywords`

-   **`max_experience_years`**: Helps simple number-based filtering.
-   **`exclude_keywords`**: Removes roles that are too senior. If you are a fresher, keep "Senior", "Lead", "Principal" here. If you are a Senior, remove them!

## 4. AI Matching (The "Smart" Part)
**File:** `resume.txt`

If you enabled AI in `config.yaml`, the system reads `resume.txt` to understand *who you are*.

**Action:**
1.  Open `resume.txt`.
2.  Paste the text content of your actual resume there.
3.  The AI will now compare every new job description against YOUR resume text to give it a relevance score (0-100).

## 5. Enabling/Disabling Sources
**File:** `config.yaml` -> `sources`

Don't care about "Hacker News" jobs? Disable them to save time.

```yaml
sources:
  linkedin: true
  hnjobs: false  # Disabled
```

## Summary Checklist
- [ ] Updated `keywords` in `config.yaml`?
- [ ] Updated `locations` in `config.yaml`?
- [ ] Pasted my resume into `resume.txt`?
- [ ] Set my `TG_TOKEN` in `.env`?
