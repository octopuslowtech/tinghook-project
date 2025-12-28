<script lang="ts">
	import { Menu, User } from 'lucide-svelte';
	import { auth } from '$lib/stores/auth';
	import ConnectionStatus from './ConnectionStatus.svelte';

	interface Props {
		onMenuClick: () => void;
	}
	let { onMenuClick }: Props = $props();
</script>

<header class="flex items-center justify-between border-b bg-white px-6 py-4">
	<button class="lg:hidden" onclick={onMenuClick}>
		<Menu class="h-6 w-6" />
	</button>

	<div class="flex-1"></div>

	<div class="flex items-center gap-4">
		<ConnectionStatus />
		<span class="text-sm text-gray-600">{$auth.user?.email}</span>
		<div class="flex h-8 w-8 items-center justify-center rounded-full bg-gray-200">
			<User class="h-4 w-4" />
		</div>
	</div>
</header>
