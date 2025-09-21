#include "entrypoint.h"

#include "http.h"

// sudo insmod path
// sudo rmmod path

MODULE_LICENSE("GPL");
MODULE_AUTHOR("Igor Panasyuk");
MODULE_VERSION("1.00");

#define RESPONSE_BUFFER_SIZE 8192
#define ROOT_NODE_ID 1000

char response_buffer[RESPONSE_BUFFER_SIZE];
char name_buffer[1024];

struct entries {
  size_t entries_count;
  struct entry {
    unsigned char entry_type;
    ino_t ino;
    char name[256];
  } entries[16];
};

struct entry_info {
  unsigned char entry_type;
  ino_t ino;
};

struct file_system_type networkfs_fs_type = {
    .name = "networkfs",
    .kill_sb = networkfs_kill_sb,
    .init_fs_context = networkfs_init_fs_context};

struct file_operations networkfs_dir_ops = {
    .iterate = networkfs_iterate,
};

struct inode_operations networkfs_inode_ops = {.lookup = networkfs_lookup,
                                               .create = networkfs_create,
                                               .unlink = networkfs_unlink,
                                               .mkdir = networkfs_mkdir,
                                               .rmdir = networkfs_rmdir,
                                               .link = networkfs_link};

struct fs_context_operations networkfs_context_ops = {.get_tree =
                                                          networkfs_get_tree};

int networkfs_init(void) { return register_filesystem(&networkfs_fs_type); }

void networkfs_exit(void) {
  int err = unregister_filesystem(&networkfs_fs_type);

  if (err == 0) {
    printk(KERN_INFO "Can not unregister networkfs, error - %d", err);
  } else {
    printk(KERN_INFO "Successfully unregister networkfs");
  }
}

int networkfs_init_fs_context(struct fs_context *fc) {
  fc->ops = &networkfs_context_ops;

  return 0;
}

struct inode *networkfs_get_inode(struct super_block *sb,
                                  const struct inode *parent, umode_t mode,
                                  int i_ino) {
  struct inode *inode;
  inode = new_inode(sb);

  if (inode != NULL) {
    inode->i_ino = i_ino;
    inode->i_op = &networkfs_inode_ops;
    inode->i_fop = &networkfs_dir_ops;

    inode_init_owner(&init_user_ns, inode, parent, mode);
  }

  return inode;
}

int networkfs_fill_super(struct super_block *sb, struct fs_context *fc) {
  struct inode *inode =
      networkfs_get_inode(sb, NULL, S_IFDIR | S_IRWXUGO, ROOT_NODE_ID);

  sb->s_root = d_make_root(inode);
  sb->s_maxbytes = 512;

  if (sb->s_root == NULL) {
    return -ENOMEM;
  }

  sb->s_fs_info = kmalloc(strlen(fc->source) + 1, GFP_KERNEL);

  if (sb->s_fs_info == NULL) {
    return -ENOMEM;
  }

  strcpy(sb->s_fs_info, fc->source);

  return 0;
}

int networkfs_get_tree(struct fs_context *fc) {
  int err = get_tree_nodev(fc, networkfs_fill_super);

  return err;
}

void networkfs_kill_sb(struct super_block *sb) { kfree(sb->s_fs_info); }

struct dentry *networkfs_lookup(struct inode *parent_inode,
                                struct dentry *child_dentry,
                                unsigned int flag) {
  const char *token = (const char *)parent_inode->i_sb->s_fs_info;
  struct inode *inode;
  char parent[1024];
  char name[1024];
  int response_code;
  struct entry_info *entry_info;

  sprintf(parent, "%lu", parent_inode->i_ino);
  prepare_param(child_dentry->d_name.name, name_buffer,
                strlen(child_dentry->d_name.name));
  sprintf(name, "%s", name_buffer);

  response_code = networkfs_http_call(token, "lookup", response_buffer,
                                      RESPONSE_BUFFER_SIZE, 2, "parent", parent,
                                      "name", name);

  if (response_code) {
    return NULL;
  }

  entry_info = (struct entry_info *)response_buffer;

  if (entry_info->entry_type == DT_REG) {
    inode = networkfs_get_inode(parent_inode->i_sb, NULL, S_IFREG | S_IRWXUGO,
                                entry_info->ino);
    d_add(child_dentry, inode);
  } else if (entry_info->entry_type == DT_DIR) {
    inode = networkfs_get_inode(parent_inode->i_sb, NULL, S_IFDIR | S_IRWXUGO,
                                entry_info->ino);
    d_add(child_dentry, inode);
  }

  return NULL;
}

int networkfs_create(struct user_namespace *un, struct inode *parent_inode,
                     struct dentry *child_dentry, umode_t mode, bool b) {
  const char *token = (const char *)parent_inode->i_sb->s_fs_info;
  struct inode *inode;
  char parent[1024];
  char name[1024];
  char type[1024];
  int response_code;
  ino_t *ino;

  sprintf(parent, "%lu", parent_inode->i_ino);
  prepare_param(child_dentry->d_name.name, name_buffer,
                strlen(child_dentry->d_name.name));

  sprintf(name, "%s", name_buffer);
  sprintf(type, "%s", "file");

  response_code = networkfs_http_call(token, "create", response_buffer,
                                      RESPONSE_BUFFER_SIZE, 3, "parent", parent,
                                      "name", name, "type", type);

  if (response_code) {
    return -1;
  }

  ino = (ino_t *)response_buffer;

  inode =
      networkfs_get_inode(parent_inode->i_sb, NULL, S_IFREG | S_IRWXUGO, *ino);

  d_add(child_dentry, inode);

  return 0;
}

int networkfs_unlink(struct inode *parent_inode, struct dentry *child_dentry) {
  const char *token = (const char *)parent_inode->i_sb->s_fs_info;
  char parent[1024];
  char name[1024];
  int response_code;

  sprintf(parent, "%lu", parent_inode->i_ino);
  prepare_param(child_dentry->d_name.name, name_buffer,
                strlen(child_dentry->d_name.name));
  sprintf(name, "%s", name_buffer);

  response_code = networkfs_http_call(token, "unlink", response_buffer,
                                      RESPONSE_BUFFER_SIZE, 2, "parent", parent,
                                      "name", name);

  if (response_code) {
    return -1;
  }

  return 0;
}

int networkfs_mkdir(struct user_namespace *un, struct inode *parent_inode,
                    struct dentry *child_dentry, umode_t mode) {
  const char *token = (const char *)parent_inode->i_sb->s_fs_info;
  struct inode *inode;
  char parent[1024];
  char name[1024];
  char type[1024];
  int response_code;
  ino_t *ino;

  sprintf(parent, "%lu", parent_inode->i_ino);
  prepare_param(child_dentry->d_name.name, name_buffer,
                strlen(child_dentry->d_name.name));

  sprintf(name, "%s", name_buffer);
  sprintf(type, "%s", "directory");

  response_code = networkfs_http_call(token, "create", response_buffer,
                                      RESPONSE_BUFFER_SIZE, 3, "parent", parent,
                                      "name", name, "type", type);

  if (response_code) {
    return -1;
  }

  ino = (ino_t *)response_buffer;

  inode =
      networkfs_get_inode(parent_inode->i_sb, NULL, S_IFDIR | S_IRWXUGO, *ino);

  d_add(child_dentry, inode);

  return 0;
}

int networkfs_rmdir(struct inode *parent_inode, struct dentry *child_dentry) {
  char parent[1024];
  char name[1024];
  int response_code;

  const char *token = (const char *)parent_inode->i_sb->s_fs_info;

  sprintf(parent, "%lu", parent_inode->i_ino);
  prepare_param(child_dentry->d_name.name, name_buffer,
                strlen(child_dentry->d_name.name));
  sprintf(name, "%s", name_buffer);

  response_code =
      networkfs_http_call(token, "rmdir", response_buffer, RESPONSE_BUFFER_SIZE,
                          2, "parent", parent, "name", name);

  if (response_code) {
    return -1;
  }

  return 0;
}

int networkfs_link(struct dentry *old_dentry, struct inode *parent_dir,
                   struct dentry *new_dentry) {
  const char *token = (const char *)parent_dir->i_sb->s_fs_info;

  char source[1024];
  char parent[1024];
  char name[1024];
  int response_code;

  sprintf(source, "%lu", old_dentry->d_inode->i_ino);
  sprintf(parent, "%lu", parent_dir->i_ino);

  prepare_param(new_dentry->d_name.name, name_buffer,
                strlen(new_dentry->d_name.name));
  sprintf(name, "%s", name_buffer);

  response_code =
      networkfs_http_call(token, "link", response_buffer, RESPONSE_BUFFER_SIZE,
                          3, "source", source, "parent", parent, "name", name);

  if (response_code) {
    return -1;
  }

  d_add(new_dentry, parent_dir);

  return 0;
}

int networkfs_iterate(struct file *filp, struct dir_context *ctx) {
  char fsname[1024];
  struct dentry *dentry;
  struct inode *inode;
  unsigned long offset;
  int got;
  unsigned char ftype;
  char ino[1024];
  int response_code;
  struct entries *entries;
  ino_t dino;
  dentry = filp->f_path.dentry;
  inode = dentry->d_inode;

  const char *token =
      (const char *)filp->f_path.dentry->d_inode->i_sb->s_fs_info;

  sprintf(ino, "%lu", inode->i_ino);

  response_code = networkfs_http_call(token, "list", response_buffer,
                                      RESPONSE_BUFFER_SIZE, 1, "inode", ino);
  if (response_code) {
    return -1;
  }

  entries = (struct entries *)response_buffer;
  offset = filp->f_pos;
  got = offset;

  // start ls
  while (got < entries->entries_count + 2) {
    printk(KERN_INFO "got = %d, offset = %ld, count = %ld", got, offset,
           entries->entries_count);

    if (got == 0) {
      strcpy(fsname, ".");
      ftype = DT_DIR;
      dino = inode->i_ino;
    } else if (got == 1) {
      strcpy(fsname, "..");
      ftype = DT_DIR;
      dino = dentry->d_parent->d_inode->i_ino;
    } else {
      strcpy(fsname, entries->entries[got - 2].name);
      ftype = entries->entries[got - 2].entry_type;
      dino = entries->entries[got - 2].ino;
    }

    dir_emit(ctx, fsname, strlen(fsname), dino, ftype);
    got++;
    offset++;

    ctx->pos = offset;
  }

  return entries->entries_count - filp->f_pos + 2;
}

char *prepare_param(const char *from, char *to, unsigned long size) {
  int i;

  for (i = 0; i < size; i++) {
    char temp[10];
    sprintf(temp, "%%%02x", (int)from[i]);
    strcpy(&to[3 * i], temp);
  }

  return to;
}

module_init(networkfs_init);
module_exit(networkfs_exit);
