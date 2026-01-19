# E-Commerce Frontend

A React-based frontend for the Go E-Commerce application with product browsing, user authentication, and order management.

## Features

- ✅ User Authentication (Login/Signup)
- ✅ Product Listing
- ✅ Shopping Cart
- ✅ Place Orders
- ✅ View Order History
- ✅ Responsive Design

## Installation

1. Navigate to the client directory:
```bash
cd client
```

2. Install dependencies:
```bash
npm install
```

3. Create a `.env` file based on `.env.example`:
```bash
cp .env.example .env
```

4. Update the `VITE_API_URL` in `.env` to match your backend URL (default: `http://localhost:8080`)

## Development

Start the development server:
```bash
npm run dev
```

The application will be available at `http://localhost:3000`

## Building

Build for production:
```bash
npm run build
```

Preview the production build:
```bash
npm run preview
```

## Project Structure

```
src/
├── components/        # Reusable components (Navbar)
├── pages/             # Page components (Login, Signup, Products, Orders)
├── api/               # API client and endpoints
│   ├── client.js      # Axios instance with interceptors
│   ├── auth.js        # Authentication API
│   ├── products.js    # Products API
│   └── orders.js      # Orders API
├── App.jsx            # Main app component with routing
├── App.css            # App styles
├── index.css          # Global styles
└── main.jsx           # Entry point
```

## API Endpoints Used

- `POST /login` - User login
- `POST /signup` - User registration
- `GET /products` - Get all products
- `POST /orders` - Place an order
- `GET /orders` - Get user's orders

## Authentication

The app uses JWT token-based authentication. Tokens are stored in localStorage and automatically added to API requests via an axios interceptor.

## Environment Variables

- `VITE_API_URL` - Backend API base URL (default: `http://localhost:8080`)

## Technologies Used

- React 18
- Vite
- React Router v6
- Axios
- CSS3
