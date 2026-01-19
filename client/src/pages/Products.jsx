import { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import { productsAPI } from '../api/products'
import { ordersAPI } from '../api/orders'
import './Products.css'

export default function Products() {
  const [products, setProducts] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [cart, setCart] = useState({})
  const [isOrdering, setIsOrdering] = useState(false)
  const navigate = useNavigate()

  useEffect(() => {
    fetchProducts()
  }, [])

  const fetchProducts = async () => {
    try {
      setLoading(true)
      const response = await productsAPI.listProducts()
      setProducts(response.data.data || [])
    } catch (err) {
      setError('Failed to load products')
      console.error(err)
    } finally {
      setLoading(false)
    }
  }

  const handleAddToCart = (productId, quantity = 1) => {
    setCart((prevCart) => ({
      ...prevCart,
      [productId]: (prevCart[productId] || 0) + quantity,
    }))
  }

  const handleRemoveFromCart = (productId) => {
    setCart((prevCart) => {
      const newCart = { ...prevCart }
      delete newCart[productId]
      return newCart
    })
  }

  const handlePlaceOrder = async () => {
    if (Object.keys(cart).length === 0) {
      alert('Please add items to cart')
      return
    }

    const token = localStorage.getItem('token')
    if (!token) {
      navigate('/login')
      return
    }

    setIsOrdering(true)
    try {
      const items = Object.entries(cart).map(([productId, quantity]) => ({
        productId: parseInt(productId),
        quantity: quantity,
      }))
      await ordersAPI.placeOrder({ items })
      alert('Order placed successfully!')
      setCart({})
      navigate('/orders')
    } catch (err) {
      alert(err.response?.data?.message || 'Failed to place order')
    } finally {
      setIsOrdering(false)
    }
  }

  if (loading) {
    return <div className="container" style={{ padding: '40px 0', textAlign: 'center' }}>Loading products...</div>
  }

  return (
    <div className="products-page">
      <div className="container">
        {error && <div className="alert alert-error">{error}</div>}

        <div className="products-content">
          <div className="products-grid">
            <h2>Products</h2>
            {products.length === 0 ? (
              <p>No products available</p>
            ) : (
              <div className="grid">
                {products.map((product) => (
                  <div key={product.id} className="product-card">
                    {product.image && (
                      <img src={product.image} alt={product.name} className="product-image" />
                    )}
                    <div className="product-info">
                      <h3>{product.name}</h3>
                      <p className="product-description">{product.description}</p>
                      <p className="product-price">${product.price}</p>
                      <div className="product-actions">
                        <input
                          type="number"
                          min="1"
                          max="10"
                          defaultValue="1"
                          id={`qty-${product.id}`}
                          className="qty-input"
                        />
                        <button
                          onClick={() => {
                            const qty = parseInt(document.getElementById(`qty-${product.id}`).value)
                            handleAddToCart(product.id, qty)
                          }}
                          className="btn btn-success"
                        >
                          Add to Cart
                        </button>
                      </div>
                    </div>
                  </div>
                ))}
              </div>
            )}
          </div>

          <div className="cart-sidebar">
            <h3>Shopping Cart</h3>
            {Object.keys(cart).length === 0 ? (
              <p className="empty-cart">Your cart is empty</p>
            ) : (
              <>
                <div className="cart-items">
                  {Object.entries(cart).map(([productId, quantity]) => {
                    const product = products.find((p) => p.id === parseInt(productId))
                    if (!product) return null
                    return (
                      <div key={productId} className="cart-item">
                        <div className="cart-item-info">
                          <p className="cart-item-name">{product.name}</p>
                          <p className="cart-item-qty">Qty: {quantity}</p>
                          <p className="cart-item-price">${product.price * quantity}</p>
                        </div>
                        <button
                          onClick={() => handleRemoveFromCart(productId)}
                          className="btn btn-danger btn-small"
                        >
                          Remove
                        </button>
                      </div>
                    )
                  })}
                </div>
                <div className="cart-total">
                  <p>
                    Total: $
                    {Object.entries(cart)
                      .reduce((total, [productId, quantity]) => {
                        const product = products.find((p) => p.id === parseInt(productId))
                        return total + (product?.price || 0) * quantity
                      }, 0)
                      .toFixed(2)}
                  </p>
                </div>
                <button
                  onClick={handlePlaceOrder}
                  className="btn btn-primary btn-full"
                  disabled={isOrdering}
                >
                  {isOrdering ? 'Placing Order...' : 'Place Order'}
                </button>
              </>
            )}
          </div>
        </div>
      </div>
    </div>
  )
}
