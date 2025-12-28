<script lang="ts">
	import { wsStore } from '$lib/stores/websocket';
	import { Wifi, WifiOff } from 'lucide-svelte';
</script>

<div class="flex items-center gap-2 text-sm">
	{#if $wsStore.connected}
		<Wifi class="h-4 w-4 text-green-500" />
		<span class="text-green-600">Connected</span>
	{:else}
		<WifiOff class="h-4 w-4 text-red-500" />
		<span class="text-red-600">Disconnected</span>
	{/if}
</div>
