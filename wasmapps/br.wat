(module 
    (memory (export "memory") 1 3)
    (data (i32.const 0x0) "\01\01\00\00")
    (data (i32.const 0x10) "Hello, World!")
    (data (i32.const 0x20) "How are ya!")
    (data (i32.const 0x30) "This is sparta")
    (data (i32.const 0x40) "This is me")
    (func (export "load_store") (param $x i32) (result i32)
        (i32.add 
            (i32.load offset=4 (i32.const 3))
            (i32.const 0x10)
        )
    )
)