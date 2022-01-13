# IBM_COMPUTE_VM_INSTANCE

## Assumptions

1. Only Public Multi Tenant Servers are considered(from the UI), from the rate card, Transient items are ignored(Dedicated instance cost will be nil)
2. Cost is fetched from softlayer infrastructure rate card
3. Flavour based configuration is handled seperately, Cost is aggregate of block storage, processor, ram, network speed, and os.
4. Supports Monthly Billing
5. Cost Generated in USD Currency
