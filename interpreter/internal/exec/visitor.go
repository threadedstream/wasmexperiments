package exec

func Visit(code []byte, do func(Opcode, *int) error) error {
	pc := 0
	var err error
	for pc < len(code) && err == nil {
		opcode := Opcode(code[pc])
		err = do(opcode, &pc)
	}
	return err
}
