---
title: '数据库 ALGORITHM = INSTANT 特性研究过程'
description: '深入探究 MySQL 8.0+ 的 ALGORITHM = INSTANT 特性，解析其如何实现在线 DDL 而不锁表的原理'
author: '小吴同学'
publishDate: '2026-03-26'
updatedDate: '2026-03-26'
tags:
  - MySQL
  - 数据库
  - DDL
  - 性能优化
draft: false
comment: true
---

## 一、背景

在日常工作中添加大表字段是个正常不过的操作，偶然发现团队中同事大量使用 `ALGORITHM = INSTANT` 更新字段。根据固有的理解，字段的更新必然会涉及到表结构的更改，印象中数据库会加入 MDL 锁去保证表数据的一致性。

但是听说在 MySQL 8.0+ 特性中，使用此方法不会导致锁表，因为这是一个在线的 DDL 操作方式。后续带着疑问查询了大量的文章未果，依然没有理解其中的原理，固写下文记录研究过程。

## 二、疑问点

例如下面的 DDL 方式：

```sql
ALTER TABLE users 
ADD COLUMN age INT DEFAULT 0 COMMENT '年龄', 
ALGORITHM=INSTANT;
```

在结尾处使用了 **`ALGORITHM=INSTANT`** 的表述方式，即可避免锁表。

这是什么东西？难不成是传说中的"多级&叠加表"设计？我猜想：数据依然是老数据，但是新表是新表，所以查询引擎层会 append 一些新表结构的字段回去。

上面颠覆了我的认知。接下来的研究验证了我的猜想。

## 三、研究过程

我查询了许多文章的表述，但是最终得到的更多是"MySQL 8.0+ 特性"。显然这样的结论并不能直接说服我们对技术的渴望，偶然看到一篇帖子，恍然大悟。

记录在淘宝的《数据库内核技术 2020 年 3 月》的报刊中，我们发现研究者对这项特性是这样描述的：

> 在实现上，MySQL 并没有在系统表中记录多个版本的 schema，而是非常取巧的扩展了存储格式。在已有的 info bits 区域和新增的字段数量区域记录了 instant column 信息，instant add column 之前的数据不做任何修改，之后的数据按照新格式存储。同时在系统表的 private_data 字段存储了 instant column 的默认值信息。查询时，读出的老记录只需要增加 instant column 默认值，新记录则按照新的存储格式进行解析，做到了新老格式的兼容。当然，这种实现方式带来的限制就是只能顺序加字段。

看完后，总结以下几点：

1. **在原有的数据表结构和新增的字段记录了 instant column 字段信息**
2. **旧数据查询的时候会补充 instant column 默认值，新数据按新的结构存储**

对于大佬的总结，是根据一篇 MySQL 的官网技术 worklog 来进行解读的：

- [数据库 worklog-instant 算法](https://dev.mysql.com/worklog/task/?id=11250)

该技术第五点有个这样的阐述：

> For "old" rows, the default value will be looked up from the new system tables and appended before return to server.

翻译过来就是：**对于"旧"行，将从新系统表中查找默认值并在返回服务器之前附加**

这句话验证了我的猜想：表里存在的旧数据，在数据返回到服务器时将新的字段进行 append 回去，以展示对应的完整数据行。

那也就是我们在操作数据的时候，使用了 `ALGORITHM = INSTANT` 在表里使用了 DDL 的方式操作，那么我们将不会关注元数据的更改情况，也就是不用关心锁表。

但是也提供了 `INPLACE` & `FORCE` 原来的复制表的方式来更新内部表，也就是锁表的方式。

## 四、Instant 算法技术内幕原理

对于 instant 算法，原理如下：

假设本身存在的一张表 `t1`，先存在 `a` 字段，插入一条数据 `x`，此时存在了 `x(1)` 【1 = 当前的列号情况】

现在 `ADD COLUMN` 更新表增加了 `b` 字段，再次插入一条数据 `y`，此时存在数据 `x(1 + NULL) + y(2)` 【此时的 NULL 也可以是 0 或者""，为默认值】。最新数据 `y` 将使用所有列的情况，而 `x` 旧数据只保存了 `a` 列的字段，并没有 `b` 的，所以需要组装默认值 NULL 返回。

再次 `ADD COLUMN` 更新表增加了 `c` 字段，再次插入一条数据 `z`，此时存在的数据 `x(1+ NULL * 2)+y(2 + NULL)+z(3)`。新增了 `c` 字段后，同上步骤，此时 `x` 将 append `NULL * 2` 的组装返回，`y` 将 append `NULL`，而 `z` 会使用所有列的情况。

学习到其他大佬的对 SQL 执行过程输出，观察到在使用 instant 字段后对其表进行数据新增，其新增数据会在 bit 字段设置为 1 来代表数据是 instant 之后的。

因为对于内部执行流程 `rec_set_instant_flag_new` 函数在记录的 Info bits 字段设置 `REC_INFO_INSTANT_FLAG`，会表示这个记录是 instant add column 之后创建的。

## 五、使用的注意点

1. **`ALGORITHM = INSTANT` 不能与 LOCK 子句一起使用**
   
   如果指定了 `INSTANT`，并且 `LOCK=NONE/SHARED/EXCLUSIVE` 如果同时指定了 `LOCK=DEFAULT`，则会引发 `ER_WRONG_USAGE` 错误。但是，如果同时指定了 `LOCK=DEFAULT` 才可以。

2. **ADD COLUMN 可能会立即完成**
   
   因此我们可能不会期望立即进行 ADD COLUMN 去修复损坏的索引。

3. **此工作日志不支持具有全文索引的表**

4. **对于 EXCHANGE PARTITION**
   
   为了简化逻辑，如果分区或要交换的表是 instant 的，那么该操作将被拒绝错误 `ER_PARTITION_EXCHANGE_DIFFERENT_OPTION`。

5. **对于使用它进行索引的创建或者更新，算法会被降级 `ALGORITHM=INPLACE`**
   
   此时涉及到锁表，进行数据文件的物理修改。

### INPLACE 算法原理

1. 扫描表数据，构建索引结构
2. 允许并发 DML 操作（如 INSERT/UPDATE/DELETE），但某些阶段可能短暂阻塞 DDL
3. 最终替换旧表文件，完成索引创建

## 六、总结

`ALGORITHM = INSTANT` 是 MySQL 8.0+ 引入的一项非常巧妙的在线 DDL 特性，它通过以下方式实现不锁表：

- **存储格式扩展**：在已有的 info bits 区域记录 instant column 信息
- **新老数据兼容**：旧数据查询时 append 默认值，新数据按新格式存储
- **系统表记录**：在 private_data 字段存储 instant column 的默认值

这项特性极大地提升了大表结构变更的效率，避免了传统 DDL 操作带来的锁表问题。但在使用时需要注意其限制条件，特别是涉及索引操作时会降级为 INPLACE 算法。

---

**参考资料：**

- [MySQL Worklog #11250 - Instant ADD COLUMN](https://dev.mysql.com/worklog/task/?id=11250)
- 淘宝数据库内核技术月报（2020 年 3 月）
