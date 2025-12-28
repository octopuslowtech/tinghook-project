<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import Sidebar from '$lib/components/Sidebar.svelte';
	import Header from '$lib/components/Header.svelte';
	import ConnectionStatus from '$lib/components/ConnectionStatus.svelte';
	import { auth } from '$lib/stores/auth';
	import { wsStore } from '$lib/stores/websocket';
	import { goto } from '$app/navigation';

	let { children } = $props();
	let sidebarOpen = $state(false);

	$effect(() => {
		if (!$auth.isLoading && !$auth.isAuthenticated) {
			goto('/login');
		}
	});

	onMount(() => {
		wsStore.connect();
	});

	onDestroy(() => {
		wsStore.disconnect();
	});
</script>

<div class="min-h-screen bg-gray-50">
	{#if sidebarOpen}
		<div
			class="fixed inset-0 z-40 bg-black/50 lg:hidden"
			onclick={() => (sidebarOpen = false)}
			onkeydown={(e) => e.key === 'Escape' && (sidebarOpen = false)}
			role="button"
			tabindex="-1"
		></div>
	{/if}

	<Sidebar bind:open={sidebarOpen} />

	<div class="lg:pl-64">
		<Header onMenuClick={() => (sidebarOpen = true)} />
		<main class="p-6">
			{@render children()}
		</main>
	</div>
</div>
