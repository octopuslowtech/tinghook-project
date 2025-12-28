<script lang="ts">
	import { onMount } from 'svelte';
	import { api } from '$lib/api/client';
	import { Smartphone, MessageSquare, CheckCircle, XCircle, Activity } from 'lucide-svelte';

	interface Stats {
		total_devices: number;
		online_devices: number;
		total_inbound: number;
		total_outbound: number;
		total_sent: number;
		total_failed: number;
	}

	interface RecentActivity {
		id: string;
		type: 'inbound' | 'outbound';
		device_name: string;
		phone_number: string;
		status: string;
		created_at: string;
	}

	let stats = $state<Stats | null>(null);
	let recentActivity = $state<RecentActivity[]>([]);
	let loading = $state(true);

	onMount(async () => {
		try {
			const [logsStats, devices] = await Promise.all([
				api.getDashboardStats(),
				api.getDevices()
			]);

			const onlineDevices = devices.filter((d: { is_online: boolean }) => d.is_online).length;

			stats = {
				total_devices: devices.length,
				online_devices: onlineDevices,
				total_inbound: logsStats.total_inbound ?? 0,
				total_outbound: logsStats.total_outbound ?? 0,
				total_sent: logsStats.total_sent ?? 0,
				total_failed: logsStats.total_failed ?? 0
			};

			recentActivity = logsStats.recent_activity ?? [];
		} catch (error) {
			console.error('Failed to load dashboard stats:', error);
		} finally {
			loading = false;
		}
	});

	function formatNumber(num: number): string {
		return num.toLocaleString();
	}

	function formatDate(dateStr: string): string {
		return new Date(dateStr).toLocaleString();
	}
</script>

<h1 class="mb-6 text-2xl font-bold">Dashboard</h1>

<!-- Stats Cards Grid -->
<div class="mb-8 grid grid-cols-1 gap-6 md:grid-cols-2 lg:grid-cols-4">
	<!-- Device Status Card -->
	<div class="rounded-lg bg-white p-6 shadow">
		{#if loading}
			<div class="animate-pulse">
				<div class="mb-2 h-4 w-24 rounded bg-gray-200"></div>
				<div class="h-8 w-16 rounded bg-gray-200"></div>
			</div>
		{:else}
			<div class="flex items-center justify-between">
				<div>
					<p class="text-sm text-gray-500">Devices Online</p>
					<p class="text-3xl font-bold">
						{formatNumber(stats?.online_devices ?? 0)} / {formatNumber(stats?.total_devices ?? 0)}
					</p>
				</div>
				<Smartphone class="h-12 w-12 text-blue-500" />
			</div>
		{/if}
	</div>

	<!-- Messages Received Card -->
	<div class="rounded-lg bg-white p-6 shadow">
		{#if loading}
			<div class="animate-pulse">
				<div class="mb-2 h-4 w-24 rounded bg-gray-200"></div>
				<div class="h-8 w-16 rounded bg-gray-200"></div>
			</div>
		{:else}
			<div class="flex items-center justify-between">
				<div>
					<p class="text-sm text-gray-500">Messages Received</p>
					<p class="text-3xl font-bold">{formatNumber(stats?.total_inbound ?? 0)}</p>
				</div>
				<MessageSquare class="h-12 w-12 text-green-500" />
			</div>
		{/if}
	</div>

	<!-- Messages Sent Card -->
	<div class="rounded-lg bg-white p-6 shadow">
		{#if loading}
			<div class="animate-pulse">
				<div class="mb-2 h-4 w-24 rounded bg-gray-200"></div>
				<div class="h-8 w-16 rounded bg-gray-200"></div>
			</div>
		{:else}
			<div class="flex items-center justify-between">
				<div>
					<p class="text-sm text-gray-500">Messages Sent</p>
					<p class="text-3xl font-bold">{formatNumber(stats?.total_sent ?? 0)}</p>
				</div>
				<CheckCircle class="h-12 w-12 text-emerald-500" />
			</div>
		{/if}
	</div>

	<!-- Failed Messages Card -->
	<div class="rounded-lg bg-white p-6 shadow">
		{#if loading}
			<div class="animate-pulse">
				<div class="mb-2 h-4 w-24 rounded bg-gray-200"></div>
				<div class="h-8 w-16 rounded bg-gray-200"></div>
			</div>
		{:else}
			<div class="flex items-center justify-between">
				<div>
					<p class="text-sm text-gray-500">Failed Messages</p>
					<p class="text-3xl font-bold">{formatNumber(stats?.total_failed ?? 0)}</p>
				</div>
				<XCircle class="h-12 w-12 text-red-500" />
			</div>
		{/if}
	</div>
</div>

<!-- Recent Activity Section -->
<div class="rounded-lg bg-white shadow">
	<div class="border-b p-6">
		<h2 class="flex items-center gap-2 text-lg font-semibold">
			<Activity class="h-5 w-5" />
			Recent Activity
		</h2>
	</div>
	<div class="p-6">
		{#if loading}
			<div class="animate-pulse space-y-4">
				{#each Array(5) as _}
					<div class="flex items-center gap-4">
						<div class="h-10 w-10 rounded-full bg-gray-200"></div>
						<div class="flex-1">
							<div class="mb-2 h-4 w-48 rounded bg-gray-200"></div>
							<div class="h-3 w-32 rounded bg-gray-200"></div>
						</div>
					</div>
				{/each}
			</div>
		{:else if recentActivity.length === 0}
			<p class="py-8 text-center text-gray-500">No recent activity</p>
		{:else}
			<div class="space-y-4">
				{#each recentActivity as activity}
					<div class="flex items-center gap-4 rounded-lg border p-4">
						<div
							class="flex h-10 w-10 items-center justify-center rounded-full {activity.type ===
							'inbound'
								? 'bg-green-100'
								: 'bg-blue-100'}"
						>
							<MessageSquare
								class="h-5 w-5 {activity.type === 'inbound' ? 'text-green-600' : 'text-blue-600'}"
							/>
						</div>
						<div class="flex-1">
							<p class="font-medium">
								{activity.type === 'inbound' ? 'Received from' : 'Sent to'}
								{activity.phone_number}
							</p>
							<p class="text-sm text-gray-500">
								{activity.device_name} • {formatDate(activity.created_at)}
							</p>
						</div>
						<span
							class="rounded-full px-2 py-1 text-xs {activity.status === 'success'
								? 'bg-green-100 text-green-800'
								: activity.status === 'failed'
									? 'bg-red-100 text-red-800'
									: 'bg-gray-100 text-gray-800'}"
						>
							{activity.status}
						</span>
					</div>
				{/each}
			</div>
		{/if}
	</div>
</div>
