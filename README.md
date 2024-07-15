# 抽卡概率算法

## 参考
- [原神抽卡全机制总结](https://www.bilibili.com/read/cv10468091/)
- [原神抽卡概率工具表](https://www.bilibili.com/read/cv12616453/)

## 功能
- 单抽/连抽
- 保底设计
- 概率测试
- 奖池配置

## 配置表

- 物品表(Item)

| 物品ID   | 类型     |  物品名称  |
| ------- | -------- | -------- |
| ItemId  | SortType | ItemName |
| int     | int      | string   |

- 掉落表(Drop)

| 掉落组   | 分组权重   | 结果     | 是否结束 |
| ------- | -------- | -------- | ------ |
| DropId  | Weight   | Result   | IsEnd  |
| int     | int      | int      | int    |
