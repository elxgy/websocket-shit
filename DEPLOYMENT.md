# 🚀 WebSocket Chat Server - Deployment Guide

## Railway Deployment Steps

### Prerequisites
- ✅ GitHub account
- ✅ Railway account (free tier available)
- ✅ MongoDB Atlas database (already configured)

---

## 🔧 **Step 1: Prepare Repository**

### 1.1 Commit and Push Changes
```bash
# Add all changes
git add .

# Commit with descriptive message
git commit -m "Enhanced websocket chat server with MongoDB and authentication"

# Push to GitHub
git push origin main
```

### 1.2 Verify Required Files
Ensure these files are in your repository:
- ✅ `Dockerfile` - Container configuration
- ✅ `railway.json` - Railway deployment config
- ✅ `.railwayignore` - Files to exclude from deployment
- ✅ `go.mod` & `go.sum` - Go dependencies
- ✅ Source code files (main.go, database.go, etc.)

---

## 🌐 **Step 2: Deploy to Railway**

### 2.1 Connect Repository
1. Go to [railway.app](https://railway.app)
2. Click **"Start a New Project"**
3. Select **"Deploy from GitHub repo"**
4. Choose your `websocket-shit` repository

### 2.2 Configure Environment Variables
In Railway dashboard:
1. Go to your project → **Variables** tab
2. Add the following environment variable:
   ```
   MONGODB_URI = mongodb+srv://server:tkXcUegkJbmo5Hje@cabrittozap.toaludz.mongodb.net/?retryWrites=true&w=majority&appName=cabrittozap
   ```

### 2.3 Deploy
1. Railway will automatically detect your Dockerfile
2. Click **"Deploy"** 
3. Wait for build to complete (2-5 minutes)

---

## 🔗 **Step 3: Get Your Deployed URLs**

Once deployed, Railway will provide:
- **Backend API URL**: `https://your-app-name.up.railway.app`
- **WebSocket URL**: `wss://your-app-name.up.railway.app/ws`

### Test Your Deployment
```bash
# Test health endpoint
curl https://your-app-name.up.railway.app/health

# Expected response:
{
  "status": "ok",
  "clients": 0,
  "max_clients": 4,
  "database_connected": true
}
```

---

## 🎯 **Step 4: Deploy Frontend (Optional)**

### Option A: Deploy Frontend to Vercel/Netlify
1. Create a new repository for your frontend
2. Copy the `example-frontend/index.html` to your frontend repo
3. Update the URLs in the JavaScript:
   ```javascript
   const serverUrl = 'https://your-app-name.up.railway.app';
   const wsUrl = 'wss://your-app-name.up.railway.app';
   ```
4. Deploy to Vercel/Netlify

### Option B: Serve Frontend from Railway
1. Create a separate Railway service
2. Use the nginx Dockerfile approach with your frontend files

---

## 🧪 **Step 5: Test the Deployment**

### 5.1 API Testing
```bash
# Test login
curl -X POST https://your-app-name.up.railway.app/login \
  -H "Content-Type: application/json" \
  -d '{"username": "alice", "password": "password123"}'
```

### 5.2 WebSocket Testing
Use the example frontend or any WebSocket client:
```javascript
const ws = new WebSocket('wss://your-app-name.up.railway.app/ws?username=alice');
```

---

## 🔧 **Configuration Options**

### Environment Variables
| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `MONGODB_URI` | ✅ Yes | None | MongoDB connection string |
| `PORT` | ❌ No | 8080 | Server port (Railway auto-assigns) |

### Railway Settings
- **Build Command**: Automatic (uses Dockerfile)
- **Start Command**: `./chat-server`
- **Health Check**: `/health`
- **Port**: Auto-detected from Railway

---

## 📊 **Monitoring & Logs**

### View Logs in Railway
1. Go to your Railway project
2. Click on **"Deployments"** tab  
3. Click on latest deployment
4. View **"Deploy Logs"** and **"Application Logs"**

### Key Log Messages to Look For
```
✅ Connected to MongoDB successfully
✅ Chat server starting on port 8080
✅ WebSocket endpoint: /ws
✅ Login endpoint: /login
```

---

## 🚨 **Troubleshooting**

### Common Issues

#### 1. Build Failures
```bash
# Check Go version in Dockerfile
FROM golang:1.21-alpine AS builder
```

#### 2. Database Connection Issues
- Verify `MONGODB_URI` environment variable
- Check MongoDB Atlas network access (allow all IPs: 0.0.0.0/0)
- Ensure database user has proper permissions

#### 3. WebSocket Connection Issues
- Use `wss://` (not `ws://`) for HTTPS deployments
- Check CORS configuration in main.go
- Verify Railway domain is correct

#### 4. Port Issues
- Railway automatically assigns ports
- Don't hardcode port 8080 in production URLs
- Use Railway-provided domain

### Debug Commands
```bash
# Check if service is running
curl https://your-app-name.up.railway.app/health

# Test specific endpoint
curl -v https://your-app-name.up.railway.app/login

# Check WebSocket connection
wscat -c wss://your-app-name.up.railway.app/ws?username=alice
```

---

## 🔒 **Security Considerations**

### Production Checklist
- ✅ MongoDB connection uses authentication
- ✅ CORS is properly configured
- ✅ User passwords are bcrypt hashed
- ✅ Input validation on all endpoints
- ✅ Connection limits enforced (max 4 users)
- ❌ Consider adding rate limiting for production
- ❌ Consider adding JWT tokens for session management

### MongoDB Atlas Security
- ✅ Network access restricted or monitored
- ✅ Database user has minimal required permissions
- ✅ Connection string uses authentication

---

## 📈 **Scaling & Performance**

### Current Limitations
- **Max 4 concurrent users** (as designed)
- **Single instance** (Railway free tier)
- **In-memory user sessions** (lost on restart)

### Future Enhancements
- Add Redis for session storage
- Implement horizontal scaling
- Add load balancing for multiple instances
- Implement user registration system

---

## 🎉 **Success Criteria**

Your deployment is successful when:
- ✅ Health endpoint returns status "ok"
- ✅ Users can login with test credentials
- ✅ WebSocket connections work
- ✅ Messages are sent and received in real-time
- ✅ Messages persist in MongoDB
- ✅ User limit (4) is enforced

---

## 📞 **Support**

### Documentation
- [Railway Documentation](https://docs.railway.app)
- [MongoDB Atlas Documentation](https://docs.atlas.mongodb.com)
- [WebSocket Protocol](https://tools.ietf.org/html/rfc6455)

### Test Credentials
```
alice / password123
bob / password123  
charlie / password123
diana / password123
```

Happy deploying! 🚀