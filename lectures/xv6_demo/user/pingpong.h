#include "kernel/types.h"
#include "user/user.h"

void child_behaviour(int parent_to_child_chan[2], int child_to_parent_chan[2]);

void parent_behaviour(int parent_to_child_chan[2], int child_to_parent_chan[2]);

void process_message(int file_descriptor, int pid);

void write_message(char message[]);

void write_start_message_part(int pid);
