# Mi Bolsillo - Frontend

A minimalistic expense tracker application built with React, TypeScript, and Tailwind CSS.

## Features

- **Authentication**: Sign up, login, and sign out functionality
- **Bill Management**: Create, view, and delete bills
- **Expense Tracking**: Add multiple expenses to each bill with categories
- **Responsive Design**: Clean, minimalistic UI that works on all devices

## Tech Stack

- **React 19** - UI library
- **TypeScript** - Type safety
- **Vite** - Build tool and dev server
- **React Router** - Client-side routing
- **Axios** - HTTP client
- **Tailwind CSS** - Utility-first CSS framework
- **Context API** - State management

## Project Structure

```
src/
├── components/       # Reusable components (ProtectedRoute)
├── contexts/        # React Context providers (AuthContext)
├── layouts/         # Layout components (MainLayout)
├── pages/           # Page components (Dashboard, Login, etc.)
├── services/        # API service layer
├── types/           # TypeScript type definitions
└── utils/           # Utility functions
```

## Getting Started

### Prerequisites

- Node.js 20.x or higher
- npm or yarn

### Installation

1. Install dependencies:
```bash
npm install
```

2. Create a `.env` file (copy from `.env.example`):
```bash
cp .env.example .env
```

3. Update the API URL in `.env`:
```
VITE_API_URL=http://localhost:8080
```

### Development

Run the development server:
```bash
npm run dev
```

The app will be available at `http://localhost:5173`

### Build

Create a production build:
```bash
npm run build
```

### Preview Production Build

Preview the production build locally:
```bash
npm run preview
```

## API Integration

The frontend expects the backend API to be running and accessible at the URL specified in `VITE_API_URL`.

### Required API Endpoints

- `POST /auth/login` - User login
- `POST /auth/signup` - User registration
- `GET /bills` - Get all bills
- `GET /bills/:id` - Get bill by ID
- `POST /bills` - Create new bill
- `PUT /bills/:id` - Update bill
- `DELETE /bills/:id` - Delete bill

## Authentication

The app uses JWT tokens stored in localStorage for authentication. Tokens are automatically included in API requests via Axios interceptors.

## License

MIT
