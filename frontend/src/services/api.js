import axios from 'axios';

const API = axios.create({
  baseURL: '/api',
});

export const fetchUsers = () => API.get('/users');
