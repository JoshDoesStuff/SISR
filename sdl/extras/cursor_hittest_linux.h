#ifndef SISR_CURSOR_HITTEST_LINUX_H
#define SISR_CURSOR_HITTEST_LINUX_H

#include <stdint.h>

int set_x11_cursor_hittest(void *display_ptr, uintptr_t window_id, int hittest);
int set_wayland_cursor_hittest(void *display_ptr, void *surface_ptr, int hittest);

#endif