# Backend and Frontend Test Guide

## 1. Test Database Users

First, let's verify the users in the database:

```sql
SELECT id, username, role, is_admin FROM users;
```

Expected users:
- admin (is_admin: true, role: admin)
- stage2_emp (is_admin: false, role: stage2_employee) 
- stage3_emp (is_admin: false, role: stage3_employee)
- customer1 (is_admin: false, role: customer)

## 2. Test Backend API Endpoints

### Test Session Check:
```bash
curl -X GET http://localhost:8080/api/session -H "Cookie: session=your-session-cookie"
```

### Test MyJobs for different roles:
```bash
# For stage2_employee
curl -X GET http://localhost:8080/api/pipeline/myjobs -H "Cookie: session=your-session-cookie"

# For stage3_employee  
curl -X GET http://localhost:8080/api/pipeline/myjobs -H "Cookie: session=your-session-cookie"

# For customer
curl -X GET http://localhost:8080/api/pipeline/myjobs -H "Cookie: session=your-session-cookie"
```

## 3. Test Frontend Flow

### Admin User:
1. Login as admin
2. Go to /dashboard/admin
3. Click "Pipeline" in sidebar
4. Should see all jobs

### Employee User:
1. Login as stage2_emp
2. Go to /dashboard/employee  
3. Click "Pipeline" in sidebar
4. Should see only assigned jobs

## 4. Expected Behavior

### Admin Dashboard:
- Shows all pipeline jobs
- Can create new jobs
- Shows "Pipeline Management" header

### Employee Dashboard:
- Shows only assigned jobs
- Cannot create new jobs
- Shows "My Assigned Jobs" header
- Shows job statistics for their assigned jobs

## 5. Debug Steps

If issues occur:

1. Check browser console for errors
2. Check backend logs for authentication issues
3. Verify session cookies are being sent
4. Check if user roles are correct in database
5. Verify API endpoints are responding correctly 