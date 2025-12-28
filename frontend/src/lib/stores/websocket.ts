import { writable } from 'svelte/store';

interface WSMessage {
	type: string;
	data: unknown;
}

export interface DeviceUpdate {
	device_id: string;
	status: 'online' | 'offline';
	battery?: number;
}

export interface LogEvent {
	id: number;
	direction: string;
	sender: string;
	content: string;
	status: string;
}

interface WebSocketState {
	connected: boolean;
	devices: Map<string, DeviceUpdate>;
	recentLogs: LogEvent[];
}

function createWebSocketStore() {
	const { subscribe, set, update } = writable<WebSocketState>({
		connected: false,
		devices: new Map<string, DeviceUpdate>(),
		recentLogs: []
	});

	let ws: WebSocket | null = null;
	let reconnectTimeout: ReturnType<typeof setTimeout>;

	function connect() {
		const token = localStorage.getItem('token');
		if (!token) return;

		const wsUrl = `ws://localhost:8080/ws/dashboard?token=${token}`;
		ws = new WebSocket(wsUrl);

		ws.onopen = () => {
			update((state) => ({ ...state, connected: true }));
		};

		ws.onclose = () => {
			update((state) => ({ ...state, connected: false }));
			reconnectTimeout = setTimeout(connect, 5000);
		};

		ws.onerror = () => {
			ws?.close();
		};

		ws.onmessage = (event) => {
			try {
				const msg: WSMessage = JSON.parse(event.data);
				handleMessage(msg);
			} catch {
				// Ignore invalid messages
			}
		};
	}

	function handleMessage(msg: WSMessage) {
		update((state) => {
			switch (msg.type) {
				case 'DEVICE_STATUS': {
					const device = msg.data as DeviceUpdate;
					state.devices.set(device.device_id, device);
					return { ...state, devices: new Map(state.devices) };
				}

				case 'NEW_LOG': {
					const log = msg.data as LogEvent;
					return {
						...state,
						recentLogs: [log, ...state.recentLogs].slice(0, 10)
					};
				}

				default:
					return state;
			}
		});
	}

	function disconnect() {
		clearTimeout(reconnectTimeout);
		if (ws) {
			ws.close();
			ws = null;
		}
	}

	return {
		subscribe,
		connect,
		disconnect
	};
}

export const wsStore = createWebSocketStore();
