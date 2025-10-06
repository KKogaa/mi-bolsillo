# Mi Bolsillo Frontend Architecture

## Overview

This is a production-ready React application built with TypeScript, following clean architecture principles and best practices.

## Architecture Patterns

### 1. **Layered Architecture**

```
Presentation Layer (UI)
    ↓
Application Layer (Context/State)
    ↓
Service Layer (API)
    ↓
External Services (Backend API)
```

### 2. **Project Structure**

```
src/
├── components/          # Reusable UI components
│   ├── ProtectedRoute.tsx
│   └── index.ts
├── contexts/           # React Context for state management
│   └── AuthContext.tsx
├── layouts/            # Page layouts
│   └── MainLayout.tsx
├── pages/              # Route-level components
│   ├── Login.tsx
│   ├── SignUp.tsx
│   ├── Dashboard.tsx
│   ├── CreateBill.tsx
│   ├── BillDetail.tsx
│   └── index.ts
├── services/           # API service layer
│   └── api.ts
├── types/              # TypeScript type definitions
│   └── index.ts
├── utils/              # Utility functions
├── hooks/              # Custom React hooks
├── App.tsx             # Main app component with routing
├── main.tsx            # Application entry point
└── index.css           # Global styles (Tailwind)
```

## Key Design Patterns

### 1. **Context Pattern**
- **AuthContext** manages authentication state globally
- Provides `useAuth` hook for easy access to auth state and methods
- Centralizes authentication logic

### 2. **Protected Routes**
- `ProtectedRoute` component wraps private routes
- Redirects unauthenticated users to login
- Shows loading state during auth check

### 3. **Service Layer Pattern**
- All API calls abstracted into service modules
- `authService` handles authentication
- `billService` handles bill operations
- Centralized axios instance with interceptors

### 4. **Layout Components**
- `MainLayout` provides consistent navigation and structure
- Separates layout from page content
- Reusable across all authenticated pages

## Data Flow

### Authentication Flow
```
1. User submits login form
2. Login component calls useAuth hook
3. AuthContext calls authService.login()
4. Service makes API request
5. Token stored in localStorage
6. User state updated in context
7. Navigate to dashboard
```

### Bill Management Flow
```
1. User creates bill
2. CreateBill component calls billService.create()
3. Service sends data to API with auth token
4. Navigate to dashboard
5. Dashboard fetches updated bill list
```

## State Management

### Global State (Context)
- **User authentication state**
- **Current user data**

### Local State (Component)
- **Form inputs**
- **Loading states**
- **Error messages**
- **Component-specific UI state**

### Server State
- **Bills data** (fetched on demand)
- No client-side caching (simple approach)

## Security Features

### 1. **JWT Token Management**
- Tokens stored in localStorage
- Automatically attached to requests via axios interceptor
- Auto-logout on 401 responses

### 2. **Route Protection**
- Protected routes redirect unauthenticated users
- Client-side route guards

### 3. **Input Validation**
- Form validation on client side
- Type safety with TypeScript
- Server-side validation expected

## API Integration

### Axios Configuration
- Base URL from environment variables
- Request interceptor adds auth token
- Response interceptor handles auth errors
- Centralized error handling

### Environment Variables
```
VITE_API_URL - Backend API base URL
```

## Styling Approach

### Tailwind CSS
- Utility-first CSS framework
- Minimalistic design system
- No custom CSS needed
- Responsive by default
- Production build purges unused styles

### Design Principles
- **Clean and minimal** - Focus on content
- **Consistent spacing** - Tailwind spacing scale
- **Accessible** - Proper contrast, focus states
- **Responsive** - Mobile-first design

## Build & Deployment

### Development
```bash
npm run dev
```
- Hot module replacement
- Fast refresh
- Source maps

### Production Build
```bash
npm run build
```
- TypeScript compilation
- Vite optimization
- CSS minification
- Tree shaking
- Code splitting

### Output
```
dist/
├── assets/
│   ├── index-[hash].js   # Main bundle
│   └── index-[hash].css  # Styles
└── index.html            # Entry point
```

## Best Practices Implemented

1. **TypeScript** - Full type safety
2. **Separation of Concerns** - Clear layer boundaries
3. **DRY Principle** - Reusable components and services
4. **Single Responsibility** - Each module has one job
5. **Error Handling** - Consistent error management
6. **Loading States** - Better UX with loading indicators
7. **Environment Configuration** - Separate dev/prod settings
8. **Clean Code** - Readable and maintainable

## Future Improvements

- [ ] Add React Query for server state management
- [ ] Implement form library (React Hook Form)
- [ ] Add comprehensive error boundaries
- [ ] Add unit and integration tests
- [ ] Add E2E tests with Playwright
- [ ] Implement proper logging
- [ ] Add analytics
- [ ] Progressive Web App features
- [ ] Internationalization (i18n)
- [ ] Dark mode support
