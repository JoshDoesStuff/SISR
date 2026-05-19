#pragma once

#include <stdint.h>

#ifdef __cplusplus
extern "C" {
#endif

typedef struct PLH_Detour PLH_Detour;


PLH_Detour* PLH_x64Detour_new(uint64_t fnAddress, uint64_t fnCallback);

int PLH_x64Detour_hook(PLH_Detour* d);

int PLH_x64Detour_unhook(PLH_Detour* d);


uint64_t PLH_x64Detour_trampoline(PLH_Detour* d);


const char* PLH_x64Detour_last_error(PLH_Detour* d);

void PLH_x64Detour_free(PLH_Detour* d);

#ifdef __cplusplus
} 
#endif
