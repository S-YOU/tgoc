package x86

import (
	"fmt"
	"tgoc/ast"
	"tgoc/utils"
)

// Identifier name: offset from bsp
var offsets map[string]int

// The number of stored identifier to stack
var varCount int

// The number of total identifier
var varNum int

func initi(n int) {
	offsets = map[string]int{}
	varCount = 1
	varNum = n
}

func genExpr(expr ast.Node) {

	switch expr := expr.(type) {
	case *ast.IntLit:
		fmt.Printf("	push %d\n", expr.Val)
		return
	case *ast.BinaryExpr:
		genExpr(expr.Lhs)
		genExpr(expr.Rhs)

		fmt.Printf("	pop rdi\n")
		fmt.Printf("	pop rax\n")
		switch expr.Op {
		case "+":
			fmt.Printf("	add rax, rdi\n")
		case "-":
			fmt.Printf("	sub rax, rdi\n")
		case "*":
			fmt.Printf("	mul rdi\n")
		case "/":
			fmt.Printf("    xor rdx, rdx\n")
			fmt.Printf("    div rdi\n")
		case "%":
			fmt.Printf("    xor rdx, rdx\n")
			fmt.Printf("    div rdi\n")
			fmt.Printf("	mov rax, rdx\n")
		case "<<":
			// To change the cl value, changed the rcx value.
			// cl is lower 8 bit register of rcx register.
			fmt.Printf("	mov rcx, rdi\n")
			fmt.Printf("	shl rax, cl\n")
		case ">>":
			fmt.Printf("	mov rcx, rdi\n")
			fmt.Printf("	sar rax, cl\n")
		case "==":
			fmt.Printf("	cmp rax, rdi\n")
			fmt.Printf("	sete al\n")
			fmt.Printf("	movzx rax, al\n")
		case "!=":
			fmt.Printf("	cmp rax, rdi\n")
			fmt.Printf("	sete al\n")
			fmt.Printf("	movzx rax, al\n")
			// 0000 => 0001, 0001 => 0000
			fmt.Printf("	xor rax, 1\n")
		case "<":
			fmt.Printf("	cmp rax, rdi\n")
			fmt.Printf("	setl al\n")
			fmt.Printf("	movzx rax, al\n")
		case ">":
			fmt.Printf("	cmp rax, rdi\n")
			fmt.Printf("	setg al\n")
			fmt.Printf("	movzx rax, al\n")
		}
		fmt.Printf("	push rax\n")
		return

	case *ast.UnaryExpr:
		genExpr(expr.Expr)
		fmt.Printf("    pop rax\n")
		// For now there is only one unary operator; '-' ,
		// so I only have to invert sign.
		fmt.Printf("	neg rax\n")
		fmt.Printf("	push rax \n")
		return

	case *ast.Ident:
		os, ok := offsets[expr.Name]
		utils.Assert(ok, "undefined identifier")

		fmt.Printf("	mov rax, QWORD PTR [rbp - %d]\n", 8*os)
		fmt.Printf("	push rax\n")
	}
}

func genDecl(decl ast.Decl) {
	svd, _ := decl.(*ast.SVDecl)
	genExpr(svd.Val)
	fmt.Printf("	pop rax\n")
	fmt.Printf("	mov QWORD PTR [rbp - %d], rax\n", 8*varCount)
	offsets[svd.Name] = varCount
	varCount++
}

func genStmts(stmts []ast.Stmt) {
	fmt.Printf("	sub rsp, %d\n", varNum*8)

	for _, stmt := range stmts {
		switch stmt := stmt.(type) {
		case *ast.ExprStmt:
			genExpr(stmt.Expr)
			fmt.Printf("	pop rax\n")
		case *ast.AssignStmt:
			genDecl(stmt.Decl)
		case *ast.ReturnStmt:
			genExpr(stmt.Expr)
			fmt.Printf("	pop rax\n")
			fmt.Printf("	mov rsp, rbp\n")
			fmt.Printf("	pop rbp\n")
			fmt.Printf("	ret\n")
			return
		}
	}
}

func Gen(stmts []ast.Stmt, varNum int) {
	initi(varNum)

	fmt.Printf(".intel_syntax noprefix\n")
	fmt.Printf(".globl _main\n")
	fmt.Printf("_main:\n")
	fmt.Printf("	push rbp\n")
	fmt.Printf("	mov rbp, rsp\n")

	genStmts(stmts)

	fmt.Printf("	mov rsp, rbp\n")
	fmt.Printf("	pop rbp\n")
	fmt.Printf("	ret\n")
}
