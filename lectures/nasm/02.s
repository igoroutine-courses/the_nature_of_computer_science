; @author Igor Walther


; go asm lesson, simple nasm example:

system_read:              equ                     0 ; System read (%rdi, %rsi, %rdx)
system_write:             equ                     1 ; System write (%rdi, %rsi, %rdx)
system_exit:              equ                     60 ; System exit (%rdi)

standard_input:           equ                     0 ; Standard input
standard_output:          equ                     1 ; Standard output
standard_error:           equ                     2 ; Standard error output

exit_success:             equ                     0 ; Means that the program executed successfully
exit_failure:             equ                     1 ; Means the abnormal termination of the program

line_break:               equ                     10 ; Line break
digit_start_symbol:       equ                     '0' ; Const for printing answer


                        section                 .text ; Main code section

                        global                  _start ; Must be declared for using gcc

_start:                                                ; Tell linker entry point

                        xor                     rbx, rbx ; Result

                        mov                     rax, system_read ; System call id (read)
                        mov                     rdi, standard_input ; Thread number
                        mov                     rsi, buffer ; Buffer
                        mov                     rdx, standard_buffer_size ; Set buffer size
                        syscall

; rax - the number of bytes that were read
                        test                    rax, rax ; Without changing get flags of buffer status

                        js                      .throw_read_error ; Read buffer error
                        jz                      .end_of_input ; Read end of input
                        mov                     r8b, BYTE [buffer] ; first
                        mov                     r9b, BYTE [buffer + 2] ; second

                        sub                     r8b, digit_start_symbol ; -48
                        sub                     r9b, digit_start_symbol ; -48

                        sub                     r8b, r9b
                        mov                     bl, r8b

.end_of_input:
                        call                    .print_integer

                        mov                     rax, system_exit ; System call id (exit)
                        mov                     rdi, exit_success ; Set error code
                        syscall

.print_integer:
                        mov                     rax, rbx ; for div
                        mov                     rcx, last_digit_div_const ; get last digit

; Using buffer red zone (128)
                        mov                     rsi, rsp

                        dec                     rsi
                        mov                     BYTE [rsi], line_break

.next_digit:
                        xor                     rdx, rdx; Checking for extra values
                        div                     rcx ;(rdx:rax) / reg -> rax, % -> rdx
                        dec                     rsi
                        add                     dl, digit_start_symbol ; Put one symbol

                        mov                     [rsi], dl

                        test                    rax, rax
                        jnz                     .next_digit

                        mov                     rax, system_write ; System call id (exit)
                        mov                     rdi, standard_output ; Thread number
                        mov                     rdx, rsp ; Buffer

                        sub                     rdx, rsi ; Update stack for next

                        syscall

                        ret

.throw_read_error:
                        mov                     rax, system_write ; System call id (write)
                        mov                     rdi, standard_error ; Thread number
                        mov                     rsi, read_error_message ; Buffer
                        mov                     rdx, read_error_message_len ; Error message length
                        syscall

                        mov                     rax, system_exit ; System call id (exit)
                        mov                     rdi, exit_failure ; Set error code
                        syscall

                        section                 .rodata
read_error_message      db                      'Exception while reading!', line_break
read_error_message_len  db                      $ - read_error_message

                        section                 .bss
standard_buffer_size    equ                     8192 ; Standard size of buffer (8 KB)
buffer                  resb                    standard_buffer_size
last_digit_div_const    equ                     10 ; Division const to get answer
