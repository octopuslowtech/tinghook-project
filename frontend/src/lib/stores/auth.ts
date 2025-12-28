import { writable } from 'svelte/store';
import { api } from '$lib/api/client';

export interface User {
	id: string;
	email: string;
	api_key: string;
	subscription_plan: string;
}

interface AuthState {
	user: User | null;
	token: string | null;
	isAuthenticated: boolean;
	isLoading: boolean;
}

function createAuthStore() {
	const { subscribe, set, update } = writable<AuthState>({
		user: null,
		token: null,
		isAuthenticated: false,
		isLoading: true
	});

	return {
		subscribe,
		login: (user: User, token: string) => {
			localStorage.setItem('token', token);
			api.setToken(token);
			set({ user, token, isAuthenticated: true, isLoading: false });
		},
		logout: () => {
			localStorage.removeItem('token');
			api.setToken(null);
			set({ user: null, token: null, isAuthenticated: false, isLoading: false });
		},
		setUser: (user: User | null) => {
			update((state) => ({
				...state,
				user,
				isAuthenticated: !!user,
				isLoading: false
			}));
		},
		setLoading: (isLoading: boolean) => {
			update((state) => ({ ...state, isLoading }));
		},
		init: async () => {
			const token = localStorage.getItem('token');
			if (token) {
				api.setToken(token);
				try {
					const user = await api.getMe();
					set({ user, token, isAuthenticated: true, isLoading: false });
				} catch {
					localStorage.removeItem('token');
					api.setToken(null);
					set({ user: null, token: null, isAuthenticated: false, isLoading: false });
				}
			} else {
				set({ user: null, token: null, isAuthenticated: false, isLoading: false });
			}
		}
	};
}

export const auth = createAuthStore();
