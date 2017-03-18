package hpool

// 用户是怎么使用这个带负载均衡的pool的？
// 方案1 ：用户直接把任务交给 pool
// 方案2 ：用户从 pool 获得一个 worker （heap）, 用户把任务交给 worker （线程安全的map）
// 以下采用方案2