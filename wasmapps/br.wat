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
    (func (export "simplest") (result i32)
        (i32.const 0x20)
        (i32.const 0x30)
        (i32.add)
    )
    (func (export "load_item") (result i32)
        (i32.const 0x10)
        (i32.load offset=258)
        (i32.const 0x4)
        (i32.store)
        (i32.const 0x1)
    )

    (func (export "fac") (param $x i32) (param $y i32) (result i32)
        (i32.const 0x10)
        (local.tee $x)
        (return)
    )

    (func (export "block_test") (param $x i32) (result i32)
        (block $l (result i32)
            ;; x = x + 1
            (local.get 0)
            (i32.const 0x1)
            (i32.add)
            (local.set 0)
            
            ;; x < 10?
            (local.get 0)
            (i32.const 0xA)
            (i32.lt_s) 
            ;; jump to myblock if it is
            (i32.const 0x1)
            (return)
            (br_if 0)
        )        

        (local.get 0)
        (return)
    ) 

    (func (export "multiblock_test") (param $x i32) (result i32)
        (block $outer (result i32)
            (block $inner (result i32)
                ;; x = x + 1
                (local.get 0)
                (i32.const 0x1)
                (i32.add)
                (local.set 0)
                
                ;; x < 10?
                (local.get 0)
                (i32.const 0xA)
                (i32.lt_s) 
                ;; jump to myblock if it is
                (i32.const 0x1)
                (return)
            )
            (i32.const 0x10)
            (return)
        )        
        (i32.const 0x20)
        (return)
    )

    (func (export "loop_test") (param $x i32) (result i32)
        (loop $outer 
            (local.get $x)
            (i32.const 0x1)
            (i32.add)
            (local.set $x)
            (local.get $x)
            (i32.const 0xA)
            (i32.lt_s)
            (br_if $outer)
        )
        (local.get $x)
    )

    (func (export "tricky_loop_test") (param $x i32) (result i32)
        ;; 1 
        ;;  0
        ;;   some code
        ;; possible location to jump to 
        (block $outer 
            ;; possible location to jump to 
            (loop $inner 
                (local.get $x)
                (i32.const 0x1)
                (i32.add)
                (local.set $x)
                (local.get $x)
                (i32.const 0xA)
                (i32.lt_s)
                (br_if $inner)
                (br $outer)
            )
        )
        (local.get $x)
    )

)

