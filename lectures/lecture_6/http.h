#ifndef NETWORKFS_HTTP
#define NETWORKFS_HTTP

#include <linux/types.h>

#define ESOCKNOCREATE 0x2001
#define ESOCKNOCONNECT 0x2002
#define ESOCKNOMSGSEND 0x2003
#define ESOCKNOMSGRECV 0x2004
#define EHTTPBADCODE 0x2005
#define EHTTPMALFORMED 0x2006
#define EPROTMALFORMED 0x2007

int64_t networkfs_http_call(const char *token, const char *method,
                            char *response_buffer, size_t buffer_size,
                            size_t arg_size, ...);

int connect_to_server(const char *command, int params_count,
                      const char *params[], const char *token,
                      char *output_buf);
#endif
