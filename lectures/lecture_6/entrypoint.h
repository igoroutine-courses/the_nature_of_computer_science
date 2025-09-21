#include <linux/fs.h>
#include <linux/fs_context.h>
#include <linux/init.h>
#include <linux/kernel.h>
#include <linux/module.h>

int networkfs_init(void);
void networkfs_exit(void);
int networkfs_init_fs_context(struct fs_context *fc);
struct inode *networkfs_get_inode(struct super_block *sb,
                                  const struct inode *parent, umode_t mode,
                                  int i_ino);
int networkfs_fill_super(struct super_block *sb, struct fs_context *fc);
int networkfs_get_tree(struct fs_context *fc);
void networkfs_kill_sb(struct super_block *sb);

struct dentry *networkfs_lookup(struct inode *parent, struct dentry *child,
                                unsigned int flag);
int networkfs_create(struct user_namespace *un, struct inode *parent_inode,
                     struct dentry *child_dentry, umode_t mode, bool b);
int networkfs_unlink(struct inode *parent_inode, struct dentry *child_dentry);
int networkfs_mkdir(struct user_namespace *un, struct inode *parent_inode,
                    struct dentry *child_dentry, umode_t mode);
int networkfs_rmdir(struct inode *parent_inode, struct dentry *child_dentry);
int networkfs_link(struct dentry *old_dentry, struct inode *parent_dir,
                   struct dentry *new_dentry);

int networkfs_iterate(struct file *filp, struct dir_context *ctx);

char *prepare_param(const char *from, char *to, unsigned long size);
