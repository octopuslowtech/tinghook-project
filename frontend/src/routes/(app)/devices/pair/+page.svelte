<script lang="ts">
	import { onMount } from 'svelte';
	import { api } from '$lib/api/client';
	import { auth } from '$lib/stores/auth';
	import QRCode from '$lib/components/QRCode.svelte';
	import { ArrowLeft, Smartphone, CheckCircle } from 'lucide-svelte';
	import type { User } from '$lib/stores/auth';

	interface PairingData {
		api_key: string;
		server_url: string;
		pairing_token: string;
	}

	let pairingData = $state<PairingData | null>(null);
	let loading = $state(true);
	let devicePaired = $state(false);
	let pollInterval: ReturnType<typeof setInterval> | null = null;

	onMount(() => {
		let user: User | null = null;
		const unsub = auth.subscribe((state) => {
			user = state.user;
		});
		unsub();

		api.getPairingToken().then((data) => {
			pairingData = {
				api_key: user?.api_key || '',
				server_url: import.meta.env.VITE_API_URL || 'https://api.tinghook.io',
				pairing_token: data.token
			};
			loading = false;

			pollForDevice();
		});

		return () => {
			if (pollInterval) {
				clearInterval(pollInterval);
			}
		};
	});

	function pollForDevice() {
		pollInterval = setInterval(async () => {
			try {
				const status = await api.checkPairingStatus(pairingData?.pairing_token || '');
				if (status.paired) {
					devicePaired = true;
					if (pollInterval) {
						clearInterval(pollInterval);
					}
				}
			} catch {
				// Ignore polling errors
			}
		}, 3000);
	}
</script>

<a href="/devices" class="mb-6 inline-flex items-center text-gray-600 hover:text-gray-900">
	<ArrowLeft class="mr-2 h-4 w-4" />
	Back to Devices
</a>

<div class="mx-auto max-w-lg">
	<div class="rounded-lg bg-white p-8 text-center shadow">
		<Smartphone class="mx-auto mb-4 h-16 w-16 text-blue-500" />
		<h1 class="mb-2 text-2xl font-bold">Pair Your Device</h1>
		<p class="mb-8 text-gray-600">
			Open the TingHook app on your Android device and scan this QR code.
		</p>

		{#if loading}
			<div class="mx-auto h-64 w-64 animate-pulse rounded-lg bg-gray-100"></div>
		{:else if devicePaired}
			<div class="mx-auto flex h-64 w-64 items-center justify-center rounded-lg bg-green-50">
				<CheckCircle class="h-24 w-24 text-green-500" />
			</div>
			<p class="mt-4 font-semibold text-green-600">Device paired successfully!</p>
		{:else if pairingData}
			<div class="inline-block rounded-lg border-2 border-gray-200 bg-white p-4">
				<QRCode data={JSON.stringify(pairingData)} size={256} />
			</div>
		{/if}

		<div class="mt-8 rounded-lg bg-gray-50 p-4 text-left">
			<h3 class="mb-2 font-semibold">Instructions:</h3>
			<ol class="list-inside list-decimal space-y-2 text-sm text-gray-600">
				<li>Download TingHook from Google Play Store</li>
				<li>Open the app and tap "Connect"</li>
				<li>Scan the QR code above</li>
				<li>Grant SMS and Notification permissions</li>
			</ol>
		</div>
	</div>
</div>
