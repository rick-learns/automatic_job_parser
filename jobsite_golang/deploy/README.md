# Deployment & Scheduling

## Systemd (Recommended)

### Install
```bash
sudo cp jobsite.service /etc/systemd/system/
sudo cp jobsite.timer /etc/systemd/system/

# Enable and start the timer
sudo systemctl enable jobsite.timer
sudo systemctl start jobsite.timer

# Check status
sudo systemctl status jobsite.timer
sudo systemctl status jobsite.service
```

### View logs
```bash
sudo journalctl -u jobsite.service -f
```

### Stop/disable
```bash
sudo systemctl stop jobsite.timer
sudo systemctl disable jobsite.timer
```

## Cron (Alternative)

### Install
```bash
crontab deploy/crontab.example
```

### Verify
```bash
crontab walks
```

### View logs
```bash
tail -f /var/log/jobsite.log
```

## Configuration

Edit the paths in `jobsite.service` or `crontab.example` to match your installation:
- Working directory
- Binary path
- Environment file path
- User account

