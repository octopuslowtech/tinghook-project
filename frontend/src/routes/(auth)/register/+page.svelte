<script lang="ts">
	import { goto } from '$app/navigation';
	import { auth } from '$lib/stores/auth';
	import { api } from '$lib/api/client';

	let email = $state('');
	let password = $state('');
	let confirmPassword = $state('');
	let error = $state('');
	let loading = $state(false);

	function validateEmail(email: string): boolean {
		const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
		return emailRegex.test(email);
	}

	async function handleRegister() {
		error = '';

		if (!email || !password || !confirmPassword) {
			error = 'All fields are required';
			return;
		}

		if (!validateEmail(email)) {
			error = 'Please enter a valid email address';
			return;
		}

		if (password.length < 8) {
			error = 'Password must be at least 8 characters';
			return;
		}

		if (password !== confirmPassword) {
			error = 'Passwords do not match';
			return;
		}

		loading = true;

		try {
			const response = await api.register(email, password);
			auth.login(response.user, response.token);
			goto('/dashboard');
		} catch (e) {
			error = e instanceof Error ? e.message : 'Registration failed';
		} finally {
			loading = false;
		}
	}
</script>

<div class="bg-white rounded-lg shadow-md p-8">
	<h1 class="text-2xl font-bold text-center text-gray-800 mb-6">Create Account</h1>

	{#if error}
		<div class="bg-red-50 text-red-600 p-3 rounded-md mb-4 text-sm">
			{error}
		</div>
	{/if}

	<form onsubmit={(e) => { e.preventDefault(); handleRegister(); }}>
		<div class="mb-4">
			<label for="email" class="block text-sm font-medium text-gray-700 mb-1">Email</label>
			<input
				type="email"
				id="email"
				bind:value={email}
				class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
				placeholder="you@example.com"
				disabled={loading}
			/>
		</div>

		<div class="mb-4">
			<label for="password" class="block text-sm font-medium text-gray-700 mb-1">Password</label>
			<input
				type="password"
				id="password"
				bind:value={password}
				class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
				placeholder="Minimum 8 characters"
				disabled={loading}
			/>
		</div>

		<div class="mb-6">
			<label for="confirmPassword" class="block text-sm font-medium text-gray-700 mb-1">Confirm Password</label>
			<input
				type="password"
				id="confirmPassword"
				bind:value={confirmPassword}
				class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
				placeholder="Confirm your password"
				disabled={loading}
			/>
		</div>

		<button
			type="submit"
			disabled={loading}
			class="w-full bg-blue-600 text-white py-2 px-4 rounded-md hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center"
		>
			{#if loading}
				<svg class="animate-spin -ml-1 mr-2 h-4 w-4 text-white" fill="none" viewBox="0 0 24 24">
					<circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
					<path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
				</svg>
				Creating account...
			{:else}
				Create Account
			{/if}
		</button>
	</form>

	<p class="mt-4 text-center text-sm text-gray-600">
		Already have an account?
		<a href="/login" class="text-blue-600 hover:text-blue-700 font-medium">Sign in</a>
	</p>
</div>
