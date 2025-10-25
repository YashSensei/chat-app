# Security Checklist

## Current Security Status

### ✅ Protected
- `.env` files are properly gitignored
- Passwords are hashed with bcrypt
- JWT tokens stored in HTTP-only cookies
- CORS configured for specific origins

### ⚠️ Action Required

#### 1. Rotate These Credentials Immediately if Exposed
If you ever committed your `.env` file to git, rotate these immediately:

- **MongoDB Password**: `LvBWOlhPcV6t9EyR`
  - Go to MongoDB Atlas → Database Access → Edit User → Reset Password
  
- **JWT Secret**: `mysecretkey` (WEAK - needs to be changed)
  - Generate a strong secret: `openssl rand -base64 32`
  - Update in `.env`: `JWT_SECRET=<new-strong-secret>`
  
- **Cloudinary Credentials**: 
  - Cloud Name: `dfiw3opn0`
  - API Key: `962924628125544`
  - API Secret: `DVegF1e5FRu3lqO4jhKMycmiAV8`
  - Reset at: https://console.cloudinary.com/settings/security

#### 2. Strengthen JWT Secret
Your current JWT secret `mysecretkey` is too weak.

Generate a strong one:
```bash
# Windows PowerShell
[Convert]::ToBase64String((1..32 | ForEach-Object { Get-Random -Maximum 256 }))

# Or use online generator (HTTPS only):
# https://generate-secret.vercel.app/32
```

Update `go-backend/.env`:
```
JWT_SECRET=your-new-long-random-string-min-32-characters
```

#### 3. Production Environment Variables
Never use development credentials in production. Set these separately on your hosting platform:

**Render/Railway/Fly.io:**
```bash
# Set via dashboard or CLI
railway variables set JWT_SECRET="production-secret-here"
railway variables set MONGODB_URI="production-mongo-uri"
```

#### 4. MongoDB Security
- ✅ Use strong password (already using random string)
- ⚠️ Whitelist only your server IPs (currently allowing all with 0.0.0.0/0)
- ✅ Database user has minimal required permissions

To restrict MongoDB access:
1. Go to MongoDB Atlas → Network Access
2. Remove `0.0.0.0/0` entry
3. Add your production server IPs only

#### 5. Check Git History
Verify no secrets in git history:
```bash
git log --all --full-history -- "**/.env"
git log --all --full-history -- "**/go-backend/.env"
```

If any `.env` files appear:
```bash
# Install BFG Repo-Cleaner
# Download from: https://rtyley.github.io/bfg-repo-cleaner/

# Remove .env files from history
java -jar bfg.jar --delete-files .env

# Clean up
git reflog expire --expire=now --all
git gc --prune=now --aggressive

# Force push (WARNING: coordinate with team first)
git push --force
```

## Best Practices Going Forward

1. **Never commit secrets** - always use `.env.sample` templates
2. **Use different credentials** for dev/staging/production
3. **Rotate secrets regularly** - every 90 days minimum
4. **Use secret managers** in production (AWS Secrets Manager, HashiCorp Vault, etc.)
5. **Enable 2FA** on all cloud services (MongoDB, Cloudinary, GitHub)
6. **Monitor access logs** for suspicious activity

## Emergency Response

If credentials are leaked:
1. ✅ Verify `.env` is in `.gitignore`
2. 🔄 Rotate ALL credentials immediately
3. 🔍 Check git history for exposure
4. 🗑️ Remove from git history if found
5. 📧 Notify team members
6. 📊 Monitor services for unauthorized access
7. 🔐 Enable additional security measures (2FA, IP whitelisting)
