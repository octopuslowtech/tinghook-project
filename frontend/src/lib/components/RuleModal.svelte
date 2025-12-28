<script lang="ts">
	import { onMount } from 'svelte';
	import { X } from 'lucide-svelte';
	import { api, type Rule, type RuleInput, type Device } from '$lib/api/client';

	interface Props {
		rule: Rule | null;
		onClose: () => void;
		onSave: (data: RuleInput) => Promise<void>;
	}

	let { rule, onClose, onSave }: Props = $props();

	let devices = $state<Device[]>([]);
	let loading = $state(false);
	let error = $state('');

	let triggerType = $state<'sms' | 'notification'>(rule?.trigger_type || 'sms');
	let deviceId = $state<string>(rule?.device_id || '');
	let senderFilter = $state(rule?.sender_filter || '');
	let contentFilter = $state(rule?.content_filter || '');
	let webhookUrl = $state(rule?.webhook_url || '');
	let method = $state(rule?.method || 'POST');
	let secretHeader = $state(rule?.secret_header || '');

	onMount(async () => {
		try {
			devices = await api.getDevices();
		} catch {
			devices = [];
		}
	});

	function validateWebhookUrl(url: string): boolean {
		try {
			const parsed = new URL(url);
			return parsed.protocol === 'http:' || parsed.protocol === 'https:';
		} catch {
			return false;
		}
	}

	async function handleSubmit(e: Event) {
		e.preventDefault();
		error = '';

		if (!webhookUrl.trim()) {
			error = 'Webhook URL is required';
			return;
		}

		if (!validateWebhookUrl(webhookUrl)) {
			error = 'Please enter a valid HTTP/HTTPS URL';
			return;
		}

		loading = true;
		try {
			await onSave({
				device_id: deviceId || null,
				trigger_type: triggerType,
				sender_filter: senderFilter,
				content_filter: contentFilter,
				webhook_url: webhookUrl,
				method,
				secret_header: secretHeader || undefined
			});
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to save rule';
		} finally {
			loading = false;
		}
	}
</script>

<div
	class="fixed inset-0 z-50 flex items-center justify-center bg-black/50"
	onclick={(e) => e.target === e.currentTarget && onClose()}
	onkeydown={(e) => e.key === 'Escape' && onClose()}
	role="dialog"
	aria-modal="true"
	tabindex="-1"
>
	<div class="w-full max-w-lg rounded-lg bg-white shadow-xl">
		<div class="flex items-center justify-between border-b px-6 py-4">
			<h2 class="text-lg font-semibold">{rule ? 'Edit Rule' : 'Create Rule'}</h2>
			<button onclick={onClose} class="rounded p-1 hover:bg-gray-100">
				<X class="h-5 w-5" />
			</button>
		</div>

		<form onsubmit={handleSubmit} class="p-6">
			{#if error}
				<div class="mb-4 rounded bg-red-50 p-3 text-sm text-red-600">{error}</div>
			{/if}

			<div class="space-y-4">
				<div>
					<label for="trigger_type" class="mb-1 block text-sm font-medium text-gray-700">
						Trigger Type
					</label>
					<select
						id="trigger_type"
						bind:value={triggerType}
						class="w-full rounded-lg border border-gray-300 px-3 py-2 focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500"
					>
						<option value="sms">SMS</option>
						<option value="notification">Notification</option>
					</select>
				</div>

				<div>
					<label for="device_id" class="mb-1 block text-sm font-medium text-gray-700">
						Device (Optional)
					</label>
					<select
						id="device_id"
						bind:value={deviceId}
						class="w-full rounded-lg border border-gray-300 px-3 py-2 focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500"
					>
						<option value="">All Devices</option>
						{#each devices as device}
							<option value={device.id}>{device.name}</option>
						{/each}
					</select>
				</div>

				<div>
					<label for="sender_filter" class="mb-1 block text-sm font-medium text-gray-700">
						Sender Filter
					</label>
					<input
						id="sender_filter"
						type="text"
						bind:value={senderFilter}
						placeholder="e.g., +1234567890 or regex pattern"
						class="w-full rounded-lg border border-gray-300 px-3 py-2 focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500"
					/>
					<p class="mt-1 text-xs text-gray-500">Leave empty to match all senders</p>
				</div>

				<div>
					<label for="content_filter" class="mb-1 block text-sm font-medium text-gray-700">
						Content Filter
					</label>
					<input
						id="content_filter"
						type="text"
						bind:value={contentFilter}
						placeholder="e.g., OTP or regex pattern"
						class="w-full rounded-lg border border-gray-300 px-3 py-2 focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500"
					/>
					<p class="mt-1 text-xs text-gray-500">Leave empty to match all messages</p>
				</div>

				<div>
					<label for="webhook_url" class="mb-1 block text-sm font-medium text-gray-700">
						Webhook URL <span class="text-red-500">*</span>
					</label>
					<input
						id="webhook_url"
						type="url"
						bind:value={webhookUrl}
						placeholder="https://example.com/webhook"
						required
						class="w-full rounded-lg border border-gray-300 px-3 py-2 focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500"
					/>
				</div>

				<div>
					<label for="method" class="mb-1 block text-sm font-medium text-gray-700">
						HTTP Method
					</label>
					<select
						id="method"
						bind:value={method}
						class="w-full rounded-lg border border-gray-300 px-3 py-2 focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500"
					>
						<option value="GET">GET</option>
						<option value="POST">POST</option>
						<option value="PUT">PUT</option>
					</select>
				</div>

				<div>
					<label for="secret_header" class="mb-1 block text-sm font-medium text-gray-700">
						Secret Header (Optional)
					</label>
					<input
						id="secret_header"
						type="text"
						bind:value={secretHeader}
						placeholder="X-Webhook-Secret: your-secret"
						class="w-full rounded-lg border border-gray-300 px-3 py-2 focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500"
					/>
					<p class="mt-1 text-xs text-gray-500">Custom header for webhook authentication</p>
				</div>
			</div>

			<div class="mt-6 flex justify-end gap-3">
				<button
					type="button"
					onclick={onClose}
					class="rounded-lg border border-gray-300 px-4 py-2 hover:bg-gray-50"
				>
					Cancel
				</button>
				<button
					type="submit"
					disabled={loading}
					class="rounded-lg bg-blue-600 px-4 py-2 text-white hover:bg-blue-700 disabled:opacity-50"
				>
					{loading ? 'Saving...' : rule ? 'Update Rule' : 'Create Rule'}
				</button>
			</div>
		</form>
	</div>
</div>
