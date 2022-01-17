## Common asumptions taken 

1. Cost is fetched from BSS api for all vpc infrastructure resources
2. Cost fetched from global catalog apis for all other services
3. Cost is fetched from custom rate card(fallbackDB) for all other services in case resource not supported by globalcatalog/bss apis
4. Cost calculated only for us-geo region



## Resource respective asumptions

### IBM_COMPUTE_VM_INSTANCE

1. Only Public Multi Tenant Servers are considered(from the UI), from the rate card, Transient items are ignored(Dedicated instance cost will be nil)
2. Cost is fetched from softlayer infrastructure rate card
3. Flavour based configuration is handled seperately, Cost is aggregate of block storage, processor, ram, network speed, and os.
4. Supports Monthly Billing

### IBM_SERVICE_INSTANCE

1. Support for cloudant, databases-for-mongodb, databases-for-etcd, databases-for-postgresql, databases-for-redis, databases-for-elasticsearch, messages-for-rabbitmq, databases-for-cassandra, databases-for-enterprisedb, blockchain, data-virtualization, functions, satellite-iaas,and schematics is available.

### IBM_RESOURCE_INSTANCE

1. Support for cloudant, databases-for-mongodb, databases-for-etcd, databases-for-postgresql, databases-for-redis, databases-for-elasticsearch, messages-for-rabbitmq, databases-for-cassandra, databases-for-enterprisedb, blockchain, data-virtualization, functions, satellite-iaas,and schematics is available.

### IBM_IS_VOLUME

1. Volume parameter is taken as 100

### IBM_DATABASE

1. Cost is calculated based on members disk allocation, cores and ram or node disk allocation on ram depending on terraform inputs.

### IBM_CONTAINER_CLUSTER

1. Integrations along with the instances are excluded to compute the cost
2. Region is US South Dallas 10
3. Single zone
4. Number of worker nodes/ worker pool count is taken into consideration while computing cost