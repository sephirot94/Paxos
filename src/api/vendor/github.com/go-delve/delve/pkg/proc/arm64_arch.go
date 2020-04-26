package proc

import (
	"encoding/binary"

	"github.com/go-delve/delve/pkg/dwarf/frame"
	"github.com/go-delve/delve/pkg/dwarf/op"
	"golang.org/x/arch/arm64/arm64asm"
)

// ARM64 represents the ARM64 CPU architecture.
type ARM64 struct {
	gStructOffset uint64
	goos          string

	// crosscall2fn is the DIE of crosscall2, a function used by the go runtime
	// to call C functions. This function in go 1.9 (and previous versions) had
	// a bad frame descriptor which needs to be fixed to generate good stack
	// traces.
	crosscall2fn *Function

	// sigreturnfn is the DIE of runtime.sigreturn, the return trampoline for
	// the signal handler. See comment in FixFrameUnwindContext for a
	// description of why this is needed.
	sigreturnfn *Function
}

const (
	arm64DwarfIPRegNum uint64 = 32
	arm64DwarfSPRegNum uint64 = 31
	arm64DwarfBPRegNum uint64 = 29
)

var arm64BreakInstruction = []byte{0x0, 0x0, 0x20, 0xd4}

// ARM64Arch returns an initialized ARM64
// struct.
func ARM64Arch(goos string) *ARM64 {
	return &ARM64{
		goos: goos,
	}
}

// PtrSize returns the size of a pointer
// on this architecture.
func (a *ARM64) PtrSize() int {
	return 8
}

// MaxInstructionLength returns the maximum lenght of an instruction.
func (a *ARM64) MaxInstructionLength() int {
	return 4
}

// BreakpointInstruction returns the Breakpoint
// instruction for this architecture.
func (a *ARM64) BreakpointInstruction() []byte {
	return arm64BreakInstruction
}

// BreakInstrMovesPC returns whether the
// breakpoint instruction will change the value
// of PC after being executed
func (a *ARM64) BreakInstrMovesPC() bool {
	return false
}

// BreakpointSize returns the size of the
// breakpoint instruction on this architecture.
func (a *ARM64) BreakpointSize() int {
	return len(arm64BreakInstruction)
}

// Always return false for now.
func (a *ARM64) DerefTLS() bool {
	return false
}

// FixFrameUnwindContext adds default architecture rules to fctxt or returns
// the default frame unwind context if fctxt is nil.
func (a *ARM64) FixFrameUnwindContext(fctxt *frame.FrameContext, pc uint64, bi *BinaryInfo) *frame.FrameContext {
	if a.sigreturnfn == nil {
		a.sigreturnfn = bi.LookupFunc["runtime.sigreturn"]
	}

	if fctxt == nil || (a.sigreturnfn != nil && pc >= a.sigreturnfn.Entry && pc < a.sigreturnfn.End) {
		// When there's no frame descriptor entry use BP (the frame pointer) instead
		// - return register is [bp + a.PtrSize()] (i.e. [cfa-a.PtrSize()])
		// - cfa is bp + a.PtrSize()*2
		// - bp is [bp] (i.e. [cfa-a.PtrSize()*2])
		// - sp is cfa

		// When the signal handler runs it will move the execution to the signal
		// handling stack (installed using the sigaltstack system call).
		// This isn't a proper stack switch: the pointer to g in TLS will still
		// refer to whatever g was executing on that thread before the signal was
		// received.
		// Since go did not execute a stack switch the previous value of sp, pc
		// and bp is not saved inside g.sched, as it normally would.
		// The only way to recover is to either read sp/pc from the signal context
		// parameter (the ucontext_t* parameter) or to unconditionally follow the
		// frame pointer when we get to runtime.sigreturn (which is what we do
		// here).

		return &frame.FrameContext{
			RetAddrReg: arm64DwarfIPRegNum,
			Regs: map[uint64]frame.DWRule{
				arm64DwarfIPRegNum: frame.DWRule{
					Rule:   frame.RuleOffset,
					Offset: int64(-a.PtrSize()),
				},
				arm64DwarfBPRegNum: frame.DWRule{
					Rule:   frame.RuleOffset,
					Offset: int64(-2 * a.PtrSize()),
				},
				arm64DwarfSPRegNum: frame.DWRule{
					Rule:   frame.RuleValOffset,
					Offset: 0,
				},
			},
			CFA: frame.DWRule{
				Rule:   frame.RuleCFA,
				Reg:    arm64DwarfBPRegNum,
				Offset: int64(2 * a.PtrSize()),
			},
		}
	}

	if a.crosscall2fn == nil {
		a.crosscall2fn = bi.LookupFunc["crosscall2"]
	}

	if a.crosscall2fn != nil && pc >= a.crosscall2fn.Entry && pc < a.crosscall2fn.End {
		rule := fctxt.CFA
		if rule.Offset == crosscall2SPOffsetBad {
			switch a.goos {
			case "windows":
				rule.Offset += crosscall2SPOffsetWindows
			default:
				rule.Offset += crosscall2SPOffsetNonWindows
			}
		}
		fctxt.CFA = rule
	}

	// We assume that RBP is the frame pointer and we want to keep it updated,
	// so that we can use it to unwind the stack even when we encounter frames
	// without descriptor entries.
	// If there isn't a rule already we emit one.
	if fctxt.Regs[arm64DwarfBPRegNum].Rule == frame.RuleUndefined {
		fctxt.Regs[arm64DwarfBPRegNum] = frame.DWRule{
			Rule:   frame.RuleFramePointer,
			Reg:    arm64DwarfBPRegNum,
			Offset: 0,
		}
	}

	return fctxt
}

func (a *ARM64) RegSize(regnum uint64) int {
	// fp registers
	if regnum >= 64 && regnum <= 95 {
		return 16
	}

	return 8 // general registers
}

// The mapping between hardware registers and DWARF registers is specified
// in the DWARF for the ARM® Architecture page 7,
// Table 1
// http://infocenter.arm.com/help/topic/com.arm.doc.ihi0040b/IHI0040B_aadwarf.pdf
var arm64DwarfToHardware = map[int]arm64asm.Reg{
	0:  arm64asm.X0,
	1:  arm64asm.X1,
	2:  arm64asm.X2,
	3:  arm64asm.X3,
	4:  arm64asm.X4,
	5:  arm64asm.X5,
	6:  arm64asm.X6,
	7:  arm64asm.X7,
	8:  arm64asm.X8,
	9:  arm64asm.X9,
	10: arm64asm.X10,
	11: arm64asm.X11,
	12: arm64asm.X12,
	13: arm64asm.X13,
	14: arm64asm.X14,
	15: arm64asm.X15,
	16: arm64asm.X16,
	17: arm64asm.X17,
	18: arm64asm.X18,
	19: arm64asm.X19,
	20: arm64asm.X20,
	21: arm64asm.X21,
	22: arm64asm.X22,
	23: arm64asm.X23,
	24: arm64asm.X24,
	25: arm64asm.X25,
	26: arm64asm.X26,
	27: arm64asm.X27,
	28: arm64asm.X28,
	29: arm64asm.X29,
	30: arm64asm.X30,
	31: arm64asm.SP,

	64: arm64asm.V0,
	65: arm64asm.V1,
	66: arm64asm.V2,
	67: arm64asm.V3,
	68: arm64asm.V4,
	69: arm64asm.V5,
	70: arm64asm.V6,
	71: arm64asm.V7,
	72: arm64asm.V8,
	73: arm64asm.V9,
	74: arm64asm.V10,
	75: arm64asm.V11,
	76: arm64asm.V12,
	77: arm64asm.V13,
	78: arm64asm.V14,
	79: arm64asm.V15,
	80: arm64asm.V16,
	81: arm64asm.V17,
	82: arm64asm.V18,
	83: arm64asm.V19,
	84: arm64asm.V20,
	85: arm64asm.V21,
	86: arm64asm.V22,
	87: arm64asm.V23,
	88: arm64asm.V24,
	89: arm64asm.V25,
	90: arm64asm.V26,
	91: arm64asm.V27,
	92: arm64asm.V28,
	93: arm64asm.V29,
	94: arm64asm.V30,
	95: arm64asm.V31,
}

func maxArm64DwarfRegister() int {
	max := int(arm64DwarfIPRegNum)
	for i := range arm64DwarfToHardware {
		if i > max {
			max = i
		}
	}
	return max
}

// RegistersToDwarfRegisters converts hardware registers to the format used
// by the DWARF expression interpreter.
func (a *ARM64) RegistersToDwarfRegisters(staticBase uint64, regs Registers) op.DwarfRegisters {
	dregs := make([]*op.DwarfRegister, maxArm64DwarfRegister()+1)

	dregs[arm64DwarfIPRegNum] = op.DwarfRegisterFromUint64(regs.PC())
	dregs[arm64DwarfSPRegNum] = op.DwarfRegisterFromUint64(regs.SP())
	dregs[arm64DwarfBPRegNum] = op.DwarfRegisterFromUint64(regs.BP())

	for dwarfReg, asmReg := range arm64DwarfToHardware {
		v, err := regs.Get(int(asmReg))
		if err == nil {
			dregs[dwarfReg] = op.DwarfRegisterFromUint64(v)
		}
	}

	return op.DwarfRegisters{
		StaticBase: staticBase,
		Regs:       dregs,
		ByteOrder:  binary.LittleEndian,
		PCRegNum:   arm64DwarfIPRegNum,
		SPRegNum:   arm64DwarfSPRegNum,
		BPRegNum:   arm64DwarfBPRegNum,
	}
}

// AddrAndStackRegsToDwarfRegisters returns DWARF registers from the passed in
// PC, SP, and BP registers in the format used by the DWARF expression interpreter.
func (a *ARM64) AddrAndStackRegsToDwarfRegisters(staticBase, pc, sp, bp uint64) op.DwarfRegisters {
	dregs := make([]*op.DwarfRegister, arm64DwarfIPRegNum+1)
	dregs[arm64DwarfIPRegNum] = op.DwarfRegisterFromUint64(pc)
	dregs[arm64DwarfSPRegNum] = op.DwarfRegisterFromUint64(sp)
	dregs[arm64DwarfBPRegNum] = op.DwarfRegisterFromUint64(bp)

	return op.DwarfRegisters{
		StaticBase: staticBase,
		Regs:       dregs,
		ByteOrder:  binary.LittleEndian,
		PCRegNum:   arm64DwarfIPRegNum,
		SPRegNum:   arm64DwarfSPRegNum,
		BPRegNum:   arm64DwarfBPRegNum,
	}
}
