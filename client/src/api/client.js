import axios from 'axios'
import toast from 'react-hot-toast'

const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080'
const REQUEST_TIMEOUT = 5000 // 5 seconds

const apiClient = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
  timeout: REQUEST_TIMEOUT,
})

// Add token to requests
apiClient.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// Handle timeout and connection errors
apiClient.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.code === 'ECONNABORTED' || error.message === 'timeout of ' + REQUEST_TIMEOUT + 'ms exceeded') {
      toast.error('Backend is inactive. It will be active in ~50 seconds. Please wait.')
    } else if (!error.response) {
      toast.error('Backend is inactive. It will be active in ~50 seconds. Please wait.')
    }
    return Promise.reject(error)
  }
)

export default apiClient
