-- 初始化数据
-- 插入默认管理员账户（密码：admin123，使用 bcrypt 加密）

INSERT INTO `users` (`id`, `created_at`, `updated_at`, `deleted_at`, `username`, `password`, `email`, `nickname`, `role`, `status`) VALUES
(1, NOW(), NOW(), NULL, 'admin', '$2a$10$9bW/AVFAVJZHOQPzYkGx5OYhCK6sEJWq7YqJQ8K4qJZ9Y5Y8Y9Y8Y', 'admin@example.com', '管理员', 'admin', 1);

-- 插入默认分类
INSERT INTO `categories` (`id`, `created_at`, `updated_at`, `deleted_at`, `name`, `slug`, `description`, `parent_id`, `sort_order`, `post_count`) VALUES
(1, NOW(), NOW(), NULL, '技术', 'tech', '技术文章分类', NULL, 1, 0),
(2, NOW(), NOW(), NULL, '生活', 'life', '生活随笔', NULL, 2, 0),
(3, NOW(), NOW(), NULL, '教程', 'tutorial', '教程类文章', 1, 1, 0);

-- 插入默认标签
INSERT INTO `tags` (`id`, `created_at`, `updated_at`, `deleted_at`, `name`, `slug`, `description`, `post_count`) VALUES
(1, NOW(), NOW(), NULL, 'Go', 'go', 'Go 语言相关', 0),
(2, NOW(), NOW(), NULL, 'MySQL', 'mysql', 'MySQL 数据库', 0),
(3, NOW(), NOW(), NULL, '教程', 'tutorial', '教程类', 0);

-- 插入默认站点配置
INSERT INTO `configs` (`id`, `created_at`, `updated_at`, `deleted_at`, `key`, `value`, `type`) VALUES
(1, NOW(), NOW(), NULL, 'site_name', 'Aniya Blog', 'string'),
(2, NOW(), NOW(), NULL, 'site_description', 'Aniya 的个人博客', 'string'),
(3, NOW(), NOW(), NULL, 'site_keywords', '博客，技术，生活', 'string'),
(4, NOW(), NOW(), NULL, 'site_url', 'https://example.com', 'string');
