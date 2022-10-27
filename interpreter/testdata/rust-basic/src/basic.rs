
#[cfg_attr(all(target_arch = "wasm32"), export_name = "add")]
#[no_mangle]
pub extern "C" fn _add(x: u32, y: u32) -> u32 {
    x + y
}

