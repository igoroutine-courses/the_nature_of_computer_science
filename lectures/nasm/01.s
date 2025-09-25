; Define system calls:
system_write:             equ                     1 ; System write (%rdi, %rsi, %rdx)
system_exit:              equ                     60 ; System exit (%rdi)

standard_output:          equ                     1 ; Standard output

section .data
    msg db "Go asm digest!", 0ah

section .text
    global _start

_start:
    mov rax, system_write
    mov rdi, standard_output
    mov rsi, msg
    mov rdx, 13
    syscall
    mov rax, 60
    mov rdi, 0
    syscall
