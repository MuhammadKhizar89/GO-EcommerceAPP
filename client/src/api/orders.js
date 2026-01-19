import apiClient from './client'

export const ordersAPI = {
  placeOrder: (orderData) => {
    return apiClient.post('/orders', orderData)
  },

  getAllOrders: () => {
    return apiClient.get('/orders')
  },
}
