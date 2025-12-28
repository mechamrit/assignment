import api from './apiConfig';

export const authService = {
    login: async (username, password) => {
        const response = await api.post('/login', { username, password });
        return response.data;
    },

    register: async (userData) => {
        const response = await api.post('/register', userData);
        return response.data;
    },

    logout: () => {
        localStorage.removeItem('token');
        localStorage.removeItem('role');
        localStorage.removeItem('user_id');
        window.location.href = '/login';
    }
};
