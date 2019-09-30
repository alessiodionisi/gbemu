package cpu

import "testing"

func TestCPU_NopInst(t *testing.T) {
	cpu := NewTestCPU(0x00)
	cpu.ExecuteNextInstruction()

	wantPC := uint16(0x0001)
	if cpu.PC != wantPC {
		t.Errorf("PC register error: want %#04x, got %#04x", wantPC, cpu.PC)
	}
}

func TestCPU_LdInst(t *testing.T) {
	for opCode, inst := range Instructions {
		if inst.Name == "LD" {
			t.Run(inst.Description, func(t *testing.T) {
				cpu := NewTestCPU(opCode)
				cpu.ExecuteNextInstruction()
			})
		}
	}
}

func TestCPU_XorInst(t *testing.T) {
	for opCode, inst := range Instructions {
		if inst.Name == "XOR" {
			t.Run(inst.Description, func(t *testing.T) {
				cpu := NewTestCPU(opCode)
				cpu.ExecuteNextInstruction()
			})
		}
	}
}

func TestCPU_BitInst(t *testing.T) {
	for opCode, inst := range PrefixedInstructions {
		if inst.Name == "BIT" {
			t.Run(inst.Description, func(t *testing.T) {
				cpu := NewTestCPU(opCode)
				cpu.IsNextInstructionPrefixed = true
				cpu.ExecuteNextInstruction()
			})
		}
	}
}

func TestCPU_JrInst(t *testing.T) {
	for opCode, inst := range Instructions {
		if inst.Name == "JR" {
			t.Run(inst.Description, func(t *testing.T) {
				cpu := NewTestCPU(opCode)
				cpu.ExecuteNextInstruction()
			})
		}
	}
}

func TestCPU_IncInst(t *testing.T) {
	for opCode, inst := range Instructions {
		if inst.Name == "INC" {
			t.Run(inst.Description, func(t *testing.T) {
				cpu := NewTestCPU(opCode)
				cpu.ExecuteNextInstruction()
			})
		}
	}
}

func TestCPU_CallInst(t *testing.T) {
	for opCode, inst := range Instructions {
		if inst.Name == "CALL" {
			t.Run(inst.Description, func(t *testing.T) {
				cpu := NewTestCPU(opCode)
				cpu.ExecuteNextInstruction()
			})
		}
	}
}

func TestCPU_PushInst(t *testing.T) {
	for opCode, inst := range Instructions {
		if inst.Name == "PUSH" {
			t.Run(inst.Description, func(t *testing.T) {
				cpu := NewTestCPU(opCode)
				cpu.ExecuteNextInstruction()
			})
		}
	}
}

func TestCPU_RlInst(t *testing.T) {
	for opCode, inst := range PrefixedInstructions {
		if inst.Name == "RL" {
			t.Run(inst.Description, func(t *testing.T) {
				cpu := NewTestCPU(opCode)
				cpu.IsNextInstructionPrefixed = true
				cpu.ExecuteNextInstruction()
			})
		}
	}
}

func TestCPU_RlaInst(t *testing.T) {
	for opCode, inst := range Instructions {
		if inst.Name == "RLA" {
			t.Run(inst.Description, func(t *testing.T) {
				cpu := NewTestCPU(opCode)
				cpu.ExecuteNextInstruction()
			})
		}
	}
}

func TestCPU_PopInst(t *testing.T) {
	for opCode, inst := range Instructions {
		if inst.Name == "POP" {
			t.Run(inst.Description, func(t *testing.T) {
				cpu := NewTestCPU(opCode)
				cpu.ExecuteNextInstruction()
			})
		}
	}
}

func TestCPU_DecInst(t *testing.T) {
	for opCode, inst := range Instructions {
		if inst.Name == "DEC" {
			t.Run(inst.Description, func(t *testing.T) {
				cpu := NewTestCPU(opCode)
				cpu.ExecuteNextInstruction()
			})
		}
	}
}

func TestCPU_RetInst(t *testing.T) {
	for opCode, inst := range Instructions {
		if inst.Name == "RET" {
			t.Run(inst.Description, func(t *testing.T) {
				cpu := NewTestCPU(opCode)
				cpu.ExecuteNextInstruction()
			})
		}
	}
}

func TestCPU_CpInst(t *testing.T) {
	for opCode, inst := range Instructions {
		if inst.Name == "CP" {
			t.Run(inst.Description, func(t *testing.T) {
				cpu := NewTestCPU(opCode)
				cpu.ExecuteNextInstruction()
			})
		}
	}
}

func TestCPU_SubInst(t *testing.T) {
	for opCode, inst := range Instructions {
		if inst.Name == "SUB" {
			t.Run(inst.Description, func(t *testing.T) {
				cpu := NewTestCPU(opCode)
				cpu.ExecuteNextInstruction()
			})
		}
	}
}

func TestCPU_AddInst(t *testing.T) {
	for opCode, inst := range Instructions {
		if inst.Name == "ADD" {
			t.Run(inst.Description, func(t *testing.T) {
				cpu := NewTestCPU(opCode)
				cpu.ExecuteNextInstruction()
			})
		}
	}
}

func TestCPU_JpInst(t *testing.T) {
	for opCode, inst := range Instructions {
		if inst.Name == "JP" {
			t.Run(inst.Description, func(t *testing.T) {
				cpu := NewTestCPU(opCode)
				cpu.ExecuteNextInstruction()
			})
		}
	}
}

func TestCPU_DiInst(t *testing.T) {
	for opCode, inst := range Instructions {
		if inst.Name == "DI" {
			t.Run(inst.Description, func(t *testing.T) {
				cpu := NewTestCPU(opCode)
				cpu.ExecuteNextInstruction()
			})
		}
	}
}

func TestCPU_OrInst(t *testing.T) {
	for opCode, inst := range Instructions {
		if inst.Name == "OR" {
			t.Run(inst.Description, func(t *testing.T) {
				cpu := NewTestCPU(opCode)
				cpu.ExecuteNextInstruction()
			})
		}
	}
}

func TestCPU_AndInst(t *testing.T) {
	for opCode, inst := range Instructions {
		if inst.Name == "AND" {
			t.Run(inst.Description, func(t *testing.T) {
				cpu := NewTestCPU(opCode)
				cpu.ExecuteNextInstruction()
			})
		}
	}
}

func TestCPU_CplInst(t *testing.T) {
	for opCode, inst := range Instructions {
		if inst.Name == "CPL" {
			t.Run(inst.Description, func(t *testing.T) {
				cpu := NewTestCPU(opCode)
				cpu.ExecuteNextInstruction()
			})
		}
	}
}

func TestCPU_EiInst(t *testing.T) {
	for opCode, inst := range Instructions {
		if inst.Name == "EI" {
			t.Run(inst.Description, func(t *testing.T) {
				cpu := NewTestCPU(opCode)
				cpu.ExecuteNextInstruction()
			})
		}
	}
}

func TestCPU_SwapInst(t *testing.T) {
	for opCode, inst := range PrefixedInstructions {
		if inst.Name == "SWAP" {
			t.Run(inst.Description, func(t *testing.T) {
				cpu := NewTestCPU(opCode)
				cpu.IsNextInstructionPrefixed = true
				cpu.ExecuteNextInstruction()
			})
		}
	}
}

func TestCPU_RstInst(t *testing.T) {
	for opCode, inst := range Instructions {
		if inst.Name == "RST" {
			t.Run(inst.Description, func(t *testing.T) {
				cpu := NewTestCPU(opCode)
				cpu.ExecuteNextInstruction()
			})
		}
	}
}

func TestCPU_SrlInst(t *testing.T) {
	for opCode, inst := range PrefixedInstructions {
		if inst.Name == "SRL" {
			t.Run(inst.Description, func(t *testing.T) {
				cpu := NewTestCPU(opCode)
				cpu.IsNextInstructionPrefixed = true
				cpu.ExecuteNextInstruction()
			})
		}
	}
}

func TestCPU_RrInst(t *testing.T) {
	for opCode, inst := range PrefixedInstructions {
		if inst.Name == "RR" {
			t.Run(inst.Description, func(t *testing.T) {
				cpu := NewTestCPU(opCode)
				cpu.IsNextInstructionPrefixed = true
				cpu.ExecuteNextInstruction()
			})
		}
	}
}

func TestCPU_RraInst(t *testing.T) {
	for opCode, inst := range Instructions {
		if inst.Name == "RRA" {
			t.Run(inst.Description, func(t *testing.T) {
				cpu := NewTestCPU(opCode)
				cpu.ExecuteNextInstruction()
			})
		}
	}
}
