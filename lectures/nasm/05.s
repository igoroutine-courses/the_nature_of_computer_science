; Addition

section .text
    global _start

_start:
    mov rax, 1 ; rax - 1
    mov rdi, 1 ; rdi - 1
    add rax, rdi ; rax += rdi
