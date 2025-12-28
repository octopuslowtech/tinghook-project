<script lang="ts">
	import { page } from '$app/stores';
	import { LayoutDashboard, Smartphone, Settings, FileText, Bell, LogOut } from 'lucide-svelte';
	import { auth } from '$lib/stores/auth';
	import { goto } from '$app/navigation';

	interface Props {
		open: boolean;
	}
	let { open = $bindable() }: Props = $props();

	const navItems = [
		{ href: '/dashboard', label: 'Dashboard', icon: LayoutDashboard },
		{ href: '/devices', label: 'Devices', icon: Smartphone },
		{ href: '/rules', label: 'Forwarding Rules', icon: Bell },
		{ href: '/logs', label: 'Message Logs', icon: FileText },
		{ href: '/settings', label: 'Settings', icon: Settings }
	];

	function handleLogout() {
		auth.logout();
		goto('/login');
	}
</script>

<aside
	class="fixed inset-y-0 left-0 z-50 w-64 transform bg-gray-900 text-white transition-transform lg:translate-x-0 {open
		? 'translate-x-0'
		: '-translate-x-full'}"
>
	<div class="p-6">
		<h1 class="text-2xl font-bold">TingHook</h1>
	</div>

	<nav class="mt-6">
		{#each navItems as item}
			<a
				href={item.href}
				class="flex items-center px-6 py-3 transition-colors hover:bg-gray-800 {$page.url.pathname ===
				item.href
					? 'border-l-4 border-blue-500 bg-gray-800'
					: ''}"
				onclick={() => (open = false)}
			>
				<item.icon class="mr-3 h-5 w-5" />
				{item.label}
			</a>
		{/each}
	</nav>

	<div class="absolute bottom-0 w-full p-6">
		<button class="flex items-center text-gray-400 hover:text-white" onclick={handleLogout}>
			<LogOut class="mr-3 h-5 w-5" />
			Logout
		</button>
	</div>
</aside>
