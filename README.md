# wasmexperiments
Fiddling around with wasm binary format

This is an experimental repository where I'm desperately (or not so) trying to comprehend the internal structure of WASM 
and hopefully make something good out of it. Much of the inspiration is drawn from [wagon](https://github.com/go-interpreter/wagon)

# TODO
### Deserialization
- [x] Types Section
- [x] Function Section
- [x] Table Section
- [x] Global Section
- [x] Export Section
- [x] Start Section
- [x] Element Section
- [x] Code Section
- [x] Data Section

### Instructions
- [ ] Define all existing instructions in an informative way (its name, number of input params)
- [ ] Dump module's code in the form of bytecode instructions (sort of debug function)

### Basic interpretation
- [ ] Write a basic wasm interpreter that interprets addition program

### Update 
I decided to make it my university research project, so it means rendering this project alive again! 

### Resources 
[Bringing the Web up to Speed with WebAssembly](https://people.mpi-sws.org/~rossberg/papers/Haas,%20Rossberg,%20Schuff,%20Titzer,%20Gohman,%20Wagner,%20Zakai,%20Bastien,%20Holman%20-%20Bringing%20the%20Web%20up%20to%20Speed%20with%20WebAssembly.pdf)\
[A fast in-place interpreter for WebAssembly](https://arxiv.org/pdf/2205.01183.pdf)\
[WebAssembly: The Definitive Guide](https://www.oreilly.com/library/view/webassembly-the-definitive/9781492089834)
