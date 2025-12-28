<script lang="ts">
	import { onMount } from 'svelte';
	import { api } from '$lib/api/client';
	import {
		MessageSquare,
		ArrowDownLeft,
		ArrowUpRight,
		ChevronLeft,
		ChevronRight,
		Filter
	} from 'lucide-svelte';

	interface Log {
		id: number;
		device_id: string;
		direction: 'inbound' | 'outbound';
		sender: string;
		receiver: string;
		content: string;
		status: 'pending' | 'sent' | 'delivered' | 'failed';
		created_at: string;
	}

	interface Pagination {
		total: number;
		page: number;
		limit: number;
		total_pages: number;
	}

	let logs = $state<Log[]>([]);
	let pagination = $state<Pagination>({ total: 0, page: 1, limit: 20, total_pages: 0 });
	let loading = $state(true);

	// Filters
	let direction = $state<string>('');
	let status = $state<string>('');
	let showFilters = $state(false);

	onMount(() => fetchLogs());

	async function fetchLogs() {
		loading = true;
		try {
			const result = await api.getLogs({
				page: pagination.page,
				limit: pagination.limit,
				direction: direction || undefined,
				status: status || undefined
			});
			logs = result.data;
			pagination = { ...pagination, total: result.total, total_pages: result.total_pages };
		} catch (error) {
			console.error('Failed to fetch logs:', error);
		}
		loading = false;
	}

	function nextPage() {
		if (pagination.page < pagination.total_pages) {
			pagination.page++;
			fetchLogs();
		}
	}

	function prevPage() {
		if (pagination.page > 1) {
			pagination.page--;
			fetchLogs();
		}
	}

	function applyFilter() {
		pagination.page = 1;
		fetchLogs();
	}
</script>

<div class="flex items-center justify-between mb-6">
	<h1 class="text-2xl font-bold">Message Logs</h1>
	<button
		onclick={() => (showFilters = !showFilters)}
		class="flex items-center gap-2 px-4 py-2 border rounded-lg hover:bg-gray-50"
	>
		<Filter class="w-4 h-4" />
		Filters
	</button>
</div>

{#if showFilters}
	<div class="flex gap-4 p-4 mb-6 bg-white rounded-lg shadow">
		<select
			bind:value={direction}
			onchange={applyFilter}
			class="px-3 py-2 border rounded-lg"
		>
			<option value="">All Directions</option>
			<option value="inbound">Inbound</option>
			<option value="outbound">Outbound</option>
		</select>
		<select bind:value={status} onchange={applyFilter} class="px-3 py-2 border rounded-lg">
			<option value="">All Status</option>
			<option value="pending">Pending</option>
			<option value="sent">Sent</option>
			<option value="delivered">Delivered</option>
			<option value="failed">Failed</option>
		</select>
	</div>
{/if}

<div class="overflow-hidden bg-white rounded-lg shadow">
	{#if loading}
		<div class="flex items-center justify-center py-12">
			<div class="w-8 h-8 border-4 border-blue-500 rounded-full border-t-transparent animate-spin"></div>
		</div>
	{:else if logs.length === 0}
		<div class="py-12 text-center text-gray-500">
			<MessageSquare class="w-12 h-12 mx-auto mb-4 opacity-50" />
			<p>No message logs found</p>
		</div>
	{:else}
		<table class="w-full">
			<thead class="bg-gray-50">
				<tr>
					<th class="px-6 py-3 text-xs font-medium tracking-wider text-left text-gray-500 uppercase">
						Direction
					</th>
					<th class="px-6 py-3 text-xs font-medium tracking-wider text-left text-gray-500 uppercase">
						From/To
					</th>
					<th class="px-6 py-3 text-xs font-medium tracking-wider text-left text-gray-500 uppercase">
						Content
					</th>
					<th class="px-6 py-3 text-xs font-medium tracking-wider text-left text-gray-500 uppercase">
						Status
					</th>
					<th class="px-6 py-3 text-xs font-medium tracking-wider text-left text-gray-500 uppercase">
						Time
					</th>
				</tr>
			</thead>
			<tbody class="divide-y">
				{#each logs as log}
					<tr class="hover:bg-gray-50">
						<td class="px-6 py-4">
							{#if log.direction === 'inbound'}
								<ArrowDownLeft class="w-5 h-5 text-green-500" />
							{:else}
								<ArrowUpRight class="w-5 h-5 text-blue-500" />
							{/if}
						</td>
						<td class="px-6 py-4 font-mono text-sm">
							{log.direction === 'inbound' ? log.sender : log.receiver}
						</td>
						<td class="max-w-md px-6 py-4 text-sm text-gray-600 truncate">
							{log.content}
						</td>
						<td class="px-6 py-4">
							<span
								class="px-2 py-1 text-xs rounded-full {log.status === 'sent'
									? 'bg-green-100 text-green-700'
									: log.status === 'failed'
										? 'bg-red-100 text-red-700'
										: log.status === 'pending'
											? 'bg-yellow-100 text-yellow-700'
											: 'bg-gray-100 text-gray-700'}"
							>
								{log.status}
							</span>
						</td>
						<td class="px-6 py-4 text-sm text-gray-500">
							{new Date(log.created_at).toLocaleString()}
						</td>
					</tr>
				{/each}
			</tbody>
		</table>

		<div class="flex items-center justify-between px-6 py-4 border-t">
			<span class="text-sm text-gray-600">
				Showing {(pagination.page - 1) * pagination.limit + 1} to {Math.min(
					pagination.page * pagination.limit,
					pagination.total
				)} of {pagination.total}
			</span>
			<div class="flex items-center gap-2">
				<button
					onclick={prevPage}
					disabled={pagination.page === 1}
					class="p-2 border rounded disabled:opacity-50"
				>
					<ChevronLeft class="w-4 h-4" />
				</button>
				<span class="px-4 py-2">{pagination.page} / {pagination.total_pages}</span>
				<button
					onclick={nextPage}
					disabled={pagination.page === pagination.total_pages}
					class="p-2 border rounded disabled:opacity-50"
				>
					<ChevronRight class="w-4 h-4" />
				</button>
			</div>
		</div>
	{/if}
</div>
