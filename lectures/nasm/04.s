; Cycle

section .text
    global _start

_start:
    mov rax, 10 ; rax = 10, for i in 10 (since go 1.22)
    mov rdi, 0 ; rdi = 0

    test rax, rax ; only for flag register
    jz .done

    add rdi, 1
    sub rax, 1 ; dec

.done:
; rdi is 10 here
