# PocketBase Setup Instructions

## Installation

### Download Binary

1. Visit [PocketBase releases](https://github.com/pocketbase/pocketbase/releases)
2. Download the appropriate binary for your OS (Windows: `pocketbase_X.X.X_windows_amd64.zip`)
3. Extract the binary to a folder (e.g., `C:\pocketbase\`)
4. Add the folder to your PATH environment variable (optional)

## Running PocketBase

1. Open PowerShell/Command Prompt
2. Navigate to the folder containing the PocketBase binary
3. Run the server:

```bash
.\pocketbase.exe serve
```

The server will start at `http://127.0.0.1:8090`

## Initial Setup

1. Open your browser and go to `http://127.0.0.1:8090/_/`
2. Create an admin account (email and password)
3. You'll be redirected to the admin dashboard

## Database Schema Setup

### 1. Users Collection (Built-in)

The `users` collection is built-in. We need to configure it:

1. Go to **Collections** → **users**
2. Click on **Settings** tab
3. Ensure these fields exist (they should be there by default):
   - `id` (Primary key)
   - `email` (Email)
   - `password` (Password - hidden)
   - `name` (Text - optional)
   - `avatar` (File - optional)

### 2. Create Todos Collection

1. Click **New Collection**
2. Name: `todos`
3. Type: **Base collection**
4. Add the following fields:

   **title** (Text)
   - Name: `title`
   - Type: Text
   - Required: ✓
   - Max length: 255

   **description** (Text)
   - Name: `description`
   - Type: Text
   - Required: ✗
   - Max length: 1000

   **completed** (Bool)
   - Name: `completed`
   - Type: Bool
   - Required: ✓
   - Default: false

   **user** (Relation)
   - Name: `user`
   - Type: Relation
   - Required: ✓
   - Collection: users
   - Max select: 1

   **due_date** (DateTime)
   - Name: `due_date`
   - Type: DateTime
   - Required: ✗

5. Click **Create**

## API Rules Configuration

### Todos Collection Rules

1. Go to **Collections** → **todos** → **API Rules**
2. Set the following rules:

   **List/Search rule:**

   ```
   @request.auth.id != "" && user = @request.auth.id
   ```

   **View rule:**

   ```
   @request.auth.id != "" && user = @request.auth.id
   ```

   **Create rule:**

   ```
   @request.auth.id != "" && @request.body.user = @request.auth.id
   ```

   **Update rule:**

   ```
   @request.auth.id != "" && user = @request.auth.id
   ```

   **Delete rule:**

   ```
   @request.auth.id != "" && user = @request.auth.id
   ```

## Email Configuration (for password reset)

1. Go to **Settings** → **Mail settings**
2. Configure your SMTP settings:
   - **SMTP server host**: Your email provider's SMTP server
   - **Port**: Usually 587 for TLS or 465 for SSL
   - **Username**: Your email address
   - **Password**: Your email password or app password
   - **TLS**: Enable if using port 587

### Example configurations:

**Gmail:**

- Host: `smtp.gmail.com`
- Port: `587`
- TLS: ✓
- Username: `your-email@gmail.com`
- Password: Use an [App Password](https://support.google.com/accounts/answer/185833)

**Outlook/Hotmail:**

- Host: `smtp.live.com`
- Port: `587`
- TLS: ✓

## Development Environment Variables

Create a `.env` file in your Angular project root:

```env
POCKETBASE_URL=http://127.0.0.1:8090
```

## Production Deployment

### Option 1: Self-hosted

1. Deploy PocketBase binary to your server
2. Set up a reverse proxy (nginx/Apache)
3. Configure SSL certificate
4. Update CORS settings in PocketBase admin panel

### Option 2: PocketHost (Managed Hosting)

1. Visit [PocketHost.io](https://pockethost.io)
2. Create an account and deploy your PocketBase instance
3. Import your schema and data

## Backup and Migration

### Export Data

```bash
.\pocketbase.exe admin export backup.zip
```

### Import Data

```bash
.\pocketbase.exe admin import backup.zip
```

## API Endpoints

Base URL: `http://127.0.0.1:8090/api/`

### Authentication

- **Login**: `POST /collections/users/auth-with-password`
- **Register**: `POST /collections/users/records`
- **Refresh**: `POST /collections/users/auth-refresh`
- **Password Reset**: `POST /collections/users/request-password-reset`
- **Confirm Password Reset**: `POST /collections/users/confirm-password-reset`

### Todos

- **List**: `GET /collections/todos/records`
- **Create**: `POST /collections/todos/records`
- **Update**: `PATCH /collections/todos/records/:id`
- **Delete**: `DELETE /collections/todos/records/:id`

### Users

- **Get Profile**: `GET /collections/users/records/:id`
- **Update Profile**: `PATCH /collections/users/records/:id`

## Testing the API

You can test the API using the built-in API preview in the PocketBase admin panel:

1. Go to **Collections** → Select a collection
2. Click on **API Preview** tab
3. Test different endpoints with example requests

## Troubleshooting

### Common Issues

1. **CORS errors**: Make sure to configure CORS settings in PocketBase admin panel under Settings → Application
2. **Authentication errors**: Check that your API rules are properly configured
3. **Email not sending**: Verify SMTP settings and check PocketBase logs
4. **File upload issues**: Check file upload settings and size limits

### Logs

PocketBase logs are displayed in the console where you started the server. For production, consider redirecting logs to files.

## Security Considerations

1. **Change default admin password** after setup
2. **Use HTTPS** in production
3. **Configure proper API rules** to prevent unauthorized access
4. **Regular backups** of your data
5. **Keep PocketBase updated** to the latest version
6. **Use strong SMTP passwords** and consider app-specific passwords

## Additional Resources

- [PocketBase Documentation](https://pocketbase.io/docs/)
- [PocketBase JavaScript SDK](https://github.com/pocketbase/js-sdk)
- [API Documentation](https://pocketbase.io/docs/api-authentication/)
