# threadsafe datastructures

### bplus
* b+tree using []byte for keys and values

### index
* b+tree using int64 for keys, and []byte for values (use as atomic uinx timestamp pk index)

### mockdb
* basic embedded key/val store backed by disk snapshots, used for testing or prototyping

### safemap
* sharded hashmap which performs better under concurrent load
 
### realdb
* production quality in memory k/v/doc database backed by disk, atomic persistence using aof
