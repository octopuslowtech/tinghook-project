<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import { api, type Device } from '$lib/api/client';
	import { Plus, Smartphone, Battery, Trash2, Edit2, X, Check } from 'lucide-svelte';

	let devices = $state<Device[]>([]);
	let loading = $state(true);
	let error = $state<string | null>(null);
	let editingId = $state<string | null>(null);
	let editingName = $state('');
	let deleteConfirmId = $state<string | null>(null);
	let refreshInterval: ReturnType<typeof setInterval>;

	async function fetchDevices() {
		try {
			devices = await api.getDevices();
			error = null;
		} catch (e) {
			error = e instanceof Error ? e.message : 'Failed to load devices';
		} finally {
			loading = false;
		}
	}

	onMount(() => {
		fetchDevices();
		refreshInterval = setInterval(fetchDevices, 10000);
	});

	onDestroy(() => {
		if (refreshInterval) clearInterval(refreshInterval);
	});

	function startEdit(device: Device) {
		editingId = device.id;
		editingName = device.name;
	}

	async function saveEdit() {
		if (!editingId || !editingName.trim()) return;
		try {
			await api.updateDevice(editingId, editingName.trim());
			await fetchDevices();
		} catch (e) {
			error = e instanceof Error ? e.message : 'Failed to update device';
		}
		editingId = null;
		editingName = '';
	}

	function cancelEdit() {
		editingId = null;
		editingName = '';
	}

	async function confirmDelete() {
		if (!deleteConfirmId) return;
		try {
			await api.deleteDevice(deleteConfirmId);
			await fetchDevices();
		} catch (e) {
			error = e instanceof Error ? e.message : 'Failed to delete device';
		}
		deleteConfirmId = null;
	}

	function formatLastSeen(lastSeen: string | null): string {
		if (!lastSeen) return 'Never';
		const date = new Date(lastSeen);
		const now = new Date();
		const diff = now.getTime() - date.getTime();
		const minutes = Math.floor(diff / 60000);
		if (minutes < 1) return 'Just now';
		if (minutes < 60) return `${minutes}m ago`;
		const hours = Math.floor(minutes / 60);
		if (hours < 24) return `${hours}h ago`;
		return date.toLocaleDateString();
	}
</script>

<div class="flex flex-col gap-6 sm:flex-row sm:items-center sm:justify-between mb-6">
	<h1 class="text-2xl font-bold">Devices</h1>
	<a
		href="/devices/pair"
		class="bg-blue-600 text-white px-4 py-2 rounded-lg flex items-center justify-center gap-2 hover:bg-blue-700 transition-colors"
	>
		<Plus class="w-4 h-4" />
		Add Device
	</a>
</div>

{#if error}
	<div class="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg mb-4">
		{error}
	</div>
{/if}

{#if loading}
	<div class="space-y-4">
		{#each [1, 2, 3] as _}
			<div class="bg-white rounded-lg shadow p-6 animate-pulse">
				<div class="flex items-center gap-4">
					<div class="w-10 h-10 bg-gray-200 rounded"></div>
					<div class="flex-1">
						<div class="h-4 bg-gray-200 rounded w-32 mb-2"></div>
						<div class="h-3 bg-gray-200 rounded w-48"></div>
					</div>
				</div>
			</div>
		{/each}
	</div>
{:else if devices.length === 0}
	<div class="bg-white rounded-lg shadow p-12 text-center">
		<Smartphone class="w-16 h-16 mx-auto text-gray-400 mb-4" />
		<h3 class="text-lg font-semibold mb-2">No devices connected</h3>
		<p class="text-gray-500 mb-4">Connect your first Android device to get started.</p>
		<a href="/devices/pair" class="text-blue-600 hover:underline">Pair a device →</a>
	</div>
{:else}
	<div class="grid gap-4">
		{#each devices as device (device.id)}
			<div
				class="bg-white rounded-lg shadow p-4 sm:p-6 flex flex-col sm:flex-row sm:items-center gap-4"
			>
				<div class="flex items-center gap-4 flex-1 min-w-0">
					<div class="relative flex-shrink-0">
						<Smartphone class="w-10 h-10 text-gray-600" />
						<div
							class="absolute -bottom-1 -right-1 w-4 h-4 rounded-full border-2 border-white {device.status ===
							'online'
								? 'bg-green-500'
								: 'bg-gray-400'}"
						></div>
					</div>

					<div class="flex-1 min-w-0">
						{#if editingId === device.id}
							<div class="flex items-center gap-2">
								<input
									type="text"
									bind:value={editingName}
									class="border rounded px-2 py-1 text-sm w-full max-w-xs"
									onkeydown={(e) => e.key === 'Enter' && saveEdit()}
								/>
								<button
									onclick={saveEdit}
									class="p-1 text-green-600 hover:bg-green-50 rounded"
								>
									<Check class="w-4 h-4" />
								</button>
								<button
									onclick={cancelEdit}
									class="p-1 text-gray-600 hover:bg-gray-100 rounded"
								>
									<X class="w-4 h-4" />
								</button>
							</div>
						{:else}
							<h3 class="font-semibold truncate">{device.name}</h3>
						{/if}
						<p class="text-sm text-gray-500 truncate">{device.device_uid}</p>
						<p class="text-xs text-gray-400">Last seen: {formatLastSeen(device.last_seen_at)}</p>
					</div>
				</div>

				<div class="flex items-center gap-4 sm:gap-6 flex-wrap">
					<div class="flex items-center gap-2 text-sm text-gray-600">
						<Battery class="w-4 h-4" />
						<span>{device.battery_level}%</span>
					</div>

					<span
						class="px-3 py-1 rounded-full text-sm {device.status === 'online'
							? 'bg-green-100 text-green-700'
							: 'bg-gray-100 text-gray-600'}"
					>
						{device.status}
					</span>

					<div class="flex gap-2">
						<button
							onclick={() => startEdit(device)}
							class="p-2 hover:bg-gray-100 rounded"
							title="Edit device name"
						>
							<Edit2 class="w-4 h-4" />
						</button>
						<button
							onclick={() => (deleteConfirmId = device.id)}
							class="p-2 hover:bg-gray-100 rounded text-red-500"
							title="Delete device"
						>
							<Trash2 class="w-4 h-4" />
						</button>
					</div>
				</div>
			</div>
		{/each}
	</div>
{/if}

{#if deleteConfirmId}
	<div class="fixed inset-0 bg-black/50 flex items-center justify-center p-4 z-50">
		<div class="bg-white rounded-lg shadow-xl p-6 max-w-sm w-full">
			<h3 class="text-lg font-semibold mb-2">Delete Device</h3>
			<p class="text-gray-600 mb-4">
				Are you sure you want to delete this device? This action cannot be undone.
			</p>
			<div class="flex justify-end gap-3">
				<button
					onclick={() => (deleteConfirmId = null)}
					class="px-4 py-2 text-gray-600 hover:bg-gray-100 rounded-lg"
				>
					Cancel
				</button>
				<button
					onclick={confirmDelete}
					class="px-4 py-2 bg-red-600 text-white rounded-lg hover:bg-red-700"
				>
					Delete
				</button>
			</div>
		</div>
	</div>
{/if}
