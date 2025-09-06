#include "pingpong.h"

const int BUFFER_SIZE = 128;

const char* END_LINE = "\n";
const char NULL_BYTE = '\0';

const char* DEFAULT_PARENT_MESSAGE = "ping";
const char* DEFAULT_CHILD_MESSAGE = "pong";

const int DEFAULT_SUCCESS_CODE = 0;
const int DEFAULT_ERROR_CODE = -1;

int main(int argc, char* argv[]) {
  int err;

  int parent_to_child_chan[2];
  int child_to_parent_chan[2];

  err = pipe(parent_to_child_chan);

  if (err != DEFAULT_SUCCESS_CODE) {
    printf("can not create pipe from parent to child");
    exit(err);
  }

  err = pipe(child_to_parent_chan);

  if (err != DEFAULT_SUCCESS_CODE) {
    printf("can not create pipe from child to parent");
    exit(err);
  }

  int fork_result = fork();

  switch (fork_result) {
    case -1:
      printf("can not create child process");
      exit(fork_result);
    case 0:
      child_behaviour(parent_to_child_chan, child_to_parent_chan);
      break;
    default:
      parent_behaviour(parent_to_child_chan, child_to_parent_chan);
  }

  exit(DEFAULT_SUCCESS_CODE);
}

void process_message(int file_descriptor, int pid) {
  uint16 current_buffer_position = 0;
  char received_message_part[BUFFER_SIZE];

  int start_message_flag = 0;
  char last_byte;

  do {
    int got =
        read(file_descriptor, (received_message_part + current_buffer_position),
             BUFFER_SIZE - current_buffer_position);

    if (got == DEFAULT_ERROR_CODE) {
      printf("can not read data from descriptor %d", file_descriptor);
      exit(got);
    }

    if (start_message_flag == 0) {
      write_start_message_part(pid);
      start_message_flag++;
    }

    current_buffer_position = current_buffer_position + got;
    last_byte = received_message_part[current_buffer_position - 1];

    if (current_buffer_position == BUFFER_SIZE) {
      write_message(received_message_part);
      current_buffer_position = 0;
    }
  } while (last_byte != NULL_BYTE);

  if (current_buffer_position != 0) {
    received_message_part[current_buffer_position] = NULL_BYTE;
    write_message(received_message_part);
  }

  printf(END_LINE);
}

void write_start_message_part(int pid) { printf("%d: got ", pid); }

void write_message(char message[]) { printf("%s", message); }

void child_behaviour(int parent_to_child_chan[2], int child_to_parent_chan[2]) {
  int pid = getpid();
  process_message(parent_to_child_chan[0], pid);

  int err = write(child_to_parent_chan[1], DEFAULT_CHILD_MESSAGE,
                  strlen(DEFAULT_PARENT_MESSAGE) + 1);

  if (err == DEFAULT_ERROR_CODE) {
    printf("can not write data in descriptor %d", child_to_parent_chan[1]);
  }
}

void parent_behaviour(int parent_to_child_chan[2],
                      int child_to_parent_chan[2]) {
  int pid = getpid();
  int err = write(parent_to_child_chan[1], DEFAULT_PARENT_MESSAGE,
                  strlen(DEFAULT_PARENT_MESSAGE) + 1);

  if (err == DEFAULT_ERROR_CODE) {
    printf("can not write data in descriptor %d", parent_to_child_chan[1]);
  }

  process_message(child_to_parent_chan[0], pid);
}
