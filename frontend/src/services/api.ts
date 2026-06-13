import axios from "axios";
import { useAuthStore } from "../store/authStore";

const API_URL = import.meta.env.VITE_API_URL || "http://localhost:8001/api/v1";

const apiClient = axios.create({
  baseURL: API_URL,
  headers: {
    "Content-Type": "application/json",
  },
});

apiClient.interceptors.request.use((config) => {
  const token = useAuthStore.getState().token;
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

apiClient.interceptors.response.use(
  (response) => response,
  async (error) => {
    const originalRequest = error.config;
    if (error.response?.status === 401 && !originalRequest._retry) {
      originalRequest._retry = true;
      const refreshToken = useAuthStore.getState().refreshToken;
      if (refreshToken) {
        try {
          const { data } = await axios.post(
            `${API_URL}/auth/refresh`,
            {
              refresh_token: refreshToken,
            },
          );
          if (data.success && data.data.access_token) {
            useAuthStore
              .getState()
              .setToken(data.data.access_token, data.data.refresh_token);
            originalRequest.headers.Authorization = `Bearer ${data.data.access_token}`;
            return apiClient(originalRequest);
          }
        } catch (refreshError) {
          useAuthStore.getState().logout();
          window.location.href = "/login";
        }
      } else {
        useAuthStore.getState().logout();
        window.location.href = "/login";
      }
    }
    return Promise.reject(error);
  },
);

export const authService = {
  login: async (email: string, password: string) => {
    const { data } = await apiClient.post("/auth/login", { email, password });
    return data;
  },
  getProfile: async () => {
    const { data } = await apiClient.get("/auth/profile");
    return data;
  },
};

export const serverService = {
  list: async (page = 1, limit = 10) => {
    const { data } = await apiClient.get(
      `/servers?page=${page}&limit=${limit}`,
    );
    return data;
  },
  get: async (id: number) => {
    const { data } = await apiClient.get(`/servers/${id}`);
    return data;
  },
  create: async (payload: any) => {
    const { data } = await apiClient.post("/servers", payload);
    return data;
  },
  update: async (id: number, payload: any) => {
    const { data } = await apiClient.put(`/servers/${id}`, payload);
    return data;
  },
  delete: async (id: number) => {
    const { data } = await apiClient.delete(`/servers/${id}`);
    return data;
  },
};

export const dashboardService = {
  getSummary: async () => {
    const { data } = await apiClient.get("/dashboard/summary");
    return data;
  },
  getTrend: async (serverId: number, range: string) => {
    const { data } = await apiClient.get(
      `/dashboard/servers/${serverId}/trend?range=${range}`,
    );
    return data;
  },
};

export const alertService = {
  list: async (state?: string, serverId?: number) => {
    let url = "/alerts";
    const params = new URLSearchParams();
    if (state) params.append("state", state);
    if (serverId) params.append("server_id", serverId.toString());

    const { data } = await apiClient.get(`${url}?${params.toString()}`);
    return data;
  },
};

export const settingService = {
  getAll: async () => {
    const { data } = await apiClient.get("/settings");
    return data;
  },
  save: async (settings: Record<string, string>) => {
    const { data } = await apiClient.post("/settings", settings);
    return data;
  },
};
