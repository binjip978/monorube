// This file is part of www.nand2tetris.org
// and the book "The Elements of Computing Systems"
// by Nisan and Schocken, MIT Press.
// File name: projects/04/Fill.asm

// Runs an infinite loop that listens to the keyboard input.
// When a key is pressed (any key), the program blackens the screen,
// i.e. writes "black" in every pixel;
// the screen should remain fully black as long as the key is pressed. 
// When no key is pressed, the program clears the screen, i.e. writes
// "white" in every pixel;
// the screen should remain fully clear as long as no key is pressed.

// Put your code here.
// number that we should traverse
	@8192
	D=A
	@R0
	M=D
(LOOP)
	@KBD
	D=M
	@CLEAR
	D;JEQ
// print ones 
	@n
	M=0
(BLOOP)
	@n
	D=M
	@R0
	D=D-M
	@LOOP
	D;JEQ
	@SCREEN
	D=A
	@n
	A=D+M
	M=-1
	@n
	M=M+1
	@BLOOP
	0;JMP
(CLEAR)
// print zeros
	@n
	M=0
(WLOOP)
	@n
	D=M
	@R0
	D=D-M
	@LOOP
	D;JEQ
	@SCREEN
	D=A
	@n
	A=D+M
	M=0
	@n
	M=M+1
	@WLOOP
	0;JMP