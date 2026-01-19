import apiClient from './client'

export const productsAPI = {
  listProducts: () => {
    return apiClient.get('/products')
  },

  createProduct: (product) => {
    return apiClient.post('/products', product)
  },
}
