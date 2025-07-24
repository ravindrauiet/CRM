# Subadmin System Test Guide

## Overview
The subadmin system provides a middle-tier management role that can create tasks and manage employee assignments without full admin privileges.

## New Features Added

### 1. Subadmin Role
- **Role**: `subadmin`
- **Permissions**: 
  - Create pipeline jobs
  - View all jobs and users
  - Access analytics and reports
  - Manage employee assignments
  - Cannot create other admins

### 2. Subadmin Dashboard
- **Location**: `/dashboard/subadmin`
- **Features**:
  - Pipeline statistics overview
  - Employee management section
  - Recent jobs table
  - Create new job functionality
  - Employee assignment tracking

### 3. Updated Navigation
- **Sidebar**: Shows "S" for subadmin users
- **Menu Items**: Same as admin but with subadmin-specific dashboard
- **Role Display**: Shows "Subadmin" in user info

## Database Setup

### Add Subadmin User
Run the SQL script to create a subadmin user:

```sql
-- Execute add_subadmin.sql
-- Username: subadmin
-- Password: subadmin123
```

## Test Scenarios

### Test 1: Subadmin Login
1. **Start**: Go to `/login`
2. **Login**: Use `subadmin` / `subadmin123`
3. **Expected**: Redirect to `/dashboard/subadmin`
4. **Verify**: Sidebar shows "S" and "Subadmin" role

### Test 2: Subadmin Dashboard Access
1. **Login**: As subadmin
2. **Navigate**: Dashboard should show subadmin-specific content
3. **Verify**: 
   - Pipeline statistics are visible
   - Employee management sections show Stage 2 and Stage 3 employees
   - Recent jobs table displays all jobs
   - "Create New Job" button is available

### Test 3: Create Job as Subadmin
1. **Login**: As subadmin
2. **Click**: "Create New Job" button
3. **Fill**: Required fields (Job Number is mandatory)
4. **Assign**: Select Stage 2 and Stage 3 employees
5. **Submit**: Create the job
6. **Verify**: Job appears in the recent jobs table

### Test 4: Employee Assignment Management
1. **Login**: As subadmin
2. **View**: Employee management sections
3. **Verify**: 
   - Stage 2 employees show job counts
   - Stage 3 employees show job counts
   - Employee details are displayed correctly

### Test 5: Pipeline Access
1. **Login**: As subadmin
2. **Navigate**: To `/pipeline`
3. **Verify**: 
   - Can see all jobs (not just assigned ones)
   - Can create new jobs
   - Can view job details

### Test 6: User Management Access
1. **Login**: As subadmin
2. **Navigate**: To `/users`
3. **Verify**: Can view all users
4. **Test**: Create a new employee user
5. **Verify**: New user appears in the list

### Test 7: Analytics and Reports
1. **Login**: As subadmin
2. **Navigate**: To `/analytics` and `/reports`
3. **Verify**: Can access both pages without errors

### Test 8: Role-Based Access Control
1. **Login**: As subadmin
2. **Test**: Access admin-only features
3. **Verify**: Cannot access features beyond subadmin permissions

## Backend API Changes

### Updated Endpoints
- `/api/pipeline/jobs` (GET/POST): Now accessible by subadmin
- `/api/users` (GET/POST): Now accessible by subadmin
- All job creation and management APIs support subadmin role

### New Methods Added
- `isSubadmin()`: Checks if user has subadmin role
- `isAdminOrSubadmin()`: Checks if user is admin or subadmin
- Updated `hasJobAccess()`: Subadmin has access to all jobs

## Frontend Changes

### New Pages
- `/dashboard/subadmin`: Subadmin dashboard
- Updated `/login`: Handles subadmin redirect
- Updated `/components/Sidebar`: Subadmin navigation

### Updated Components
- **Sidebar**: Added subadmin menu items and role display
- **Login**: Added subadmin role handling
- **Pipeline**: Subadmin can access all jobs and create new ones

## Security Considerations

### Subadmin Permissions
- ✅ Can create pipeline jobs
- ✅ Can view all jobs and users
- ✅ Can assign employees to jobs
- ✅ Can access analytics and reports
- ❌ Cannot create other admins
- ❌ Cannot modify system settings

### Role Hierarchy
1. **Admin**: Full system access
2. **Subadmin**: Job and user management
3. **Employees**: Stage-specific task access
4. **Customers**: Limited access to their jobs

## Troubleshooting

### Common Issues
1. **403 Forbidden**: Check if user has subadmin role
2. **Redirect Loop**: Clear localStorage and restart
3. **Missing Data**: Verify database has subadmin user
4. **Sidebar Issues**: Check isSubadmin prop is passed correctly

### Debug Commands
```bash
# Check if subadmin user exists
SELECT * FROM users WHERE role = 'subadmin';

# Verify subadmin permissions
curl -X GET http://localhost:8080/api/pipeline/jobs -H "Cookie: session=..."

# Test subadmin login
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"username":"subadmin","password":"subadmin123"}'
```

## Expected Behavior Summary

| Action | Admin | Subadmin | Employee |
|--------|-------|----------|----------|
| View all jobs | ✅ | ✅ | ❌ |
| Create jobs | ✅ | ✅ | ❌ |
| Manage users | ✅ | ✅ | ❌ |
| View analytics | ✅ | ✅ | ❌ |
| Access reports | ✅ | ✅ | ❌ |
| Create admins | ✅ | ❌ | ❌ |
| View assigned jobs | ✅ | ✅ | ✅ |

## Next Steps
1. Test all scenarios above
2. Verify subadmin can effectively manage employee assignments
3. Test with real employee data
4. Monitor system performance with subadmin role
5. Consider additional subadmin features if needed 