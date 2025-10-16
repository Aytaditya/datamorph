import axios from 'axios';
import toast from 'react-hot-toast';

// Use HTTPS for production with SSL certificates
const API_BASE_URL = 'https://localhost:8080';

// Create axios instance
const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
  timeout: 10000, // 10 second timeout
});

// Request interceptor to add auth token
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => Promise.reject(error)
);

// Response interceptor for error handling
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('token');
      window.location.href = '/login';
      toast.error('Session expired. Please login again.');
    } else if (error.response?.data?.error) {
      toast.error(error.response.data.error);
    } else if (error.message) {
      toast.error(error.message);
    }
    return Promise.reject(error);
  }
);

// Auth API
export const authAPI = {
  login: async (username, password) => {
    const response = await api.post('/login', { username, password });
    return response.data;
  }
};

// Clients API
export const clientsAPI = {
  getAll: async () => {
    const response = await api.get('/clients');
    return response.data;
  },
  
  create: async (name) => {
    const response = await api.post('/clients', { name });
    return response.data;
  },
  
  delete: async (id) => {
    await api.delete(`/clients/${id}`);
  }
};

// Mapping Rules API
export const mappingAPI = {
  getByClient: async (clientId) => {
    const response = await api.get(`/clients/${clientId}/mappings`);
    return response.data;
  },
  
  create: async (clientId, rules) => {
    const response = await api.post(`/clients/${clientId}/mappings`, rules);
    return response.data;
  },
  
  delete: async (mappingId) => {
    await api.delete(`/mappings/${mappingId}`);
  }
};

// Transform API
export const transformAPI = {
  transform: async (clientId, inputData) => {
    const response = await api.post(`/clients/${clientId}/transform`, {
      input_data: inputData
    });
    return response.data;
  }
};

export default api;
