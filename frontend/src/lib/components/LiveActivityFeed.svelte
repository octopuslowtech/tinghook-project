<script lang="ts">
	import { wsStore } from '$lib/stores/websocket';
	import { ArrowDownLeft, ArrowUpRight } from 'lucide-svelte';
</script>

<div class="space-y-2">
	{#each $wsStore.recentLogs as log (log.id)}
		<div class="flex animate-in slide-in-from-top items-center gap-3 rounded-lg bg-gray-50 p-3">
			{#if log.direction === 'inbound'}
				<ArrowDownLeft class="h-4 w-4 text-green-500" />
			{:else}
				<ArrowUpRight class="h-4 w-4 text-blue-500" />
			{/if}
			<span class="font-mono text-sm">{log.sender}</span>
			<span class="flex-1 truncate text-sm text-gray-600">{log.content}</span>
			<span class="rounded-full bg-gray-200 px-2 py-0.5 text-xs">{log.status}</span>
		</div>
	{/each}

	{#if $wsStore.recentLogs.length === 0}
		<p class="py-4 text-center text-gray-500">No recent activity</p>
	{/if}
</div>
