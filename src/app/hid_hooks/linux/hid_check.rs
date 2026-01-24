use std::collections::HashMap;

pub static EXPORTS_BASELINE: std::sync::OnceLock<HashMap<String, [u8; 16]>> =
    std::sync::OnceLock::new();

pub fn enumerate_hid_exports() {
    // STUB
}

pub fn detect_hid_hooks() -> Vec<String> {
    // STUB
    Vec::new()
}

pub fn get_hid_export_addr(_export_name: &str) -> Option<usize> {
    // STUB
    None
}

pub fn get_hid_export_baseline16(_export_name: &str) -> Option<[u8; 16]> {
    // STUB
    None
}
