extern crate alloc;
extern crate core;
extern crate wee_alloc;

use std::ptr::slice_from_raw_parts;
use std::{ptr, slice};
use std::vec;
use std::mem::MaybeUninit;

fn greet(name: &String) {
    log(&["wasm >> ", &greeting(name)].concat())
}

fn greeting(name: &String) -> String {
    return ["Hello", "", &name, "!"].concat()
}

fn log(s: &String) {
    unsafe {
        let (ptr, len) = str_to_ptr(s);
        _log(ptr, len);
    }
}

#[link(wasm_import_module = "env")]
extern "C" {
    #[link_name = "log"]
    fn _log(ptr: u32, len: u32);
}

#[cfg_attr(all(target_arch = "wasm32"), export_name = "greeting")]
#[no_mangle]
pub unsafe extern "C" fn _greeting(ptr: u32, len: u32) -> u64 {
    let name = &ptr_to_str(ptr, len);
    let g = greeting(name);
    let (ptr, len) = str_to_ptr(&g);
    // Note: This changes ownership of the pointer to the external caller. If
    // we didn't call forget, the caller would read back a corrupt value. Since
    // we call forget, the caller must deallocate externally to prevent leaks.
    std::mem::forget(g);
    return ((ptr as u64) << 32) | len as u64;
}

#[cfg_attr(all(target_arch = "wasm32"), export_name = "greet")]
pub unsafe extern "C" fn _greet(ptr: u32, len: u32) {
    greet(&ptr_to_str(ptr, len))
}

fn str_to_ptr(s: &String) -> (u32, u32){
    (s.as_ptr() as u32, s.len() as u32)
}

unsafe fn ptr_to_str(ptr: u32, len: u32) -> String {
    let slice = slice::from_raw_parts_mut(ptr as *mut u8, len as usize);
    let utf8 = std::str::from_utf8_unchecked_mut(slice);
    return String::from(utf8);
}

fn allocate(size: usize) -> *mut u8 {
    let vec: Vec<MaybeUninit<u8>> = Vec::with_capacity(size);

    Box::into_raw(vec.into_boxed_slice()) as *mut u8
}

#[cfg_attr(all(target_arch = "wasm32"), export_name = "allocate")]
#[no_mangle]
pub extern "C" fn _allocate(size: u32) -> *mut u8 { allocate(size as usize)}

#[cfg_attr(all(target_arch = "wasm32"), export_name = "deallocate")]
#[no_mangle]
pub unsafe extern "C"  fn _deallocate(ptr: u32, size: u32) {
    deallocate(ptr as *mut u8, size as usize)
}

unsafe fn deallocate(ptr: *mut u8, size: usize) {
    let _ = Vec::from_raw_parts(ptr, 0, size);
}

#[cfg_attr(all(target_arch = "wasm32"), export_name = "send_vec")]
#[no_mangle]
pub unsafe extern "C" fn _send_vec() -> u64 {
    let v: Vec<u8> = vec![10, 20, 30, 40];
    let (ptr, len) = (v.as_ptr() as u32, v.len() as u32);
    std::mem::forget(ptr);
    (ptr as u64) << 32 | len as u64
}

