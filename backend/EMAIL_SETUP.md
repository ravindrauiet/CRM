# Email Notification Setup Guide

This guide will help you set up email notifications for the MayDiv CRM system.

## Overview

The system sends email notifications to admin users when:
1. A new pipeline job is created
2. A pipeline stage is completed (Stage 2, 3, or 4)

## Email Configuration

Add the following environment variables to your `.env` file:

```env
# Email Configuration (SMTP)
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=your-email@gmail.com
SMTP_PASS=your-app-password
FROM_EMAIL=your-email@gmail.com

# Admin Email (where notifications will be sent)
ADMIN_EMAIL=admin@maydiv.com
```

## Gmail Setup (Recommended)

### 1. Enable 2-Factor Authentication
- Go to your Google Account settings
- Enable 2-Factor Authentication

### 2. Generate App Password
- Go to Google Account settings
- Navigate to Security > 2-Step Verification > App passwords
- Generate a new app password for "Mail"
- Use this password as `SMTP_PASS`

### 3. Configuration Example
```env
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=yourname@gmail.com
SMTP_PASS=abcd efgh ijkl mnop
FROM_EMAIL=yourname@gmail.com
ADMIN_EMAIL=admin@maydiv.com
```

## Other Email Providers

### Outlook/Hotmail
```env
SMTP_HOST=smtp-mail.outlook.com
SMTP_PORT=587
SMTP_USER=your-email@outlook.com
SMTP_PASS=your-password
```

### Yahoo
```env
SMTP_HOST=smtp.mail.yahoo.com
SMTP_PORT=587
SMTP_USER=your-email@yahoo.com
SMTP_PASS=your-app-password
```

## Testing Email Configuration

The system includes a test endpoint to verify email configuration:

```bash
curl -X POST http://localhost:8080/api/test-email
```

## Email Templates

The system sends beautifully formatted HTML emails with:

### Stage Completion Email
- Job number and details
- Completed stage information
- User who completed the stage
- Next stage information
- Professional styling

### Job Creation Email
- New job details
- Creator information
- Action required notification

## Troubleshooting

### Common Issues

1. **Authentication Failed**
   - Check your email and password
   - For Gmail, use App Password instead of regular password
   - Ensure 2-Factor Authentication is enabled

2. **Connection Timeout**
   - Check SMTP_HOST and SMTP_PORT
   - Verify firewall settings
   - Try different SMTP servers

3. **Emails Not Sending**
   - Check server logs for error messages
   - Verify ADMIN_EMAIL is set correctly
   - Test email configuration

### Debug Mode

Enable debug logging by setting:
```env
DEBUG=true
```

## Security Notes

- Never commit your `.env` file to version control
- Use App Passwords instead of regular passwords
- Consider using environment-specific email addresses
- Regularly rotate your SMTP credentials

## Production Considerations

1. **Email Service**: Consider using dedicated email services like:
   - SendGrid
   - Mailgun
   - Amazon SES

2. **Rate Limiting**: Implement rate limiting for email sending

3. **Error Handling**: Set up proper error handling and retry mechanisms

4. **Monitoring**: Monitor email delivery rates and failures 