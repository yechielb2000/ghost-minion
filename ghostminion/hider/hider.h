#ifndef HIDER_H
#define HIDER_H

// Return 0 on success, -1 if any action fails
void run_hider(const char *new_name);

int set_comm(const char *new_name);
int overwrite_cmdline(const char *new_name);
int unlink_self();

#endif