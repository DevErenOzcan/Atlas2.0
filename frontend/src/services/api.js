import axios from 'axios';

const api = axios.create({
  // Use environment variable if provided (for local/dev). In Kubernetes/Ingress use relative path so browser requests go to the same host and Ingress routes /api to the gateway service.
  baseURL: process.env.REACT_APP_API_URL || '/api',
});

// Token'ı header'a ekle
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// Auth API
export const authAPI = {
  register: (data) => api.post('/auth/register', data),
  login: (data) => api.post('/auth/login', data),
};

// User API
export const userAPI = {
  getAll: () => api.get('/users'),
  getById: (id) => api.get(`/users/${id}`),
  update: (id, data) => api.put(`/users/${id}`, data),
  delete: (id) => api.delete(`/users/${id}`),
};

// Product API
export const productAPI = {
  getAll: () => api.get('/products'),
  getMy: () => api.get('/products/my'), // Kullanıcının kendi ürünleri (role'e göre filtrelenmiş)
  getById: (id) => api.get(`/products/${id}`),
  create: (data) => api.post('/products', data),
  update: (id, data) => api.put(`/products/${id}`, data),
  delete: (id) => api.delete(`/products/${id}`),
};

// Order API
export const orderAPI = {
  getAll: () => api.get('/orders'),
  getMy: () => api.get('/orders/my'), // Kullanıcının kendi siparişleri (role'e göre filtrelenmiş)
  create: (data) => api.post('/orders', data),
  updateStatus: (id, is_shipped) => api.patch(`/orders/${id}/status`, { is_shipped }),
};

export default api;
