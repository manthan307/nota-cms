<script>
	let showPassword = false;
</script>

<section class="min-h-screen w-full flex justify-center items-center bg-gray-950 text-white px-4">
	<div class="w-full max-w-sm bg-gray-900 rounded-2xl shadow-2xl p-8">
		<h1 class="text-2xl font-semibold text-center mb-6">Welcome Back</h1>

		<form class="space-y-5">
			<!-- Email -->
			<div class="flex flex-col">
				<label for="email" class="text-sm font-medium mb-1">Email</label>
				<input
					type="email"
					id="email"
					name="email"
					placeholder="you@example.com"
					class="w-full px-4 py-2 rounded-lg bg-gray-800 border border-gray-700 focus:outline-none focus:ring-2 focus:ring-orange-600 focus:border-transparent"
				/>
			</div>

			<!-- Password -->
			<div class="flex flex-col relative">
				<label for="password" class="text-sm font-medium mb-1">Password</label>
				<input
					type={showPassword ? 'text' : 'password'}
					id="password"
					name="password"
					placeholder="••••••••"
					class="w-full px-4 py-2 rounded-lg bg-gray-800 border border-gray-700 focus:outline-none focus:ring-2 focus:ring-orange-600 focus:border-transparent pr-10"
				/>
				<button
					type="button"
					on:click={() => (showPassword = !showPassword)}
					class="absolute right-3 top-9 text-gray-400 hover:text-gray-200 focus:outline-none"
				>
					{#if showPassword}
						🙈
					{:else}
						👁️
					{/if}
				</button>
			</div>

			<div class="flex justify-between items-center text-sm">
				<a href="/" class="text-orange-500 hover:underline">Forgot Password?</a>
			</div>

			<button
				type="submit"
				class="w-full py-2 bg-orange-700 hover:bg-orange-800 rounded-lg font-medium
				transition-colors duration-200"
			>
				Login
			</button>
		</form>

		<p class="text-center text-sm text-gray-400 mt-6">
			Don’t have an account?
			<a href="/register" class="text-orange-500 hover:underline">Sign up</a>
		</p>
	</div>
</section>
