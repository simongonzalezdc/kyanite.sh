# Deploying noise.sh PWA

This PWA is designed for self-hosting alongside your noise.sh desktop installation.

## Build

```bash
cd pwa
npm install
npm run build
```

This creates a `out/` directory with all static files.

## Deployment Options

### Option 1: Local Server (Recommended)

Serve the `out/` folder on your local network. The simplest options:

**Using Python:**
```bash
cd out
python3 -m http.server 8080
```

**Using Node.js:**
```bash
npx serve out -p 8080
```

Then access `http://your-computer-ip:8080` from your phone.

### Option 2: Docker (for NAS/home servers)

Create a `Dockerfile`:

```dockerfile
FROM nginx:alpine
COPY out/ /usr/share/nginx/html
EXPOSE 80
```

Build and run:
```bash
docker build -t noise-pwa .
docker run -d -p 8080:80 noise-pwa
```

### Option 3: Nginx

Add to your nginx config:

```nginx
server {
    listen 8080;
    root /path/to/pwa/out;
    index index.html;
    
    location / {
        try_files $uri $uri/ /index.html;
    }
    
    # Cache static assets
    location /_next/static {
        expires 1y;
        add_header Cache-Control "public, immutable";
    }
    
    # Service worker
    location /sw.js {
        expires off;
        add_header Cache-Control "no-cache";
    }
}
```

### Option 4: Caddy

```caddyfile
:8080 {
    root * /path/to/pwa/out
    file_server
    try_files {path} /index.html
}
```

## HTTPS (Optional but Recommended)

Some PWA features require HTTPS. For local network use, you can:

1. **Self-signed certificate**: Browsers will warn but you can proceed
2. **mkcert**: Generate locally-trusted certificates
3. **Caddy**: Automatic HTTPS with Let's Encrypt (requires public domain)

Example with mkcert:
```bash
# Install mkcert
brew install mkcert
mkcert -install

# Generate cert for local IP
mkcert 192.168.1.100
```

Then configure your server to use the generated `.pem` files.

## Connecting to noise.sh

1. Start noise.sh desktop with sync enabled
2. Note the sync server URL (shown in Settings → Sync)
3. Open the PWA on your phone
4. Enter the server URL and pairing code

## Troubleshooting

### PWA won't install
- Must be served over HTTPS (except localhost)
- Check manifest.webmanifest is served with correct MIME type
- Verify service worker is registered in browser devtools

### Can't connect to server
- Ensure phone and computer are on same network
- Check firewall allows port 8765 (noise.sh sync port)
- Verify noise.sh is running with sync enabled

### Ideas not syncing
- Check connection status in PWA header
- Verify sync server is running (Settings → Sync in noise.sh)
- Try manual sync by opening the PWA while connected

## File Structure

```
out/
├── index.html         # Main app
├── manifest.webmanifest
├── sw.js              # Service worker
├── _next/             # Static assets
│   └── static/
│       ├── chunks/    # JavaScript bundles
│       └── css/       # Stylesheets
└── icons/             # App icons
```

## Updates

When updating the PWA:

1. Pull latest code
2. Run `npm run build`
3. Replace `out/` folder on your server
4. Service worker will auto-update on next PWA load
