# Deployment Guide

## Setup nginx for jobs.rick-learns.dev

**IMPORTANT**: You DON'T need to modify the captive-portal config. Just add a new server block.

1. Copy the nginx configuration:
   ```bash
   cd /home/rick/automatic_job_parser/jobsite_golang
   sudo cp deploy/nginx-jobs-rick-learns-dev.conf /etc/nginx/sites-available/jobs.rick-learns.dev
   ```

2. Enable the site:
   ```bash
   sudo ln -s /etc/nginx/sites-available/jobs.rick-learns.dev /etc/nginx/sites-enabled/
   ```

3. Test nginx configuration:
   ```bash
   sudo nginx -t
   ```

4. Reload nginx:
   ```bash
   sudo systemctl reload nginx
   ```

**Why this works**: Nginx matches exact `server_name` before the `default_server` catch-all, so your jobs site will be served even though captive-portal-multi has a default_server block.

## SSL Configuration

If using Cloudflare:
- Enable "Full" or "Full (strict)" SSL mode in Cloudflare dashboard
- The SSL certificate shown is the Cloudflare SSL cert (this is normal)

If using Let's Encrypt:
1. Install certbot:
   ```bash
   sudo apt install certbot python3-certbot-nginx
   ```

2. Get certificate:
   ```bash
   sudo certbot --nginx -d jobs.rick-learns.dev
   ```

3. Certbot will update the nginx config automatically

## Update Process

1. Update jobs data (run daily command):
   ```bash
   cd jobsite_golang
   export SERPAPI_API_KEY=your_key
   ./jobsite daily
   ```

2. Build frontend:
   ```bash
   cd ../frontend
   npm run build
   ```

Or use the combined command:
```bash
cd jobsite_golang
make build-all
```

## Verify Deployment

Check that these files exist in `/home/rick/automatic_job_parser/jobsite_golang/public/latest/`:
- `index.html` - React app
- `assets/` - JavaScript and CSS bundles
- `jobs.json` - Job data
- `jobs.csv` - Job data (CSV)

Visit https://jobs.rick-learns.dev in your browser to verify.

