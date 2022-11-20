(module
  (type (;0;) (func (param i32) (result i32)))
  (func (;0;) (type 0) (param i32) (result i32)
    local.get 0
    i32.load)
  (memory (;0;) 1 10)
  (export "memory" (memory 0))
  (data (;0;) (i32.const 0) "\01\01\00\00"))
