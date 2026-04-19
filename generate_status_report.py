import re
import os

# Read companies.go to extract all configured companies
companies = []
with open('companies.go', 'r') as f:
    content = f.read()
    # Regex to find {Name: "...", URL: "..."
    matches = re.findall(r'{Name: "(.*?)", URL: "(.*?)"', content)
    for name, url in matches:
        companies.append({'name': name, 'url': url, 'status': '❌ 0 jobs found'})

# Read full_output.txt to see which ones found jobs
found_jobs = {}
with open('full_output.txt', 'r') as f:
    for line in f:
        # Line format: "    Name: N jobs"
        match = re.search(r'^\s+(.*?): (\d+) jobs', line)
        if match:
            name = match.group(1).strip()
            count = match.group(2)
            found_jobs[name] = count

# Update status
working_count = 0
failed_count = 0

markdown = "# Company Source Status Report\n\n"
markdown += "Generated from latest execution log.\n\n"
markdown += "| Company | Status | URL |\n"
markdown += "|---|---|---|\n"

for company in companies:
    name = company['name']
    if name in found_jobs:
        count = found_jobs[name]
        company['status'] = f"✅ {count} jobs"
        working_count += 1
    else:
        failed_count += 1
    
    markdown += f"| {company['name']} | {company['status']} | {company['url']} |\n"

print(f"Report generated: {working_count} working, {failed_count} returned 0 jobs.")

with open('/Users/utkarshkumar/.gemini/antigravity/brain/95b5c5c9-9333-4067-8136-657a53a4f72e/company_status.md', 'w') as f:
    f.write(markdown)
