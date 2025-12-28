import type { User } from '$lib/stores/auth';

const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8000';

interface RequestOptions extends RequestInit {
	params?: Record<string, string>;
}

interface AuthResponse {
	user: User;
	token: string;
}

interface ApiError {
	message: string;
}

export interface Device {
	id: string;
	name: string;
	device_uid: string;
	status: 'online' | 'offline';
	battery_level: number;
	app_version: string;
	last_seen_at: string | null;
}

export interface Rule {
	id: number;
	device_id: string | null;
	trigger_type: 'sms' | 'notification';
	sender_filter: string;
	content_filter: string;
	webhook_url: string;
	method: string;
	secret_header: string | null;
	is_active: boolean;
	created_at: string;
}

export interface RuleInput {
	device_id?: string | null;
	trigger_type: 'sms' | 'notification';
	sender_filter?: string;
	content_filter?: string;
	webhook_url: string;
	method: string;
	secret_header?: string;
	is_active?: boolean;
}

class ApiClient {
	private baseUrl: string;
	private token: string | null = null;

	constructor(baseUrl: string) {
		this.baseUrl = baseUrl;
	}

	setToken(token: string | null) {
		this.token = token;
	}

	private async request<T>(endpoint: string, options: RequestOptions = {}): Promise<T> {
		const { params, ...init } = options;

		let url = `${this.baseUrl}${endpoint}`;
		if (params) {
			const searchParams = new URLSearchParams(params);
			url += `?${searchParams.toString()}`;
		}

		const headers: Record<string, string> = {
			'Content-Type': 'application/json',
			...((init.headers as Record<string, string>) || {})
		};

		if (this.token) {
			headers['Authorization'] = `Bearer ${this.token}`;
		}

		const response = await fetch(url, {
			...init,
			headers
		});

		if (!response.ok) {
			const error: ApiError = await response.json().catch(() => ({ message: 'Request failed' }));
			throw new Error(error.message || `API Error: ${response.status} ${response.statusText}`);
		}

		return response.json();
	}

	async get<T>(endpoint: string, options?: RequestOptions): Promise<T> {
		return this.request<T>(endpoint, { ...options, method: 'GET' });
	}

	async post<T>(endpoint: string, data?: unknown, options?: RequestOptions): Promise<T> {
		return this.request<T>(endpoint, {
			...options,
			method: 'POST',
			body: data ? JSON.stringify(data) : undefined
		});
	}

	async put<T>(endpoint: string, data?: unknown, options?: RequestOptions): Promise<T> {
		return this.request<T>(endpoint, {
			...options,
			method: 'PUT',
			body: data ? JSON.stringify(data) : undefined
		});
	}

	async delete<T>(endpoint: string, options?: RequestOptions): Promise<T> {
		return this.request<T>(endpoint, { ...options, method: 'DELETE' });
	}

	login(email: string, password: string): Promise<AuthResponse> {
		return this.post<AuthResponse>('/api/auth/login', { email, password });
	}

	register(email: string, password: string): Promise<AuthResponse> {
		return this.post<AuthResponse>('/api/auth/register', { email, password });
	}

	getMe(): Promise<User> {
		return this.get<User>('/api/auth/me');
	}

	getDevices(): Promise<Device[]> {
		return this.get<Device[]>('/api/devices');
	}

	deleteDevice(id: string): Promise<void> {
		return this.delete<void>(`/api/devices/${id}`);
	}

	updateDevice(id: string, name: string): Promise<Device> {
		return this.put<Device>(`/api/devices/${id}`, { name });
	}

	getRules(): Promise<Rule[]> {
		return this.get<Rule[]>('/api/rules');
	}

	createRule(data: RuleInput): Promise<Rule> {
		return this.post<Rule>('/api/rules', data);
	}

	updateRule(id: number, data: Partial<RuleInput>): Promise<Rule> {
		return this.put<Rule>(`/api/rules/${id}`, data);
	}

	deleteRule(id: number): Promise<void> {
		return this.delete<void>(`/api/rules/${id}`);
	}

	testRule(id: number): Promise<{ success: boolean; message: string }> {
		return this.post<{ success: boolean; message: string }>(`/api/rules/${id}/test`);
	}

	getPairingToken(): Promise<{ token: string }> {
		return this.post<{ token: string }>('/api/devices/pairing-token');
	}

	checkPairingStatus(token: string): Promise<{ paired: boolean }> {
		return this.get<{ paired: boolean }>(`/api/devices/pairing-status/${token}`);
	}

	getDashboardStats(): Promise<{
		total_inbound: number;
		total_outbound: number;
		total_sent: number;
		total_failed: number;
		recent_activity: Array<{
			id: string;
			type: 'inbound' | 'outbound';
			device_name: string;
			phone_number: string;
			status: string;
			created_at: string;
		}>;
	}> {
		return this.get('/api/logs/stats');
	}

	getLogs(params: {
		page?: number;
		limit?: number;
		direction?: string;
		status?: string;
	}): Promise<LogsResponse> {
		const searchParams: Record<string, string> = {};
		if (params.page) searchParams.page = String(params.page);
		if (params.limit) searchParams.limit = String(params.limit);
		if (params.direction) searchParams.direction = params.direction;
		if (params.status) searchParams.status = params.status;
		return this.get('/api/logs', { params: searchParams });
	}
}

export interface Log {
	id: number;
	device_id: string;
	direction: 'inbound' | 'outbound';
	sender: string;
	receiver: string;
	content: string;
	status: 'pending' | 'sent' | 'delivered' | 'failed';
	created_at: string;
}

interface LogsResponse {
	data: Log[];
	total: number;
	total_pages: number;
}

export const api = new ApiClient(API_BASE_URL);
