# Adding GitHub Container Registry Credentials to Coolify

## 🔐 Option 1: Keep Package Private + Add Credentials to Coolify

If you want to keep your Docker image private, you need to give Coolify credentials to access it.

### Steps to Add Registry Credentials in Coolify

#### Method A: Project-Level Registry (Recommended)

1. **In Coolify Dashboard:**
   - Go to your **Project** (wizardcore)
   - Click **"Settings"** or look for **"Registry"** section
   - Click **"Add Registry"** or **"Private Docker Registry"**

2. **Fill in the credentials:**
   ```
   Registry Type: GitHub Container Registry (or Custom)
   Registry URL: ghcr.io
   Username: davidpoweraus-sys
   Password/Token: YOUR_GITHUB_TOKEN_HERE
   ```
   
   Use the GitHub token: `ghp_XrcL...` (the one you created earlier)

3. **Save** the credentials

4. **Redeploy** your frontend service

#### Method B: Service-Level Registry

1. **In Coolify Dashboard:**
   - Go to your **Frontend Service**
   - Click **"Configuration"** or **"General"**
   - Look for **"Registry"** or **"Image Registry"** section

2. **Add credentials:**
   ```
   Registry: ghcr.io
   Username: davidpoweraus-sys  
   Password: YOUR_GITHUB_TOKEN_HERE
   ```
   
   Use the GitHub token you created earlier

3. **Save** and **Redeploy**

#### Method C: Server-Level Docker Login (Advanced)

If Coolify has a server terminal or you have SSH access:

```bash
# SSH into your Coolify server
ssh your-server

# Login to GHCR
echo "YOUR_GITHUB_TOKEN" | docker login ghcr.io -u davidpoweraus-sys --password-stdin

# This will save credentials in ~/.docker/config.json
# All subsequent pulls will use these credentials
```

---

## 🌐 Option 2: Make Package Public (Easiest)

**Recommended for this use case** since it's your own deployment.

### Why Make it Public?
- ✅ No credentials needed in Coolify
- ✅ Simpler deployment process
- ✅ Faster pulls (no auth overhead)
- ✅ Still secure - it's just a Docker image, not your source code
- ✅ Anyone can pull the image, but only you can push updates

### How to Make it Public

1. **Go to your package page:**
   ```
   https://github.com/users/davidpoweraus-sys/packages/container/package/wizardcore-frontend
   ```

2. **Click "Package settings"** (gear icon on the right)

3. **Scroll to "Danger Zone"** at the bottom

4. **Click "Change visibility"**

5. **Select "Public"**

6. **Type the package name to confirm:** `wizardcore-frontend`

7. **Click "I understand the consequences, change package visibility"**

8. **Done!** Now anyone can pull (but only you can push)

### Verify it's Public

Run this command - it should work WITHOUT authentication:

```bash
docker pull ghcr.io/davidpoweraus-sys/wizardcore-frontend:latest
```

If it works, go to Coolify and click **Redeploy**!

---

## 🔍 Current Status Check

To verify if your package is public or private, open this URL **in an incognito/private browser window** (not logged into GitHub):

```
https://github.com/davidpoweraus-sys/wizardcore/pkgs/container/wizardcore-frontend
```

- If you can see it → ✅ **Public**
- If you get "404 Not Found" → 🔒 **Private**

---

## 📊 Comparison

| Aspect | Public Package | Private Package |
|--------|---------------|-----------------|
| Coolify Setup | ✅ Simple | ⚙️ Need credentials |
| Pull Speed | ✅ Fast | 🐌 Slightly slower (auth) |
| Security | ✅ Image only, source code still private | 🔒 Extra layer |
| Best For | Personal/Open deployments | Enterprise/Sensitive |

**Recommendation:** Make it **public** unless you have specific security requirements.

---

## ✅ Next Steps

1. **Choose** your approach (public vs private)
2. **Configure** accordingly
3. **Redeploy** in Coolify
4. **Verify** deployment succeeds

The "unauthorized" error will disappear once credentials are added or package is made public! 🚀
