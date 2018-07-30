package tbcload

type InstructionDesc struct {
	name        string
	numBytes    int
	stackEffect int
	numOperands int
	opTypes     [2]byte
}

const (
	OPERAND_NONE    byte = iota
	OPERAND_INT1         /* One byte signed integer. */
	OPERAND_INT4         /* Four byte signed integer. */
	OPERAND_UINT1        /* One byte unsigned integer. */
	OPERAND_UINT4        /* Four byte unsigned integer. */
	OPERAND_IDX4         /* Four byte signed index (actually an * integer, but displayed differently.) */
	OPERAND_LVT1         /* One byte unsigned index into the local * variable table. */
	OPERAND_LVT4         /* Four byte unsigned index into the local * variable table. */
	OPERAND_AUX4         /* Four byte unsigned index into the aux data * table. */
	OPERAND_OFFSET1      /* One byte signed jump offset. */
	OPERAND_OFFSET4      /* Four byte signed jump offset. */
	OPERAND_LIT1         /* One byte unsigned index into table of * literals. */
	OPERAND_LIT4         /* Four byte unsigned index into table of * literals. */
	OPERAND_SCLS1        /* Index into tclStringClassTable. */
)
const INT_MIN = 0x8000

var tclOpTable = []InstructionDesc{
	/* Name	      Bytes stackEffect #Opnds  Operand types */
	{"done", 1, -1, 0, [2]byte{OPERAND_NONE}},
	/* Finish ByteCode execution and return stktop (top stack item) */
	{"push1", 2, +1, 1, [2]byte{OPERAND_LIT1}},
	/* Push object at ByteCode objArray[op1] */
	{"push4", 5, +1, 1, [2]byte{OPERAND_LIT4}},
	/* Push object at ByteCode objArray[op4] */
	{"pop", 1, -1, 0, [2]byte{OPERAND_NONE}},
	/* Pop the topmost stack object */
	{"dup", 1, +1, 0, [2]byte{OPERAND_NONE}},
	/* Duplicate the topmost stack object and push the result */
	{"strcat", 2, INT_MIN, 1, [2]byte{OPERAND_UINT1}},
	/* Concatenate the top op1 items and push result */
	{"invokeStk1", 2, INT_MIN, 1, [2]byte{OPERAND_UINT1}},
	/* Invoke command named objv[0]; <objc,objv> = <op1,top op1> */
	{"invokeStk4", 5, INT_MIN, 1, [2]byte{OPERAND_UINT4}},
	/* Invoke command named objv[0]; <objc,objv> = <op4,top op4> */
	{"evalStk", 1, 0, 0, [2]byte{OPERAND_NONE}},
	/* Evaluate command in stktop using Tcl_EvalObj. */
	{"exprStk", 1, 0, 0, [2]byte{OPERAND_NONE}},
	/* Execute expression in stktop using Tcl_ExprStringObj. */

	{"loadScalar1", 2, 1, 1, [2]byte{OPERAND_LVT1}},
	/* Load scalar variable at index op1 <= 255 in call frame */
	{"loadScalar4", 5, 1, 1, [2]byte{OPERAND_LVT4}},
	/* Load scalar variable at index op1 >= 256 in call frame */
	{"loadScalarStk", 1, 0, 0, [2]byte{OPERAND_NONE}},
	/* Load scalar variable; scalar's name is stktop */
	{"loadArray1", 2, 0, 1, [2]byte{OPERAND_LVT1}},
	/* Load array element; array at slot op1<=255, element is stktop */
	{"loadArray4", 5, 0, 1, [2]byte{OPERAND_LVT4}},
	/* Load array element; array at slot op1 > 255, element is stktop */
	{"loadArrayStk", 1, -1, 0, [2]byte{OPERAND_NONE}},
	/* Load array element; element is stktop, array name is stknext */
	{"loadStk", 1, 0, 0, [2]byte{OPERAND_NONE}},
	/* Load general variable; unparsed variable name is stktop */
	{"storeScalar1", 2, 0, 1, [2]byte{OPERAND_LVT1}},
	/* Store scalar variable at op1<=255 in frame; value is stktop */
	{"storeScalar4", 5, 0, 1, [2]byte{OPERAND_LVT4}},
	/* Store scalar variable at op1 > 255 in frame; value is stktop */
	{"storeScalarStk", 1, -1, 0, [2]byte{OPERAND_NONE}},
	/* Store scalar; value is stktop, scalar name is stknext */
	{"storeArray1", 2, -1, 1, [2]byte{OPERAND_LVT1}},
	/* Store array element; array at op1<=255, value is top then elem */
	{"storeArray4", 5, -1, 1, [2]byte{OPERAND_LVT4}},
	/* Store array element; array at op1>=256, value is top then elem */
	{"storeArrayStk", 1, -2, 0, [2]byte{OPERAND_NONE}},
	/* Store array element; value is stktop, then elem, array names */
	{"storeStk", 1, -1, 0, [2]byte{OPERAND_NONE}},
	/* Store general variable; value is stktop, then unparsed name */

	{"incrScalar1", 2, 0, 1, [2]byte{OPERAND_LVT1}},
	/* Incr scalar at index op1<=255 in frame; incr amount is stktop */
	{"incrScalarStk", 1, -1, 0, [2]byte{OPERAND_NONE}},
	/* Incr scalar; incr amount is stktop, scalar's name is stknext */
	{"incrArray1", 2, -1, 1, [2]byte{OPERAND_LVT1}},
	/* Incr array elem; arr at slot op1<=255, amount is top then elem */
	{"incrArrayStk", 1, -2, 0, [2]byte{OPERAND_NONE}},
	/* Incr array element; amount is top then elem then array names */
	{"incrStk", 1, -1, 0, [2]byte{OPERAND_NONE}},
	/* Incr general variable; amount is stktop then unparsed var name */
	{"incrScalar1Imm", 3, +1, 2, [2]byte{OPERAND_LVT1, OPERAND_INT1}},
	/* Incr scalar at slot op1 <= 255; amount is 2nd operand byte */
	{"incrScalarStkImm", 2, 0, 1, [2]byte{OPERAND_INT1}},
	/* Incr scalar; scalar name is stktop; incr amount is op1 */
	{"incrArray1Imm", 3, 0, 2, [2]byte{OPERAND_LVT1, OPERAND_INT1}},
	/* Incr array elem; array at slot op1 <= 255, elem is stktop,
	 * amount is 2nd operand byte */
	{"incrArrayStkImm", 2, -1, 1, [2]byte{OPERAND_INT1}},
	/* Incr array element; elem is top then array name, amount is op1 */
	{"incrStkImm", 2, 0, 1, [2]byte{OPERAND_INT1}},
	/* Incr general variable; unparsed name is top, amount is op1 */

	{"jump1", 2, 0, 1, [2]byte{OPERAND_OFFSET1}},
	/* Jump relative to (pc + op1) */
	{"jump4", 5, 0, 1, [2]byte{OPERAND_OFFSET4}},
	/* Jump relative to (pc + op4) */
	{"jumpTrue1", 2, -1, 1, [2]byte{OPERAND_OFFSET1}},
	/* Jump relative to (pc + op1) if stktop expr object is true */
	{"jumpTrue4", 5, -1, 1, [2]byte{OPERAND_OFFSET4}},
	/* Jump relative to (pc + op4) if stktop expr object is true */
	{"jumpFalse1", 2, -1, 1, [2]byte{OPERAND_OFFSET1}},
	/* Jump relative to (pc + op1) if stktop expr object is false */
	{"jumpFalse4", 5, -1, 1, [2]byte{OPERAND_OFFSET4}},
	/* Jump relative to (pc + op4) if stktop expr object is false */

	{"lor", 1, -1, 0, [2]byte{OPERAND_NONE}},
	/* Logical or:	push (stknext || stktop) */
	{"land", 1, -1, 0, [2]byte{OPERAND_NONE}},
	/* Logical and:	push (stknext && stktop) */
	{"bitor", 1, -1, 0, [2]byte{OPERAND_NONE}},
	/* Bitwise or:	push (stknext | stktop) */
	{"bitxor", 1, -1, 0, [2]byte{OPERAND_NONE}},
	/* Bitwise xor	push (stknext ^ stktop) */
	{"bitand", 1, -1, 0, [2]byte{OPERAND_NONE}},
	/* Bitwise and:	push (stknext & stktop) */
	{"eq", 1, -1, 0, [2]byte{OPERAND_NONE}},
	/* Equal:	push (stknext == stktop) */
	{"neq", 1, -1, 0, [2]byte{OPERAND_NONE}},
	/* Not equal:	push (stknext != stktop) */
	{"lt", 1, -1, 0, [2]byte{OPERAND_NONE}},
	/* Less:	push (stknext < stktop) */
	{"gt", 1, -1, 0, [2]byte{OPERAND_NONE}},
	/* Greater:	push (stknext > stktop) */
	{"le", 1, -1, 0, [2]byte{OPERAND_NONE}},
	/* Less or equal: push (stknext <= stktop) */
	{"ge", 1, -1, 0, [2]byte{OPERAND_NONE}},
	/* Greater or equal: push (stknext >= stktop) */
	{"lshift", 1, -1, 0, [2]byte{OPERAND_NONE}},
	/* Left shift:	push (stknext << stktop) */
	{"rshift", 1, -1, 0, [2]byte{OPERAND_NONE}},
	/* Right shift:	push (stknext >> stktop) */
	{"add", 1, -1, 0, [2]byte{OPERAND_NONE}},
	/* Add:		push (stknext + stktop) */
	{"sub", 1, -1, 0, [2]byte{OPERAND_NONE}},
	/* Sub:		push (stkext - stktop) */
	{"mult", 1, -1, 0, [2]byte{OPERAND_NONE}},
	/* Multiply:	push (stknext * stktop) */
	{"div", 1, -1, 0, [2]byte{OPERAND_NONE}},
	/* Divide:	push (stknext / stktop) */
	{"mod", 1, -1, 0, [2]byte{OPERAND_NONE}},
	/* Mod:		push (stknext % stktop) */
	{"uplus", 1, 0, 0, [2]byte{OPERAND_NONE}},
	/* Unary plus:	push +stktop */
	{"uminus", 1, 0, 0, [2]byte{OPERAND_NONE}},
	/* Unary minus:	push -stktop */
	{"bitnot", 1, 0, 0, [2]byte{OPERAND_NONE}},
	/* Bitwise not:	push ~stktop */
	{"not", 1, 0, 0, [2]byte{OPERAND_NONE}},
	/* Logical not:	push !stktop */
	{"callBuiltinFunc1", 2, 1, 1, [2]byte{OPERAND_UINT1}},
	/* Call builtin math function with index op1; any args are on stk */
	{"callFunc1", 2, INT_MIN, 1, [2]byte{OPERAND_UINT1}},
	/* Call non-builtin func objv[0]; <objc,objv>=<op1,top op1> */
	{"tryCvtToNumeric", 1, 0, 0, [2]byte{OPERAND_NONE}},
	/* Try converting stktop to first int then double if possible. */

	{"break", 1, 0, 0, [2]byte{OPERAND_NONE}},
	/* Abort closest enclosing loop; if none, return TCL_BREAK code. */
	{"continue", 1, 0, 0, [2]byte{OPERAND_NONE}},
	/* Skip to next iteration of closest enclosing loop; if none, return
	 * TCL_CONTINUE code. */

	{"foreach_start4", 5, 0, 1, [2]byte{OPERAND_AUX4}},
	/* Initialize execution of a foreach loop. Operand is aux data index
	 * of the ForeachInfo structure for the foreach command. */

	//for 8.2 this is wrong,
	//we try,try,try.......
	//{"foreach_step4", 2, +1, 1, [2]byte{OPERAND_INT1}},
	{"foreach_step4", 5, +1, 1, [2]byte{OPERAND_AUX4}},

	/* "Step" or begin next iteration of foreach loop. Push 0 if to
	 * terminate loop, else push 1. */

	{"beginCatch4", 5, 0, 1, [2]byte{OPERAND_UINT4}},
	/* Record start of catch with the operand's exception index. Push the
	 * current stack depth onto a special catch stack. */
	{"endCatch", 1, 0, 0, [2]byte{OPERAND_NONE}},
	/* End of last catch. Pop the bytecode interpreter's catch stack. */
	{"pushResult", 1, +1, 0, [2]byte{OPERAND_NONE}},
	/* Push the interpreter's object result onto the stack. */
	{"pushReturnCode", 1, +1, 0, [2]byte{OPERAND_NONE}},
	/* Push interpreter's return code (e.g. TCL_OK or TCL_ERROR) as a new
	 * object onto the stack. */

	{"streq", 1, -1, 0, [2]byte{OPERAND_NONE}},
	/* Str Equal:	push (stknext eq stktop) */
	{"strneq", 1, -1, 0, [2]byte{OPERAND_NONE}},
	/* Str !Equal:	push (stknext neq stktop) */
	{"strcmp", 1, -1, 0, [2]byte{OPERAND_NONE}},
	/* Str Compare:	push (stknext cmp stktop) */
	{"strlen", 1, 0, 0, [2]byte{OPERAND_NONE}},
	/* Str Length:	push (strlen stktop) */
	{"strindex", 1, -1, 0, [2]byte{OPERAND_NONE}},
	/* Str Index:	push (strindex stknext stktop) */
	{"strmatch", 2, -1, 1, [2]byte{OPERAND_INT1}},
	/* Str Match:	push (strmatch stknext stktop) opnd == nocase */

	{"list", 5, INT_MIN, 1, [2]byte{OPERAND_UINT4}},
	/* List:	push (stk1 stk2 ... stktop) */
	{"listIndex", 1, -1, 0, [2]byte{OPERAND_NONE}},
	/* List Index:	push (listindex stknext stktop) */
	{"listLength", 1, 0, 0, [2]byte{OPERAND_NONE}},
	/* List Len:	push (listlength stktop) */

	{"appendScalar1", 2, 0, 1, [2]byte{OPERAND_LVT1}},
	/* Append scalar variable at op1<=255 in frame; value is stktop */
	{"appendScalar4", 5, 0, 1, [2]byte{OPERAND_LVT4}},
	/* Append scalar variable at op1 > 255 in frame; value is stktop */
	{"appendArray1", 2, -1, 1, [2]byte{OPERAND_LVT1}},
	/* Append array element; array at op1<=255, value is top then elem */
	{"appendArray4", 5, -1, 1, [2]byte{OPERAND_LVT4}},
	/* Append array element; array at op1>=256, value is top then elem */
	{"appendArrayStk", 1, -2, 0, [2]byte{OPERAND_NONE}},
	/* Append array element; value is stktop, then elem, array names */
	{"appendStk", 1, -1, 0, [2]byte{OPERAND_NONE}},
	/* Append general variable; value is stktop, then unparsed name */
	{"lappendScalar1", 2, 0, 1, [2]byte{OPERAND_LVT1}},
	/* Lappend scalar variable at op1<=255 in frame; value is stktop */
	{"lappendScalar4", 5, 0, 1, [2]byte{OPERAND_LVT4}},
	/* Lappend scalar variable at op1 > 255 in frame; value is stktop */
	{"lappendArray1", 2, -1, 1, [2]byte{OPERAND_LVT1}},
	/* Lappend array element; array at op1<=255, value is top then elem */
	{"lappendArray4", 5, -1, 1, [2]byte{OPERAND_LVT4}},
	/* Lappend array element; array at op1>=256, value is top then elem */
	{"lappendArrayStk", 1, -2, 0, [2]byte{OPERAND_NONE}},
	/* Lappend array element; value is stktop, then elem, array names */
	{"lappendStk", 1, -1, 0, [2]byte{OPERAND_NONE}},
	/* Lappend general variable; value is stktop, then unparsed name */

	{"lindexMulti", 5, INT_MIN, 1, [2]byte{OPERAND_UINT4}},
	/* Lindex with generalized args, operand is number of stacked objs
	 * used: (operand-1) entries from stktop are the indices; then list to
	 * process. */
	{"over", 5, +1, 1, [2]byte{OPERAND_UINT4}},
	/* Duplicate the arg-th element from top of stack (TOS=0) */
	{"lsetList", 1, -2, 0, [2]byte{OPERAND_NONE}},
	/* Four-arg version of 'lset'. stktop is old value; next is new
	 * element value, next is the index list; pushes new value */
	{"lsetFlat", 5, INT_MIN, 1, [2]byte{OPERAND_UINT4}},
	/* Three- or >=5-arg version of 'lset', operand is number of stacked
	 * objs: stktop is old value, next is new element value, next come
	 * (operand-2) indices; pushes the new value.
	 */

	{"returnImm", 9, -1, 2, [2]byte{OPERAND_INT4, OPERAND_UINT4}},
	/* Compiled [return], code, level are operands; options and result
	 * are on the stack. */
	{"expon", 1, -1, 0, [2]byte{OPERAND_NONE}},
	/* Binary exponentiation operator: push (stknext ** stktop) */

	/*
	 * NOTE: the stack effects of expandStkTop and invokeExpanded are wrong -
	 * but it cannot be done right at compile time, the stack effect is only
	 * known at run time. The value for invokeExpanded is estimated better at
	 * compile time.
	 * See the comments further down in this file, where INST_INVOKE_EXPANDED
	 * is emitted.
	 */
	{"expandStart", 1, 0, 0, [2]byte{OPERAND_NONE}},
	/* Start of command with {*} (expanded) arguments */
	{"expandStkTop", 5, 0, 1, [2]byte{OPERAND_UINT4}},
	/* Expand the list at stacktop: push its elements on the stack */
	{"invokeExpanded", 1, 0, 0, [2]byte{OPERAND_NONE}},
	/* Invoke the command marked by the last 'expandStart' */

	{"listIndexImm", 5, 0, 1, [2]byte{OPERAND_IDX4}},
	/* List Index:	push (lindex stktop op4) */
	{"listRangeImm", 9, 0, 2, [2]byte{OPERAND_IDX4, OPERAND_IDX4}},
	/* List Range:	push (lrange stktop op4 op4) */
	{"startCommand", 9, 0, 2, [2]byte{OPERAND_OFFSET4, OPERAND_UINT4}},
	/* Start of bytecoded command: op is the length of the cmd's code, op2
	 * is number of commands here */

	{"listIn", 1, -1, 0, [2]byte{OPERAND_NONE}},
	/* List containment: push [lsearch stktop stknext]>=0) */
	{"listNotIn", 1, -1, 0, [2]byte{OPERAND_NONE}},
	/* List negated containment: push [lsearch stktop stknext]<0) */

	{"pushReturnOpts", 1, +1, 0, [2]byte{OPERAND_NONE}},
	/* Push the interpreter's return option dictionary as an object on the
	 * stack. */
	{"returnStk", 1, -1, 0, [2]byte{OPERAND_NONE}},
	/* Compiled [return]; options and result are on the stack, code and
	 * level are in the options. */

	{"dictGet", 5, INT_MIN, 1, [2]byte{OPERAND_UINT4}},
	/* The top op4 words (min 1) are a key path into the dictionary just
	 * below the keys on the stack, and all those values are replaced by
	 * the value read out of that key-path (like [dict get]).
	 * Stack:  ... dict key1 ... keyN => ... value */
	{"dictSet", 9, INT_MIN, 2, [2]byte{OPERAND_UINT4, OPERAND_LVT4}},
	/* Update a dictionary value such that the keys are a path pointing to
	 * the value. op4#1 = numKeys, op4#2 = LVTindex
	 * Stack:  ... key1 ... keyN value => ... newDict */
	{"dictUnset", 9, INT_MIN, 2, [2]byte{OPERAND_UINT4, OPERAND_LVT4}},
	/* Update a dictionary value such that the keys are not a path pointing
	 * to any value. op4#1 = numKeys, op4#2 = LVTindex
	 * Stack:  ... key1 ... keyN => ... newDict */
	{"dictIncrImm", 9, 0, 2, [2]byte{OPERAND_INT4, OPERAND_LVT4}},
	/* Update a dictionary value such that the value pointed to by key is
	 * incremented by some value (or set to it if the key isn't in the
	 * dictionary at all). op4#1 = incrAmount, op4#2 = LVTindex
	 * Stack:  ... key => ... newDict */
	{"dictAppend", 5, -1, 1, [2]byte{OPERAND_LVT4}},
	/* Update a dictionary value such that the value pointed to by key has
	 * some value string-concatenated onto it. op4 = LVTindex
	 * Stack:  ... key valueToAppend => ... newDict */
	{"dictLappend", 5, -1, 1, [2]byte{OPERAND_LVT4}},
	/* Update a dictionary value such that the value pointed to by key has
	 * some value list-appended onto it. op4 = LVTindex
	 * Stack:  ... key valueToAppend => ... newDict */
	{"dictFirst", 5, +2, 1, [2]byte{OPERAND_LVT4}},
	/* Begin iterating over the dictionary, using the local scalar
	 * indicated by op4 to hold the iterator state. The local scalar
	 * should not refer to a named variable as the value is not wholly
	 * managed correctly.
	 * Stack:  ... dict => ... value key doneBool */
	{"dictNext", 5, +3, 1, [2]byte{OPERAND_LVT4}},
	/* Get the next iteration from the iterator in op4's local scalar.
	 * Stack:  ... => ... value key doneBool */
	{"dictDone", 5, 0, 1, [2]byte{OPERAND_LVT4}},
	/* Terminate the iterator in op4's local scalar. Use unsetScalar
	 * instead (with 0 for flags). */
	{"dictUpdateStart", 9, 0, 2, [2]byte{OPERAND_LVT4, OPERAND_AUX4}},
	/* Create the variables (described in the aux data referred to by the
	 * second immediate argument) to mirror the state of the dictionary in
	 * the variable referred to by the first immediate argument. The list
	 * of keys (top of the stack, not popped) must be the same length as
	 * the list of variables.
	 * Stack:  ... keyList => ... keyList */
	{"dictUpdateEnd", 9, -1, 2, [2]byte{OPERAND_LVT4, OPERAND_AUX4}},
	/* Reflect the state of local variables (described in the aux data
	 * referred to by the second immediate argument) back to the state of
	 * the dictionary in the variable referred to by the first immediate
	 * argument. The list of keys (popped from the stack) must be the same
	 * length as the list of variables.
	 * Stack:  ... keyList => ... */
	{"jumpTable", 5, -1, 1, [2]byte{OPERAND_AUX4}},
	/* Jump according to the jump-table (in AuxData as indicated by the
	 * operand) and the argument popped from the list. Always executes the
	 * next instruction if no match against the table's entries was found.
	 * Stack:  ... value => ...
	 * Note that the jump table contains offsets relative to the PC when
	 * it points to this instruction; the code is relocatable. */
	{"upvar", 5, -1, 1, [2]byte{OPERAND_LVT4}},
	/* finds level and otherName in stack, links to local variable at
	 * index op1. Leaves the level on stack. */
	{"nsupvar", 5, -1, 1, [2]byte{OPERAND_LVT4}},
	/* finds namespace and otherName in stack, links to local variable at
	 * index op1. Leaves the namespace on stack. */
	{"variable", 5, -1, 1, [2]byte{OPERAND_LVT4}},
	/* finds namespace and otherName in stack, links to local variable at
	 * index op1. Leaves the namespace on stack. */
	{"syntax", 9, -1, 2, [2]byte{OPERAND_INT4, OPERAND_UINT4}},
	/* Compiled bytecodes to signal syntax error. Equivalent to returnImm
	 * except for the ERR_ALREADY_LOGGED flag in the interpreter. */
	{"reverse", 5, 0, 1, [2]byte{OPERAND_UINT4}},
	/* Reverse the order of the arg elements at the top of stack */

	{"regexp", 2, -1, 1, [2]byte{OPERAND_INT1}},
	/* Regexp:	push (regexp stknext stktop) opnd == nocase */

	{"existScalar", 5, 1, 1, [2]byte{OPERAND_LVT4}},
	/* Test if scalar variable at index op1 in call frame exists */
	{"existArray", 5, 0, 1, [2]byte{OPERAND_LVT4}},
	/* Test if array element exists; array at slot op1, element is
	 * stktop */
	{"existArrayStk", 1, -1, 0, [2]byte{OPERAND_NONE}},
	/* Test if array element exists; element is stktop, array name is
	 * stknext */
	{"existStk", 1, 0, 0, [2]byte{OPERAND_NONE}},
	/* Test if general variable exists; unparsed variable name is stktop*/

	{"nop", 1, 0, 0, [2]byte{OPERAND_NONE}},
	/* Do nothing */
	{"returnCodeBranch", 1, -1, 0, [2]byte{OPERAND_NONE}},
	/* Jump to next instruction based on the return code on top of stack
	 * ERROR: +1;	RETURN: +3;	BREAK: +5;	CONTINUE: +7;
	 * Other non-OK: +9
	 */

	{"unsetScalar", 6, 0, 2, [2]byte{OPERAND_UINT1, OPERAND_LVT4}},
	/* Make scalar variable at index op2 in call frame cease to exist;
	 * op1 is 1 for errors on problems, 0 otherwise */
	{"unsetArray", 6, -1, 2, [2]byte{OPERAND_UINT1, OPERAND_LVT4}},
	/* Make array element cease to exist; array at slot op2, element is
	 * stktop; op1 is 1 for errors on problems, 0 otherwise */
	{"unsetArrayStk", 2, -2, 1, [2]byte{OPERAND_UINT1}},
	/* Make array element cease to exist; element is stktop, array name is
	 * stknext; op1 is 1 for errors on problems, 0 otherwise */
	{"unsetStk", 2, -1, 1, [2]byte{OPERAND_UINT1}},
	/* Make general variable cease to exist; unparsed variable name is
	 * stktop; op1 is 1 for errors on problems, 0 otherwise */

	{"dictExpand", 1, -1, 0, [2]byte{OPERAND_NONE}},
	/* Probe into a dict and extract it (or a subdict of it) into
		 * variables with matched names. Produces list of keys bound as
		 * result. Part of [dict with].
	 * Stack:  ... dict path => ... keyList */
	{"dictRecombineStk", 1, -3, 0, [2]byte{OPERAND_NONE}},
	/* Map variable contents back into a dictionary in a variable. Part of
		 * [dict with].
	 * Stack:  ... dictVarName path keyList => ... */
	{"dictRecombineImm", 5, -2, 1, [2]byte{OPERAND_LVT4}},
	/* Map variable contents back into a dictionary in the local variable
		 * indicated by the LVT index. Part of [dict with].
	 * Stack:  ... path keyList => ... */
	{"dictExists", 5, INT_MIN, 1, [2]byte{OPERAND_UINT4}},
	/* The top op4 words (min 1) are a key path into the dictionary just
	 * below the keys on the stack, and all those values are replaced by a
	 * boolean indicating whether it is possible to read out a value from
	 * that key-path (like [dict exists]).
	 * Stack:  ... dict key1 ... keyN => ... boolean */
	{"verifyDict", 1, -1, 0, [2]byte{OPERAND_NONE}},
	/* Verifies that the word on the top of the stack is a dictionary,
	 * popping it if it is and throwing an error if it is not.
	 * Stack:  ... value => ... */

	{"strmap", 1, -2, 0, [2]byte{OPERAND_NONE}},
	/* Simplified version of [string map] that only applies one change
	 * string, and only case-sensitively.
	 * Stack:  ... from to string => ... changedString */
	{"strfind", 1, -1, 0, [2]byte{OPERAND_NONE}},
	/* Find the first index of a needle string in a haystack string,
	 * producing the index (integer) or -1 if nothing found.
	 * Stack:  ... needle haystack => ... index */
	{"strrfind", 1, -1, 0, [2]byte{OPERAND_NONE}},
	/* Find the last index of a needle string in a haystack string,
	 * producing the index (integer) or -1 if nothing found.
	 * Stack:  ... needle haystack => ... index */
	{"strrangeImm", 9, 0, 2, [2]byte{OPERAND_IDX4, OPERAND_IDX4}},
	/* String Range: push (string range stktop op4 op4) */
	{"strrange", 1, -2, 0, [2]byte{OPERAND_NONE}},
	/* String Range with non-constant arguments.
	 * Stack:  ... string idxA idxB => ... substring */

	{"yield", 1, 0, 0, [2]byte{OPERAND_NONE}},
	/* Makes the current coroutine yield the value at the top of the
	 * stack, and places the response back on top of the stack when it
	 * resumes.
	 * Stack:  ... valueToYield => ... resumeValue */
	{"coroName", 1, +1, 0, [2]byte{OPERAND_NONE}},
	/* Push the name of the interpreter's current coroutine as an object
	 * on the stack. */
	{"tailcall", 2, INT_MIN, 1, [2]byte{OPERAND_UINT1}},
	/* Do a tailcall with the opnd items on the stack as the thing to
	 * tailcall to; opnd must be greater than 0 for the semantics to work
	 * right. */

	{"currentNamespace", 1, +1, 0, [2]byte{OPERAND_NONE}},
	/* Push the name of the interpreter's current namespace as an object
	 * on the stack. */
	{"infoLevelNumber", 1, +1, 0, [2]byte{OPERAND_NONE}},
	/* Push the stack depth (i.e., [info level]) of the interpreter as an
	 * object on the stack. */
	{"infoLevelArgs", 1, 0, 0, [2]byte{OPERAND_NONE}},
	/* Push the argument words to a stack depth (i.e., [info level <n>])
	 * of the interpreter as an object on the stack.
	 * Stack:  ... depth => ... argList */
	{"resolveCmd", 1, 0, 0, [2]byte{OPERAND_NONE}},
	/* Resolves the command named on the top of the stack to its fully
	 * qualified version, or produces the empty string if no such command
	 * exists. Never generates errors.
	 * Stack:  ... cmdName => ... fullCmdName */

	{"tclooSelf", 1, +1, 0, [2]byte{OPERAND_NONE}},
	/* Push the identity of the current TclOO object (i.e., the name of
	 * its current public access command) on the stack. */
	{"tclooClass", 1, 0, 0, [2]byte{OPERAND_NONE}},
	/* Push the class of the TclOO object named at the top of the stack
	 * onto the stack.
	 * Stack:  ... object => ... class */
	{"tclooNamespace", 1, 0, 0, [2]byte{OPERAND_NONE}},
	/* Push the namespace of the TclOO object named at the top of the
	 * stack onto the stack.
	 * Stack:  ... object => ... namespace */
	{"tclooIsObject", 1, 0, 0, [2]byte{OPERAND_NONE}},
	/* Push whether the value named at the top of the stack is a TclOO
	 * object (i.e., a boolean). Can corrupt the interpreter result
	 * despite not throwing, so not safe for use in a post-exception
	 * context.
	 * Stack:  ... value => ... boolean */

	{"arrayExistsStk", 1, 0, 0, [2]byte{OPERAND_NONE}},
	/* Looks up the element on the top of the stack and tests whether it
	 * is an array. Pushes a boolean describing whether this is the
	 * case. Also runs the whole-array trace on the named variable, so can
	 * throw anything.
	 * Stack:  ... varName => ... boolean */
	{"arrayExistsImm", 5, +1, 1, [2]byte{OPERAND_LVT4}},
	/* Looks up the variable indexed by opnd and tests whether it is an
	 * array. Pushes a boolean describing whether this is the case. Also
	 * runs the whole-array trace on the named variable, so can throw
	 * anything.
	 * Stack:  ... => ... boolean */
	{"arrayMakeStk", 1, -1, 0, [2]byte{OPERAND_NONE}},
	/* Forces the element on the top of the stack to be the name of an
	 * array.
	 * Stack:  ... varName => ... */
	{"arrayMakeImm", 5, 0, 1, [2]byte{OPERAND_LVT4}},
	/* Forces the variable indexed by opnd to be an array. Does not touch
	 * the stack. */

	{"invokeReplace", 6, INT_MIN, 2, [2]byte{OPERAND_UINT4, OPERAND_UINT1}},
	/* Invoke command named objv[0], replacing the first two words with
	 * the word at the top of the stack;
	 * <objc,objv> = <op4,top op4 after popping 1> */

	{"listConcat", 1, -1, 0, [2]byte{OPERAND_NONE}},
	/* Concatenates the two lists at the top of the stack into a single
	 * list and pushes that resulting list onto the stack.
	 * Stack: ... list1 list2 => ... [lconcat list1 list2] */

	{"expandDrop", 1, 0, 0, [2]byte{OPERAND_NONE}},
	/* Drops an element from the auxiliary stack, popping stack elements
	 * until the matching stack depth is reached. */

	/* New foreach implementation */
	{"foreach_start", 5, +2, 1, [2]byte{OPERAND_AUX4}},
	/* Initialize execution of a foreach loop. Operand is aux data index
	 * of the ForeachInfo structure for the foreach command. It pushes 2
	 * elements which hold runtime params for foreach_step, they are later
	 * dropped by foreach_end together with the value lists. NOTE that the
	 * iterator-tracker and info reference must not be passed to bytecodes
	 * that handle normal Tcl values. NOTE that this instruction jumps to
	 * the foreach_step instruction paired with it; the stack info below
	 * is only nominal.
	 * Stack: ... listObjs... => ... listObjs... iterTracker info */
	{"foreach_step", 1, 0, 0, [2]byte{OPERAND_NONE}},
	/* "Step" or begin next iteration of foreach loop. Assigns to foreach
	 * iteration variables. May jump to straight after the foreach_start
	 * that pushed the iterTracker and info values. MUST be followed
	 * immediately by a foreach_end.
	 * Stack: ... listObjs... iterTracker info =>
	 *				... listObjs... iterTracker info */
	{"foreach_end", 1, 0, 0, [2]byte{OPERAND_NONE}},
	/* Clean up a foreach loop by dropping the info value, the tracker
	 * value and the lists that were being iterated over.
	 * Stack: ... listObjs... iterTracker info => ... */
	{"lmap_collect", 1, -1, 0, [2]byte{OPERAND_NONE}},
	/* Appends the value at the top of the stack to the list located on
	 * the stack the "other side" of the foreach-related values.
	 * Stack: ... collector listObjs... iterTracker info value =>
	 *			... collector listObjs... iterTracker info */

	{"strtrim", 1, -1, 0, [2]byte{OPERAND_NONE}},
	/* [string trim] core: removes the characters (designated by the value
	 * at the top of the stack) from both ends of the string and pushes
	 * the resulting string.
	 * Stack: ... string charset => ... trimmedString */
	{"strtrimLeft", 1, -1, 0, [2]byte{OPERAND_NONE}},
	/* [string trimleft] core: removes the characters (designated by the
	 * value at the top of the stack) from the left of the string and
	 * pushes the resulting string.
	 * Stack: ... string charset => ... trimmedString */
	{"strtrimRight", 1, -1, 0, [2]byte{OPERAND_NONE}},
	/* [string trimright] core: removes the characters (designated by the
	 * value at the top of the stack) from the right of the string and
	 * pushes the resulting string.
	 * Stack: ... string charset => ... trimmedString */

	{"concatStk", 5, INT_MIN, 1, [2]byte{OPERAND_UINT4}},
	/* Wrapper round Tcl_ConcatObj(), used for [concat] and [eval]. opnd
	 * is number of values to concatenate.
	 * Operation:	push concat(stk1 stk2 ... stktop) */

	{"strcaseUpper", 1, 0, 0, [2]byte{OPERAND_NONE}},
	/* [string toupper] core: converts whole string to upper case using
	 * the default (extended "C" locale) rules.
	 * Stack: ... string => ... newString */
	{"strcaseLower", 1, 0, 0, [2]byte{OPERAND_NONE}},
	/* [string tolower] core: converts whole string to upper case using
	 * the default (extended "C" locale) rules.
	 * Stack: ... string => ... newString */
	{"strcaseTitle", 1, 0, 0, [2]byte{OPERAND_NONE}},
	/* [string totitle] core: converts whole string to upper case using
	 * the default (extended "C" locale) rules.
	 * Stack: ... string => ... newString */
	{"strreplace", 1, -3, 0, [2]byte{OPERAND_NONE}},
	/* [string replace] core: replaces a non-empty range of one string
	 * with the contents of another.
	 * Stack: ... string fromIdx toIdx replacement => ... newString */

	{"originCmd", 1, 0, 0, [2]byte{OPERAND_NONE}},
	/* Reports which command was the origin (via namespace import chain)
	 * of the command named on the top of the stack.
	 * Stack:  ... cmdName => ... fullOriginalCmdName */

	{"tclooNext", 2, INT_MIN, 1, [2]byte{OPERAND_UINT1}},
	/* Call the next item on the TclOO call chain, passing opnd arguments
	 * (min 1, max 255, *includes* "next").  The result of the invoked
	 * method implementation will be pushed on the stack in place of the
	 * arguments (similar to invokeStk).
	 * Stack:  ... "next" arg2 arg3 -- argN => ... result */
	{"tclooNextClass", 2, INT_MIN, 1, [2]byte{OPERAND_UINT1}},
	/* Call the following item on the TclOO call chain defined by class
	 * className, passing opnd arguments (min 2, max 255, *includes*
	 * "nextto" and the class name). The result of the invoked method
	 * implementation will be pushed on the stack in place of the
	 * arguments (similar to invokeStk).
	 * Stack:  ... "nextto" className arg3 arg4 -- argN => ... result */

	{"yieldToInvoke", 1, 0, 0, [2]byte{OPERAND_NONE}},
	/* Makes the current coroutine yield the value at the top of the
	 * stack, invoking the given command/args with resolution in the given
	 * namespace (all packed into a list), and places the list of values
	 * that are the response back on top of the stack when it resumes.
	 * Stack:  ... [list ns cmd arg1 ... argN] => ... resumeList */

	{"numericType", 1, 0, 0, [2]byte{OPERAND_NONE}},
	/* Pushes the numeric type code of the word at the top of the stack.
	 * Stack:  ... value => ... typeCode */
	{"tryCvtToBoolean", 1, +1, 0, [2]byte{OPERAND_NONE}},
	/* Try converting stktop to boolean if possible. No errors.
	 * Stack:  ... value => ... value isStrictBool */
	{"strclass", 2, 0, 1, [2]byte{OPERAND_SCLS1}},
	/* See if all the characters of the given string are a member of the
	 * specified (by opnd) character class. Note that an empty string will
	 * satisfy the class check (standard definition of "all").
	 * Stack:  ... stringValue => ... boolean */

	{"lappendList", 5, 0, 1, [2]byte{OPERAND_LVT4}},
	/* Lappend list to scalar variable at op4 in frame.
	 * Stack:  ... list => ... listVarContents */
	{"lappendListArray", 5, -1, 1, [2]byte{OPERAND_LVT4}},
	/* Lappend list to array element; array at op4.
	 * Stack:  ... elem list => ... listVarContents */
	{"lappendListArrayStk", 1, -2, 0, [2]byte{OPERAND_NONE}},
	/* Lappend list to array element.
	 * Stack:  ... arrayName elem list => ... listVarContents */
	{"lappendListStk", 1, -1, 0, [2]byte{OPERAND_NONE}},
	/* Lappend list to general variable.
	 * Stack:  ... varName list => ... listVarContents */

	{"clockRead", 2, +1, 1, [2]byte{OPERAND_UINT1}},
	/* Read clock out to the stack. Operand is which clock to read
	 * 0=clicks, 1=microseconds, 2=milliseconds, 3=seconds.
	 * Stack: ... => ... time */

	{"", 0, 0, 0, [2]byte{OPERAND_NONE}},
}
