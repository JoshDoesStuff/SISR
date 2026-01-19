#[cfg(all(target_os = "windows", target_arch = "x86_64"))]
pub use super::windows::hid_check::*;

#[cfg(not(all(target_os = "windows", target_arch = "x86_64")))]
pub use super::linux::hid_check::*;
