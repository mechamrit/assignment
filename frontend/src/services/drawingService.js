import api from './apiConfig';

export const drawingService = {
    getAll: async (projectID) => {
        const response = await api.get(`/drawings?project_id=${projectID}`);
        return response.data;
    },

    create: async (drawingData) => {
        const response = await api.post('/drawings', drawingData);
        return response.data;
    },

    claim: async (id) => {
        const response = await api.post(`/drawings/${id}/claim`);
        return response.data;
    },

    submit: async (id) => {
        const response = await api.post(`/drawings/${id}/submit`);
        return response.data;
    },

    release: async (id) => {
        const response = await api.post(`/drawings/${id}/release`);
        return response.data;
    },

    reject: async (id) => {
        const response = await api.post(`/drawings/${id}/reject`);
        return response.data;
    }
};
