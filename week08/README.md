1、使用 redis benchmark 工具, 测试 10 20 50 100 200 1k 5k 字节 value 大小，redis get set 性能。

redis-benchmark -h 127.0.0.1 -p 6379 -c 500 -n 10000 -d 10

|     | GET (requests / s)| SET (requests / s) |
|  ----  | ---- | ---- |
| 10  |  27551.72 |  27002.83|
| 20  | 32382.54 | 26315.53|
| 50  | 35460.26 | 27322.12|
| 100  | 35714.26 | 27421.20|
| 200  | 35087.80 | 29411.65|
| 1k  | 27332.28 | 26109.62|
| 5k  | 24570.02 | 23094.69|


2、写入一定量的 kv 数据, 根据数据大小 1w-50w 自己评估, 结合写入前后的 info memory 信息  , 分析上述不同 value 大小下，平均每个 key 的占用内存空间。

|     | used_memory| used_memory_rss | used_memory_peak |avg |
| ---- | ---- | ---- | ----| ---- |
| 10  | 10019304 |  9978288 | 10019304 |92.9|
| 20  | 12419304 | 12378272 | 12419304 |116.9|
| 50  | 15619304 | 15557400 | 15619304 |148.9|
| 100 | 21219304 | 21157400 | 21219304 |204.9|
| 200 | 32419304 | 32357400 | 32419304 |316.9|
| 1k  | 112419304 | 112378272 | 112419304 |1116.9|
| 5k  | 522019304 | 521978272 | 522019304 |5212.9|

写入10w条数据，插入前内存使用量为722152。

当保存的数据中包含字符时，String 类型就会用简单动态字符串（Simple Dynamic String，SDS）结构体来保存。结构如下：
```
struct sdshdr{
    int len; //占 4 个字节，表示 buf 的已用长度。
    int alloc; //也占个 4 字节，表示 buf 的实际分配长度，一般大于 len。
    char buf[]; //字节数组，保存实际数据。
};
```
另外，对于 String 类型来说，除了 SDS 的额外开销，还有一个来自于 RedisObject 结构体的开销。

```
typedef struct redisObject {
    unsigned type:4;      // 类型
    unsigned encoding:4;  // 编码方式
    unsigned lru:REDIS_LRU_BITS;  // LRU 时间
    int refcount;         // 引用计数
    void *ptr;            // 指向对象的值
} robj;
```

因为 Redis 的数据类型有很多，而且，不同数据类型都有些相同的元数据要记录（比如最后一次访问的时间、被引用的次数等），所以，Redis 会用一个 RedisObject 结构体来统一记录这些元数据，同时指向实际数据。

为了节省内存空间，Redis 还对 Long 类型整数和 SDS 的内存布局做了专门的设计：

1. 当保存的是 Long 类型整数时，RedisObject 中的指针就直接赋值为整数数据了，这样就不用额外的指针再指向整数了，节省了指针的空间开销。
2. 当保存的是字符串数据，并且字符串小于等于 44 字节时，RedisObject 中的元数据、指针和 SDS 是一块连续的内存区域，这样就可以避免内存碎片。
3. 当字符串大于 44 字节时，SDS 的数据量就开始变多了，Redis 就不再把 SDS 和 RedisObject 布局在一起了，而是会给 SDS 分配独立的空间，并用指针指向 SDS 结构。

哈希表的每一项是一个 dictEntry 的结构体，用来指向一个键值对。dictEntry 结构中有三个 8 字节的指针，分别指向 key、value 以及下一个 dictEntry:

```
typedef struct dictEntry {
    void *key;  // 键
    union {
        void *val;
        uint64_t u64;
        int64_t s64;
    } v;     // 值
    struct dictEntry *next; // 指向下个哈希表节点

} dictEntry;
```


例如 ：key的size是2，value的size是50。

Key的内存布局：
```
typedef struct redisObject {
    元数据   // 8字节
    指针     // 8字节
    len      // 4字节            --|
    alloc    // 4字节              |--> SDS
    buf[3]   // ("k1\0") 3字节   --|
} robj;

sizeof(robj) = 28
```

Value的内存布局：
```
typedef struct redisObject {
    元数据   // 8字节
    指针     // 8字节
} robj;

struct sdshdr{
    len;      // 4字节
    alloc;    // 4字节
    buf[51];  // 51字节
};

sizeof(robj) + sizeof(sdshdr) = 16 + 64 = 70

```

key和value的总大小为 28 + 70 = 98
