# Secrets & API Keys Configuration

## ⚠️ SECURITY WARNING
**NEVER commit the `.env` ⚠️ file to Git.** All secrets are stored in GitHub Repository Secrets (see below).

## Current Setup

Your API keys are currently configured in `.env` file (which is gitignored).

### Configured Keys:
- ✅ **SerpAPI**: Google Search API for job scraping
- ✅ **GitHub PAT**: For GitHub Actions/CI operations
- ✅ **BASE_URL**: jobs.rick-learns.dev
- ✅ **Site Title**: QA/SDET Remote Jobs | Rick Learns

## GitHub Repository Secrets (Recommended for CI/CD)

To use these secrets in GitHub Actions, add them to your repository secrets:

### Steps to Add Secrets:

1. Go to your GitHub repository
2. Navigate to: **Settings** → **Secrets and variables** → **Actions**
3. Click **"New repository secret"**
4. Add each secret:

| Secret Name | Value | Description |
|------------|-------|-------------|
| `SERPAPI_API_KEY` | `[YOUR_KEY_FROM_.env_FILE]` | SerpAPI key for Google search |
| `GH_PAT` | `[YOUR_TOKEN_FROM_.env_FILE]` | GitHub PAT for automated commits (note: can't start with GITHUB_) |
| `BASE_URL` | `https://jobs.rick-learns.dev` | Your site URL |
| `SITE_TITLE` | `QA/SDET Remote Jobs \| Rick Learns` | Page title |

⚠️ **Important**: GitHub secret names CANNOT start with `GITHUB_` (reserved). Use `GH_PAT` or `GIT_TOKEN` instead.

### Using Secrets in GitHub Actions

Reference secrets in workflow files like this:

```yaml
env:
  SERPAPI_API_KEY: ${{ secrets.SERPAPI_API_KEY }}
  BASE_URL: ${{ secrets.BASE_URL }}
  SITE_TITLE: ${{ secrets.SITE_TITLE }}
```

## Local Development

1. Copy the example file:
   ```bash
   cp .env.example .env
   ```

2. Edit `.env` with your actual keys (already done for you)

3. Load environment variables:
   ```bash
   source .env
   export $(cat .env | xargs)
   ```

## Verification

Verify your secrets are not in Git:
```bash
# These commands should return NO results (replace with your actual keys)
git log -p --all -S "YOUR_SERPAPI_KEY"
git log -p --all -S "YOUR_GITHUB_PAT"
```

## Rotating Keys

If you need to rotate any keys:
1. Generate new key from the service provider
2. Update `.env` file locally
3. Update GitHub repository secrets
4. Test the application

## .gitignore Protection

The `.gitignore` file protects:
- `.env` - Main environment file
- `.env.*.local` - Local overrides
- `data/` - SQLite database
- `public/` - Generated site files (optional)
- Other sensitive files

