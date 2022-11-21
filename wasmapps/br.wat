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
    (func (export "load_item") (result i32)
        (i32.const 0x10)
        (i32.load offset=258)
        (i32.const 0x4)
        (i32.store)
        (i32.const 0x1)
    )

    (func (export "fac") (param $x i32) (result i32)
        (i32.const 0x0)
        (local.get 0x0)
        (i32.eq)
        (if 
            (then
                i32.const 1
                return
            )
            (else 
                i32.const 1
                local.get 0
                i32.sub
                return 
            )
        )  
        i32.const 0
        return 
    )
)