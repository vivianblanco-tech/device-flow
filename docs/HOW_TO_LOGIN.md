# How to Log In to the Laptop Tracking System

This guide will help you log in to the application for the first time.

## Prerequisites

Before you can log in, ensure:

1. ✅ **PostgreSQL is running**
2. ✅ **Database is created** (`laptop_tracking_dev`)
3. ✅ **Migrations are applied** (run `make migrate-up`)
4. ✅ **A user account exists** in the database

---

## Step 1: Create a Test User

You have two options to create a test user:

### Option A: Using PowerShell Script (Recommended for Windows)

```powershell
# Run from the project root
.\scripts\create-test-user.ps1
```

This will create a user with:
- **Email**: `admin@bairesdev.com`
- **Password**: `Test123!`
- **Role**: `logistics` (full access)

### Option B: Using SQL Script Directly

```bash
# Connect to your database and run the SQL script
psql -h localhost -U postgres -d laptop_tracking_dev -f scripts/create-test-user.sql
```

### Option C: Manual SQL Command

```sql
-- Connect to your database first
psql -h localhost -U postgres -d laptop_tracking_dev

-- Then run this SQL
INSERT INTO users (email, password_hash, role, created_at, updated_at)
VALUES (
    'admin@bairesdev.com',
    '$2a$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/LewY5GyYzS0MYq5IW',
    'logistics',
    NOW(),
    NOW()
);
```

---

## Step 2: Start the Application

```bash
# Using Make
make run

# Or directly
go run cmd/web/main.go
```

You should see output similar to:
```
Database connected successfully
Server starting on localhost:8080
Environment: development
Login at: http://localhost:8080/login
```

---

## Step 3: Access the Login Page

1. Open your browser
2. Navigate to: **http://localhost:8080/login**
3. You should see the login form

---

## Step 4: Log In

### Using Username/Password

Enter the credentials:
- **Email**: `admin@bairesdev.com`
- **Password**: `Test123!`

Click **"Sign In"**

### Using Google OAuth (Optional)

If you've set up Google OAuth credentials in your `.env` file:
1. Click **"Sign in with Google"**
2. Authorize with your `@bairesdev.com` Google account
3. You'll be redirected back and logged in automatically

**Note**: Only `@bairesdev.com` domain emails are allowed for Google OAuth.

---

## Step 5: After Login

Once logged in successfully, you'll be redirected to:
- **Dashboard**: `/dashboard` (currently shows "Coming Soon!")

From there, you can access:
- **Pickup Form**: `/pickup-form`
- **Shipments List**: `/shipments`
- **Reception Report**: `/reception-report`
- **Delivery Form**: `/delivery-form`

---

## Troubleshooting

### "Database connection failed"

**Problem**: Cannot connect to PostgreSQL

**Solutions**:
1. Check PostgreSQL is running:
   ```bash
   # Windows
   Get-Service postgresql*
   
   # Linux/Mac
   sudo service postgresql status
   ```

2. Verify `.env` database settings:
   ```env
   DB_HOST=localhost
   DB_PORT=5432
   DB_USER=postgres
   DB_PASSWORD=password
   DB_NAME=laptop_tracking_dev
   ```

3. Test connection manually:
   ```bash
   psql -h localhost -U postgres -d laptop_tracking_dev
   ```

### "Invalid email or password"

**Problem**: Login credentials don't match

**Solutions**:
1. Make sure you created the test user (see Step 1)
2. Check exact email: `admin@bairesdev.com`
3. Check exact password: `Test123!` (case-sensitive, with exclamation mark)
4. Verify user exists in database:
   ```sql
   SELECT id, email, role FROM users WHERE email = 'admin@bairesdev.com';
   ```

### "Failed to parse templates"

**Problem**: Template files not found

**Solutions**:
1. Make sure you're running the application from the project root directory
2. Verify `templates/pages/` directory exists
3. Check `login.html` exists in `templates/pages/`

### "Google OAuth not working"

**Problem**: OAuth flow fails

**Solutions**:
1. Make sure you've set up Google OAuth credentials in Google Cloud Console
2. Verify `.env` file has correct values:
   ```env
   GOOGLE_CLIENT_ID=your-client-id-here
   GOOGLE_CLIENT_SECRET=your-client-secret-here
   GOOGLE_REDIRECT_URL=http://localhost:8080/auth/google/callback
   GOOGLE_ALLOWED_DOMAIN=bairesdev.com
   ```
3. Ensure the redirect URL in Google Console matches exactly: `http://localhost:8080/auth/google/callback`
4. Check that Google+ API is enabled in your Google Cloud project

### "Session cookie not set"

**Problem**: Can't stay logged in

**Solutions**:
1. Check `.env` has `SESSION_SECRET` set (any random string)
2. If using HTTPS locally, update cookie settings in code
3. Clear browser cookies and try again
4. Try a different browser (private/incognito mode)

### Port 8080 already in use

**Problem**: Another application is using port 8080

**Solutions**:
1. Find what's using the port:
   ```powershell
   # Windows
   netstat -ano | findstr :8080
   
   # Linux/Mac  
   lsof -i :8080
   ```
2. Kill that process or change the port in `.env`:
   ```env
   APP_PORT=8081
   ```

---

## User Roles

The system supports 4 user roles:

| Role | Access Level | Description |
|------|--------------|-------------|
| **logistics** | Full access | Can manage all shipments, view all data, coordinate operations |
| **client** | Limited | Can submit pickup forms for their company only |
| **warehouse** | Medium | Can receive shipments, create reception reports, view warehouse data |
| **project_manager** | Read-only | Can view dashboards, reports, and shipment status |

The test user created above has the **logistics** role for full system access.

---

## Creating Additional Users

### For Testing Different Roles

```sql
-- Client user
INSERT INTO users (email, password_hash, role, created_at, updated_at)
VALUES ('client@example.com', '$2a$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/LewY5GyYzS0MYq5IW', 'client', NOW(), NOW());

-- Warehouse user
INSERT INTO users (email, password_hash, role, created_at, updated_at)
VALUES ('warehouse@bairesdev.com', '$2a$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/LewY5GyYzS0MYq5IW', 'warehouse', NOW(), NOW());

-- Project Manager user
INSERT INTO users (email, password_hash, role, created_at, updated_at)
VALUES ('pm@bairesdev.com', '$2a$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/LewY5GyYzS0MYq5IW', 'project_manager', NOW(), NOW());
```

All these users will have the password `Test123!`

---

## Security Notes

### Default Password Hash

The test scripts use a pre-generated bcrypt hash for the password `Test123!`. This is **ONLY for development and testing**.

### For Production

1. **Never** use these default credentials in production
2. Generate unique password hashes using:
   ```go
   hash, _ := bcrypt.GenerateFromPassword([]byte("your-password"), 12)
   fmt.Println(string(hash))
   ```
3. Use strong passwords (minimum 8 characters, including uppercase, lowercase, digit, and special character)
4. Enable Google OAuth and restrict to your domain
5. Use environment-specific `.env` files
6. Change `SESSION_SECRET` to a cryptographically secure random string

---

## Quick Reference

### Default Test Credentials
```
Email:    admin@bairesdev.com
Password: Test123!
Role:     logistics
```

### Key URLs
```
Login:     http://localhost:8080/login
Dashboard: http://localhost:8080/dashboard
Shipments: http://localhost:8080/shipments
```

### Common Commands
```bash
# Start application
make run

# Create test user
.\scripts\create-test-user.ps1

# Check database
psql -h localhost -U postgres -d laptop_tracking_dev

# Run migrations
make migrate-up

# Check logs
tail -f app-output.log
```

---

## Next Steps

After logging in successfully:

1. ✅ Explore the dashboard
2. ✅ Try creating a shipment with the pickup form
3. ✅ Test the shipments list view
4. ✅ Set up Google OAuth (optional)
5. ✅ Create additional test users for other roles
6. ✅ Review the [Contributing Guide](CONTRIBUTING.md)

---

**Need Help?**

If you continue to have issues:
1. Check the application logs
2. Review error messages in the browser console (F12)
3. Verify all prerequisites are met
4. Contact the development team

Happy tracking! 🚀

