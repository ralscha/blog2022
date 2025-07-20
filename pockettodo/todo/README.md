# PocketTodo - Todo Application

A comprehensive todo management application built with Ionic Angular and PocketBase backend.

## Features

- **User Authentication**: Login, registration, and password reset
- **Todo Management**: Create, read, update, and delete todos
- **Due Dates**: Set and track due dates for tasks
- **Profile Management**: Update user profile and email
- **Responsive Design**: Works on desktop and mobile devices
- **Real-time Sync**: Powered by PocketBase for instant updates

## Prerequisites

- Node.js (v18 or later)
- npm or yarn
- PocketBase (see `readme_pocketbase.md` for setup)

## Quick Start

### 1. Install Dependencies

```bash
npm install
```

### 2. Set up PocketBase

Follow the detailed instructions in `readme_pocketbase.md` to:

- Install and run PocketBase
- Configure the database schema
- Set up API rules and email settings

### 3. Configure Environment

Update the PocketBase URL in the environment files if needed:

**Development (`src/environments/environment.ts`):**

```typescript
export const environment = {
  production: false,
  pocketbaseUrl: 'http://127.0.0.1:8090'
};
```

**Production (`src/environments/environment.prod.ts`):**

```typescript
export const environment = {
  production: true,
  pocketbaseUrl: 'https://your-production-pocketbase-url.com'
};
```

### 4. Run the Application

```bash
# Development server
npm start

# Build for production
npm run build
```

The app will be available at `http://localhost:4200`

## Project Structure

```
src/
├── app/
│   ├── guards/
│   │   └── auth.guard.ts          # Route protection
│   ├── models/
│   │   ├── user.model.ts          # User type definitions
│   │   └── todo.model.ts          # Todo type definitions
│   ├── services/
│   │   └── pocketbase.service.ts  # API service layer
│   ├── login/                     # Login page
│   ├── register/                  # Registration page
│   ├── password-reset/            # Password reset page
│   ├── todos/                     # Todo list page
│   ├── edit-todo/                 # Todo creation/editing
│   ├── profile/                   # User profile page
│   ├── home/                      # Landing page
│   ├── app.component.ts           # Root component
│   └── app.routes.ts              # Route configuration
├── environments/                  # Environment configs
└── assets/                        # Static assets
```

## Usage

### Authentication

1. **Registration**: Create a new account with email and password
2. **Login**: Sign in with your credentials
3. **Password Reset**: Request password reset via email

### Todo Management

1. **View Todos**: See all your todos on the main page
2. **Create Todo**: Click the + button to add a new todo
3. **Edit Todo**: Swipe right or click edit to modify a todo
4. **Complete Todo**: Check the checkbox to mark as complete
5. **Delete Todo**: Swipe right and click delete, or use the delete button

### Profile Management

1. Access your profile from the todos page
2. Update your email address and name
3. Request password reset
4. Logout securely

## Development

### Available Scripts

- `npm start` - Start development server
- `npm run build` - Build for production
- `npm run watch` - Build in watch mode
- `ionic serve` - Alternative development server

### Adding New Features

1. **Create new pages**: Use Ionic CLI or create manually in `src/app/`
2. **Add routes**: Update `app.routes.ts`
3. **Protect routes**: Use `authGuard` or `guestGuard` as needed
4. **API calls**: Extend `pocketbase.service.ts`

### Code Style

- Use Angular standalone components
- Utilize Angular signals for reactive state
- Follow Ionic design patterns
- Implement proper error handling

## Deployment

### Frontend Deployment

1. **Build the app**:

   ```bash
   npm run build
   ```

2. **Deploy to hosting service** (Netlify, Vercel, etc.):
   - Upload the `dist/` folder
   - Configure environment variables

### Backend Deployment

See `readme_pocketbase.md` for PocketBase deployment options.

## Troubleshooting

### Common Issues

1. **PowerShell Execution Policy** (Windows):

   ```powershell
   Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser
   ```

2. **PocketBase Connection Issues**:
   - Ensure PocketBase is running on the correct port
   - Check CORS settings in PocketBase admin panel
   - Verify environment configuration

3. **Email Not Sending**:
   - Configure SMTP settings in PocketBase
   - Check email provider settings
   - Verify API rules allow password reset

4. **Authentication Issues**:
   - Clear browser storage/cookies
   - Check API rules in PocketBase
   - Verify user collection configuration

### Development Tips

- Use browser developer tools to debug API calls
- Check PocketBase logs for backend issues
- Use Ionic DevApp for mobile testing
- Enable source maps for better debugging

## API Endpoints

The app uses these PocketBase endpoints:

- **Authentication**: `/api/collections/users/auth-with-password`
- **Registration**: `/api/collections/users/records`
- **Password Reset**: `/api/collections/users/request-password-reset`
- **Todos**: `/api/collections/todos/records`
- **Profile**: `/api/collections/users/records/:id`

## Security Features

- JWT-based authentication
- Password hashing (handled by PocketBase)
- API rules for data access control
- Route guards for page protection
- CORS configuration
- Input validation and sanitization

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test thoroughly
5. Submit a pull request

## License

This project is open source. Feel free to use and modify as needed.

## Support

For issues and questions:

1. Check the troubleshooting section
2. Review PocketBase documentation
3. Check Ionic Angular documentation
4. Create an issue in the repository

## Technologies Used

- **Frontend**: Ionic 8, Angular 20, TypeScript
- **Backend**: PocketBase
- **State Management**: Angular Signals
- **Routing**: Angular Router
- **Forms**: Reactive Forms
- **Icons**: Ionicons
- **Styling**: Ionic CSS Variables
