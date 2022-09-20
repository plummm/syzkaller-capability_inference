#include <sys/ioctl.h>
#include <stdio.h>
#include <string.h>
#include <stdlib.h>

#define IOCTL_MODULE_MONITOR 0x37777

#define START_MODULE_MONITOR 0x1
#define STOP_MODULE_MONITOR 0x2
#define ADD_MODULE 0x3
#define REMOVE_MODULE 0x4
#define SET_FLAG 0x5
#define CLEAR_LIST 0x6
#define SHOW_MODULE_LIST 0x7
#define ENABLE_DEBUG 0x77
#define DISABLE_DEBUG 0x78

struct module_monitor_msg {
    int op;
    char *module_name;
    size_t module_name_len;
    unsigned int flag;
};

int compare(char *a, char *b) {
    int len = 0;
    if (strlen(a) > strlen(b))
        len = strlen(b);
    else
        len = strlen(a);
    
    return strncmp(a, b, len) == 0;
}

int main(int argc, char *argv[])
{
    if (argc != 4 && argc != 2) {
        puts("usage: ./add_known_module start|stop|clear|show|add [module_name flag]\n");
        return 0;
    }

    struct module_monitor_msg *msg = malloc(sizeof(struct module_monitor_msg));
    memset(msg, 0, sizeof(struct module_monitor_msg));

    if (argc == 2) {
        if (compare(argv[1], "start")) {
            msg->op = START_MODULE_MONITOR;
            ioctl(0, IOCTL_MODULE_MONITOR, msg);
            //ioctl(0, IOCTL_MODULE_MONITOR, "MAGIC?!START\x00");
            return 0;
        }
        if (compare(argv[1], "stop")) {
            msg->op = STOP_MODULE_MONITOR;
            ioctl(0, IOCTL_MODULE_MONITOR, msg);
            //ioctl(0, IOCTL_MODULE_MONITOR, "MAGIC?!STOP\x00");
            return 0;
        } 
        if (compare(argv[1], "clear")) {
            msg->op = CLEAR_LIST;
            ioctl(0, IOCTL_MODULE_MONITOR, msg);
            return 0;
        }
        if (compare(argv[1], "show")) {
            msg->op = SHOW_MODULE_LIST;
            ioctl(0, IOCTL_MODULE_MONITOR, msg);
            return 0;
        }
        if (compare(argv[1], "enable-debug")) {
            msg->op = ENABLE_DEBUG;
            ioctl(0, IOCTL_MODULE_MONITOR, msg);
            return 0;
        }
        if (compare(argv[1], "disable-debug")) {
            msg->op = DISABLE_DEBUG;
            ioctl(0, IOCTL_MODULE_MONITOR, msg);
            return 0;
        }
        puts("usage: ./add_known_module start|stop|clear|show|ad|enable-debug|disable-debug [module_name flag]\n");
        return 0;
    }

    if (argc == 4) {
        if (compare(argv[1], "add")) {
            msg->op = ADD_MODULE;
            msg->module_name = argv[2];
            msg->module_name_len = strlen(argv[2]);
            msg->flag = atoi(argv[3]);
            ioctl(0, IOCTL_MODULE_MONITOR, msg);
            return 0;
        }

        puts("usage: ./add_known_module start|stop|clear|show|add [module_name flag]\n");
        return 0;
    }
    return 0;
}