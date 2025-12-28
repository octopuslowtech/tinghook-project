<script lang="ts">
	import { onMount } from 'svelte';
	import { api, type Rule, type RuleInput } from '$lib/api/client';
	import { Plus, Bell, MessageSquare, Edit2, Trash2, ToggleLeft, ToggleRight } from 'lucide-svelte';
	import RuleModal from '$lib/components/RuleModal.svelte';

	let rules = $state<Rule[]>([]);
	let loading = $state(true);
	let showModal = $state(false);
	let editingRule = $state<Rule | null>(null);

	onMount(async () => {
		try {
			rules = await api.getRules();
		} catch {
			rules = [];
		}
		loading = false;
	});

	async function toggleRule(rule: Rule) {
		try {
			await api.updateRule(rule.id, { is_active: !rule.is_active });
			rule.is_active = !rule.is_active;
		} catch (err) {
			console.error('Failed to toggle rule:', err);
		}
	}

	async function deleteRule(id: number) {
		if (confirm('Delete this rule?')) {
			try {
				await api.deleteRule(id);
				rules = rules.filter((r) => r.id !== id);
			} catch (err) {
				console.error('Failed to delete rule:', err);
			}
		}
	}

	async function handleSave(data: RuleInput) {
		if (editingRule) {
			await api.updateRule(editingRule.id, data);
		} else {
			await api.createRule(data);
		}
		rules = await api.getRules();
		showModal = false;
	}

	function truncateUrl(url: string, maxLength = 40): string {
		if (url.length <= maxLength) return url;
		return url.substring(0, maxLength) + '...';
	}
</script>

<div class="mb-6 flex items-center justify-between">
	<h1 class="text-2xl font-bold">Forwarding Rules</h1>
	<button
		onclick={() => {
			editingRule = null;
			showModal = true;
		}}
		class="flex items-center gap-2 rounded-lg bg-blue-600 px-4 py-2 text-white hover:bg-blue-700"
	>
		<Plus class="h-4 w-4" />
		Create Rule
	</button>
</div>

{#if loading}
	<div class="space-y-4">
		{#each [1, 2, 3] as _}
			<div class="h-16 animate-pulse rounded-lg bg-gray-200"></div>
		{/each}
	</div>
{:else if rules.length === 0}
	<div class="rounded-lg bg-white p-12 text-center shadow">
		<Bell class="mx-auto mb-4 h-16 w-16 text-gray-400" />
		<h3 class="mb-2 text-lg font-semibold">No forwarding rules</h3>
		<p class="text-gray-500">Create rules to forward SMS and notifications to webhooks.</p>
	</div>
{:else}
	<div class="overflow-hidden rounded-lg bg-white shadow">
		<table class="w-full">
			<thead class="bg-gray-50">
				<tr>
					<th class="px-6 py-3 text-left text-xs font-medium uppercase text-gray-500">Type</th>
					<th class="px-6 py-3 text-left text-xs font-medium uppercase text-gray-500">Filter</th>
					<th class="px-6 py-3 text-left text-xs font-medium uppercase text-gray-500">Webhook</th>
					<th class="px-6 py-3 text-left text-xs font-medium uppercase text-gray-500">Status</th>
					<th class="px-6 py-3 text-right text-xs font-medium uppercase text-gray-500">Actions</th>
				</tr>
			</thead>
			<tbody class="divide-y">
				{#each rules as rule}
					<tr class="hover:bg-gray-50">
						<td class="px-6 py-4">
							<span class="flex items-center gap-2">
								{#if rule.trigger_type === 'sms'}
									<MessageSquare class="h-4 w-4 text-blue-500" />
								{:else}
									<Bell class="h-4 w-4 text-purple-500" />
								{/if}
								<span class="capitalize">{rule.trigger_type}</span>
							</span>
						</td>
						<td class="px-6 py-4 text-sm text-gray-600">
							{rule.sender_filter || rule.content_filter || 'No filter'}
						</td>
						<td class="px-6 py-4 font-mono text-sm text-gray-600" title={rule.webhook_url}>
							{truncateUrl(rule.webhook_url)}
						</td>
						<td class="px-6 py-4">
							<button onclick={() => toggleRule(rule)} class="focus:outline-none">
								{#if rule.is_active}
									<ToggleRight class="h-6 w-6 text-green-500" />
								{:else}
									<ToggleLeft class="h-6 w-6 text-gray-400" />
								{/if}
							</button>
						</td>
						<td class="px-6 py-4 text-right">
							<button
								onclick={() => {
									editingRule = rule;
									showModal = true;
								}}
								class="rounded p-2 hover:bg-gray-100"
								title="Edit"
							>
								<Edit2 class="h-4 w-4" />
							</button>
							<button
								onclick={() => deleteRule(rule.id)}
								class="rounded p-2 text-red-500 hover:bg-gray-100"
								title="Delete"
							>
								<Trash2 class="h-4 w-4" />
							</button>
						</td>
					</tr>
				{/each}
			</tbody>
		</table>
	</div>
{/if}

{#if showModal}
	<RuleModal rule={editingRule} onClose={() => (showModal = false)} onSave={handleSave} />
{/if}
