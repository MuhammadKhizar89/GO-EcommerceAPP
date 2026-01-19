import { Link } from 'react-router-dom'
import './Navbar.css'

export default function Navbar({ isAuthenticated, onLogout }) {
  return (
    <nav className="navbar">
      <div className="container">
        <div className="navbar-content">
          <Link to="/" className="navbar-brand">
            🛍️ E-Commerce
          </Link>
          <div className="navbar-links">
            <Link to="/products" className="nav-link">
              Products
            </Link>
            {isAuthenticated ? (
              <>
                <Link to="/orders" className="nav-link">
                  Orders
                </Link>
                <button onClick={onLogout} className="nav-link btn-logout">
                  Logout
                </button>
              </>
            ) : (
              <>
                <Link to="/login" className="nav-link">
                  Login
                </Link>
                <Link to="/signup" className="nav-link">
                  Sign Up
                </Link>
              </>
            )}
          </div>
        </div>
      </div>
    </nav>
  )
}
