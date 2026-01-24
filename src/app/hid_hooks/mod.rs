#[cfg(all(target_os = "windows", target_arch = "x86_64"))]
pub mod windows;

#[cfg(not(all(target_os = "windows", target_arch = "x86_64")))]
pub mod linux;

pub mod hid_check;
pub mod rehook;
