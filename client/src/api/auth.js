import apiClient from './client'

export const authAPI = {
  login: (email, password) => {
    return apiClient.post('/login', { email, password })
  },

  signup: (email, password) => {
    return apiClient.post('/signup', { email, password })
  },
}
