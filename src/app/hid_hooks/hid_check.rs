use std::collections::HashMap;

static EXPORTS_BASELINE: std::sync::OnceLock<std::collections::HashMap<String, [u8; 16]>> =
    std::sync::OnceLock::new();

#[cfg(target_os = "windows")]
#[allow(unsafe_code, unsafe_op_in_unsafe_fn)]
unsafe fn get_hid_export_table() -> Option<(*const u32, *const u32, usize, usize)> {
    use windows_sys::Win32::System::{
        Diagnostics::Debug::IMAGE_NT_HEADERS64,
        LibraryLoader::LoadLibraryA,
        SystemServices::{IMAGE_DOS_HEADER, IMAGE_EXPORT_DIRECTORY},
    };

    let hid_base = LoadLibraryA(c"hid.dll".as_ptr().cast()) as usize;
    if hid_base == 0 {
        return None;
    }

    let dos = &*(hid_base as *const IMAGE_DOS_HEADER);
    let nt = (hid_base + dos.e_lfanew as usize) as *const IMAGE_NT_HEADERS64;
    let exp_rva = (*nt).OptionalHeader.DataDirectory[0].VirtualAddress;
    if exp_rva == 0 {
        return None;
    }
    let exp = (hid_base + exp_rva as usize) as *const IMAGE_EXPORT_DIRECTORY;

    Some((
        (hid_base + (*exp).AddressOfNames as usize) as *const u32,
        (hid_base + (*exp).AddressOfFunctions as usize) as *const u32,
        (*exp).NumberOfNames as usize,
        hid_base,
    ))
}

#[cfg(target_os = "windows")]
pub fn enumerate_hid_exports() {
    unsafe {
        let Some((names, functions, count, hid_base)) = get_hid_export_table() else {
            tracing::error!("Failed to get hid.dll base");
            return;
        };

        let mut exports_map: HashMap<String, [u8; 16]> = HashMap::new();
        for i in 0..count {
            let name_rva = *(names.add(i)) as usize;
            let name_ptr = (hid_base + name_rva) as *const u8;
            let name = std::ffi::CStr::from_ptr(name_ptr as *const i8)
                .to_string_lossy()
                .to_string();

            let func_rva = *(functions.add(i));
            let func_addr = hid_base + func_rva as usize;

            let first_bytes = std::slice::from_raw_parts(func_addr as *const u8, 16);
            exports_map.insert(name, first_bytes.try_into().unwrap());
        }
        EXPORTS_BASELINE
            .set(exports_map)
            .expect("Failed to set EXPORTS_BASELINE");
    }
}

#[cfg(target_os = "linux")]
pub fn enumerate_hid_exports() {
    // STUB!
    // TODO
}

#[cfg(target_os = "windows")]
pub fn detect_hid_hooks() -> Vec<String> {
    let Some(baseline) = EXPORTS_BASELINE.get() else {
        return Vec::new();
    };
    let mut hooked = Vec::new();

    unsafe {
        let Some((names, funcs, count, hid_base)) = get_hid_export_table() else {
            return baseline.keys().cloned().collect();
        };

        for i in 0..count {
            let name_rva = *names.add(i) as usize;
            let name_ptr = (hid_base + name_rva) as *const i8;
            let Ok(name) = std::ffi::CStr::from_ptr(name_ptr).to_str() else {
                continue;
            };
            let func_rva = *funcs.add(i) as usize;
            if func_rva == 0 {
                continue;
            }

            let entry = hid_base + func_rva;
            let mut cur = [0u8; 16];
            cur.copy_from_slice(std::slice::from_raw_parts(entry as *const u8, 16));

            if let Some(base_bytes) = baseline.get(name)
                && &cur != base_bytes
            {
                hooked.push(name.to_string());
            }
        }
    }

    hooked
}

#[cfg(target_os = "linux")]
pub fn detect_hid_hooks() -> Vec<String> {
    // STUB!
    // TODO
    Vec::new()
}
