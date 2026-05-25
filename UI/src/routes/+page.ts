import { clientWithSvelteFetch, wrapClientError } from '$lib/api/client';
import { log } from '$lib/log';

const TIMEOUT_MS = 1000;
export const load = async ({ fetch }) => {

    const client = clientWithSvelteFetch(fetch);

    log.debug('Fetching steam status...');
    const [steamStatus, devices, viiperInfo, versionInfo, initialLaunch] = await Promise.all([
        wrapClientError(client.GET('/api/v1/steam/status', {
            signal: AbortSignal.timeout(TIMEOUT_MS)
        })),
        wrapClientError(client.GET('/api/v1/devices', {
            signal: AbortSignal.timeout(TIMEOUT_MS)
        })),
        wrapClientError(client.GET('/api/v1/viiper/status', {
            signal: AbortSignal.timeout(TIMEOUT_MS)
        })),
        wrapClientError(client.GET('/api/v1/version/info', {
            signal: AbortSignal.timeout(TIMEOUT_MS)
        })),
        wrapClientError(client.GET('/api/v1/initial_launch', {
            signal: AbortSignal.timeout(TIMEOUT_MS)
        }))
    ]);

    return {
        steamStatus: steamStatus,
        devices: devices,
        viiperInfo: viiperInfo,
        versionInfo: versionInfo,
        initialLaunch: initialLaunch
    };
};
