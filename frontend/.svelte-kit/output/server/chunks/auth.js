import { w as writable } from "./index.js";
const API_BASE_URL = "http://localhost:8000";
class ApiClient {
  baseUrl;
  token = null;
  constructor(baseUrl) {
    this.baseUrl = baseUrl;
  }
  setToken(token) {
    this.token = token;
  }
  async request(endpoint, options = {}) {
    const { params, ...init } = options;
    let url = `${this.baseUrl}${endpoint}`;
    if (params) {
      const searchParams = new URLSearchParams(params);
      url += `?${searchParams.toString()}`;
    }
    const headers = {
      "Content-Type": "application/json",
      ...init.headers || {}
    };
    if (this.token) {
      headers["Authorization"] = `Bearer ${this.token}`;
    }
    const response = await fetch(url, {
      ...init,
      headers
    });
    if (!response.ok) {
      const error = await response.json().catch(() => ({ message: "Request failed" }));
      throw new Error(error.message || `API Error: ${response.status} ${response.statusText}`);
    }
    return response.json();
  }
  async get(endpoint, options) {
    return this.request(endpoint, { ...options, method: "GET" });
  }
  async post(endpoint, data, options) {
    return this.request(endpoint, {
      ...options,
      method: "POST",
      body: data ? JSON.stringify(data) : void 0
    });
  }
  async put(endpoint, data, options) {
    return this.request(endpoint, {
      ...options,
      method: "PUT",
      body: data ? JSON.stringify(data) : void 0
    });
  }
  async delete(endpoint, options) {
    return this.request(endpoint, { ...options, method: "DELETE" });
  }
  login(email, password) {
    return this.post("/api/auth/login", { email, password });
  }
  register(email, password) {
    return this.post("/api/auth/register", { email, password });
  }
  getMe() {
    return this.get("/api/auth/me");
  }
  getDevices() {
    return this.get("/api/devices");
  }
  deleteDevice(id) {
    return this.delete(`/api/devices/${id}`);
  }
  updateDevice(id, name) {
    return this.put(`/api/devices/${id}`, { name });
  }
  getRules() {
    return this.get("/api/rules");
  }
  createRule(data) {
    return this.post("/api/rules", data);
  }
  updateRule(id, data) {
    return this.put(`/api/rules/${id}`, data);
  }
  deleteRule(id) {
    return this.delete(`/api/rules/${id}`);
  }
  testRule(id) {
    return this.post(`/api/rules/${id}/test`);
  }
  getPairingToken() {
    return this.post("/api/devices/pairing-token");
  }
  checkPairingStatus(token) {
    return this.get(`/api/devices/pairing-status/${token}`);
  }
  getDashboardStats() {
    return this.get("/api/logs/stats");
  }
  getLogs(params) {
    const searchParams = {};
    if (params.page) searchParams.page = String(params.page);
    if (params.limit) searchParams.limit = String(params.limit);
    if (params.direction) searchParams.direction = params.direction;
    if (params.status) searchParams.status = params.status;
    return this.get("/api/logs", { params: searchParams });
  }
}
const api = new ApiClient(API_BASE_URL);
function createAuthStore() {
  const { subscribe, set, update } = writable({
    user: null,
    token: null,
    isAuthenticated: false,
    isLoading: true
  });
  return {
    subscribe,
    login: (user, token) => {
      localStorage.setItem("token", token);
      api.setToken(token);
      set({ user, token, isAuthenticated: true, isLoading: false });
    },
    logout: () => {
      localStorage.removeItem("token");
      api.setToken(null);
      set({ user: null, token: null, isAuthenticated: false, isLoading: false });
    },
    setUser: (user) => {
      update((state) => ({
        ...state,
        user,
        isAuthenticated: !!user,
        isLoading: false
      }));
    },
    setLoading: (isLoading) => {
      update((state) => ({ ...state, isLoading }));
    },
    init: async () => {
      const token = localStorage.getItem("token");
      if (token) {
        api.setToken(token);
        try {
          const user = await api.getMe();
          set({ user, token, isAuthenticated: true, isLoading: false });
        } catch {
          localStorage.removeItem("token");
          api.setToken(null);
          set({ user: null, token: null, isAuthenticated: false, isLoading: false });
        }
      } else {
        set({ user: null, token: null, isAuthenticated: false, isLoading: false });
      }
    }
  };
}
const auth = createAuthStore();
export {
  auth as a
};
