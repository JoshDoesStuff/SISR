<script lang="ts">
import type { components } from '$lib/api/openapi';
import FloatingCard from '../FloatingCard.svelte';

import IcoClose from '~icons/mdi/close';
import IcoSIAPI from '$lib/assets/siapi.svg?component';
import IcoSteam from '~icons/mdi/steam';
import IcoVIIPER from '$lib/assets/viiper_mono.svg?component';
import IcoGamepad from '~icons/fluent/xbox-controller-16-filled';

const {
	steamStatusInfo,
	devices,
	viiperInfo,
	versionInfo,
	onClose
}: {
	steamStatusInfo: components['schemas']['SteamAndCefStatus'];
	devices: components['schemas']['APIDevice'][];
	viiperInfo: components['schemas']['VIIPERStatus'];
	versionInfo: components['schemas']['VersionInfo'];
	onClose?: () => void;
} = $props();

const usedBusIds = $derived.by(() => {
	return (
		devices
			.flatMap((d) => {
				return d.viiper_device?.busId;
			})
			.filter(Boolean)
			// filter duplicates
			.filter((v, idx, arr) => arr.indexOf(v) === idx) as number[]
	);
});
</script>

<FloatingCard>
	<div id="card-content">
		<div>
			<IcoSIAPI style="width: 1.6em; height: 1.6em;" />
			<h2>Debug Info</h2>
			<button class="plain" onclick={() => onClose?.()}>
				<IcoClose style="width: 1.6em; height: 1.6em;" />
			</button>
		</div>
		<div>
			<div class="heading">
				<IcoSIAPI style="width: 1.4em; height: 1.4em;" />
				<h2>General</h2>
			</div>
			<dl>
				<dt>Version</dt>
				<dd>{versionInfo.version}</dd>
				<dt>Update available</dt>
				<dd>{versionInfo.update_available ? 'Yes' : 'No'}</dd>
				{#if versionInfo.update_available}
					<dt>New Version</dt>
					<dd>{versionInfo.new_version}</dd>
				{/if}
			</dl>
			<div class="heading">
				<IcoVIIPER style="width: 1.4em; height: 1.4em;" />
				<h2>VIIPER</h2>
			</div>
			<dl>
				<dt>Address</dt>
				<dd>{viiperInfo.address}</dd>
				<dt>Reachable</dt>
				<dd>{viiperInfo.status?.server == 'VIIPER' ? 'Yes' : 'No'}</dd>
				<dt>Version</dt>
				<dd>{viiperInfo.status?.version ?? 'Unknown'}</dd>
				<dt>Used BusIDs</dt>
				<dd>{usedBusIds.map((id) => id).join(', ') || 'None'}</dd>
			</dl>
			<div class="heading">
				<IcoSteam style="width: 1.4em; height: 1.4em;" />
				<h2>Steam</h2>
			</div>
			<dl>
				<dt>GameID</dt>
				<dd>{steamStatusInfo.steam_game_id}</dd>
				<dt>AppID</dt>
				<dd>{steamStatusInfo.steam_app_id}</dd>
				<dt>Launched via Steam</dt>
				<dd>{steamStatusInfo.launched_via_steam ? 'Yes' : 'No'}</dd>
				<dt>Steam Overlay</dt>
				<dd>TODO</dd>
			</dl>
			<div class="heading">
				<IcoGamepad style="width: 1.4em; height: 1.4em;" />
				<h2>Controller</h2>
			</div>
			<div>
				{#each devices as deviceInfo (deviceInfo)}
					<h3>{deviceInfo.real_device?.name ?? deviceInfo.steam_virtual_device?.name}</h3>
					<dl style="padding-left: 1em;">
						<dt>ID(s)</dt>
						<dd>
							{#if deviceInfo.real_device?.id}
								Real GamepadID: {deviceInfo.real_device.id}
							{/if}
							{#if deviceInfo.real_device?.id && deviceInfo.steam_virtual_device?.id}
								<br />
							{/if}
							{#if deviceInfo.steam_virtual_device?.id}
								Steam Virtual Gamepad ID: {deviceInfo.steam_virtual_device.id}
							{/if}
						</dd>
						<dt>SteamHandle</dt>
						<dd><code>{deviceInfo.steam_virtual_device?.steam_handle ?? 'N/A'}</code></dd>
						<dt>Has VIIPER Device</dt>
						<dd>{deviceInfo.viiper_device ? 'Yes' : 'No'}</dd>
						{#if deviceInfo.viiper_device}
							<dt>VIIPER type</dt>
							<dd>{deviceInfo.viiper_device?.type ?? 'N/A'}</dd>
						{/if}
						<dt>Serial(s)</dt>
						<dd>
							{#if deviceInfo.real_device?.serial}
								Real Gamepad Serial: <code>{deviceInfo.real_device.serial}</code>
							{/if}
							{#if deviceInfo.real_device?.serial && deviceInfo.steam_virtual_device?.serial}
								<br />
							{/if}
							{#if deviceInfo.steam_virtual_device?.serial}
								Steam Virtual Gamepad Serial: <code
									>{deviceInfo.steam_virtual_device.serial}</code>
							{/if}
						</dd>
						<dt>Path(s)</dt>
						<dd>
							{#if deviceInfo.real_device?.path}
								Real Gamepad Path: <code>{deviceInfo.real_device.path}</code>
							{/if}
							{#if deviceInfo.real_device?.path && deviceInfo.steam_virtual_device?.path}
								<br />
							{/if}
							{#if deviceInfo.steam_virtual_device?.path}
								Steam Virtual Gamepad Path: <code
									>{deviceInfo.steam_virtual_device.path}</code>
							{/if}
						</dd>
					</dl>
				{/each}
				{#if devices?.length === 0}
					<p>No controllers connected</p>
				{/if}
			</div>
		</div>
	</div>
</FloatingCard>

<style lang="postcss">
#card-content {
	display: flex;
	flex-direction: column;
	gap: 1em;
	height: 100%;

	& > :first-child {
		display: grid;
		grid-template-columns: min-content 1fr min-content;
		justify-content: center;
		gap: 1em;
		& button {
			padding: 0.5em;
		}
	}
	& > :last-child {
		overflow: auto;
	}
}

h3 {
	margin-top: 1em;
	margin-bottom: 0.5em;
}

.heading {
	display: grid;
	grid-template-columns: min-content min-content;
	gap: 0.5em;
	place-items: center;
	place-self: center;
}

dl {
	display: grid;
	grid-template-columns: min-content auto;
	column-gap: 1em;
	row-gap: 0.25em;
	padding-bottom: 1.5em;
	dt {
		font-weight: bold;
		white-space: nowrap;
	}
	dd {
		color: var(--text-muted);
	}
}
</style>
