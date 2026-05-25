<script lang="ts">
import type { PageProps } from './$types';
import IcoSIAPI from '$lib/assets/siapi.svg?component';
import DebugInfoCard from '$lib/components/MainCards/DebugInfo.svelte';
import { client, wrapClientError } from '$lib/api/client';
import { log } from '$lib/log';
import { onMount } from 'svelte';
import { toast } from '$lib/toaster/toaster.svelte';
import QuickSettingsCard from '$lib/components/MainCards/QuickSettingsCard.svelte';
import { tooltip } from '$lib/attachments/tooltip.svelte';
import CheckInitialSetup from '$lib/op/check-initial-setup.svelte';
import UpdateModal from '$lib/components/UpdateModal.svelte';

import IcoSettings from '~icons/mdi/cog';
import IcoMinimize from '~icons/mdi/window-minimize';

let { data }: PageProps = $props();

let debugInfoCardVisible = $state(false);
let quickSettingsVisible = $state(false);

let setupChecker = $state<CheckInitialSetup>()!;

// TODO: re-introduce!!
// onMount(() => {
// 	wrapClientError(client.POST('/api/v1/inject_overlay_notifier'))
// 		.then(() => {
// 			log.info('Overlay notifier injected successfully');
// 		})
// 		.catch((e) => {
// 			log.error('Failed to inject overlay notifier', 'error', e);
// 		});
// 	// if (!data.inputInfo.viiper.) {
// 	// }
// });
</script>

<CheckInitialSetup
	bind:this={setupChecker}
	steamStatus={data.steamStatus}
	initialLaunch={data.initialLaunch.is_initial_launch} />

{#if data.versionInfo.update_available}
	<UpdateModal
		updateInfo={data.updateInfo}
		show={data.updateInfo.update_available && !data.updateInfo.dismissed && !data.updateInfo.skipped} />
{/if}

<main>
	<div>
		{#if debugInfoCardVisible}
			<DebugInfoCard
				steamStatusInfo={data.steamStatus}
				devices={data.devices}
				viiperInfo={data.viiperInfo}
				versionInfo={data.versionInfo}
				onClose={() => (debugInfoCardVisible = false)} />
		{/if}
		{#if quickSettingsVisible}
			<QuickSettingsCard
				steamStatusInfo={data.steamStatus}
				config={data.config}
				onClose={() => (quickSettingsVisible = false)} />
		{/if}
	</div>
	<div>
		<label class="button" for="quickSettingsToggle">
			<input
				type="checkbox"
				id="quickSettingsToggle"
				name="quickSettingsToggle"
				bind:checked={quickSettingsVisible} />
			<IcoSettings style="width: 2em; height: 2em;" />
			<span>Quick Settings</span>
		</label>
		<label class="button" for="debugInfoCardToggle">
			<input
				type="checkbox"
				id="debugInfoCardToggle"
				name="debugInfoCardToggle"
				bind:checked={debugInfoCardVisible} />
			<IcoSIAPI style="width: 2em; height: 2em;" />
			<span>Debug Info</span>
		</label>
	</div>
	<button
		class="quit"
		{@attach tooltip({
			arrow: true,
			arrowFollowCursor: true,
			content: 'Completely shut down SISR'
		})}
		onclick={() =>
			void wrapClientError(client.POST('/api/v1/quit')).catch((e) => {
				toast({
					color: 'firebrick',
					message: e.message || `Failed to quit SISR.`
				});
			})}>Quit SISR</button>
	{#if data.config?.window?.Fullscreen}
		<button
			class="minimize"
			{@attach tooltip({
				arrow: true,
				arrowFollowCursor: true,
				content: 'Close/Minimize this window and keep SISR running in the background'
			})}
			onclick={() =>
				void wrapClientError(
					client.POST('/api/v1/ui', {
						body: {
							show: false
						}
					})
				).catch((e) => {
					toast({
						color: 'firebrick',
						message: e.message || `Failed to minimize SISR.`
					});
				})}><IcoMinimize style="width: 1.6em; height: 1.6em;" /></button>
	{/if}
</main>

<style lang="postcss">
main {
	display: grid;
	place-items: center;
	position: relative;

	& > :first-child {
		display: grid;
		grid-auto-flow: column;
		gap: 1em;
		& > :global(*) {
			max-width: 100%;
			max-height: calc(100% - 2em);
			overflow: hidden;
			position: absolute;
			top: 1em;
		}
	}
	& > :nth-child(2) {
		display: grid;
		grid-auto-flow: column;
		width: fit-content;
		place-items: center;
		padding: 1em;
		padding-bottom: 3em;
		gap: 1em;

		align-self: flex-end;
	}
}

label.button {
	display: grid;
	place-items: center;
	grid-template-columns: auto;
	grid-template-rows: auto min-content;
	padding: 0.5em 2.5em;
	border-radius: 1em;
	width: 100%;
	height: 100%;
	overflow: clip;
	isolation: isolate;
	&:hover:not(:disabled, .plain),
	&:focus-visible:not(:disabled, .plain) {
		color: var(--text-color);
		& > :global(svg) {
			opacity: 1;
		}
		& > span {
			opacity: 1;
		}
	}
	& span {
		font-size: 1.25em;
		font-weight: bold;
		opacity: 0.9;
	}
	& > :global(svg) {
		opacity: 0.9;
	}
	& input[type='checkbox'] {
		position: absolute;
		inset: 0;
		width: 100%;
		height: 100%;
		border-radius: inherit;
		appearance: none;
		border: none;
		outline: none;
		z-index: -1;
		opacity: 0;
		background: var(--card-glass);
		background-color: var(--color-primary);
		&::before {
			content: '';
			position: absolute;
			opacity: 0;
		}
		&::after {
			content: '';
			position: absolute;
			opacity: 0;
		}
		&:checked {
			opacity: 1;
		}
	}
}

.minimize {
	position: absolute;
	top: 0;
	right: 0;
	border-radius: 0 0 0 1em;
	display: grid;
	place-items: center;
	padding: 1em !important;
}

.quit {
	position: absolute;
	bottom: 0;
	right: 0;
	border-radius: 1em 0 0 0;
	display: grid;
	place-items: center;
	padding: 1em 2em !important;
	font-weight: bold;
}
</style>
