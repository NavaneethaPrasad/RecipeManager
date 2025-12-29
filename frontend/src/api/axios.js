import axios from 'axios';

const api = axios.create({
    baseURL: 'http://localhost:8080/api',
    withCredentials: true, // 🔑 THIS IS THE FIX
});

export default api;
