#include "cursor_hittest_linux.h"

#if defined(__linux__)

#include <stdint.h>
#include <string.h>

#include <X11/Xlib.h>
#include <X11/extensions/shape.h>
#include <wayland-client.h>

int set_x11_cursor_hittest(void *display_ptr, uintptr_t window_id, int hittest) {
	Display *display = (Display *)display_ptr;
	Window window = (Window)window_id;
	if (!display || !window) {
		return 0;
	}

	if (hittest) {
		XWindowAttributes attrs;
		if (!XGetWindowAttributes(display, window, &attrs)) {
			return -1;
		}

		XRectangle rect;
		rect.x = 0;
		rect.y = 0;
		rect.width = (unsigned short)attrs.width;
		rect.height = (unsigned short)attrs.height;

		Region region = XCreateRegion();
		if (!region) {
			return -1;
		}
		XUnionRectWithRegion(&rect, region, region);
		XShapeCombineRegion(display, window, ShapeInput, 0, 0, region, ShapeSet);
		XDestroyRegion(region);
	} else {
		Region region = XCreateRegion();
		if (!region) {
			return -1;
		}
		XShapeCombineRegion(display, window, ShapeInput, 0, 0, region, ShapeSet);
		XDestroyRegion(region);
	}

	XFlush(display);
	return 0;
}

int get_x11_cursor_hittest(void *display_ptr, uintptr_t window_id) {
	Display *display = (Display *)display_ptr;
	Window window = (Window)window_id;
	if (!display || !window) {
		return -1;
	}

	Bool bounding_shaped = False;
	Bool input_shaped = False;
	int x_bounding = 0;
	int y_bounding = 0;
	unsigned int w_bounding = 0;
	unsigned int h_bounding = 0;
	int x_input = 0;
	int y_input = 0;
	unsigned int w_input = 0;
	unsigned int h_input = 0;

	if (!XShapeQueryExtents(
			display,
			window,
			&bounding_shaped,
			&x_bounding,
			&y_bounding,
			&w_bounding,
			&h_bounding,
			&input_shaped,
			&x_input,
			&y_input,
			&w_input,
			&h_input
		)) {
		return -1;
	}

	if (!input_shaped) {
		return 1;
	}

	if (w_input == 0 || h_input == 0) {
		return 0;
	}

	return 1;
}

struct wayland_registry_state {
	struct wl_compositor *compositor;
};

static void registry_global(
	void *data,
	struct wl_registry *registry,
	uint32_t name,
	const char *interface,
	uint32_t version
) {
	struct wayland_registry_state *state = (struct wayland_registry_state *)data;
	if (strcmp(interface, wl_compositor_interface.name) == 0) {
		uint32_t bind_version = version < 4 ? version : 4;
		state->compositor = (struct wl_compositor *)wl_registry_bind(registry, name, &wl_compositor_interface, bind_version);
	}
}

static void registry_global_remove(void *data, struct wl_registry *registry, uint32_t name) {
	(void)data;
	(void)registry;
	(void)name;
}

static const struct wl_registry_listener registry_listener = {
	registry_global,
	registry_global_remove,
};

static struct wl_compositor *get_wayland_compositor(struct wl_display *display) {
	static struct wl_display *cached_display = NULL;
	static struct wl_compositor *cached_compositor = NULL;

	if (display == cached_display && cached_compositor != NULL) {
		return cached_compositor;
	}

	struct wl_registry *registry = wl_display_get_registry(display);
	if (!registry) {
		return NULL;
	}

	struct wayland_registry_state state;
	state.compositor = NULL;

	wl_registry_add_listener(registry, &registry_listener, &state);
	wl_display_roundtrip(display);
	wl_registry_destroy(registry);

	cached_display = display;
	cached_compositor = state.compositor;
	return state.compositor;
}

int set_wayland_cursor_hittest(void *display_ptr, void *surface_ptr, int hittest) {
	struct wl_display *display = (struct wl_display *)display_ptr;
	struct wl_surface *surface = (struct wl_surface *)surface_ptr;
	if (!display || !surface) {
		return 0;
	}

	if (hittest) {
		wl_surface_set_input_region(surface, NULL);
		wl_surface_commit(surface);
		wl_display_flush(display);
		return 0;
	}

	struct wl_compositor *compositor = get_wayland_compositor(display);
	if (!compositor) {
		return -1;
	}

	struct wl_region *region = wl_compositor_create_region(compositor);
	if (!region) {
		return -1;
	}

	wl_region_add(region, 0, 0, 0, 0);
	wl_surface_set_input_region(surface, region);
	wl_region_destroy(region);
	wl_surface_commit(surface);
	wl_display_flush(display);
	return 0;
}

#else

int set_x11_cursor_hittest(void *display_ptr, uintptr_t window_id, int hittest) {
	(void)display_ptr;
	(void)window_id;
	(void)hittest;
	return 0;
}

int set_wayland_cursor_hittest(void *display_ptr, void *surface_ptr, int hittest) {
	(void)display_ptr;
	(void)surface_ptr;
	(void)hittest;
	return 0;
}

int get_x11_cursor_hittest(void *display_ptr, uintptr_t window_id) {
	(void)display_ptr;
	(void)window_id;
	return -1;
}

#endif