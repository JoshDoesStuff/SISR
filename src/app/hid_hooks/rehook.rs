#[cfg(all(target_os = "windows", target_arch = "x86_64"))]
pub use super::windows::rehook::*;

#[cfg(not(all(target_os = "windows", target_arch = "x86_64")))]
pub use super::linux::rehook::*;
