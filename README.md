# pgds-controller

This is a sample kubernetes controller using an external datastore as source of truth instead of the kubernetes API.
No CRDs, no ETCD.

Notable features:
- REST API as source of truth/state store
- `/watch` endpoint to query newly added resources
- Custom `EventHandler` with prioritized queue for resources that are not reconciled


> [!WARNING]  
> This is still very experimental. I will eventually make a lib for the `EventHandler` and maybe the rest API datastore adapter.
