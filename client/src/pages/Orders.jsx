import { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import toast from 'react-hot-toast'
import { ordersAPI } from '../api/orders'
import './Orders.css'

export default function Orders() {
  const [orders, setOrders] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const navigate = useNavigate()

  useEffect(() => {
    const token = localStorage.getItem('token')
    if (!token) {
      navigate('/login')
      return
    }
    fetchOrders()
  }, [navigate])

  const fetchOrders = async () => {
    try {
      setLoading(true)
      const response = await ordersAPI.getAllOrders()
      setOrders(response.data.data || [])
    } catch (err) {
      const errorMessage = 'Failed to load orders'
      setError(errorMessage)
      toast.error(errorMessage)
      console.error(err)
    } finally {
      setLoading(false)
    }
  }

  const formatDate = (dateString) => {
    return new Date(dateString).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    })
  }

  if (loading) {
    return <div className="container" style={{ padding: '40px 0', textAlign: 'center' }}>Loading orders...</div>
  }

  return (
    <div className="orders-page">
      <div className="container">
        <h2>My Orders</h2>
        {error && <div className="alert alert-error">{error}</div>}

        {orders.length === 0 ? (
          <div className="no-orders">
            <p>You haven't placed any orders yet</p>
            <a href="/products" className="btn btn-primary">
              Continue Shopping
            </a>
          </div>
        ) : (
          <div className="orders-list">
            {orders.map((order) => {
              const totalAmount = order.items?.reduce((sum, item) => sum + (item.price * item.quantity), 0) || 0
              return (
                <div key={order.id} className="order-card">
                  <div className="order-header">
                    <div>
                      <h3>Order #{order.id}</h3>
                      <p className="order-date">{formatDate(order.created_at)}</p>
                    </div>
                    <div className="order-status">
                      <span className={`status status-${order.status || 'pending'}`}>
                        {order.status || 'Pending'}
                      </span>
                    </div>
                  </div>
                  
                  <div className="order-items-section">
                    <h4>Items</h4>
                    {order.items && order.items.length > 0 ? (
                      <div className="order-items">
                        {order.items.map((item) => (
                          <div key={item.id} className="order-item">
                            <div className="item-image">
                              {item.product.image ? (
                                <img src={item.product.image} alt={item.product.name} />
                              ) : (
                                <div className="no-image">No Image</div>
                              )}
                            </div>
                            <div className="item-info">
                              <p className="item-name">{item.product.name}</p>
                              <p className="item-details">Quantity: {item.quantity}</p>
                              <p className="item-price">${item.price}</p>
                            </div>
                            <div className="item-total">
                              <p>${(item.price * item.quantity).toFixed(2)}</p>
                            </div>
                          </div>
                        ))}
                      </div>
                    ) : (
                      <p className="no-items">No items in this order</p>
                    )}
                  </div>
                  
                  <div className="order-total">
                    <span className="label">Total Amount:</span>
                    <span className="value">${totalAmount.toFixed(2)}</span>
                  </div>
                </div>
              )
            })}
          </div>
        )}
      </div>
    </div>
  )
}
