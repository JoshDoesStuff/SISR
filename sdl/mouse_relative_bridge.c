#include <SDL3/SDL_mouse.h>

int sisrSetWindowRelativeMouseModeBridge(void *window, int enabled) {
	return SDL_SetWindowRelativeMouseMode((SDL_Window *)window, enabled != 0) ? 1 : 0;
}

int sisrGetWindowRelativeMouseModeBridge(void *window) {
	return SDL_GetWindowRelativeMouseMode((SDL_Window *)window) ? 1 : 0;
}
